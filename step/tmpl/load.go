package tmpl

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"

	"github.com/Masterminds/sprig"
	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/step/common"
)

//Load data from step file using its viper representation a return priority and list of steps
func Load(ctx context.Context, log *logrus.Entry, recipePath string, filename string, v *viper.Viper, dss datasource.Datasourcers) (priority uint, steps []common.Steper, err error) {
	steps = make([]common.Steper, 0, 1)

	priority = v.GetUint("priority")
	tags := v.GetStringSlice("tags")
	if len(tags) == 0 {
		tags = []string{""}
	}

	name := v.GetString("name")
	logStep := log.WithField("name", name).WithField("type", "shell")

	templateFile := v.GetString("template")
	if templateFile == "" {
		logStep.Error("No template filename provided")
		return 0, nil, fmt.Errorf("the step %s must have a template to render", name)
	}
	if !filepath.IsAbs(templateFile) {
		templateFile = filepath.Join(recipePath, templateFile)
	}

	template, err := template.New("template").Funcs(sprig.FuncMap()).ParseFiles(templateFile)
	if err != nil {
		logStep.Error("Parsing the template failed:")
		logStep.Error(err)
		return 0, nil, fmt.Errorf("error parsing the template file %s of %s step: %v", templateFile, name, err)
	}

	destinationTmpl := v.GetString("destination")
	if destinationTmpl == "" {
		logStep.Error("No destination filename provided")
		return 0, nil, fmt.Errorf("the step %s must have a destination to be rendered", name)
	}
	tdestination, err := template.New("destination").Funcs(sprig.FuncMap()).Parse(destinationTmpl)
	if err != nil {
		logStep.Error("Parsing the destination filename template failed:")
		logStep.Error(err)
		return 0, nil, fmt.Errorf("error parsing the destination of %s step: %v", name, err)
	}
	renderedDestination := bytes.NewBuffer(make([]byte, 0, 1024))

	var mode Mode
	switch strings.ToLower(v.GetString("mode")) {
	case "replace":
		mode = Replace
	case "append":
		mode = Append
	case "replaceappend":
		mode = ReplaceAppend
	default:
		mode = Replace
	}

	engines := v.GetStringSlice("engines")
	e, err := datasource.StringsToEngines(engines)
	if err != nil {
		logStep.Error(err)
		return 0, nil, err
	}

	for index, datasource := range dss.Lookup(log, tags, []datasource.Type{datasource.Database}, e) {
		var step Step

		step.Name = fmt.Sprintf("%s:%d", name, index)
		step.templateFile = templateFile
		step.template = template
		step.datasource = datasource
		step.mode = mode

		tdestination.Execute(renderedDestination, datasource.FillTmplValues())
		step.destination = renderedDestination.String()

		steps = append(steps, &step)

	}

	return priority, steps, nil
}
