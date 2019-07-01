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
func NewSaver(ctx context.Context, config *config.Config, saverConfig config.DestinationConfig, environment string, instances []string) ([]Saver, error) {
	if saverConfig.Type == "" {
		return nil, fmt.Errorf("the configuration block for this destination does not provide the type")
	}

	var ss []Saver

	switch saverConfig.Type {
	case "database":
		ds, err := database.NewSaver(ctx, config, saverConfig, environment, instances)
		for _, d := range ds {
			ss = append(ss, d)
		}
		return ss, err
	case "csv":
		writer, tmpName, err := common.OpenWriter(saverConfig)
		if err != nil {
			return nil, err
		}
		s, err := csv.NewSaver(ctx, saverConfig, tmpName, writer)
		return append(ss, s), err
	case "json":
		writer, tmpName, err := common.OpenWriter(saverConfig)
		if err != nil {
			return nil, err
		}
		s, err := json.NewSaver(ctx, saverConfig, tmpName, writer)
		return append(ss, s), err
	case "yaml":
		writer, tmpName, err := common.OpenWriter(saverConfig)
		if err != nil {
			return nil, err
		}
		s, err := yaml.NewSaver(ctx, saverConfig, tmpName, writer)
		return append(ss, s), err
	default:
		return nil, fmt.Errorf("don't know how to manage %s", saverConfig.Type)
	}
}
