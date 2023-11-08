package main

import (
	"context"
	"database/sql"
	"embed"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fredbi/gooseplus"
	"golang.org/x/sync/errgroup"

	// initialize DB driver
	_ "github.com/mattn/go-sqlite3"

	// register go migrations
	_ "github.com/fredbi/gooseplus/examples/multi-lanes/sql/long-running"
)

//go:embed sql/*/*.sql
//go:embed sql/*/*.go
var embedMigrations embed.FS

func main() {
	db, e := createTestDB()
	if e != nil {
		log.Fatalln(e)
	}

	ctx := context.Background()

	// start fast lane, then move on to regular processing
	if err := fastLaneMigration(ctx, db); err != nil {
		log.Fatalln(err)
	}

	grp, groupCtx := errgroup.WithContext(ctx)

	doneMigrating := make(chan struct{})

	// start slow lane in the background, running additional migrations, signal to the main app when it is complete
	grp.Go(func() error {
		if err := slowLaneMigration(groupCtx, db); err != nil {
			return nil
		}

		// signals main that this is done
		log.Printf("background migrations completed")
		close(doneMigrating)

		return nil
	})
	defer func() {
		if err := grp.Wait(); err != nil {
			log.Printf("background migrations failed with: %v", err)
		}
	}()

	// Main app may poll until the background job is completed.
	//
	// Alternatively, the app can check at any time is the slow migrations are done by checking something like:
	//
	// func DoneWithMigration() bool {
	//   select {
	//    case <- doneMigrating:
	//				return true
	//		default:
	//				return false
	// }
	ticker := time.NewTicker(time.Second)
POLL:
	for {
		select {
		case <-doneMigrating:
			log.Println("app can fully proceed now that long running migrations are passed")

			break POLL
		case <-ticker.C:
			log.Println("app can do some work, polling until long-running migrations are complete")
		}
	}

	log.Println("app exited gracefully")
}

func fastLaneMigration(parentCtx context.Context, db *sql.DB) error {
	ctx, cancel := context.WithTimeout(parentCtx, 2*time.Second)
	defer cancel()

	migrator := gooseplus.New(
		db,
		gooseplus.WithDialect("sqlite3"),
		gooseplus.WithFS(embedMigrations),
		gooseplus.WithMigrationTimeout(time.Second),
	)

	return migrator.Migrate(ctx)
}

func slowLaneMigration(parentCtx context.Context, db *sql.DB) error {
	ctx, cancel := context.WithTimeout(parentCtx, 10*time.Minute)
	defer cancel()

	// NOTE: we can run migrations in parallel, but we can't use different goose options
	// for parallel instances, such as WithFS() or WithDialect().
	migrator := gooseplus.New(
		db,
		gooseplus.WithDialect("sqlite3"),
		gooseplus.WithFS(embedMigrations),
		// We can do that only sequentially (goose maintains a global state).
		// With this option, long-running migrations are versioned separately.
		gooseplus.WithVersionTable("goose_db_long_running"),
		gooseplus.SetEnvironments([]string{"long-running"}),
		gooseplus.WithMigrationTimeout(5*time.Minute),
	)

	return migrator.Migrate(ctx)
}

func createTestDB() (*sql.DB, error) {
	dbDir := filepath.Join("testdata", "db")
	if err := os.MkdirAll(dbDir, 0700); err != nil {
		return nil, err
	}

	dir := filepath.Join(dbDir, "example.db")
	db, err := sql.Open("sqlite3", dir)
	if err != nil {
		return nil, err
	}

	return db, nil
}
