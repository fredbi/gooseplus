package gooseplus

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

var gooseMx sync.Mutex

// Migrator knows how to apply changes (migrations) to a versioned database schema.
//
// By default, the migrator will run migrations from the "base" environment folder.
//
// If extra environments are added using options, the migrator will merge
// the migrations with the other folders corresponding to these environments.
type Migrator struct {
	DB *sql.DB
	options
}

// New migrator with options.
func New(db *sql.DB, opts ...Option) *Migrator {
	return &Migrator{
		DB:      db,
		options: applyOptionsWithDefaults(opts),
	}
}

// Migrate applies a sequence database migrations, provided either as SQL scripts or go migration functions.
//
// Whenever a migration fails, Migrate applies rollbacks back to the initial state before returning an error.
func (m Migrator) Migrate(parentCtx context.Context) error {
	lg := m.logger.With(zap.String("operation", "migratedb"))
	if err := m.setGooseGlobals(lg); err != nil {
		return err
	}

	// leave any context with a deadline set unchanged. Otherwise, apply timeout from options.
	ctx, cancel := m.withOptionalTimeout(parentCtx)
	defer cancel()

	db := m.DB

	lg.Info("applying db migrations")

	initialVersion, err := goose.EnsureDBVersionContext(ctx, db)
	if err != nil {
		return errors.Join(ErrMigrationTable, err)
	}

	mergedMigrations, err := m.mergeMigrations()
	if err != nil {
		lg.Error("could not merge migrations", zap.Error(err))

		return errors.Join(ErrMergeMigrations, err)
	}

	if len(mergedMigrations) == 0 {
		lg.Info("no db migrations to be applied")

		return nil
	}

	var tx *sql.Tx
	defer func() {
		if tx != nil {
			_ = tx.Commit()
		}
	}()

	if m.withGlobalLock {
		if err = m.ensureLockTable(ctx); err != nil {
			lg.Error("could not ensure the lock table exists", zap.Error(err))

			return err
		}

		tx, err = m.acquireLock(ctx)
		if err != nil {
			lg.Error("lock could not be acquired prior to running migrations", zap.Error(err))

			return err
		}
	}

	rollForward := m.rollForwardFunc(db, mergedMigrations)
	rollbackTo := m.rollbackToFunc(db, mergedMigrations)

	if rollForwardErr := rollForward(ctx); rollForwardErr != nil {
		// rollback a failed release back to when the deployment started
		lg.Error("failure during rollforward",
			zap.String("action", "rollbacking to the initial state of deployment"),
			zap.Error(rollForwardErr),
		)

		if rollBackErr := rollbackTo(ctx, initialVersion); rollBackErr != nil {
			lg.Error("encountered again an error while rollbacking",
				zap.String("action", "bailed"),
				zap.String("status", "this might leave your database in an inconsistent state"),
				zap.Error(rollBackErr),
			)

			return errors.Join(
				ErrRollBack,
				fmt.Errorf("rollback error: %w", rollBackErr),
				fmt.Errorf("rollforward error: %w", rollForwardErr),
			)
		}

		return errors.Join(
			ErrRollForward,
			rollForwardErr,
		)
	}

	return nil
}

func (m Migrator) mergeMigrations() (goose.Migrations, error) {
	lg := m.logger.With(zap.String("operation", "merge-migrations"))
	var mergedMigrations goose.Migrations
	uniqueIndex := make(map[int64]*goose.Migration)

	for _, env := range m.envs {
		dir := filepath.Join(m.base, env)

		// intercept goose handling of empty folders: this should not be blocking
		if m.notExists(dir) {
			lg.Warn("no migrations for env", zap.String("env", env))

			continue
		}

		// assess if there is any go migration for this env
		hasAnyGoMigration, err := m.hasGoMigrations(dir)
		if err != nil {
			return nil, err
		}

		migrationsForEnv, err := goose.CollectMigrations(dir, 0, goose.MaxVersion)
		if err != nil && !errors.Is(err, goose.ErrNoMigrationFiles) {
			return nil, fmt.Errorf("could not collect migrations for env: %s: %w", env, err)
		}

		for _, migration := range migrationsForEnv {
			// for folders without go migrations, don't add up the globally registered ones
			if !hasAnyGoMigration && migration.Registered {
				continue
			}

			if existing, ok := uniqueIndex[migration.Version]; ok {
				if existing.Registered && migration == existing {
					// globally registered go versions may produce duplicates across envs: skip dupes
					continue
				}

				return nil, fmt.Errorf(
					"duplicate versions found in migrations: %v in %v and %v",
					migration.Version, existing.Source, migration.Source,
				)
			}

			uniqueIndex[migration.Version] = migration
			mergedMigrations = append(mergedMigrations, migration)
		}
	}

	// ensuring deduped migrations alleviates the need to handle the panic there caused by goose
	sort.Sort(mergedMigrations)

	return mergedMigrations, nil
}

func (m Migrator) rollForwardFunc(db *sql.DB, migrations goose.Migrations) func(context.Context) error {
	return func(ctx context.Context) error {
		var n int
		defer func() {
			m.logger.Info("completed", zap.Int("migrations", n))
		}()

		for ; ; n++ {
			currentDBVersion, err := goose.GetDBVersion(db)
			if err != nil {
				return fmt.Errorf("could not determine the current schema version during rollforward: %w", err)
			}

			nextMigration, err := migrations.Next(currentDBVersion)
			if err != nil {
				if errors.Is(err, goose.ErrNoNextVersion) {
					// we're done here
					return nil
				}

				return fmt.Errorf("rollforward could not retrieve the next migration: %w", err)
			}

			migrationCtx, cancel := m.withMigrationTimeout(ctx)
			defer cancel()

			if err = nextMigration.UpContext(migrationCtx, db); err != nil {
				return fmt.Errorf("rollforward error in migration %v: %w", nextMigration.Source, err)
			}
		}
	}
}

func (m Migrator) rollbackToFunc(db *sql.DB, migrations goose.Migrations) func(context.Context, int64) error {
	return func(ctx context.Context, toVersion int64) error {
		var n int
		defer func() {
			m.logger.Info("rollbacked", zap.Int("migrations", n))
		}()

		for ; ; n++ {
			currentDBVersion, err := goose.GetDBVersion(db)
			if err != nil {
				return fmt.Errorf("could not determine the current schema version during rollback: %w", err)
			}

			currentMigration, err := migrations.Current(currentDBVersion)
			if err != nil {
				if errors.Is(err, goose.ErrNoCurrentVersion) {
					// we're done here
					return nil
				}

				return fmt.Errorf("rollback could not retrieve the current migration: %w", err)
			}

			if currentMigration.Version <= toVersion {
				return nil
			}

			migrationCtx, cancel := m.withMigrationTimeout(ctx)
			defer cancel()

			if err = currentMigration.DownContext(migrationCtx, db); err != nil {
				return fmt.Errorf("rollback error in migration %v: %w", currentMigration.Source, err)
			}
		}
	}
}

func (m Migrator) notExists(dir string) bool {
	_, err := fs.Stat(m.fsys, dir)

	return errors.Is(err, fs.ErrNotExist)
}

func (m Migrator) hasGoMigrations(dir string) (bool, error) {
	dir = filepath.ToSlash(dir)
	matches, err := fs.Glob(m.fsys, path.Join(dir, "*.go")) // TODO: test that this works on windows
	if err != nil {
		return false, err
	}

	var hasAtLeastOneMatch bool
	for _, match := range matches {
		if strings.HasSuffix(match, "_test.go") {
			continue
		}

		hasAtLeastOneMatch = true
		break
	}

	return hasAtLeastOneMatch, nil
}
