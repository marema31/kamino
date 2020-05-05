package provider

import (
	"context"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/provider/common"
	"github.com/marema31/kamino/provider/csv"
	"github.com/marema31/kamino/provider/database"
	"github.com/marema31/kamino/provider/json"
	"github.com/marema31/kamino/provider/types"
	"github.com/marema31/kamino/provider/yaml"
)

//Loader provides way to load record by record.
type Loader interface {
	Next() bool
	Load(*logrus.Entry) (types.Record, error)
	Close(*logrus.Entry) error
	Name() string
}

//NewLoader analyze the datasource and return object implementing Loader of the asked type.
func (p *KaminoProvider) NewLoader(ctx context.Context, log *logrus.Entry, ds datasource.Datasourcer, table string, where string) (Loader, error) {
	engine := ds.GetEngine()

	switch engine {
	case datasource.Mysql, datasource.Postgres:
		return database.NewLoader(ctx, log, ds, table, where)
	case datasource.CSV:
		return csv.NewLoader(ctx, log, ds)
	case datasource.JSON:
		return json.NewLoader(ctx, log, ds)
	case datasource.YAML:
		return yaml.NewLoader(ctx, log, ds)
	default:
		return nil, fmt.Errorf("don't know how to manage this datasource engine: %w", common.ErrWrongParameterValue)
	}
}
