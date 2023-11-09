# gooseplus
![Lint](https://github.com/fredbi/gooseplus/actions/workflows/01-golang-lint.yaml/badge.svg)
![CI](https://github.com/fredbi/gooseplus/actions/workflows/02-test.yaml/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/fredbi/gooseplus/badge.svg?branch=master)](https://coveralls.io/github/fredbi/gooseplus?branch=master)
![Vulnerability Check](https://github.com/fredbi/gooseplus/actions/workflows/03-govulncheck.yaml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/fredbi/gooseplus)](https://goreportcard.com/report/github.com/fredbi/gooseplus)

![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/fredbi/gooseplus)
[![Go Reference](https://pkg.go.dev/badge/github.com/fredbi/gooseplus.svg)](https://pkg.go.dev/github.com/fredbi/gooseplus)
[![license](http://img.shields.io/badge/license/License-Apache-yellow.svg)](https://raw.githubusercontent.com/fredbi/gooseplus/master/LICENSE.md)

Goose DB migrations, on steroids.

## Purpose

`gooseplus` extends the great DB migration tool [`goose`](https://github.com/pressly/goose) to support
a few advanced use cases:

1. Leaves a failed deployment in its initial state: upon failure, rollbacks migrations back to when the deployment started
2. Support environment-specific migrations, so we can add migrations for tests, etc.
3. More options: structured zap logger, fined-grained timeouts ...

`gooseplus` is primarily intended to be used as a library, and does not come with a CLI command.

## Usage

```go
    db, _ := sql.Open("postgres", "test")
    migrator := New(db)

    err := migrator.Migrate(context.Background())
    ...
```

Feel free to look at the various [`examples`](examples/README.md).

## Features

* Everything `goose/v3` does out of the box.
* Rollback to the state at the start of the call to `Migrate()` after a failure.
* Environment-specific migration folders

## Concepts

### Defaults

I've tried to define sensible defaults as follows:
* default DB driver: `postgres` (like `goose`)
* default base path for migrations: `sql`
* default FS: `os.DirFS(".")`
* default timeout on the whole migration process: 5m
* default timeout on any single migration: 1m

### Environment-specific folders

Migrations are stored in a base directory as a linear sequence of SQL scripts or go migration programs.

In this directory, the `base` folder contains migrations that apply to all environments.

Additional folders may be defined to run migrations for specific environments (i.e. specific deployment contexts).

> This comes in handy in situations where we want data initialization scripts 
> (not just schema changes) to run under different environments.

Example:
```
sql/base/
sql/base/20231103204811_populate_example.sql
sql/base/20231102204811_create_example.sql

sql/production/
sql/production/20231103204911_populate_prod.sql

sql/test/
sql/test/20231103204911_populate_test.sql
```

You can change the `base` folder by setting the new list of folders: `SetEnvironments([]string{"default", "production", "test")`.

If you don't want to manage sub-folders _at all_, you can disable it with the option `SetEnvironments(nil)`.
In this case, no `base` folder will be used.

> *Attention point*: if you use go migrations these folders become go packages,
> and folder names should not be reserved names with a special meaning for golang.
> Hence `default`, `xxx_test` are names to be avoided for package names.

### Embedded file system

`goose/v3` supports embedded file systems at build time.

You can use it with `gooseplus` like so:
```go
	//go:embed sql/*/*.sql
	var embedMigrations embed.FS

	db, _ := sql.Open("postgres", "test")

	migrator := gooseplus.New(
		db,
		gooseplus.WithFS(embedMigrations),
	)
```

### Logging

`gooseplus` injects a structured `zap` logger from `go.uber.org/zap`

### Caveats

* Concurrent usage is not supported: `goose/v3` relies on a lot of globals. Migrations should normally run once.
* Minimal locking has ben added so you can run your tests with `-race`
