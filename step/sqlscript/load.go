package sqlscript

import (
	"bytes"
	"context"
	"fmt"
	"html/template"

	"github.com/Masterminds/sprig"
	"github.com/spf13/viper"

	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/step/common"
)

//Load data from step file using its viper representation a return priority and list of steps
func Load(ctx context.Context, filename string, v *viper.Viper, dss datasource.Datasourcers) (priority uint, steps []common.Steper, err error) {
	steps = make([]common.Steper, 0, 1)

	priority = v.GetUint("priority")
	tags := v.GetStringSlice("tags")
	if len(tags) == 0 {
		tags = []string{""}
	}

	name := v.GetString("name")
	templateFile := v.GetString("template")
	if templateFile == "" {
		return 0, nil, fmt.Errorf("the step %s must have a template to render", name)
	}
	template, err := template.New("template").Funcs(sprig.FuncMap()).ParseFiles(templateFile)
	if err != nil {
		return 0, nil, fmt.Errorf("error parsing the template file %s of %s step: %v", templateFile, name, err)
	}

	queryTmpl := v.GetString("query")
	if queryTmpl == "" {
		return 0, nil, fmt.Errorf("the step %s must have a query to be executed", name)
	}
	tquery, err := template.New("query").Funcs(sprig.FuncMap()).Parse(queryTmpl)
	if err != nil {
		return 0, nil, fmt.Errorf("error parsing the query of %s step: %v", name, err)
	}
	renderedQuery := bytes.NewBuffer(make([]byte, 0, 1024))

	admin := v.GetBool("admin")
	noDb := v.GetBool("noDb")

	engines := v.GetStringSlice("engines")
	e, err := datasource.StringsToEngines(engines)
	if err != nil {
		return 0, nil, err
	}

	for _, datasource := range dss.Lookup(tags, []datasource.Type{datasource.Database}, e) {
		var step Step

		step.Name = name
		step.templateFile = templateFile
		step.template = template
		step.datasource = datasource
		step.admin = admin
		step.noDb = noDb

		tquery.Execute(renderedQuery, datasource.FillTmplValues())
		step.query = renderedQuery.String()

		steps = append(steps, &step)

	}

	return priority, steps, nil
}
