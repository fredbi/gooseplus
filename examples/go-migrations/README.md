# go-migrations

Use programmatic migration scripts, rather than plain SQL.

Use-case: populating data from a file, applying any complex logic difficult to achieve in plain SQL.

Avoids using procedural extensions such as PL/pgSQL (postgres).

> Reminder: SQL migrations are plain SQL _statements_, and not SQL scripts for frontends tools
> like `psql`. Those tools may have some advanced capabilities (variables, file loading...) 
> which are not available from plain SQL.

In the following example, we run a data-initialization script that loads a plain CSV file into a table.

## Usage

```sh
go run migrate.go
{"level":"info","msg":"applying db migrations","operation":"migratedb"}
{"level":"warn","msg":"no migrations for env","operation":"merge-migrations","env":"test"}
{"level":"info","msg":"OK   20231102204811_create_example.sql (2.67ms)","operation":"migratedb"}
{"level":"info","msg":"OK   20231103204811_populate_example.sql (2.65ms)","operation":"migratedb"}
2023/11/07 18:57:09 2 values inserted from file
{"level":"info","msg":"OK   20231104120000_load_example.go (3.14ms)","operation":"migratedb"}
```
