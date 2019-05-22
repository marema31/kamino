package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/marema31/kamino/provider/common"
)

//DbSaver specifc state for database Saver provider
type DbSaver struct {
	kaminoDb
	stmt     *sql.Stmt
	colNames []string
}

//NewSaver open the database connection, prepare the insert statement and return a Saver compatible object
func NewSaver(ctx context.Context, config map[string]string) (*DbSaver, error) {
	k, err := newKaminoDb(config)
	if err != nil {
		return nil, err
	}
	rows, err := k.db.QueryContext(ctx, fmt.Sprintf("SELECT * from %s LIMIT 1", k.table)) // We don't need data, we only needs the column names
	if err != nil {
		return nil, err
	}

	columns, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	columnsname := make([]string, len(columns))
	questionmark := make([]string, len(columns))
	for i, col := range columns {
		questionmark[i] = "?"
		columnsname[i] = col.Name()
		//		fmt.Println(col.Name(), col.DatabaseTypeName(), col.ScanType())
	}

	stmt, err := k.db.Prepare(fmt.Sprintf("INSERT INTO %s ( %s) VALUES ( %s )", k.table, strings.Join(columnsname[:], ","), strings.Join(questionmark[:], ",")))
	if err != nil {
		return nil, err
	}

	return &DbSaver{*k, stmt, columnsname}, nil
}

//Save writes the record to the destination
func (ds *DbSaver) Save(record common.Record) error {
	row := make([]interface{}, len(ds.colNames))
	for i, col := range ds.colNames {
		row[i] = record[col]
	}

	_, err := ds.stmt.Exec(row...)
	return err
}

//Close closes the destination
func (ds *DbSaver) Close() error {
	ds.db.Close()
	return nil
}

//Name give the name of the destination
func (ds *DbSaver) Name() string {
	return ds.database + "_" + ds.table
}
