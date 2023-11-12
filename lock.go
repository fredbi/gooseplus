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
	db := m.DB

	if _, err := db.ExecContext(ctx, fmt.Sprintf(
		`CREATE TABLE IF NOT EXISTS %s (status INTEGER NOT NULL DEFAULT 0)`, tableName,
	)); err != nil {
		return err
	}

	cancellable, cancel := context.WithCancel(ctx)
	defer cancel()

	tx, err := db.BeginTx(cancellable, nil)
	if err != nil {
		return err
	}

	if _, err = tx.ExecContext(cancellable, fmt.Sprintf(
		`DELETE FROM %s`, tableName,
	)); err != nil {
		return err
	}

	if _, err = tx.ExecContext(ctx, fmt.Sprintf(
		`INSERT INTO %s(status) VALUES(0)`, tableName,
	)); err != nil {
		return err
	}

	return tx.Commit()
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
			`SELECT status FROM %s LIMIT 1 FOR UPDATE`, tableName,
		),
	)

	var n int

	return tx, row.Scan(&n)
}
