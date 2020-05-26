package common

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"strconv"
	"strings"

	"github.com/Masterminds/sprig/v3"
	"github.com/marema31/kamino/datasource"

	"github.com/Sirupsen/logrus"
)

// ToSkipDatabase run the query (likely a SELECT COUNT) on the datasource
// return true if the query return a non-zero value in the only column of the only row.
func ToSkipDatabase(ctx context.Context, log *logrus.Entry, ds datasource.Datasourcer, admin bool, nodb bool, queries []SkipQuery) (bool, error) {
	var needskip int

	if len(queries) == 0 {
		return false, nil
	}

	db, err := ds.OpenDatabase(log, admin, nodb)
	if err != nil {
		return false, err
	}

	for _, query := range queries {
		err = db.QueryRowContext(ctx, query.query).Scan(&needskip)
		if err != nil {
			log.Errorf("Query of skip phase failed : %v", err)
			return false, err
		}

		log.Debugf("Skip query: %s (%d ~ %d)", query.query, needskip, query.compareValue)
		// If the query returned a 0 result we shoud do it, do not test other queries (AND) to allow test table exists and table contains something
		if needskip == query.compareValue {
			if !query.inverted {
				log.Debugf("Skip query not skipped: %s", query.query)
				return false, nil
			}
		} else {
			if query.inverted {
				log.Debugf("Skip query not skipped: %s", query.query)
				return false, nil
			}
		}
	}

	return true, nil
}

// ParseQueries parse the queries templates from step configuration file.
func ParseQueries(log *logrus.Entry, queriesTmpl []string) ([]TemplateSkipQuery, error) {
	var err error

	tqueries := make([]TemplateSkipQuery, 0, len(queriesTmpl))

	if len(queriesTmpl) == 0 {
		log.Warning("No SQL queries provided")
	}

	for _, queryTmpl := range queriesTmpl {
		cmpValue := 0
		inverted := false

		if queryTmpl == "" {
			log.Error("No SQL query provided")
			return nil, fmt.Errorf("the step cannot have an empty query to be executed: %w", ErrMissingParameter)
		}

		if strings.HasPrefix(queryTmpl, "!") {
			inverted = true

			queryTmpl = strings.TrimPrefix(queryTmpl, "!")
		}

		if strings.HasPrefix(queryTmpl, "=") {
			endValue := strings.Index(queryTmpl, ":")
			if endValue == -1 {
				log.Errorf("Query %s has comparison but value has no terminator", queryTmpl)
				return nil, fmt.Errorf("the step cannot have an incorrectly formatted query to be executed: %w", ErrWrongParameterValue)
			}

			cmpValue, err = strconv.Atoi(queryTmpl[1:endValue])
			if err != nil {
				log.Errorf("Query %s has comparison but incorrect value has been provided", queryTmpl)
				return nil, fmt.Errorf("the step cannot have an incorrectly formatted query to be executed: %w", ErrWrongParameterValue)
			}

			queryTmpl = strings.TrimPrefix(queryTmpl, queryTmpl[0:endValue+1])
		}

		tquery, err := template.New("query").Funcs(sprig.FuncMap()).Parse(queryTmpl)
		if err != nil {
			log.Errorf("Parsing the SQL query template failed: %v", err)
			return nil, fmt.Errorf("error parsing the query of step: %w", err)
		}

		tqueries = append(tqueries, TemplateSkipQuery{tquery: tquery, compareValue: cmpValue, inverted: inverted, text: queryTmpl})
	}

	return tqueries, nil
}

//RenderQueries render the templated query with parameter corresponding to the datasource.
func RenderQueries(log *logrus.Entry, tqueries []TemplateSkipQuery, tmplValues datasource.TmplValues) ([]SkipQuery, error) {
	renderedQuery := bytes.NewBuffer(make([]byte, 0, 4096))
	queries := make([]SkipQuery, 0, len(tqueries))

	for _, tquery := range tqueries {
		err := tquery.tquery.Execute(renderedQuery, tmplValues)
		if err != nil {
			log.Errorf("Rendering the '%s' query template failed :%v ", tquery.text, err)
			return nil, err
		}

		queries = append(queries, SkipQuery{query: renderedQuery.String(), compareValue: tquery.compareValue, inverted: tquery.inverted})

		renderedQuery.Reset()
	}

	return queries, nil
}
