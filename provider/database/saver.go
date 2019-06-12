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
	update      dbSaverMode = iota // Will update if line with same primary exist or skip the line
	replace     dbSaverMode = iota // Will update if line with same primary exist or insert the line
	exactCopy   dbSaverMode = iota // As replace but will remove line with primary key not present in source
	truncate    dbSaverMode = iota // As insert but will truncate the table before

)

//DbSaver specifc state for database Saver provider
type DbSaver struct {
	kaminoDb
	insertStmt *sql.Stmt
	updateStmt *sql.Stmt
	colNames   []string
	mode       dbSaverMode
	wasEmpty   bool
	key        string
	ids        map[string]bool
}

//parseConfig parse the config to extract the mode and the primary key and save them in the dbSaver instance
func (ds *DbSaver) parseConfig(config map[string]string) error {
	ds.key = config["key"]

	ds.mode = exactCopy
	modestr, ok := config["mode"]
	if ok {
		switch {
		case strings.EqualFold(modestr, "onlyifempty"):
			ds.mode = onlyIfEmpty
		case strings.EqualFold(modestr, "insert"):
			ds.mode = insert
		case strings.EqualFold(modestr, "update"):
			ds.mode = update
		case strings.EqualFold(modestr, "replace"):
			ds.mode = replace
		case strings.EqualFold(modestr, "copy"):
			ds.mode = exactCopy
		case strings.EqualFold(modestr, "truncate"):
			ds.mode = truncate
		}
	}

	if ds.key == "" && (ds.mode == update || ds.mode == replace) {
		return fmt.Errorf("mode for %s.%s is %s and no key is provided", ds.database, ds.table, modestr)
	}
	return nil
}

//NewSaver open the database connection, prepare the insert statement and return a Saver compatible object
func NewSaver(ctx context.Context, config map[string]string) (*DbSaver, error) {
	var ds DbSaver
	var err error

	k, err := newKaminoDb(config)
	if err != nil {
		return nil, err
	}

	ds.db = k.db
	ds.driver = k.driver
	ds.database = k.database
	ds.table = k.table
	ds.where = k.where
	ds.ids = make(map[string]bool)

	err = ds.parseConfig(config)
	if err != nil {
		return nil, err
	}

	err = ds.createStatement(ctx)
	if err != nil {
		return nil, err
	}

	if ds.mode == replace || ds.mode == exactCopy {
		err = ds.createIdsList(ctx)
		if err != nil {
			return nil, err
		}
	}

	return &ds, nil
}

// createStatement Query the destination table to determine the available colums, create the corresponding insert/update statement and save them in the dbSaver instance
func (ds *DbSaver) createStatement(ctx context.Context) error {
	rows, err := ds.db.QueryContext(ctx, fmt.Sprintf("SELECT * from %s LIMIT 1", ds.table)) // We don't need data, we only needs the column names
	if err != nil {
		return err
	}

	columns, err := rows.ColumnTypes()
	if err != nil {
		return err
	}

	var updateSet []string
	var questionmark []string

	for _, col := range columns {
		if !strings.EqualFold(col.Name(), ds.key) {
			questionmark = append(questionmark, "?")
			ds.colNames = append(ds.colNames, col.Name())
			updateSet = append(updateSet, fmt.Sprintf("%s=?", col.Name()))
		}
		//		fmt.Println(col.Name(), col.DatabaseTypeName(), col.ScanType())
	}
	// By doing like this we ensure the primary key will be the last of column names and this array can be use for insert and update

	if ds.key != "" {
		questionmark = append(questionmark, "?")
		ds.colNames = append(ds.colNames, ds.key)
	}

	if len(columns) != len(ds.colNames) {
		return fmt.Errorf("provided key %s is not a column of %s.%s ", ds.key, ds.database, ds.table)
	}

	ds.wasEmpty = !rows.Next()
	if ds.mode == onlyIfEmpty && !ds.wasEmpty {
		log.Printf("Warning: the table %s of database %s is not empty, I will do nothing on it", ds.table, ds.database)
	}

	if ds.mode == truncate {
		// The truncate will be done at the first record save to avoid truncate a table if there is an error on config file
		ds.wasEmpty = false
	}

	ds.insertStmt, err = ds.db.Prepare(fmt.Sprintf("INSERT INTO %s ( %s) VALUES ( %s )", ds.table, strings.Join(ds.colNames[:], ","), strings.Join(questionmark[:], ",")))
	if err != nil {
		return err
	}

	if ds.mode == replace || ds.mode == update || ds.mode == exactCopy {
		ds.updateStmt, err = ds.db.Prepare(fmt.Sprintf("UPDATE %s SET  %s WHERE %s = ?", ds.table, strings.Join(updateSet[:], ","), ds.key))
		if err != nil {
			return err
		}
	}

	return nil
}

// createStatement Query the destination table to determine the available colums, create the corresponding insert/update statement and save them in the dbSaver instance
func (ds *DbSaver) createIdsList(ctx context.Context) error {
	rows, err := ds.db.QueryContext(ctx, fmt.Sprintf("SELECT %s from %s", ds.key, ds.table)) // We don't need data, we only needs the column names
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return err
		}
		ds.ids[id] = false
	}
	return nil
}

//Save writes the record to the destination
func (ds *DbSaver) Save(record common.Record) error {
	if ds.mode == onlyIfEmpty && !ds.wasEmpty {
		return nil
	}

	if ds.mode == truncate && !ds.wasEmpty {
		_, err := ds.db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", ds.table))
		if err != nil {
			return err
		}
		ds.wasEmpty = true // Avoid truncate after inserting the first record
	}

	row := make([]interface{}, len(ds.colNames))
	for i, col := range ds.colNames {
		row[i] = record[col]
	}

	switch ds.mode {
	case onlyIfEmpty:
		_, err := ds.insertStmt.Exec(row...)
		return err

	case insert:
		_, err := ds.insertStmt.Exec(row...)
		return err

	case truncate:
		_, err := ds.insertStmt.Exec(row...)
		return err

	case update:
		_, err := ds.updateStmt.Exec(row...)
		return err

	case replace:
		_, ok := ds.ids[record[ds.key]]
		ds.ids[record[ds.key]] = true
		if ok {
			_, err := ds.updateStmt.Exec(row...)
			return err
		}
		_, err := ds.insertStmt.Exec(row...)
		return err

	case exactCopy:
		_, ok := ds.ids[record[ds.key]]
		ds.ids[record[ds.key]] = true
		if ok {
			_, err := ds.updateStmt.Exec(row...)
			return err
		}
		_, err := ds.insertStmt.Exec(row...)
		return err

	}
	return nil
}

//Close closes the destination
func (ds *DbSaver) Close() error {
	if ds.mode == exactCopy {
		for id, modified := range ds.ids {
			if !modified {
				_, err := ds.db.Exec(fmt.Sprintf("DELETE from %s WHERE %s=%s", ds.table, ds.key, id))
				if err != nil {
					log.Println(err)
					return err
				}
			}
		}
	}

	ds.db.Close()
	return nil
}

//Name give the name of the destination
func (ds *DbSaver) Name() string {
	return ds.database + "_" + ds.table
}
