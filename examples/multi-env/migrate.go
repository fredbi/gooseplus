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
	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/*/*.sql
var embedMigrations embed.FS

func main() {
	db, err := createTestDB()
	if err != nil {
		log.Fatalln(err)
	}

	// run migrations from default, plus any specified environment
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "test"
	}

	migrator := gooseplus.New(
		db,
		gooseplus.WithDialect("sqlite3"),
		gooseplus.WithFS(embedMigrations),
		gooseplus.WithBasePath("migrations"),
		gooseplus.WithEnvironments(env),
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
