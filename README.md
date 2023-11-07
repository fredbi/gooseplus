# gooseplus
![Lint](https://github.com/fredbi/gooseplus/actions/workflows/01-golang-lint.yaml/badge.svg)
![CI](https://github.com/fredbi/gooseplus/actions/workflows/02-test.yaml/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/fredbi/gooseplus/badge.svg?branch=master)](https://coveralls.io/github/fredbi/gooseplus?branch=master)
![Vulnerability Check](https://github.com/fredbi/gooseplus/actions/workflows/03-govulncheck.yaml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/fredbi/gooseplus)](https://goreportcard.com/report/github.com/fredbi/gooseplus)

![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/fredbi/gooseplus)
[![Go Reference](https://pkg.go.dev/badge/github.com/fredbi/gooseplus.svg)](https://pkg.go.dev/github.com/fredbi/gooseplus)
[![license](http://img.shields.io/badge/license/License-Apache-yellow.svg)](https://raw.githubusercontent.com/fredbi/gooseplus/master/LICENSE.md)

Goose DB migrations, extended.

## Purpose

`gooseplus` extends the great DB migration tool [`goose`](https://github.com/pressly/goose) to supports a few advanced
use cases:

1. Leaves a failed deployment in the initial state: rollbacks migrations back to when the deployment failed
2. Support environment-specific migrations, so we can add migrations for tests, etc.
3. More options: structured zap logger, fined-grained timeouts ...

`gooseplus` is primarily intended to be used as a library, and does not come with a CLI command.

## Usage

```go
    db, _ := sql.Open("postgres", "test")
    migrator := New(db)

    err := migrator.Migrate(context.Background())
```

Feel free to look at the various [`examples`](examples/README.md).

## Features

* Everything `goose` does out of the box.
* Rollback to the state at the start of the call, after a failure.
* Environment-specific migration folders

## Caveats

* Concurrent usage is not supported: `goose` relies on a lot of globals. Minimal locking has ben added so you can run
your tests with `-race`, but migrations should normally run once.
