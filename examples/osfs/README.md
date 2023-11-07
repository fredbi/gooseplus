# osfs

Migrations folders mounted from regular OS FS, not embedded FS.

Use-case: bake migrations as plain files into a docker image

> NOTE: in this example, the migrations layout has been flattened, without multi-env folders.

## Usage

```sh
go run migrate.go
{"level":"info","msg":"applying db migrations","operation":"migratedb"}
{"level":"info","msg":"OK   20231102204811_create_example.sql (2.65ms)","operation":"migratedb"}
{"level":"info","msg":"OK   20231103204811_populate_example.sql (2.59ms)","operation":"migratedb"}
```
