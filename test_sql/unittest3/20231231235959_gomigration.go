package unittest

import (
	"context"
	"database/sql"
	"errors"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(Up20231230000000, Down20231230000000)
}

func Up20231230000000(ctx context.Context, tx *sql.Tx) error {
	_, _ = tx.ExecContext(ctx,
		`CREATE TABLE go_failed(x integer)`,
	)

	return errors.New("test failure")
}

func Down20231230000000(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx,
		`DROP TABLE IF EXISTS go_failed`,
	)

	return err
}
