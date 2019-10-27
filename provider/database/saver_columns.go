package database

import (
	"fmt"
	"log"
	"strings"

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

func (saver *DbSaver) getColNames(record types.Record) ([]string, []string, error) {

	var updateSet []string
	var questionmark []string
	rows, err := saver.db.QueryContext(saver.ctx, fmt.Sprintf("SELECT * from %s LIMIT 1", saver.table)) // We don't need data, we only neesaver the column names
	if err != nil {
		return nil, nil, err
	}

	columns, err := rows.ColumnTypes()
	if err != nil {
		return nil, nil, err
	}

	saver.wasEmpty = !rows.Next()
	if saver.mode == onlyIfEmpty && !saver.wasEmpty {
		log.Printf("Warning: the table %s of database %s is not empty, I will do nothing on it", saver.table, saver.database)
	}

	keyseen := false
	for _, col := range columns {
		_, ok := record[col.Name()]
		if !ok {
			log.Printf("Warning: the colum %s does not exist in source, for table %s of %s using table default value", col.Name(), saver.table, saver.database)
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
			log.Printf("Warning: the colum %s does not exist in destination table %s of %s", colr, saver.table, saver.database)
		}
	}

	if saver.key != "" && !keyseen {
		return nil, nil, fmt.Errorf("provided key %s is not available in destination table for %s.%s ", saver.key, saver.database, saver.table)
	}

	return questionmark, updateSet, nil
}

// createStatement Query the destination table to determine the available colums, create the corresponding insert/update statement and save them in the dbSaver instance
func (saver *DbSaver) createStatement(record types.Record) error {
	questionmark, updateSet, err := saver.getColNames(record)
	if err != nil {
		return err
	}

	saver.insertString = fmt.Sprintf("INSERT INTO %s ( %s) VALUES ( %s )", saver.table, strings.Join(saver.colNames[:], ","), strings.Join(questionmark[:], ","))
	saver.updateString = fmt.Sprintf("UPDATE %s SET  %s WHERE %s = %s", saver.table, strings.Join(updateSet[:], ","), saver.key, saver.questionMarkByEngine(&updateSet))

	if saver.transaction {
		saver.insertStmt, err = saver.tx.Prepare(saver.insertString)
	} else {
		saver.insertStmt, err = saver.db.Prepare(saver.insertString)
	}
	if err != nil {
		return err
	}
	if saver.mode == replace || saver.mode == update || saver.mode == exactCopy {
		if saver.transaction {
			saver.updateStmt, err = saver.tx.Prepare(saver.updateString)
		} else {
			saver.updateStmt, err = saver.db.Prepare(saver.updateString)
		}
		if err != nil {
			return err
		}
	}

	return nil
}
