package recipe

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/provider"
	"github.com/marema31/kamino/step"
)

//TODO: Implement limitation of tag/engine from the CLI, pass them to the step that should pass them to the datastore.Lookup

//Load Lookup the provided folder for datasource configuration files
func Load(ctx context.Context, configPath string, recipes string, dss datasource.Datasourcers, prov provider.Provider) error {
	baseFolder := filepath.Join(configPath, "recipes")
	fmt.Print(baseFolder)
	files, err := ioutil.ReadDir(baseFolder)
	if err != nil {
		log.Fatal(err)
	}
	//TODO: Adapt to load the recipes from the selection recipes from CLI
	for _, file := range files {
		if file.Mode().IsRegular() {
			ext := filepath.Ext(file.Name())
			if ext == ".yml" || ext == ".yaml" || ext == ".json" || ext == ".toml" {
				name := strings.TrimSuffix(file.Name(), ext)
				//TODO: use the priority and steplist
				_, _, err := step.Load(ctx, baseFolder, name, dss, prov)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
