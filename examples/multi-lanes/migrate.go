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

	parentCtx := context.Background()

	// start fast lane, then move on to regular processing
	fastLaneCtx, cancel := context.WithTimeout(parentCtx, 2*time.Second)
	if err := fastLaneMigration(fastLaneCtx, db); err != nil {
		cancel()

		log.Fatalln(err)
	}

	cancel()

	grp, groupCtx := errgroup.WithContext(parentCtx)
	slowLaneCtx, cancel := context.WithTimeout(groupCtx, 10*time.Minute)
	defer cancel()

	doneMigrating := make(chan struct{})

	// start slow lane in the background, running additional migrations, signal to the main app when it is complete
	grp.Go(func() error {
		if err := slowLaneMigration(slowLaneCtx, db); err != nil {
			return nil
		}

		// signals main that this is done
		log.Printf("background migrations completed")
		close(doneMigrating)
		cancel()

		return nil
	})
	defer func() {
		if err := grp.Wait(); err != nil {
			log.Printf("background migrations failed with: %v", err)
		}
	}()

	// main app polls until the background job is completed
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

func fastLaneMigration(ctx context.Context, db *sql.DB) error {
	migrator := gooseplus.New(
		db,
		gooseplus.WithDialect("sqlite3"),
		gooseplus.WithFS(embedMigrations),
		gooseplus.WithMigrationTimeout(time.Second),
	)

	return migrator.Migrate(ctx)
}

func slowLaneMigration(ctx context.Context, db *sql.DB) error {
	// NOTE: we can run migrations in parallel, but we can't use different goose options
	// for parallel instances, such as WithFS() or WithDialect().
	migrator := gooseplus.New(
		db,
		gooseplus.WithDialect("sqlite3"),
		gooseplus.WithFS(embedMigrations),
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
