package gooseplus

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

const lockTableSuffix = "_lock"

func (m Migrator) ensureLockTable(ctx context.Context) error {
	tableName := fmt.Sprintf("%s%s", goose.TableName(), lockTableSuffix)

	if _, err := m.DB.ExecContext(ctx, fmt.Sprintf(
		`CREATE TABLE IF NOT EXISTS %s (status INTEGER NOT NULL DEFAULT 0)`, tableName,
	)); err != nil {
		return err
	}

	if _, err := m.DB.ExecContext(ctx, fmt.Sprintf(
		`DELETE FROM %s`, tableName,
	)); err != nil {
		return err
	}

	_, err := m.DB.ExecContext(ctx, fmt.Sprintf(
		`INSERT INTO %s(status) VALUES(0)`, tableName,
	))

	return err
}

func (m Migrator) acquireLock(ctx context.Context) (*sql.Tx, error) {
	db := m.DB
	tableName := fmt.Sprintf("%s%s", goose.TableName(), lockTableSuffix)

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	row := tx.QueryRowContext(ctx,
		fmt.Sprintf(
			`SELECT MAX(status) FROM %s`, tableName,
		),
	)

	var n int

	return tx, row.Scan(&n)
}
