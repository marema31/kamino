// Package shell manage step that runs a shell command with template
package shell

import (
	"github.com/marema31/kamino/datasource"
)

// Step informations
type Step struct {
	Name       string
	datasource datasource.Datasourcer
	arguments  string
	script     string
}
