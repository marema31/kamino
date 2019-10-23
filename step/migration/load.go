package migration

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

	folderTmpl := v.GetString("folder")
	if folderTmpl == "" {
		return 0, nil, fmt.Errorf("the step %s must have a folder that contains the migration files", name)
	}
	tfolder, err := template.New("folder").Funcs(sprig.FuncMap()).Parse(folderTmpl)
	if err != nil {
		return 0, nil, fmt.Errorf("error parsing the folder of %s step: %v", name, err)
	}
	renderedFolder := bytes.NewBuffer(make([]byte, 0, 1024))

	engines := v.GetStringSlice("engines")
	e, err := datasource.StringsToEngines(engines)
	if err != nil {
		return 0, nil, err
	}

	for _, datasource := range dss.Lookup(tags, []datasource.Type{datasource.Database}, e) {
		var step Step

		step.Name = name
		step.datasource = datasource

		tfolder.Execute(renderedFolder, datasource.FillTmplValues())
		step.folder = renderedFolder.String()
		//TODO: folder not found

		steps = append(steps, step)

	}

	return priority, steps, nil
}
