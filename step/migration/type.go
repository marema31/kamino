//Package migration define steps that manage the schema migration using sql-migrate engine
package migration

import (
	"github.com/marema31/kamino/datasource"
	migrate "github.com/rubenv/sql-migrate"
)

// Step informations
type Step struct {
	Name       string
	datasource datasource.Datasourcer
	folder     string
	noUser     bool
	noAdmin    bool
	dir        migrate.MigrationDirection
	limit      int
	tableUser  string
	tableAdmin string
	query      string
	printOnly  bool //Used by migrate status
	dialect    string
	schema     string
}
