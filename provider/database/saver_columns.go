package database

import (
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/provider/types"
)

func (saver *DbSaver) questionMarkByEngine(qm *[]string) string {

	switch saver.engine {
	case datasource.Mysql:
		return "?"
	case datasource.Postgres:
		return fmt.Sprintf("$%d", len(*qm)+1)
	}
	return ""
}

func (saver *DbSaver) getColNames(log *logrus.Entry, record types.Record) ([]string, []string, error) {

	var updateSet []string
	var questionmark []string
	log.Debug("Retrieving the column names")
	log.Debugf("SELECT * from %s LIMIT 1", saver.table)
	rows, err := saver.db.QueryContext(saver.ctx, fmt.Sprintf("SELECT * from %s LIMIT 1", saver.table)) // We don't need data, we only neesaver the column names
	if err != nil {
		log.Error("Querying for retrieving column names failed")
		log.Error(err)
		return nil, nil, err
	}

	columns, err := rows.ColumnTypes()
	if err != nil {
		log.Error("Retrieving column names failed")
		log.Error(err)
		return nil, nil, err
	}

	saver.wasEmpty = !rows.Next()
	if saver.mode == onlyIfEmpty && !saver.wasEmpty {
		log.Warnf("Table %s not empty, I will do nothing on it", saver.table)
	}

	keyseen := false
	for _, col := range columns {
		_, ok := record[col.Name()]
		if !ok {
			log.Warnf("Colum %s does not exist in source, using table default value", col.Name())
			continue
		}
		if strings.EqualFold(col.Name(), saver.key) {
			keyseen = true
			continue
		}
		saver.colNames = append(saver.colNames, col.Name())
		updateSet = append(updateSet, fmt.Sprintf("%s=%s", col.Name(), saver.questionMarkByEngine(&updateSet)))
		questionmark = append(questionmark, saver.questionMarkByEngine(&questionmark))
	}

	// By doing like this we ensure the primary key will be the last of column names and this array can be use for insert and update
	if saver.key != "" {
		saver.colNames = append(saver.colNames, saver.key)
		questionmark = append(questionmark, saver.questionMarkByEngine(&questionmark))

		if !keyseen {
			log.Errorf("Provided key %s is not a column of %s", saver.key, saver.table)
			log.Error(err)
			return nil, nil, fmt.Errorf("provided key %s is not a column of %s.%s ", saver.key, saver.database, saver.table)
		}
	}

	keyseen = false
	for colr := range record {
		if saver.key != "" && saver.key == colr {
			keyseen = true
		}
		seen := false
		for _, col := range columns {
			if col.Name() == colr {
				seen = true
			}
		}
		if !seen {
			log.Warnf("Colum %s does not exist in %s", colr, saver.table)
		}
	}

	if saver.key != "" && !keyseen {
		log.Errorf("Provided key %s is not a column of %s", saver.key, saver.table)
		log.Error(err)
		return nil, nil, fmt.Errorf("provided key %s is not available in destination table for %s.%s ", saver.key, saver.database, saver.table)
	}

	return questionmark, updateSet, nil
}

// createStatement Query the destination table to determine the available colums, create the corresponding insert/update statement and save them in the dbSaver instance
func (saver *DbSaver) createStatement(log *logrus.Entry, record types.Record) error {
	questionmark, updateSet, err := saver.getColNames(log, record)
	if err != nil {
		return err
	}

	saver.insertString = fmt.Sprintf("INSERT INTO %s ( %s) VALUES ( %s )", saver.table, strings.Join(saver.colNames[:], ","), strings.Join(questionmark[:], ","))
	saver.updateString = fmt.Sprintf("UPDATE %s SET  %s WHERE %s = %s", saver.table, strings.Join(updateSet[:], ","), saver.key, saver.questionMarkByEngine(&updateSet))

	log.Debug("Preparing Insert statement")
	log.Debug(saver.insertString)
	if saver.transaction {
		saver.insertStmt, err = saver.tx.Prepare(saver.insertString)
	} else {
		saver.insertStmt, err = saver.db.Prepare(saver.insertString)
	}
	if err != nil {
		log.Error("Preparing Insert statement failed")
		log.Error(saver.insertString)
		log.Error(err)
		return err
	}
	if saver.mode == replace || saver.mode == update || saver.mode == exactCopy {
		log.Debug("Preparing Update statement")
		log.Debug(saver.updateString)
		if saver.transaction {
			saver.updateStmt, err = saver.tx.Prepare(saver.updateString)
		} else {
			saver.updateStmt, err = saver.db.Prepare(saver.updateString)
		}
		if err != nil {
			log.Error("Preparing Update statement failed")
			log.Error(saver.updateString)
			log.Error(err)
			return err
		}
	}

	return nil
}
