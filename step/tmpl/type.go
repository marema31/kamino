//Package tmpl manage step that generate files from templates
package tmpl

import (
	"html/template"

	"github.com/marema31/kamino/datasource"
)

//Mode discriminate the result on the destination file
type Mode int

const (
	// Replace the destination file for each datasource of a step (if the destination field is a not template, file will be overwritten)
	Replace Mode = iota
	// ReplaceAppend will replace the destination file for defined step but will append for each datasource of that step(if the destination field is a not template, file will be overwritten)
	ReplaceAppend Mode = iota
	// Append to the destination
	Append Mode = iota
)

// Step informations
type Step struct {
	Name         string
	datasource   datasource.Datasourcer
	templateFile string
	template     *template.Template
	destination  string
	mode         Mode
}
