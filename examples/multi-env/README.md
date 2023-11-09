# multi-env

Migrations folders selected from environment variable `APP_ENV`.

Use-case: populating data with different settings for test and production.

In the following example, we populate some initial data with different values for the test and the production environment.

## Usage

```sh
APP_ENV=production go run migrate.go
{"level":"info","msg":"applying db migrations","operation":"migratedb"}
{"level":"info","msg":"OK   20231102204811_create_example.sql (3.52ms)","operation":"migratedb"}
{"level":"info","msg":"OK   20231103204811_populate_example.sql (3.02ms)","operation":"migratedb"}
{"level":"info","msg":"OK   20231103204911_populate_prod.sql (3.33ms)","operation":"migratedb"}
```


```sh
APP_ENV=test go run migrate.go
{"level":"info","msg":"applying db migrations","operation":"migratedb"}
{"level":"info","msg":"OK   20231102204811_create_example.sql (3.52ms)","operation":"migratedb"}
{"level":"info","msg":"OK   20231103204811_populate_example.sql (3.02ms)","operation":"migratedb"}
{"level":"info","msg":"OK   20231103204911_populate_test.sql (2.68ms)","operation":"migratedb"}
```
