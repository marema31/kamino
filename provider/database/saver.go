package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/provider/types"
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
	ds           datasource.Datasourcer
	db           *sql.DB
	tx           *sql.Tx
	database     string
	table        string
	insertString string
	insertStmt   *sql.Stmt
	updateString string
	updateStmt   *sql.Stmt
	colNames     []string
	mode         dbSaverMode
	wasEmpty     bool
	key          string
	transaction  bool
	engine       datasource.Engine
	ids          map[string]bool
	ctx          context.Context
}

//NewSaver open the database connection, prepare the insert statement and return a Saver compatible object
func NewSaver(ctx context.Context, ds datasource.Datasourcer, table string, key string, mode string) (*DbSaver, error) {
	var saver DbSaver
	tv := ds.FillTmplValues()

	saver.ds = ds
	saver.ctx = ctx
	saver.database = tv.Database
	saver.transaction = tv.Transaction
	saver.engine, _ = datasource.StringToEngine(tv.Engine)

	if table == "" {
		return nil, fmt.Errorf("destination of sync does not provided a table name")
	}
	if tv.Schema != "" {
		table = fmt.Sprintf("%s.%s", tv.Schema, table)
	}
	saver.table = table

	db, err := ds.OpenDatabase(false, false)
	if err != nil {
		return nil, fmt.Errorf("can't open %s database : %v", tv.Database, err)
	}
	saver.db = db
	saver.key = key
	saver.mode = stringToMode(mode)

	saver.ids = make(map[string]bool)

	if saver.mode == replace || saver.mode == exactCopy || saver.mode == update {
		if saver.key == "" {
			return nil, fmt.Errorf("modes replace and exactCopy need a primary key for %s.%s", saver.database, saver.table)
		}
		err = saver.createIdsList()
		if err != nil {
			return nil, err
		}
	}
	return &saver, nil
}

//Save writes the record to the destination
func (saver *DbSaver) Save(record types.Record) error {
	var err error

	// Is this method is called for the first time
	//If yes fix the column order in csv file
	if saver.colNames == nil {
		if saver.transaction {
			saver.tx, err = saver.db.Begin()
			if err != nil {
				return err
			}
		}

		err := saver.createStatement(record)
		if err != nil {
			return err
		}
		// The truncate will be done at the first record save to avoid truncate a table if there is an error on config file
		if saver.mode == truncate {
			if saver.transaction {
				_, err = saver.tx.Exec(fmt.Sprintf("TRUNCATE TABLE %s", saver.table))
			} else {
				_, err = saver.db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", saver.table))
			}
			if err != nil {
				return err
			}
			saver.wasEmpty = true // Avoid truncate after inserting the first record
		}
	}
	if saver.mode == onlyIfEmpty && !saver.wasEmpty {
		return nil
	}

	row := make([]interface{}, len(saver.colNames))
	for i, col := range saver.colNames {
		row[i] = record[col]
	}
	switch saver.mode {
	case onlyIfEmpty:
		_, err := saver.insertStmt.Exec(row...)
		return err

	case insert:
		_, err := saver.insertStmt.Exec(row...)
		return err

	case truncate:
		_, err := saver.insertStmt.Exec(row...)
		return err

	case update:
		_, err := saver.updateStmt.Exec(row...)
		return err

	case replace:
		_, ok := saver.ids[record[saver.key]]
		saver.ids[record[saver.key]] = true
		if ok {
			_, err := saver.updateStmt.Exec(row...)
			return err
		}
		_, err := saver.insertStmt.Exec(row...)
		return err

	case exactCopy:
		_, ok := saver.ids[record[saver.key]]
		saver.ids[record[saver.key]] = true
		if ok {
			_, err := saver.updateStmt.Exec(row...)
			return err
		}
		_, err := saver.insertStmt.Exec(row...)
		return err

	}
	return nil
}

//Close closes the destination
func (saver *DbSaver) Close() error {
	if saver.mode == exactCopy {
		for id, modified := range saver.ids {
			if !modified {
				var err error
				if saver.transaction {
					_, err = saver.tx.Exec(fmt.Sprintf("DELETE from %s WHERE %s=%s", saver.table, saver.key, id))
				} else {
					_, err = saver.db.Exec(fmt.Sprintf("DELETE from %s WHERE %s=%s", saver.table, saver.key, id))
				}
				if err != nil {
					return err
				}
			}
		}
	}
	if saver.transaction {
		saver.tx.Commit()
	}

	saver.db.Close()
	return nil
}

//Reset reinitialize the destination (if possible)
func (saver *DbSaver) Reset() error {
	saver.colNames = nil

	if saver.transaction && saver.tx != nil {
		saver.tx.Rollback()
	}
	return nil
}

//Name give the name of the destination
func (saver *DbSaver) Name() string {
	return saver.database + "_" + saver.table
}
