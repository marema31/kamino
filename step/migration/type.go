//Package migration define steps that manage the schema migration using sql-migrate engine
package migration

import (
	"github.com/marema31/kamino/datasource"
)

// Step informations
type Step struct {
	Name       string
	datasource datasource.Datasourcer
	folder     string
}
