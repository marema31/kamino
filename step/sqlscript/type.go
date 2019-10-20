// Package sqlscript manage the steps that runs a sql command on all destinations
package sqlscript

import (
	"html/template"

	"github.com/marema31/kamino/datasource"
)

// Step informations
type Step struct {
	Name         string
	datasource   *datasource.Datasource
	admin        bool
	noDb         bool
	query        string
	template     *template.Template
	templateFile string
}
