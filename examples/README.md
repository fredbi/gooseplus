# Examples

The following examples are provided with a plain `sqlite3` DB for demoing purpose.

Obvisouly, they translate easily into production-ready DBs such as Postgres or MySQL.

## [multi-env](multi-env/README.md)

Demonstrates how to run environment-specific migrations.

## [go-migrations](go-migrations/README.md)

Demonstrates how to run go migrations (from embedded FS).

## [osfs](osfs/README.md)

Demonstrates how to run migrations from plain OS FS (which is the default).

## [communicate-with-app](communicate-with-app/README.md)

Demonstrates how go migrations may interact with the main app to share settings via context.

## [multi-lanes](multi-lanes/README.md)

Demonstrate how to set up several migration lanes at deployment time: a fast lane and a long-running one,
which continues running after a deployment is done.
