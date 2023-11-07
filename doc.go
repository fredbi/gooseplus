// Package gooseplus provides a configurable DB Migrator.
//
// It extends github.com/pressly/goose/v3 with the following features:
// * rollbacks to the initial state of a deployment whenever a migration fails in a sequence of several migrations
// * supports multiple-environments, so that it is possible to define environment-specific migrations
// * supports options: such as structured logging with zap.
package gooseplus
