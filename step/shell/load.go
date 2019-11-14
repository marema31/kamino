package shell

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/Masterminds/sprig"
	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/step/common"
)

//Load data from step file using its viper representation a return priority and list of steps
func Load(ctx context.Context, log *logrus.Entry, recipePath string, name string, nameIndex int, v *viper.Viper, dss datasource.Datasourcers) (priority uint, steps []common.Steper, err error) {
	steps = make([]common.Steper, 0, 1)

	priority = v.GetUint("priority")
	tags := v.GetStringSlice("tags")
	if len(tags) == 0 {
		tags = []string{""}
	}

	logStep := log.WithField("name", name).WithField("type", "shell")

	script := v.GetString("script")
	if script == "" {
		logStep.Error("No script provided")
		return 0, nil, fmt.Errorf("the step %s must have a script to call", name)
	}
	if !filepath.IsAbs(script) {
		script = filepath.Join(recipePath, script)
	}

	if _, err := os.Stat(script); err != nil {
		if os.IsNotExist(err) {
			log.Warningf("The script %s does not exists at load phase, I will recheck at execution if it has not been created by a previous step execution", script)
		}
	}

	arguments := v.GetString("arguments")

	targuments, err := template.New("arguments").Funcs(sprig.FuncMap()).Parse(arguments)
	if err != nil {
		logStep.Error("Parsing the arguments list failed:")
		logStep.Error(err)
		return 0, nil, fmt.Errorf("error parsing the arguments of %s step: %v", name, err)
	}
	renderedArguments := bytes.NewBuffer(make([]byte, 0, 1024))

	engines := v.GetStringSlice("engines")
	e, err := datasource.StringsToEngines(engines)
	if err != nil {
		logStep.Error(err)
		return 0, nil, err
	}

	for index, datasource := range dss.Lookup(log, tags, []datasource.Type{datasource.Database}, e) {
		var step Step

		step.Name = fmt.Sprintf("%s:%d", name, nameIndex+index)
		step.script = script
		targuments.Execute(renderedArguments, datasource.FillTmplValues())
		step.arguments = renderedArguments.String()

		step.datasource = datasource

		steps = append(steps, &step)

	}

	return priority, steps, nil
}
