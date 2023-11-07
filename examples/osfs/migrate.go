package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fredbi/gooseplus"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := createTestDB()
	if err != nil {
		log.Fatalln(err)
	}

	migrator := gooseplus.New(
		db,
		gooseplus.WithDialect("sqlite3"),
		gooseplus.WithFS(os.DirFS(".")),
		gooseplus.SetEnvironments(nil), // disable folders: only one location for all migrations
		gooseplus.WithTimeout(5*time.Second),
	)

	ctx := context.Background()
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
