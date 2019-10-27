//Package step manage the loading of step files and the creation of a list of steps that will be runned by the recipe
package step

import (
	"context"
	"fmt"
	"strings"

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

// Load the step file and returns the priority and a list of steper for this file
func Load(ctx context.Context, path string, filename string, dss datasource.Datasourcers, prov provider.Provider) (priority uint, stepList []common.Steper, err error) {
	v := viper.New()

	v.SetConfigName(filename) // The file will be named [filename].json, [filename].yaml or [filename.toml]
	v.AddConfigPath(path)
	err = v.ReadInConfig()

	if err != nil {
		return 0, nil, err
	}

	stepType := strings.ToLower(v.GetString("type"))
	if stepType == "" {
		return 0, nil, fmt.Errorf("the step %s does not provide the type", filename)
	}

	switch stepType {
	case "sql":
		return shell.Load(ctx, filename, v, dss)
	case "migration":
		return migration.Load(ctx, filename, v, dss)
	case "sync", "synchro", "synchronization":
		return sync.Load(ctx, filename, v, dss, prov)
	case "template":
		return tmpl.Load(ctx, filename, v, dss)
	case "shell":
		return sqlscript.Load(ctx, filename, v, dss)
	default:
		return 0, nil, fmt.Errorf("does not how to manage %s step type", stepType)
	}
}
