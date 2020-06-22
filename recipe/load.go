package recipe

import (
	"context"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/provider"
	"github.com/marema31/kamino/step/common"
)

//PostLoad modify the loaded step values with the values provided in the map in argument.
func (ck *Cookbook) PostLoad(log *logrus.Entry, superseed map[string]string) error {
	for _, recipe := range ck.Recipes {
		for _, steps := range recipe.steps {
			for _, step := range steps {
				if err := step.PostLoad(log, superseed); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// Load Lookup the provided folder for recipes folder and will return a Cookbook of the selected recipes/steps.
// For each recipe, it will load all datasources and the selected steps.
func (ck *Cookbook) Load(ctx context.Context, log *logrus.Entry, configPath string, recipes []string, limitedTags []string, stepNames []string, stepTypes []string) error {
	var firstError error

	for _, rname := range recipes {
		logRecipe := log.WithField("recipe", rname)

		err := ck.loadOneRecipe(ctx, logRecipe, configPath, rname, limitedTags, stepNames, stepTypes)
		if err != nil {
			if !ck.validate {
				return err
			}

			if firstError == nil {
				firstError = err
			}
		}
	}

	return firstError
}

func (ck *Cookbook) parseStep(ctx context.Context, log *logrus.Entry, dss datasource.Datasourcers, prov provider.Provider, recipePath string, rname string, limitedTags []string, stepNames []string, stepTypes []string, filename string) error {
	var firstError error

	log.Debug("Parsing step configuration")

	priority, forceSequential, steps, err := ck.stepFactory.Load(ctx, log, recipePath, filename, dss, prov, limitedTags, stepNames, stepTypes, ck.force, ck.dryRun)
	if err != nil {
		if !ck.validate {
			return err
		}

		if firstError == nil {
			firstError = err
		}
	}

	if !ck.forcedSequential[priority] {
		ck.forcedSequential[priority] = forceSequential
	}

	if _, ok := ck.Recipes[rname]; !ok {
		s := make(map[uint][]common.Steper)
		s[priority] = make([]common.Steper, 0)
		ck.Recipes[rname] = recipe{name: rname, steps: s, currentPriority: 0, dss: dss}
	}

	ck.Recipes[rname].steps[priority] = append(ck.Recipes[rname].steps[priority], steps...)

	return firstError
}

func (ck *Cookbook) loadOneRecipe(ctx context.Context, log *logrus.Entry, configPath string, rname string, limitedTags []string, stepNames []string, stepTypes []string) error {
	var firstError error

	log.Info("Reading datasources")

	prov := &provider.KaminoProvider{}
	dss := datasource.New(ck.conTimeout, ck.conRetry)
	recipePath := filepath.Join(configPath, rname)

	if err := dss.LoadAll(recipePath, log); err != nil {
		if !ck.validate {
			return err
		}

		if firstError == nil {
			firstError = err
		}
	}

	log.Info("Reading steps")

	stepsFolder := filepath.Join(recipePath, "steps")

	files, err := ioutil.ReadDir(stepsFolder)
	if err != nil {
		return err
	}

	for _, file := range files {
		ext := filepath.Ext(file.Name())
		if file.Mode().IsRegular() && (ext == ".yml" || ext == ".yaml" || ext == ".json" || ext == ".toml") {
			name := strings.TrimSuffix(file.Name(), ext)
			logRecipe := log.WithField("step", name)

			err := ck.parseStep(ctx, logRecipe, dss, prov, recipePath, rname, limitedTags, stepNames, stepTypes, name)
			if err != nil && !ck.validate {
				return err
			}

			if err != nil && firstError == nil {
				firstError = err
			}
		}
	}

	return firstError
}
