package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/provider/common"
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

//DbSaver specifc state for database Saver provider.
type DbSaver struct {
	ds           datasource.Datasourcer
	db           *sql.DB
	tx           *sql.Tx
	database     string
	table        string
	rawtable     string
	schema       string
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

//NewSaver open the database connection, prepare the insert statement and return a Saver compatible object.
func NewSaver(ctx context.Context, log *logrus.Entry, ds datasource.Datasourcer, table string, key string, mode string) (*DbSaver, error) {
	logDb := log.WithField("datasource", ds.GetName())
	tv := ds.FillTmplValues()

	var saver DbSaver

	saver.ds = ds
	saver.ctx = ctx
	saver.database = tv.Database
	saver.transaction = tv.Transaction
	saver.engine, _ = datasource.StringToEngine(tv.Engine)

	if table == "" {
		logDb.Error("No destination table provided")
		return nil, fmt.Errorf("destination of sync does not provided a table name: %w", common.ErrMissingParameter)
	}

	saver.rawtable = table
	saver.schema = tv.Schema

	if tv.Schema != "" {
		table = fmt.Sprintf("%s.%s", tv.Schema, table)
	}

	saver.table = table

	db, err := ds.OpenDatabase(logDb, false, false)
	if err != nil {
		return nil, fmt.Errorf("can't open %s database : %w", tv.Database, err)
	}

	saver.db = db
	saver.key = key
	saver.mode = stringToMode(mode)

	saver.ids = make(map[string]bool)

	if saver.mode == replace || saver.mode == exactCopy || saver.mode == update {
		if saver.key == "" {
			logDb.Errorf("Modes replace and exactCopy need a primary key for %s.%s", saver.database, saver.table)
			return nil, fmt.Errorf("modes replace and exactCopy need a primary key for %s.%s: %w", saver.database, saver.table, common.ErrMissingParameter)
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

func (saver *DbSaver) initSave(log *logrus.Entry, record types.Record) error {
	log.Debug("First Save action, preparing the needed informations")

	var err error

	if saver.transaction {
		log.Debug("Starting transaction")

		saver.tx, err = saver.db.Begin()
		if err != nil {
			log.Error("Beginning transaction failed")
			log.Error(err)

			return err
		}
	}

	log.Debug("Preparing the statements")

	err = saver.createStatement(log, record)
	if err != nil {
		return err
	}
	// The truncate will be done at the first record save to avoid truncate a table if there is an error on config file
	if saver.mode == truncate {
		log.Debug("Truncating the destination table")

		if saver.transaction {
			_, err = saver.tx.Exec(fmt.Sprintf("TRUNCATE TABLE %s", saver.table))
		} else {
			_, err = saver.db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", saver.table))
		}

		if err != nil {
			log.Error("Truncating the destination table failed")
			log.Error(err)

			return err
		}

		saver.wasEmpty = true // Avoid truncate after inserting the first record
	}

	return nil
}

//Save writes the record to the destination.
func (saver *DbSaver) Save(log *logrus.Entry, record types.Record) error {
	logDb := log.WithField("datasource", saver.ds.GetName())

	var err error

	// Is this method is called for the first time
	if saver.colNames == nil {
		err = saver.initSave(log, record)
		if err != nil {
			return err
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
		if record[col] != types.NullValue {
			row[i] = record[col]
		} else {
			row[i] = sql.NullString{}
		}
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
		if ok {
			_, err = saver.updateStmt.Exec(row...)
		} else {
			_, err = saver.insertStmt.Exec(row...)
		}

		saver.ids[record[saver.key]] = true
	case exactCopy:
		_, ok := saver.ids[record[saver.key]]
		if ok {
			_, err = saver.updateStmt.Exec(row...)
		} else {
			_, err = saver.insertStmt.Exec(row...)
		}

		saver.ids[record[saver.key]] = true
	}

	if err != nil {
		logDb.Error("Saving row failed")
		logDb.Error(err)
	}

	return err
}

func (saver *DbSaver) removeNonSynchronized(log *logrus.Entry) error {
	log.Debug("Deleting non synchronized rows")

	for id, modified := range saver.ids {
		if !modified {
			var err error

			if saver.transaction {
				_, err = saver.tx.Exec("DELETE from ? WHERE ?=?", saver.table, saver.key, id)
			} else {
				_, err = saver.db.Exec("DELETE from ? WHERE ?=?", saver.table, saver.key, id)
			}

			if err != nil {
				log.Error("Deleting non synchronized rows failed")
				log.Error(err)

				return err
			}
		}
	}

	return nil
}

//Close closes the destination.
func (saver *DbSaver) Close(log *logrus.Entry) error {
	logDb := log.WithField("datasource", saver.ds.GetName())

	if saver.mode == exactCopy {
		err := saver.removeNonSynchronized(logDb)
		if err != nil {
			return err
		}
	}

	if saver.transaction && saver.tx != nil {
		logDb.Debug("Committing transaction")

		err := saver.tx.Commit()
		if err != nil {
			logDb.Error("Committing transaction failed")
			logDb.Error(err)
		}
	}

	err := saver.ds.CloseDatabase(logDb, false, false)
	if err != nil {
		logDb.Error("Close database failed")
		logDb.Error(err)
	}

	return err
}

//Reset reinitialize the destination (if possible).
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

//Name give the name of the destination.
func (saver *DbSaver) Name() string {
	return saver.database + "_" + saver.table
}
