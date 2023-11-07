package gooseplus

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"

	// init driver
	_ "github.com/mattn/go-sqlite3"
)

const (
	testDir      = "testdata"
	testDBDriver = "sqlite3"
	testSQL      = "test_sql"
)

func TestMigrator(t *testing.T) {
	_ = os.RemoveAll(testDir)
	require.NoError(t, os.MkdirAll(testDir, 0700))
	t.Cleanup(func() {
		_ = os.RemoveAll(testDir)
	})

	t.Run("with merged migrations", func(t *testing.T) {
		db, clean := createTestDB(t)
		t.Cleanup(clean)
		ctx := context.Background()

		migrator := New(db,
			WithDialect(testDBDriver),
			WithFS(embedMigrations),
			WithBasePath(testSQL),
			SetEnvironments([]string{ // disable default
				"unittest",
				"unittest2", // test merge
			}))

		require.NoError(t,
			migrator.Migrate(ctx),
		)
		currentDBVersion, err := goose.GetDBVersion(db)
		require.NoError(t, err)
		require.Equal(t, int64(20231230000000), currentDBVersion)

		assertTableExists(t, db, "unittest", true)
		assertTableExists(t, db, "unittest_pre", true)
		assertTableExists(t, db, "unittest_post", true)
		assertTableExists(t, db, "go", true)

		t.Run("with new lot of migrations, with failure and rollback", func(t *testing.T) {
			migrator := New(db,
				WithDialect(testDBDriver),
				WithFS(embedMigrations),
				WithBasePath(testSQL),
				SetEnvironments([]string{ // disable default
					"unittest",
					"unittest2", // test merge
					"unittest3", // test rollback
				}))

			err := migrator.Migrate(ctx)
			require.ErrorContains(t, err, "test failure")

			versionAfterRollback, err := goose.GetDBVersion(db)
			require.NoError(t, err)
			require.Equal(t, currentDBVersion, versionAfterRollback)

			assertTableExists(t, db, "go_failed", false)
		})
	})

	t.Run("with empty env for migrations", func(t *testing.T) {
		db, clean := createTestDB(t)
		t.Cleanup(clean)
		ctx := context.Background()

		migrator := New(db,
			WithDialect(testDBDriver),
			WithFS(embedMigrations),
			WithBasePath(testSQL),
			SetEnvironments([]string{ // disable default
				"unittest",
				"unittest4", // test empty merge
			}))

		require.NoError(t,
			migrator.Migrate(ctx),
		)
	})
}

func createTestDB(t testing.TB) (*sql.DB, func()) {
	tempDB, err := os.MkdirTemp(testDir, "db")
	require.NoError(t, err)
	dir := filepath.Join(tempDB, "unittest.db")

	db, err := sql.Open(testDBDriver, dir)
	require.NoError(t, err)

	return db, func() { _ = os.RemoveAll(tempDB) }
}

func assertTableExists(t testing.TB, db *sql.DB, table string, exists bool) {
	var n int
	row := db.QueryRow(fmt.Sprintf(`SELECT COUNT(1) FROM %s`, table))
	err := row.Scan(&n)
	if exists {
		require.NoError(t, err)

		return
	}

	require.Error(t, err)
}
