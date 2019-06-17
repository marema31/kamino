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

//Saver provides way to save record by record
type Saver interface {
	Save(common.Record) error
	Close() error
	Reset() error
	Name() string
}

//NewSaver analyze the config map and return object implemnting Saver of the asked type
func NewSaver(ctx context.Context, config map[string]string) (Saver, error) {
	_, ok := config["type"]
	if !ok {
		return nil, fmt.Errorf("the configuration block for this destination does not provide the type")
	}

	switch config["type"] {
	case "database":
		return database.NewSaver(ctx, config)
	case "csv":
		writer, name, err := common.OpenWriter(config)
		if err != nil {
			return nil, err
		}
		return csv.NewSaver(ctx, config, name, writer)
	case "json":
		writer, name, err := common.OpenWriter(config)
		if err != nil {
			return nil, err
		}
		return json.NewSaver(ctx, config, name, writer)
	case "yaml":
		writer, name, err := common.OpenWriter(config)
		if err != nil {
			return nil, err
		}
		return yaml.NewSaver(ctx, config, name, writer)
	default:
		return nil, fmt.Errorf("don't know how to manage %s", config["type"])
	}
}
