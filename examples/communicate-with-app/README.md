# communicate-with-app

Use programmatic migration scripts, passing parameters from the calling app using the context.

Use-case: sharing logger, config and any other settings with the migration scripts.

## Usage

```sh
 go run migrate.go
{"level":"info","logger":"demo","msg":"applying db migrations","operation":"migratedb"}
{"level":"info","logger":"demo","msg":"OK   20231102204811_create_example.sql (2.68ms)","operation":"migratedb"}
{"level":"info","logger":"demo","msg":"OK   20231103204811_populate_example.sql (2.65ms)","operation":"migratedb"}
{"level":"info","logger":"demo","msg":"values inserted from file","entries":2}
{"level":"info","logger":"demo","msg":"OK   20231104120000_load_example.go (2.74ms)","operation":"migratedb"}
```

```sh
sqlite3 testdata/db/example.db 
SQLite version 3.31.1 2020-01-27 19:55:54
Enter ".help" for usage hints.
sqlite> select * from example;
one|One
three|name: passed value- Third value
four|name: passed value- Fourth value
sqlite> 
```
