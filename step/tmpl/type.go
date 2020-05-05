//Package tmpl manage step that generate files from templates
package tmpl

import (
	"html/template"
	"io"

	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/file"
)

//Mode discriminate the result on the destination file.
type Mode int

const (
	// Replace the destination file for each datasource of a step (if the destination field is a not template, file will be a concatenation of template for each datasource)
	Replace Mode = iota
	// Append to the destination
	Append Mode = iota
	// Unique avoid duplicate the generated content in the file
	Unique Mode = iota
)

// Step informations.
type Step struct {
	Name            string
	datasources     []datasource.Datasourcer
	templateFile    string
	template        *template.Template
	destination     string
	mode            Mode
	onlyIfNotExists bool
	input           file.File
	inputHandle     io.ReadCloser
	output          file.File
	outputHandle    io.Writer
	dryRun          bool
}
