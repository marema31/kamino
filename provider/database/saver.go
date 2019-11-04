package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Sirupsen/logrus"
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
func NewSaver(ctx context.Context, log *logrus.Entry, ds datasource.Datasourcer, table string, key string, mode string) (*DbSaver, error) {
	logDb := log.WithField("datasource", ds.GetName())
	var saver DbSaver
	tv := ds.FillTmplValues()

	saver.ds = ds
	saver.ctx = ctx
	saver.database = tv.Database
	saver.transaction = tv.Transaction
	saver.engine, _ = datasource.StringToEngine(tv.Engine)

	if table == "" {
		logDb.Error("No destination table provided")
		return nil, fmt.Errorf("destination of sync does not provided a table name")
	}
	if tv.Schema != "" {
		table = fmt.Sprintf("%s.%s", tv.Schema, table)
	}
	saver.table = table

	db, err := ds.OpenDatabase(logDb, false, false)
	if err != nil {
		return nil, fmt.Errorf("can't open %s database : %v", tv.Database, err)
	}
	saver.db = db
	saver.key = key
	saver.mode = stringToMode(mode)

	saver.ids = make(map[string]bool)

	if saver.mode == replace || saver.mode == exactCopy || saver.mode == update {
		if saver.key == "" {
			logDb.Errorf("Modes replace and exactCopy need a primary key for %s.%s", saver.database, saver.table)
			return nil, fmt.Errorf("modes replace and exactCopy need a primary key for %s.%s", saver.database, saver.table)
		}
		logDb.Debug("Create current IDs list")
		err = saver.createIdsList(logDb)
		if err != nil {
			logDb.Error("Getting current IDs in destination table failed")
			logDb.Error(err)
			return nil, err
		}
	}
	return &saver, nil
}

//Save writes the record to the destination
func (saver *DbSaver) Save(log *logrus.Entry, record types.Record) error {
	logDb := log.WithField("datasource", saver.ds.GetName())
	var err error

	// Is this method is called for the first time
	//If yes fix the column order in csv file
	if saver.colNames == nil {
		logDb.Debug("First Save action, preparing the needed informations")
		if saver.transaction {
			logDb.Debug("Starting transaction")
			saver.tx, err = saver.db.Begin()
			if err != nil {
				logDb.Error("Beginning transaction failed")
				logDb.Error(err)
				return err
			}
		}

		logDb.Debug("Preparing the statements")
		err := saver.createStatement(logDb, record)
		if err != nil {
			return err
		}
		// The truncate will be done at the first record save to avoid truncate a table if there is an error on config file
		if saver.mode == truncate {
			logDb.Debug("Truncating the destination table")
			if saver.transaction {
				_, err = saver.tx.Exec(fmt.Sprintf("TRUNCATE TABLE %s", saver.table))
			} else {
				_, err = saver.db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", saver.table))
			}
			if err != nil {
				logDb.Error("Truncating the destination table failed")
				logDb.Error(err)
				return err
			}
			saver.wasEmpty = true // Avoid truncate after inserting the first record
		}
	}
	if saver.mode == onlyIfEmpty && !saver.wasEmpty {
		if saver.colNames == nil {
			logDb.Warn("Destination table is not empty, we will do nothing")
		}
		return nil
	}

	row := make([]interface{}, len(saver.colNames))
	for i, col := range saver.colNames {
		row[i] = record[col]
	}
	switch saver.mode {
	case onlyIfEmpty:
		_, err = saver.insertStmt.Exec(row...)

	case insert:
		_, err = saver.insertStmt.Exec(row...)

	case truncate:
		_, err = saver.insertStmt.Exec(row...)

	case update:
		_, err = saver.updateStmt.Exec(row...)

	case replace:
		_, ok := saver.ids[record[saver.key]]
		saver.ids[record[saver.key]] = true
		if ok {
			_, err = saver.updateStmt.Exec(row...)
		} else {

			_, err = saver.insertStmt.Exec(row...)
		}

	case exactCopy:
		_, ok := saver.ids[record[saver.key]]
		saver.ids[record[saver.key]] = true
		if ok {
			_, err = saver.updateStmt.Exec(row...)
		} else {
			_, err = saver.insertStmt.Exec(row...)
		}
	}
	if err != nil {
		logDb.Error("Saving row failed")
		logDb.Error(err)
	}
	return err
}

//Close closes the destination
func (saver *DbSaver) Close(log *logrus.Entry) error {
	logDb := log.WithField("datasource", saver.ds.GetName())
	if saver.mode == exactCopy {
		logDb.Debug("Deleting non synchronized rows")
		for id, modified := range saver.ids {
			if !modified {
				var err error
				if saver.transaction {
					_, err = saver.tx.Exec(fmt.Sprintf("DELETE from %s WHERE %s=%s", saver.table, saver.key, id))
				} else {
					_, err = saver.db.Exec(fmt.Sprintf("DELETE from %s WHERE %s=%s", saver.table, saver.key, id))
				}
				if err != nil {
					logDb.Error("Deleting non synchronized rows failed")
					logDb.Error(err)
					return err
				}
			}
		}
	}
	if saver.transaction {
		logDb.Debug("Committing transaction")
		err := saver.tx.Commit()
		if err != nil {
			logDb.Error("Committing transaction failed")
			logDb.Error(err)
		}
	}

	logDb.Debug("Closing database")
	err := saver.db.Close()
	if err != nil {
		logDb.Error("Close database failed")
		logDb.Error(err)
	}
	return err
}

//Reset reinitialize the destination (if possible)
func (saver *DbSaver) Reset(log *logrus.Entry) (err error) {
	logDb := log.WithField("datasource", saver.ds.GetName())
	saver.colNames = nil

	if saver.transaction && saver.tx != nil {
		logDb.Debug("Rollbacking transaction")
		err = saver.tx.Rollback()
		if err != nil {
			logDb.Error("Rollbacking transaction failed")
			logDb.Error(err)
		}
	} else {
		logDb.Debug("Reset needed but without transaction we can not do anything")
	}
	return err
}

//Name give the name of the destination
func (saver *DbSaver) Name() string {
	return saver.database + "_" + saver.table
}
