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

//Saver provides way to save record by record
type Saver interface {
	Save(common.Record) error
	Close() error
	Reset() error
	Name() string
}

//NewSaver analyze the config map and return object implemnting Saver of the asked type
func NewSaver(ctx context.Context, config *config.Config, saverConfig map[string]string) (Saver, error) {
	_, ok := saverConfig["type"]
	if !ok {
		return nil, fmt.Errorf("the configuration block for this destination does not provide the type")
	}

	switch saverConfig["type"] {
	case "database":
		return database.NewSaver(ctx, config, saverConfig)
	case "csv":
		writer, name, tmpName, err := common.OpenWriter(saverConfig)
		if err != nil {
			return nil, err
		}
		return csv.NewSaver(ctx, saverConfig, name, tmpName, writer)
	case "json":
		writer, name, tmpName, err := common.OpenWriter(saverConfig)
		if err != nil {
			return nil, err
		}
		return json.NewSaver(ctx, saverConfig, name, tmpName, writer)
	case "yaml":
		writer, name, tmpName, err := common.OpenWriter(saverConfig)
		if err != nil {
			return nil, err
		}
		return yaml.NewSaver(ctx, saverConfig, name, tmpName, writer)
	default:
		return nil, fmt.Errorf("don't know how to manage %s", saverConfig["type"])
	}
}
