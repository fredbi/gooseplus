package gooseplus

import (
	"context"
	"io/fs"
	"os"
	"time"

	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

type (
	// Option for the Migrator.
	//
	// Default settings are:
	//
	//   dialect: "postgres",
	//   base:    "sql",
	//   envs: []string{"base"},
	//   timeout: 5 * time.Minute,
	//   migrationTimeout: 1 * time.Minute,
	//   logger:  zap.NewExample(),
	//   fsys:    os.DirFS("."),
	Option func(*options)

	options struct {
		dialect          string
		base             string
		fsys             fs.FS
		versionTable     string
		logger           *zap.Logger
		envs             []string
		timeout          time.Duration
		migrationTimeout time.Duration
		// gooseOptions     []goose.OptionsFunc
	}
)

var baseOptions = options{
	dialect:          "postgres",
	base:             "sql",
	envs:             []string{"base"},
	timeout:          5 * time.Minute,
	migrationTimeout: 1 * time.Minute,
	logger:           zap.NewExample(),
	fsys:             os.DirFS("."),
}

func applyOptionsWithDefaults(opts []Option) options {
	if len(opts) == 0 {
		return baseOptions
	}

	o := baseOptions

	for _, apply := range opts {
		apply(&o)
	}

	if len(o.envs) == 0 {
		o.envs = []string{""}
	}

	return o
}

func (o options) withOptionalTimeout(ctx context.Context) (context.Context, func()) {
	if o.timeout == 0 {
		return ctx, func() {}
	}

	if _, ok := ctx.Deadline(); ok {
		// parent context is already set with a deadline
		return ctx, func() {}
	}

	return context.WithTimeout(ctx, o.timeout)
}

func (o options) withMigrationTimeout(ctx context.Context) (context.Context, func()) {
	if o.migrationTimeout == 0 {
		return ctx, func() {}
	}

	return context.WithTimeout(ctx, o.migrationTimeout)
}

// setGooseGlobals sets global variables in goose: guard against race
// (but cannot be used concurrently with different settings)
func (o options) setGooseGlobals(lg *zap.Logger) error {
	gooseMx.Lock()
	defer gooseMx.Unlock()

	goose.SetBaseFS(o.fsys)
	o.setVersionTable()
	goose.SetLogger(zap.NewStdLog(lg))

	err := goose.SetDialect(o.dialect)

	return err
}

func (o options) setVersionTable() {
	if o.versionTable == "" {
		return
	}

	// sets goose's global var
	goose.SetTableName(o.versionTable)
}

// WithTimeout specifies a timeout to apply to the whole migration process.
//
// NOTE: if Migrate(ctx) is called with a context that already contains a deadline,
// that deadline will override this option.
//
// The zero value disables the timeout.
//
// Default is 5m.
func WithTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.timeout = timeout
	}
}

// WithMigrationTimeout specifies a timeout to apply for each individual migration.
//
// The zero value disables the timeout.
//
// Default is 1m.
func WithMigrationTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.migrationTimeout = timeout
	}
}

// WithDialect indicates the database SQL dialect.
//
// For details see https://pkg.go.dev/github.com/pressly/goose/v3#Dialect
func WithDialect(dialect string) Option {
	return func(o *options) {
		o.dialect = dialect
	}
}

// WithEnvironments appends environment-specific folders to merge with the migrations.
//
// The default setting is a single folder "base".
func WithEnvironments(envs ...string) Option {
	return func(o *options) {
		o.envs = append(o.envs, envs...)
	}
}

// SetEnvironments overrides environment-specific folders to merge with the migrations.
//
// Setting to nil or to an empty slice will disable folders: migrations will be searched for in the base path only.
func SetEnvironments(envs []string) Option {
	return func(o *options) {
		o.envs = envs
	}
}

// WithFS provides the file system where migrations are located.
//
// The base is os.Dir(".").
func WithFS(fsys fs.FS) Option {
	return func(o *options) {
		o.fsys = fsys
	}
}

// WithBasePath provides the root directory where migrations are located on the FS.
func WithBasePath(base string) Option {
	return func(o *options) {
		o.base = base
	}
}

// WithLogger provides a structured zap logger to the migrator.
func WithLogger(zlg *zap.Logger) Option {
	return func(o *options) {
		o.logger = zlg
	}
}

// WithVersionTable tells goose to use an non-default version table.
//
// The default is "". Setting an empty table equates to using the default.
func WithVersionTable(table string) Option {
	return func(o *options) {
		o.versionTable = table
	}
}

/* not workable for now.

// WithGooseOptions allows to inject other options available to `goose` users.
func WithGooseOptions(gooseOptions ...goose.OptionsFunc) Option {
	return func(o *options) {
		o.gooseOptions = gooseOptions
	}
}
*/
