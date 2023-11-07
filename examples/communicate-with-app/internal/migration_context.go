package internal

import (
	"context"
	"errors"

	"go.uber.org/zap"
)

// MigrationContext holds all the settings we want to share with migration scripts.
//
// NOTE: this is important to declare this in a separate package, and avoid cyclic import dependencies.
type MigrationContext struct {
	Logger *zap.Logger
	Name   string
}

type ctxKey uint8

const migrationContextKey ctxKey = iota + 1

// FromContext retrieves the app context from the context.
func FromContext(ctx context.Context) (MigrationContext, error) {
	v := ctx.Value(migrationContextKey)
	m, ok := v.(MigrationContext)
	if !ok {
		return MigrationContext{}, errors.New("no migration context passed")
	}

	return m, nil
}

// ToContext pushes the app context into a child context.
func (m MigrationContext) ToContext(parentCtx context.Context) context.Context {
	return context.WithValue(parentCtx, migrationContextKey, m)
}
