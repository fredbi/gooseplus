# Examples

The following examples are provided with a plain `sqlite3` DB for demoing purpose.

Of course, they translate easily with production-ready DB such as Postgres or MySQL.

## [multi-env](multi-env/README.md)

Demonstrates how to run environment-specific migrations.

## [go-migrations](go-migrations/README.md)

Demonstrates how to run go migrations (from embedded FS).

## [osfs](osfs/README.md)

Demonstrate how to run migrations from plain OS FS (which is the default).

## [communicate-with-app](communicate-with-app/README.md)

Demonstrate how go migrations may interact with the main app to share settings via context.

## [multi-lanes](multi-lanes/README.md)

Demonstrate how to setup several migration lanes at deployment time: a fast lane and a long-running one,
that continues running after a deployment is done.
