package gooseplus

import "errors"

// Errors returned by the Migrator
var (
	ErrMigrationTable  = errors.New("could not ensure goose migration table")
	ErrMergeMigrations = errors.New("error merging migrations")
	ErrRollForward     = errors.New("error rolling forward migrations. Recovered error: the db has been rollbacked to its initial state")
	ErrRollBack        = errors.New("error rolling back migrations")
)
