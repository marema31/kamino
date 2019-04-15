package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/marema31/kamino/provider"
)

type dbLoader struct {
	kaminoDb
	rows     *sql.Rows
	scanned  []interface{}
	rawBytes []sql.RawBytes
	colNames []string
}

func NewLoader(ctx context.Context, c *ConnectionInfo, table string) (*dbLoader, error) {
	k, err := New(c, table)
	if err != nil {
		return nil, err
	}
	rows, err := k.db.QueryContext(ctx, fmt.Sprintf("SELECT * from %s", table))
	if err != nil {
		return nil, err
	}

	columns, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	columnsname := make([]string, len(columns))
	for i, col := range columns {
		columnsname[i] = col.Name()
	}

	rawBytes := make([]sql.RawBytes, len(columns)) // Buffers for each column
	scanned := make([]interface{}, len(columns))   // Adress of each Buffers since sql.QueryContext needs pointer to each column buffers
	for i := range rawBytes {
		scanned[i] = &rawBytes[i]
	}

	return &dbLoader{*k, rows, scanned, rawBytes, columnsname}, nil
}

func (dl *dbLoader) Next() bool {
	return dl.rows.Next()
}

func (dl *dbLoader) Load() (provider.Record, error) {
	err := dl.rows.Scan(dl.scanned...)
	if err != nil {
		return nil, err
	}

	record := make(provider.Record, len(dl.colNames))
	for i, col := range dl.colNames {
		record[col] = string(dl.rawBytes[i])
	}

	return record, nil

}

func (dl *dbLoader) Close() {
	dl.rows.Close()
	dl.db.Close()
	return
}
