package provider

import (
	"context"
	"fmt"

	"github.com/marema31/kamino/provider/common"
	"github.com/marema31/kamino/provider/database"
)

//Saver provides way to save record by record
type Saver interface {
	Save(common.Record) error
	Close()
}

//Loader provides way to load record by record
type Loader interface {
	Next() bool
	Load() (common.Record, error)
	Close()
}

//NewLoader analyze the config map and return object implemnting Loader of the asked type
func NewLoader(ctx context.Context, config map[string]string) (Loader, error) {
	_, ok := config["type"]
	if !ok {
		return nil, fmt.Errorf("the configuration block for this source does not precise the type")
	}

	switch config["type"] {
	case "database":
		return database.NewLoader(ctx, config)
	default:
		return nil, fmt.Errorf("don't know how to manage %s", config["type"])
	}
}

//NewSaver analyze the config map and return object implemnting Saver of the asked type
func NewSaver(ctx context.Context, config map[string]string) (Saver, error) {
	_, ok := config["type"]
	if !ok {
		return nil, fmt.Errorf("the configuration block for this source does not precise the type")
	}

	switch config["type"] {
	case "database":
		return database.NewSaver(ctx, config)
	default:
		return nil, fmt.Errorf("don't know how to manage %s", config["type"])
	}
}
