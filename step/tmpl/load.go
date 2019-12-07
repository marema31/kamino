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

//PostLoad modify the loaded step values with the values provided in the map in argument
func (st *Step) PostLoad(log *logrus.Entry, superseed map[string]string) error {
	// Nothing to do
	return nil
}

//Load data from step file using its viper representation a return priority and list of steps
func Load(ctx context.Context, log *logrus.Entry, recipePath string, name string, nameIndex int, v *viper.Viper, dss datasource.Datasourcers, force bool) (priority uint, steps []common.Steper, err error) {
	steps = make([]common.Steper, 0, 1)

	priority = v.GetUint("priority")
	tags := v.GetStringSlice("tags")
	if len(tags) == 0 {
		tags = []string{""}
	}

	logStep := log.WithField("name", name).WithField("type", "tmpl")

	templateFile := v.GetString("template")
	if templateFile == "" {
		logStep.Error("No template filename provided")
		return 0, nil, fmt.Errorf("the step %s must have a template to render", name)
	}
	if !filepath.IsAbs(templateFile) {
		templateFile = filepath.Join(recipePath, templateFile)
	}

	template, err := template.New(filepath.Base(templateFile)).Funcs(sprig.FuncMap()).ParseFiles(templateFile)
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
	case "unique":
		mode = Unique
	default:
		mode = Replace
	}

	onlyIfNotExists := v.GetBool("onlyifnotexists")
	zip := v.GetBool("zip")
	gzip := v.GetBool("gzip")

	engines := v.GetStringSlice("engines")
	e, err := datasource.StringsToEngines(engines)
	if err != nil {
		logStep.Error(err)
		return 0, nil, err
	}

	datasourcesByDestinations := make(map[string][]datasource.Datasourcer)
	for _, ds := range dss.Lookup(log, tags, []datasource.Type{datasource.Database}, e) {

		tdestination.Execute(renderedDestination, ds.FillTmplValues())
		destination := renderedDestination.String()
		renderedDestination.Reset()

		if !filepath.IsAbs(destination) {
			destination = filepath.Join(recipePath, destination)
		}

		if _, ok := datasourcesByDestinations[destination]; !ok {
			datasourcesByDestinations[destination] = make([]datasource.Datasourcer, 0, 1)
		}
		datasourcesByDestinations[destination] = append(datasourcesByDestinations[destination], ds)
	}

	index := 0
	for destination, datasources := range datasourcesByDestinations {
		logStep.Debugf("creating step for %s", destination)
		var step Step

		step.Name = fmt.Sprintf("%s:%d", name, nameIndex+index)
		step.templateFile = templateFile
		step.template = template
		step.datasources = datasources
		step.mode = mode
		step.onlyIfNotExists = onlyIfNotExists

		step.destination = destination
		if step.mode == Unique || step.mode == Append {
			step.input.FilePath = destination
			step.input.Gzip = gzip
			step.input.Zip = zip
		}
		step.output.FilePath = destination
		step.output.Gzip = gzip
		step.output.Zip = zip
		steps = append(steps, &step)
		index++
	}

	return priority, steps, nil
}
