package shell

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
	script := v.GetString("script")
	if script == "" {
		return 0, nil, fmt.Errorf("the step %s must have a script to call", name)
	}

	//TODO: script not found

	arguments := v.GetString("arguments")

	targuments, err := template.New("arguments").Funcs(sprig.FuncMap()).Parse(arguments)
	if err != nil {
		return 0, nil, fmt.Errorf("error parsing the arguments of %s step: %v", name, err)
	}
	renderedArguments := bytes.NewBuffer(make([]byte, 0, 1024))

	engines := v.GetStringSlice("engines")
	e, err := datasource.StringsToEngines(engines)
	if err != nil {
		return 0, nil, err
	}

	for _, datasource := range dss.Lookup(tags, []datasource.Type{datasource.Database}, e) {
		var step Step

		step.Name = name
		step.script = script
		targuments.Execute(renderedArguments, datasource.FillTmplValues())
		step.arguments = renderedArguments.String()

		step.datasource = datasource

		steps = append(steps, &step)

	}

	return priority, steps, nil
}
