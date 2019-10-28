package recipe

import (
	"context"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/provider"
	"github.com/marema31/kamino/step/common"
)

//TODO: Implement limitation of tag/engine from the CLI, pass them to the step that should pass them to the datastore.Lookup

//Load Lookup the provided folder for datasource configuration files
func (ck *Cookbook) Load(ctx context.Context, configPath string, recipes []string, stepNames []string, stepTypes []string) error {
	for _, rname := range recipes {
		err := ck.loadOneRecipe(ctx, configPath, rname, stepNames, stepTypes)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ck *Cookbook) loadOneRecipe(ctx context.Context, configPath string, rname string, stepNames []string, stepTypes []string) error {
	prov := &provider.KaminoProvider{}
	dss := datasource.New()
	baseFolder := filepath.Join(configPath, rname)
	if err := dss.LoadAll(baseFolder); err != nil {
		return err
	}
	stepsFolder := filepath.Join(configPath, rname, "steps")
	files, err := ioutil.ReadDir(stepsFolder)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.Mode().IsRegular() {
			ext := filepath.Ext(file.Name())
			if ext == ".yml" || ext == ".yaml" || ext == ".json" || ext == ".toml" {
				name := strings.TrimSuffix(file.Name(), ext)
				//TODO: pass the stepNames and stepTypes to step.Load that should do the filtering
				priority, steps, err := ck.stepFactory.Load(ctx, stepsFolder, name, dss, prov)
				if err != nil {
					return err
				}
				if _, ok := ck.Recipes[rname]; !ok {
					s := make(map[uint][]common.Steper)
					s[priority] = make([]common.Steper, 0)
					ck.Recipes[rname] = recipe{name: rname, steps: s, currentPriority: 0}
				}
				ck.Recipes[rname].steps[priority] = append(ck.Recipes[rname].steps[priority], steps...)
			}
		}
	}
	return nil
}
