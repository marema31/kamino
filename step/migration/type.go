//Package migration define steps that manage the schema migration using sql-migrate engine
package migration

import (
	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/step/common"
	migrate "github.com/rubenv/sql-migrate"
)

// Step informations.
type Step struct {
	Name         string
	datasource   datasource.Datasourcer
	folder       string
	noUser       bool
	noAdmin      bool
	dir          migrate.MigrationDirection
	limit        int
	tableUser    string
	tableAdmin   string
	queries      []common.SkipQuery
	printOnly    bool //Used by migrate status
	dryRun       bool
	dialect      string
	schema       string
	ignoreErrors bool
}
