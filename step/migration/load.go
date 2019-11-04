package migration

import (
	"bytes"
	"context"
	"fmt"
	"html/template"

	"github.com/Masterminds/sprig"
	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/step/common"
)

//Load data from step file using its viper representation a return priority and list of steps
func Load(ctx context.Context, log *logrus.Entry, filename string, v *viper.Viper, dss datasource.Datasourcers) (priority uint, steps []common.Steper, err error) {
	steps = make([]common.Steper, 0, 1)

	priority = v.GetUint("priority")
	tags := v.GetStringSlice("tags")
	if len(tags) == 0 {
		tags = []string{""}
	}

	name := v.GetString("name")
	logStep := log.WithField("name", name).WithField("type", "migration")

	folderTmpl := v.GetString("folder")
	if folderTmpl == "" {
		logStep.Error("No migration folder provided")
		return 0, nil, fmt.Errorf("the step %s must have a folder that contains the migration files", name)
	}
	tfolder, err := template.New("folder").Funcs(sprig.FuncMap()).Parse(folderTmpl)
	if err != nil {
		logStep.Error("Parsing the migration folder name failed:")
		logStep.Error(err)
		return 0, nil, fmt.Errorf("error parsing the folder of %s step: %v", name, err)
	}
	renderedFolder := bytes.NewBuffer(make([]byte, 0, 1024))

	engines := v.GetStringSlice("engines")
	e, err := datasource.StringsToEngines(engines)
	if err != nil {
		logStep.Error(err)
		return 0, nil, err
	}

	for index, datasource := range dss.Lookup(log, tags, []datasource.Type{datasource.Database}, e) {
		var step Step

		step.Name = fmt.Sprintf("%s:%d", name, index)
		step.datasource = datasource

		tfolder.Execute(renderedFolder, datasource.FillTmplValues())
		step.folder = renderedFolder.String()
		//TODO: folder not found

		steps = append(steps, &step)

	}

	return priority, steps, nil
}
