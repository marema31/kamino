package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/marema31/kamino/config"
	"github.com/marema31/kamino/kaminodb"
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
	kdb          *kaminodb.KaminoDb
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
	ids          map[string]bool
	ctx          context.Context
}

//NewSaver open the database connection, prepare the insert statement and return a Saver compatible object
func NewSaver(ctx context.Context, config *config.Config, saverConfig config.DestinationConfig, environment string, instances []string) ([]*DbSaver, error) {
	var dss []*DbSaver
	/*TODO: uncomment and adapt
	var err error
	if saverConfig.Database == "" {
		return nil, fmt.Errorf("destination of sync does not provided a database")
	}

	kdbs, err := config.GetDbs(saverConfig.Database, environment, instances)
	if err != nil {
		return nil, err
	}
	for _, kdb := range kdbs {
		var ds DbSaver

		if saverConfig.Table == "" {
			return nil, fmt.Errorf("destination of sync does not provided a table name")
		}

		ds.table = saverConfig.Table
		if kdb.Schema != "" {
			ds.table = fmt.Sprintf("%s.%s", kdb.Schema, saverConfig.Table)
		}

		ds.db, err = kdb.Open()
		if err != nil {
			return nil, fmt.Errorf("can't open %s database : %v", kdb.Database, err)
		}

		ds.kdb = kdb
		ds.database = kdb.Database
		ds.ids = make(map[string]bool)
		ds.ctx = ctx

		err = ds.parseConfig(saverConfig)
		if err != nil {
			return nil, err
		}

		if ds.mode == replace || ds.mode == exactCopy {
			if ds.key == "" {
				return nil, fmt.Errorf("modes replace and exactCopy need a primary key for %s.%s", ds.database, ds.table)
			}
			err = ds.createIdsList()
			if err != nil {
				return nil, err
			}
		}
		dss = append(dss, &ds)
	}
	*/
	return dss, nil
}

//Save writes the record to the destination
func (ds *DbSaver) Save(record common.Record) error {
	var err error

	// Is this method is called for the first time
	//If yes fix the column order in csv file
	if ds.colNames == nil {
		if ds.kdb.Transaction {
			ds.tx, err = ds.db.Begin()
			if err != nil {
				return err
			}
		}

		err := ds.createStatement(record)
		if err != nil {
			return err
		}
		// The truncate will be done at the first record save to avoid truncate a table if there is an error on config file
		if ds.mode == truncate {
			if ds.kdb.Transaction {
				_, err = ds.tx.Exec(fmt.Sprintf("TRUNCATE TABLE %s", ds.table))
			} else {
				_, err = ds.db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", ds.table))
			}
			if err != nil {
				return err
			}
			ds.wasEmpty = true // Avoid truncate after inserting the first record
		}
	}
	if ds.mode == onlyIfEmpty && !ds.wasEmpty {
		return nil
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
				var err error
				if ds.kdb.Transaction {
					_, err = ds.tx.Exec(fmt.Sprintf("TRUNCATE TABLE %s", ds.table))
				} else {
					_, err = ds.db.Exec(fmt.Sprintf("DELETE from %s WHERE %s=%s", ds.table, ds.key, id))
				}
				if err != nil {
					return err
				}
			}
		}
	}
	if ds.kdb.Transaction {
		ds.tx.Commit()
	}

	ds.db.Close()
	return nil
}

//Reset reinitialize the destination (if possible)
func (ds *DbSaver) Reset() error {
	ds.colNames = nil

	if ds.kdb.Transaction && ds.tx != nil {
		ds.tx.Rollback()
	}
	return nil
}

//Name give the name of the destination
func (ds *DbSaver) Name() string {
	return ds.database + "_" + ds.table
}
