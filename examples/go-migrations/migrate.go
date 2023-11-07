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

	// initialize sqlite3 driver
	_ "github.com/mattn/go-sqlite3"

	// register go migrations
	_ "github.com/fredbi/gooseplus/examples/go-migrations/db/base"
)

//go:embed db/*/*.sql
//go:embed db/*/*.go
var embedMigrations embed.FS

func main() {
	db, err := createTestDB()
	if err != nil {
		log.Fatalln(err)
	}

	migrator := gooseplus.New(
		db,
		gooseplus.WithDialect("sqlite3"),
		gooseplus.WithFS(embedMigrations),
		gooseplus.WithBasePath("db"),
		gooseplus.WithMigrationTimeout(time.Second),
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if err = migrator.Migrate(ctx); err != nil {
		log.Printf("error: %v", err)
	}
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
