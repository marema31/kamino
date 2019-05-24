package provider

import (
	"context"
	"fmt"

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
func NewLoader(ctx context.Context, config map[string]string) (Loader, error) {
	_, ok := config["type"]
	if !ok {
		return nil, fmt.Errorf("the configuration block for this source does not provide the type")
	}

	switch config["type"] {
	case "database":
		return database.NewLoader(ctx, config)
	case "csv":
		reader, name, err := common.OpenReader(config)
		if err != nil {
			return nil, err
		}
		return csv.NewLoader(ctx, config, name, reader)
	case "json":
		reader, name, err := common.OpenReader(config)
		if err != nil {
			return nil, err
		}
		return json.NewLoader(ctx, config, name, reader)
	case "yaml":
		reader, name, err := common.OpenReader(config)
		if err != nil {
			return nil, err
		}
		return yaml.NewLoader(ctx, config, name, reader)
	default:
		return nil, fmt.Errorf("don't know how to manage %s", config["type"])
	}
}
