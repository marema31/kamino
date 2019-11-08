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

// Load Lookup the provided folder for recipes folder and will return a Cookbook of the selected recipes/steps.
// For each recipe, it will load all datasources and the selected steps
func (ck *Cookbook) Load(ctx context.Context, log *logrus.Entry, configPath string, recipes []string, stepNames []string, stepTypes []string) error {
	for _, rname := range recipes {
		logRecipe := log.WithField("recipe", rname)
		err := ck.loadOneRecipe(ctx, logRecipe, configPath, rname, stepNames, stepTypes)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ck *Cookbook) loadOneRecipe(ctx context.Context, log *logrus.Entry, configPath string, rname string, stepNames []string, stepTypes []string) error {
	prov := &provider.KaminoProvider{}
	log.Info("Reading datasources")
	dss := datasource.New()
	recipePath := filepath.Join(configPath, rname)
	if err := dss.LoadAll(recipePath, log); err != nil {
		return err
	}

	log.Info("Reading steps")
	stepsFolder := filepath.Join(recipePath, "steps")
	files, err := ioutil.ReadDir(stepsFolder)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.Mode().IsRegular() {
			ext := filepath.Ext(file.Name())
			if ext == ".yml" || ext == ".yaml" || ext == ".json" || ext == ".toml" {
				name := strings.TrimSuffix(file.Name(), ext)
				logRecipe := log.WithField("step", name)
				logRecipe.Debug("Parsing step configuration")
				priority, steps, err := ck.stepFactory.Load(ctx, logRecipe, recipePath, name, dss, prov, stepNames, stepTypes)
				if err != nil {
					return err
				}
				if _, ok := ck.Recipes[rname]; !ok {
					s := make(map[uint][]common.Steper)
					s[priority] = make([]common.Steper, 0)
					ck.Recipes[rname] = recipe{name: rname, steps: s, currentPriority: 0, dss: dss}
				}
				ck.Recipes[rname].steps[priority] = append(ck.Recipes[rname].steps[priority], steps...)
			}
		}
	}
	return nil
}
