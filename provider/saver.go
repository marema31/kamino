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

//Saver provides way to save record by record.
type Saver interface {
	Save(*logrus.Entry, types.Record) error
	Close(*logrus.Entry) error
	Reset(*logrus.Entry) error
	Name() string
}

//NewSaver analyze the datasource and return object implementing Saver of the asked type.
func (p *KaminoProvider) NewSaver(ctx context.Context, log *logrus.Entry, ds datasource.Datasourcer, table string, key string, mode string) (Saver, error) {
	engine := ds.GetEngine()

	switch engine {
	case datasource.Mysql, datasource.Postgres:
		return database.NewSaver(ctx, log, ds, table, key, mode)
	case datasource.CSV:
		return csv.NewSaver(ctx, log, ds)
	case datasource.JSON:
		return json.NewSaver(ctx, log, ds)
	case datasource.YAML:
		return yaml.NewSaver(ctx, log, ds)
	default:
		return nil, fmt.Errorf("don't know how to manage this datasource engine: %w", common.ErrWrongParameterValue)
	}
}
