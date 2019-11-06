//Package step manage the loading of step files and the creation of a list of steps that will be runned by the recipe
package step

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/provider"
	"github.com/marema31/kamino/step/common"
	"github.com/marema31/kamino/step/migration"
	"github.com/marema31/kamino/step/shell"
	"github.com/marema31/kamino/step/sqlscript"
	"github.com/marema31/kamino/step/sync"
	"github.com/marema31/kamino/step/tmpl"
)

// Creater is an interface to an object able to create Steper from configuration
type Creater interface {
	Load(context.Context, *logrus.Entry, string, string, datasource.Datasourcers, provider.Provider, []string, []string) (uint, []common.Steper, error)
}

// Factory implements the StepCreated and use configuration files to create the steps
type Factory struct{}

// Load the step file and returns the priority and a list of steper for this file
func (sf Factory) Load(ctx context.Context, log *logrus.Entry, recipePath string, filename string, dss datasource.Datasourcers, prov provider.Provider, stepNames []string, stepTypes []string) (priority uint, stepList []common.Steper, err error) {
	v := viper.New()
	v.SetConfigName(filename) // The file will be named [filename].json, [filename].yaml or [filename.toml]
	stepsFolder := filepath.Join(recipePath, "steps")
	v.AddConfigPath(stepsFolder)
	err = v.ReadInConfig()
	if err != nil {
		log.Errorf("Unable to parse configuration: %v", err)
		return 0, nil, err
	}

	stepType := strings.ToLower(v.GetString("type"))
	if stepType == "" {
		log.Errorf("Step type is empty")
		return 0, nil, fmt.Errorf("the step %s does not provide the type", filename)
	}

	//TODO: Implement limitation of tag/engine from the CLI, pass them to the step that should pass them to the datastore.Lookup
	logStep := log.WithField("type", stepType)
	//by sending stepNames and  stepTypes to the Load functions
	switch stepType {
	case "shell":
		priority, stepList, err = shell.Load(ctx, logStep, recipePath, filename, v, dss)
	case "migration":
		priority, stepList, err = migration.Load(ctx, logStep, recipePath, filename, v, dss)
	case "sync", "synchro", "synchronization":
		priority, stepList, err = sync.Load(ctx, logStep, recipePath, filename, v, dss, prov)
	case "template":
		priority, stepList, err = tmpl.Load(ctx, logStep, recipePath, filename, v, dss)
	case "sql", "sqlscript":
		priority, stepList, err = sqlscript.Load(ctx, logStep, recipePath, filename, v, dss)
	default:
		log.Errorf("Do not know how to manage %s step type", stepType)
		return 0, nil, fmt.Errorf("does not how to manage %s step type", stepType)
	}
	if err != nil {
		logStep.Error("Parsing step configuration failed")
	}
	logStep.Debugf("Created %d steps at priority %d", len(stepList), priority)
	return priority, stepList, err
}
