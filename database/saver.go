package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/marema31/kamino/provider"
)

type dbSaver struct {
	kaminoDb
	stmt     *sql.Stmt
	colNames []string
}

func NewSaver(ctx context.Context, c *ConnectionInfo, table string) (*dbSaver, error) {
	k, err := New(c, table)
	if err != nil {
		return nil, err
	}
	rows, err := k.db.QueryContext(ctx, fmt.Sprintf("SELECT * from %s LIMIT 1", table)) // We don't need data, we only needs the column names
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

	stmt, err := k.db.Prepare(fmt.Sprintf("INSERT INTO %s ( %s) VALUES ( %s )", table, strings.Join(columnsname[:], ","), strings.Join(questionmark[:], ",")))
	if err != nil {
		log.Fatal(err)
	}

	return &dbSaver{*k, stmt, columnsname}, nil
}

func (ds *dbSaver) Save(record provider.Record) error {
	row := make([]interface{}, len(ds.colNames))
	for i, col := range ds.colNames {
		row[i] = record[col]
	}

	_, err := ds.stmt.Exec(row...)
	return err
}

func (ds *dbSaver) Close() {
	ds.db.Close()
	return
}
