# go-migrations

Use programmatic migration scripts, rather than plain SQL.

Use-case: populating data from a file, any complex logic difficult to achieve in plain SQL.

Avoids using procedural extensions such as PL/pgSQL (postgres).

> Reminder: SQL migrations are plain SQL, not SQL scripts for frontends like `psql`. Those tools may
> have some advanced capabilities (variables, file loading...) not available from plain SQL.
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
