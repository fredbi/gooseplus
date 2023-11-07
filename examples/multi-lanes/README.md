# multi-lanes

Proceed with 2 migration lanes: one processing fast running migrations, another one for long-running scripts.

Use-case: usual migrations can proceed swiftly, so as to startup a deployed service quickly.

However some data manipulations that we would like to keep under versioning may take some time
(e.g. rebuild index, refresh materialized view, prepare bulk data initialization...).

Allowing for a special lane for those scripts allow to start the service as soon as the database state
is workable for the app.

This example illustrate how the completion of the long-running migrations can be signaled to the main app,
for example to open some features that require that part.

## Usage

```sh
go run migrate.go
{"level":"info","msg":"applying db migrations","operation":"migratedb"}
{"level":"info","msg":"OK   20231102204811_create_example.sql (2.89ms)","operation":"migratedb"}
{"level":"info","msg":"OK   20231103204811_populate_example.sql (2.64ms)","operation":"migratedb"}
{"level":"info","msg":"applying db migrations","operation":"migratedb"}
2023/11/07 22:43:24 app can do some work, polling until long-running migrations are complete
2023/11/07 22:43:25 app can do some work, polling until long-running migrations are complete
2023/11/07 22:43:26 app can do some work, polling until long-running migrations are complete
2023/11/07 22:43:27 app can do some work, polling until long-running migrations are complete
2023/11/07 22:43:28 app can do some work, polling until long-running migrations are complete
{"level":"info","msg":"OK   20231104120000_load_example.go (5s)","operation":"migratedb"}
2023/11/07 22:43:28 background migrations completed
2023/11/07 22:43:28 app can fully proceed now that long running migrations are passed
2023/11/07 22:43:28 app exited gracefully
```
