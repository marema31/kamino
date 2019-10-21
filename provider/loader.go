package provider

import (
	"context"
	"fmt"

	"github.com/marema31/kamino/datasource"

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

//NewLoader analyze the datasource and return object implemnting Loader of the asked type
func NewLoader(ctx context.Context, ds *datasource.Datasource, table string, where string) (Loader, error) {
	switch ds.Engine {
	case datasource.Mysql, datasource.Postgres:
		return database.NewLoader(ctx, ds, table, where)
	case datasource.CSV:
		return csv.NewLoader(ctx, ds)
	case datasource.JSON:
		return json.NewLoader(ctx, ds)
	case datasource.YAML:
		return yaml.NewLoader(ctx, ds)
	default:
		return nil, fmt.Errorf("don't know how to manage this datasource engine")
	}
}
