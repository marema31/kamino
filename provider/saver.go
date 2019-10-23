package provider

import (
	"context"
	"fmt"

	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/provider/csv"
	"github.com/marema31/kamino/provider/database"
	"github.com/marema31/kamino/provider/json"
	"github.com/marema31/kamino/provider/types"
	"github.com/marema31/kamino/provider/yaml"
)

//Saver provides way to save record by record
type Saver interface {
	Save(types.Record) error
	Close() error
	Reset() error
	Name() string
}

//NewSaver analyze the datasource and return object implemnting Saver of the asked type
func NewSaver(ctx context.Context, ds datasource.Datasourcer, table string, key string, mode string) (Saver, error) {
	tv := ds.FillTmplValues()
	engine, _ := datasource.StringToEngine(tv.Engine)

	switch engine {
	case datasource.Mysql, datasource.Postgres:
		return database.NewSaver(ctx, ds, table, key, mode)
	case datasource.CSV:
		return csv.NewSaver(ctx, ds)
	case datasource.JSON:
		return json.NewSaver(ctx, ds)
	case datasource.YAML:
		return yaml.NewSaver(ctx, ds)
	default:
		return nil, fmt.Errorf("don't know how to manage this datasource engine")
	}
}
