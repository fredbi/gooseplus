package gooseplus_test

import (
	"context"
	"database/sql"
	"embed"
	"log"
	"os"
	"path/filepath"

	"github.com/fredbi/gooseplus"
	"go.uber.org/zap"

	// registers go migrations for unit tests
	_ "github.com/fredbi/gooseplus/test_sql/unittest"
	_ "github.com/fredbi/gooseplus/test_sql/unittest3"

	// init driver
	_ "github.com/mattn/go-sqlite3"
)

//go:embed test_sql/*/*.sql
//go:embed test_sql/*/*.go
var embedMigrations embed.FS

func ExampleMigrator() {
	const (
		dir            = "exampledata"
		driver         = "sqlite3"
		migrationsRoot = "test_sql"
	)

	if err := os.MkdirAll(dir, 0700); err != nil {
		log.Println(err)

		return
	}

	defer func() {
		_ = os.RemoveAll(dir)
	}()

	tempDB, err := os.MkdirTemp(dir, "db")
	if err != nil {
		log.Println(err)

		return
	}

	db, err := sql.Open("sqlite3", filepath.Join(tempDB, "example.db"))
	if err != nil {
		log.Println(err)

		return
	}

	zlg := zap.NewExample()

	migrator := gooseplus.New(db,
		gooseplus.WithDialect(driver),
		gooseplus.WithFS(embedMigrations),
		gooseplus.WithBasePath(migrationsRoot),
		gooseplus.WithLogger(zlg),
	)

	if err := migrator.Migrate(context.Background()); err != nil {
		log.Println(err)

		return
	}
}
