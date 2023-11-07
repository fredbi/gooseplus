package base

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(UpLoadFromFile, DownLoadFromFile)
}

func UpLoadFromFile(ctx context.Context, tx *sql.Tx) error {
	i, err := applyToCsv(ctx, func(ctx context.Context, record []string) error {
		_, err := tx.ExecContext(ctx, `INSERT INTO example VALUES(?,?)`, record[0], record[1])

		return err
	})

	if err != nil {
		return err
	}

	log.Printf("%d values inserted from file", i)

	return nil
}

func DownLoadFromFile(ctx context.Context, tx *sql.Tx) error {
	i, err := applyToCsv(ctx, func(ctx context.Context, record []string) error {
		_, err := tx.ExecContext(ctx, `DELETE FROM example WHERE id = ?`, record[0])

		return err
	})

	if err != nil {
		return err
	}

	log.Printf("%d values deleted from file", i)

	return nil
}

func currentDir() string {
	_, filename, _, _ := runtime.Caller(1)

	return filepath.Dir(filename)
}

func applyToCsv(ctx context.Context, fn func(context.Context, []string) error) (int, error) {
	fd, err := os.Open(filepath.Join(currentDir(), "example.csv"))
	if err != nil {
		return 0, err
	}
	r := csv.NewReader(fd)
	_, err = r.Read() // skip header
	if err == io.EOF {
		return 0, err
	}

	var i int
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}

		if len(record) < 2 {
			return 0, fmt.Errorf("invalid record, missing values: %v", record)
		}

		i++

		if err = fn(ctx, record); err != nil {
			return i, err
		}
	}

	return i, nil
}
