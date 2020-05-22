// Package sqlscript manage the steps that runs a sql command on all destinations
package sqlscript

import (
	"database/sql"

	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/step/common"
)

// Step informations.
type Step struct {
	Name         string
	datasource   datasource.Datasourcer
	admin        bool
	noDb         bool
	queries      []common.SkipQuery
	templateFile string
	sqlCmds      []string
	transaction  bool
	tx           *sql.Tx
	dryRun       bool
}
