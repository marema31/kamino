package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/marema31/kamino/provider/common"
)

type dbSaverMode int

const (
	onlyIfEmpty dbSaverMode = iota // Will insert only if database was empty
	insert      dbSaverMode = iota // Will insert all line from source (may break if primary key already present)
	replace     dbSaverMode = iota // Will update if line with same primary exist or insert the line
	exactCopy   dbSaverMode = iota // As replace but will remove line with primary key not present in source
)

//DbSaver specifc state for database Saver provider
type DbSaver struct {
	kaminoDb
	stmt     *sql.Stmt
	colNames []string
	mode     dbSaverMode
	wasEmpty bool
}

//NewSaver open the database connection, prepare the insert statement and return a Saver compatible object
func NewSaver(ctx context.Context, config map[string]string) (*DbSaver, error) {
	k, err := newKaminoDb(config)
	if err != nil {
		return nil, err
	}

	mode := exactCopy
	modestr, ok := config["mode"]
	if ok {
		switch {
		case strings.EqualFold(modestr, "onlyifempty"):
			mode = onlyIfEmpty
		case strings.EqualFold(modestr, "insert"):
			mode = insert
		case strings.EqualFold(modestr, "replace"):
			mode = replace
		case strings.EqualFold(modestr, "copy"):
			mode = exactCopy
		}
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

	wasEmpty := !rows.Next()
	if mode == onlyIfEmpty && !wasEmpty {
		log.Printf("Warning: the table %s of database %s is not empty, I will do nothing on it", k.table, k.database)
	}

	stmt, err := k.db.Prepare(fmt.Sprintf("INSERT INTO %s ( %s) VALUES ( %s )", k.table, strings.Join(columnsname[:], ","), strings.Join(questionmark[:], ",")))
	if err != nil {
		return nil, err
	}

	return &DbSaver{*k, stmt, columnsname, mode, wasEmpty}, nil
}

//Save writes the record to the destination
func (ds *DbSaver) Save(record common.Record) error {
	if ds.mode == onlyIfEmpty && !ds.wasEmpty {
		return nil
	}

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
