package provider

import (
	"context"
	"fmt"

	"github.com/marema31/kamino/config"
	"github.com/marema31/kamino/provider/common"
	"github.com/marema31/kamino/provider/csv"
	"github.com/marema31/kamino/provider/database"
	"github.com/marema31/kamino/provider/json"
	"github.com/marema31/kamino/provider/yaml"
)

//Loader provides way to load record by record
type Loader interface {
	Next() bool
	Load() (common.Record, error)
	Close()
	Name() string
}

//NewLoader analyze the config map and return object implemnting Loader of the asked type
func NewLoader(ctx context.Context, config *config.Config, loaderConfig map[string]string) (Loader, error) {
	_, ok := loaderConfig["type"]
	if !ok {
		return nil, fmt.Errorf("the configuration block for this source does not provide the type")
	}

	switch loaderConfig["type"] {
	case "database":
		return database.NewLoader(ctx, config, loaderConfig)
	case "csv":
		reader, name, err := common.OpenReader(loaderConfig)
		if err != nil {
			return nil, err
		}
		return csv.NewLoader(ctx, loaderConfig, name, reader)
	case "json":
		reader, name, err := common.OpenReader(loaderConfig)
		if err != nil {
			return nil, err
		}
		return json.NewLoader(ctx, loaderConfig, name, reader)
	case "yaml":
		reader, name, err := common.OpenReader(loaderConfig)
		if err != nil {
			return nil, err
		}
		return yaml.NewLoader(ctx, loaderConfig, name, reader)
	default:
		return nil, fmt.Errorf("don't know how to manage %s", loaderConfig["type"])
	}
}
