package base

import (
	"context"
	"database/sql"
	"time"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(UpLongRunning, DownLongRunning)
}

func UpLongRunning(_ context.Context, _ *sql.Tx) error {
	time.Sleep(5 * time.Second)

	return nil
}

func DownLongRunning(_ context.Context, _ *sql.Tx) error {
	return nil
}
