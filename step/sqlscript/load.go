package sqlscript

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"path"
	"path/filepath"
	"strings"

	"github.com/Masterminds/sprig/v3"
	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/step/common"
)

//PostLoad modify the loaded step values with the values provided in the map in argument.
func (st *Step) PostLoad(log *logrus.Entry, superseed map[string]string) error {
	// Nothing to do
	return nil
}

//Load data from step file using its viper representation a return priority and list of steps.
func Load(ctx context.Context, log *logrus.Entry, recipePath string, name string, nameIndex int, v *viper.Viper, dss datasource.Datasourcers, force bool, dryRun bool, limitedTags []string) (priority uint, steps []common.Steper, err error) { //nolint: funlen
	steps = make([]common.Steper, 0, 1)
	priority = v.GetUint("priority")
	ignoreErrors := v.GetBool("ignoreErrors")

	tags := v.GetStringSlice("tags")
	if len(tags) == 0 {
		tags = []string{""}
	}

	logStep := log.WithField("name", name).WithField("type", "sql")

	templateFile := v.GetString("template")
	if templateFile == "" {
		logStep.Error("No SQL script template filename provided")

		return 0, nil, fmt.Errorf("the step %s must have a SQL script template to render: %w", name, common.ErrMissingParameter)
	}

	if !filepath.IsAbs(templateFile) {
		templateFile = filepath.Join(recipePath, templateFile)
	}

	tsqlscript, err := template.New(path.Base(templateFile)).Funcs(sprig.FuncMap()).ParseFiles(templateFile)
	if err != nil {
		logStep.Error("Parsing the SQL script template failed:")
		logStep.Error(err)

		return 0, nil, fmt.Errorf("error parsing the template file %s of %s step: %w", templateFile, name, err)
	}

	tqueries, err := common.ParseQueries(logStep, v.GetStringSlice("queries"))
	if err != nil {
		return 0, nil, err
	}

	unique := v.GetBool("unique")
	if unique {
		logStep.Debug("Will limitate step to unique datasource")
	}

	admin := v.GetBool("admin")
	noDb := v.GetBool("noDb")
	v.SetDefault("transaction", true)
	wantTransaction := v.GetBool("transaction")

	engines := v.GetStringSlice("engines")

	e, err := datasource.StringsToEngines(engines)
	if err != nil {
		logStep.Error(err)
		return 0, nil, err
	}

	renderedSQLScript := bytes.NewBuffer(make([]byte, 0, 8192))
	datasourceHashes := map[string]bool{}

	lookedUp, _, err := dss.Lookup(logStep, tags, limitedTags, []datasource.Type{datasource.Database}, e)
	if err != nil {
		return 0, nil, err
	}

	for index, datasource := range lookedUp {
		var step Step

		if unique {
			hash := datasource.GetHash(logStep, admin, noDb)
			if _, ok := datasourceHashes[hash]; ok {
				logStep.Debugf("Skipping datasource %s: not unique", datasource.GetName())

				continue
			}

			datasourceHashes[hash] = true
		}

		step.Name = fmt.Sprintf("%s:%d", name, nameIndex+index)
		step.dryRun = dryRun
		step.templateFile = templateFile
		step.datasource = datasource
		step.admin = admin
		step.noDb = noDb
		step.ignoreErrors = ignoreErrors

		if datasource.IsTransaction() {
			step.transaction = wantTransaction
		}

		tmplValues := datasource.FillTmplValues()

		queries, err := common.RenderQueries(logStep, tqueries, tmplValues)
		if err != nil {
			return 0, nil, err
		}

		step.queries = queries

		err = tsqlscript.Execute(renderedSQLScript, tmplValues)
		if err != nil {
			logStep.Error("Rendering the sql script template failed")
			logStep.Error(err)

			return 0, nil, err
		}

		step.sqlCmds, err = splitSQLStatements(logStep, renderedSQLScript)
		if err != nil {
			return 0, nil, err
		}

		renderedSQLScript.Reset()

		steps = append(steps, &step)
	}

	return priority, steps, nil
}

// Checks the line to see if the line has a statement-ending semicolon
// or if the line contains a double-dash comment.
func endsWithSemicolon(line string) bool {
	prev := ""
	scanner := bufio.NewScanner(strings.NewReader(line))
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		word := scanner.Text()
		if strings.HasPrefix(word, "--") {
			break
		}

		prev = word
	}

	return strings.HasSuffix(prev, ";")
}

// Split the given sql script into individual statements.
//
// The base case is to simply split on semicolons, as these
// naturally terminate a statement.
//
// However, more complex cases like pl/pgsql can have semicolons
// within a statement. For these cases, we provide the explicit annotations
// 'StatementBegin' and 'StatementEnd' to allow the script to
// tell us to ignore semicolons.
func splitSQLStatements(log *logrus.Entry, r io.Reader) (stmts []string, err error) {
	// buffer for the current statement that can be in multiple lines
	var buf bytes.Buffer

	scanner := bufio.NewScanner(r)
	statementEnded := false
	ignoreSemicolons := false

	for scanner.Scan() {
		line := scanner.Text()

		// handle any goose-specific commands
		if strings.HasPrefix(line, "--") {
			cmd := strings.TrimSpace(line[2:])
			switch cmd {
			case "StatementBegin":
				ignoreSemicolons = true
			case "StatementEnd":
				if ignoreSemicolons {
					statementEnded = true
					ignoreSemicolons = false
				} else {
					statementEnded = false
				}
			}
		} else if _, err := buf.WriteString(line + "\n"); err != nil {
			// Add the line to the current statement
			log.Error("Splitting SQL script failed")
			log.Error(err)

			return nil, err
		}

		// Wrap up the two supported cases: 1) basic with semicolon; 2) psql statement
		// Lines that end with semicolon that are in a statement block
		// do not conclude statement.
		// add the statement to the slice, and empty the buffer
		if (!ignoreSemicolons && endsWithSemicolon(line)) || statementEnded {
			statementEnded = false

			stmts = append(stmts, buf.String())
			buf.Reset()
		}
	}

	if err := scanner.Err(); err != nil {
		log.Error("Splitting SQL script failed")
		log.Error(err)

		return nil, err
	}

	// diagnose likely SQL script errors
	if ignoreSemicolons {
		log.Warn("Saw '-- StatementBegin' with no matching '-- StatementEnd'")
	}

	if bufferRemaining := strings.TrimSpace(buf.String()); len(bufferRemaining) > 0 {
		log.Warnf("Unexpected unfinished SQL query: %s. Missing a semicolon?", bufferRemaining)
	}

	return stmts, err
}
