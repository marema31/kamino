package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/marema31/kamino/provider/common"
)

//DbLoader specifc state for database Loader provider
type DbLoader struct {
	kaminoDb
	rows     *sql.Rows
	scanned  []interface{}
	rawBytes []sql.RawBytes
	colNames []string
}

//NewLoader open the database connection, make the data query and return a Loader compatible object
func NewLoader(ctx context.Context, config map[string]string) (*DbLoader, error) {
	k, err := newKaminoDb(config)
	if err != nil {
		return nil, err
	}

	where := k.where
	if k.where != "" {
		where = fmt.Sprintf("WHERE %s", k.where)
	}
	rows, err := k.db.QueryContext(ctx, fmt.Sprintf("SELECT * from %s %s", k.table, where))
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

	return &DbLoader{*k, rows, scanned, rawBytes, columnsname}, nil
}

//Next moves to next record and return false if there is no more records
func (dl *DbLoader) Next() bool {
	return dl.rows.Next()
}

//Load reads the next record and return it
func (dl *DbLoader) Load() (common.Record, error) {
	err := dl.rows.Scan(dl.scanned...)
	if err != nil {
		return nil, err
	}

	// Rows.Err will report the last error encountered by Rows.Scan.
	if err = dl.rows.Err(); err != nil {
		return nil, err
	}

	record := make(common.Record, len(dl.colNames))
	for i, col := range dl.colNames {
		record[col] = string(dl.rawBytes[i])
	}

	return record, nil

}

//Close closes the datasource
func (dl *DbLoader) Close() {
	dl.rows.Close()
	dl.db.Close()
}

//Name give the name of the destination
func (dl *DbLoader) Name() string {
	return dl.database + "_" + dl.table
}
