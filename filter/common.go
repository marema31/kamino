package filter

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/Masterminds/sprig/v3"
)

type tmplEnv struct {
	Environments map[string]string
}

func parseField(fieldName string, fieldValue string, data tmplEnv) (string, error) {
	var buf bytes.Buffer

	tmpl, err := template.New(fieldName).Funcs(sprig.FuncMap()).Parse(fieldValue)
	if err != nil {
		return "", fmt.Errorf("parsing %s provided: %w", fieldName, err)
	}

	if err = tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("expanding %s provided: %w", fieldName, err)
	}

	parsed := buf.String()

	return parsed, nil
}
