package database

import (
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/provider/common"
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
	updateSet := make([]string, 0)
	questionmark := make([]string, 0)

	log.Debug("Retrieving the column names")

	query := fmt.Sprintf("SELECT column_name AS name FROM information_schema.columns WHERE table_schema = '%s' AND table_name ='%s';", saver.database, saver.table) //nolint: gosec
	log.Debug(query)

	rows, err := saver.db.QueryContext(saver.ctx, query)
	if err != nil {
		log.Error("Querying for retrieving column names failed")
		log.Error(err)

		return nil, nil, err
	}
	defer rows.Close()

	columns := make([]string, 0)

	for rows.Next() {
		var col string
		if err := rows.Scan(&col); err != nil {
			// Check for a scan error.
			// Query rows will be closed with defer.
			log.Fatal(err)
		}

		columns = append(columns, col)
	}

	keyseen := false

	for _, col := range columns {
		_, ok := record[col]
		if !ok {
			log.Warnf("Column %s does not exist in source, using table default value", col)
			continue
		}

		if strings.EqualFold(col, saver.key) {
			keyseen = true
			continue
		}

		saver.colNames = append(saver.colNames, col)
		updateSet = append(updateSet, fmt.Sprintf("%s=%s", col, saver.questionMarkByEngine(&updateSet)))
		questionmark = append(questionmark, saver.questionMarkByEngine(&questionmark))
	}

	// By doing like this we ensure the primary key will be the last of column names and this array can be use for insert and update
	if saver.key != "" {
		saver.colNames = append(saver.colNames, saver.key)
		questionmark = append(questionmark, saver.questionMarkByEngine(&questionmark))

		if !keyseen {
			log.Errorf("Provided key %s is not a column of %s", saver.key, saver.table)
			log.Error(err)

			return nil, nil, fmt.Errorf("provided key %s is not a column of %s.%s : %w", saver.key, saver.database, saver.table, common.ErrMissingParameter)
		}
	}

	for colr := range record {
		seen := false

		for _, col := range columns {
			if col == colr {
				seen = true
			}
		}

		if !seen {
			log.Warnf("Colum %s does not exist in %s", colr, saver.table)
		}
	}

	var count int

	query = fmt.Sprintf("SELECT COUNT(*) FROM %s", saver.table) //nolint: gosec
	log.Debug(query)

	row := saver.db.QueryRowContext(saver.ctx, query)

	err = row.Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	saver.wasEmpty = (count == 0)
	if saver.mode == onlyIfEmpty && !saver.wasEmpty {
		log.Warnf("Table %s not empty, I will do nothing on it", saver.table)
	}

	return questionmark, updateSet, nil
}

// createStatement Query the destination table to determine the available colums, create the corresponding insert/update statement and save them in the dbSaver instance.
func (saver *DbSaver) createStatement(log *logrus.Entry, record types.Record) error {
	questionmark, updateSet, err := saver.getColNames(log, record)
	if err != nil {
		return err
	}

	saver.insertString = fmt.Sprintf("INSERT INTO %s ( %s) VALUES ( %s )", saver.table, strings.Join(saver.colNames, ","), strings.Join(questionmark, ","))           //nolint:gosec
	saver.updateString = fmt.Sprintf("UPDATE %s SET  %s WHERE %s = %s", saver.table, strings.Join(updateSet, ","), saver.key, saver.questionMarkByEngine(&updateSet)) //nolint:gosec

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
