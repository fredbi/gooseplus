package gooseplus

// registers go migrations for unit tests
import (
	"embed"

	_ "github.com/fredbi/gooseplus/test_sql/unittest"
	_ "github.com/fredbi/gooseplus/test_sql/unittest3"
)

//go:embed test_sql/*/*.sql
//go:embed test_sql/*/*.go
var embedMigrations embed.FS
