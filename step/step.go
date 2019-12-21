//Package step manage the loading of step files and the creation of a list of steps that will be runned by the recipe
package step

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/mb0/glob"
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
	Load(context.Context, *logrus.Entry, string, string, datasource.Datasourcers, provider.Provider, []string, []string, []string, bool, bool) (uint, []common.Steper, error)
}

// Factory implements the StepCreated and use configuration files to create the steps
type Factory struct {
	indexes map[string]int
}

func normalizeStepType(stepType string) (string, error) {
	switch stepType {
	case "shell":
		return "shell", nil
	case "migration":
		return "migration", nil
	case "sync", "synchro", "synchronization":
		return "synchronization", nil
	case "template", "tmpl":
		return "template", nil
	case "sql", "sqlscript":
		return "sqlscript", nil
	default:
		return "", fmt.Errorf("does not how to manage %s step type", stepType)
	}

}

func (sf *Factory) nameInStepNames(log *logrus.Entry, name string, stepNames []string) (bool, error) {
	g, err := glob.New(glob.Default())
	if err != nil {
		log.Errorf("Unable to initialize the globbing engine: %v", err)
		return false, err
	}
	for _, testedName := range stepNames {
		n := strings.ToLower(name)
		t := strings.ToLower(testedName)
		if strings.EqualFold(t, n) {
			return true, nil
		}
		matched, err := g.Match(t, n)
		if err != nil {
			log.Errorf("Using the step name %s failed: %v", stepNames, err)
			return false, err
		}
		if matched {
			return true, nil
		}
	}
	return false, nil
}

// Load the step file and returns the priority and a list of steper for this file
func (sf *Factory) Load(ctx context.Context, log *logrus.Entry, recipePath string, filename string, dss datasource.Datasourcers, prov provider.Provider, limitedTags []string, stepNames []string, stepTypes []string, force bool, dryRun bool) (priority uint, stepList []common.Steper, err error) {
	v := viper.New()
	v.SetConfigName(filename) // The file will be named [filename].json, [filename].yaml or [filename.toml]
	stepsFolder := filepath.Join(recipePath, "steps")
	v.AddConfigPath(stepsFolder)
	err = v.ReadInConfig()
	if err != nil {
		log.Errorf("Unable to parse configuration: %v", err)
		return 0, nil, err
	}

	name := v.GetString("name")
	if sf.indexes == nil {
		sf.indexes = make(map[string]int)
	}

	if len(stepNames) != 0 {
		matched, err := sf.nameInStepNames(log, name, stepNames)
		if err != nil {
			return 0, nil, err
		}
		if !matched {
			return 0, make([]common.Steper, 0), nil
		}
	}

	nameIndex, ok := sf.indexes[name]
	if !ok {
		nameIndex = 0
	}

	stepType := strings.ToLower(v.GetString("type"))
	if stepType == "" {
		log.Errorf("Step type is empty")
		return 0, nil, fmt.Errorf("the step %s does not provide the type", filename)
	}
	if stepType, err = normalizeStepType(stepType); err != nil {
		log.Errorf("Do not know how to manage %s step type", stepType)
		return 0, nil, err
	}

	if len(stepTypes) != 0 {
		normalizedStepTypes := make([]string, 0, len(stepTypes))
		for _, testedType := range stepTypes {
			if testedType, err = normalizeStepType(testedType); err != nil {
				log.Errorf("Do not know how to filter on %s step type", testedType)
				return 0, nil, err
			}
			normalizedStepTypes = append(normalizedStepTypes, testedType)
		}
		found := false
		for _, testedType := range normalizedStepTypes {
			if strings.EqualFold(strings.ToLower(stepType), strings.ToLower(testedType)) {
				found = true
			}
		}
		if !found {
			return 0, make([]common.Steper, 0), nil
		}
	}

	logStep := log.WithField("type", stepType)
	switch stepType {
	case "shell":
		priority, stepList, err = shell.Load(ctx, logStep, recipePath, name, nameIndex, v, dss, force, dryRun, limitedTags)
	case "migration":
		priority, stepList, err = migration.Load(ctx, logStep, recipePath, name, nameIndex, v, dss, force, dryRun, limitedTags)
	case "synchronization":
		priority, stepList, err = sync.Load(ctx, logStep, recipePath, name, nameIndex, v, dss, prov, force, dryRun, limitedTags)
	case "template":
		priority, stepList, err = tmpl.Load(ctx, logStep, recipePath, name, nameIndex, v, dss, force, dryRun, limitedTags)
	case "sqlscript":
		priority, stepList, err = sqlscript.Load(ctx, logStep, recipePath, name, nameIndex, v, dss, force, dryRun, limitedTags)
	}
	if err != nil {
		logStep.Error("Parsing step configuration failed")
	}
	logStep.Debugf("Created %d steps at priority %d", len(stepList), priority)
	sf.indexes[name] = nameIndex + len(stepList)
	return priority, stepList, err
}
