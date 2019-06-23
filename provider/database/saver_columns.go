package database

import (
	"fmt"
	"log"
	"strings"

	"github.com/marema31/kamino/kaminodb"
	"github.com/marema31/kamino/provider/common"
)

func (ds *DbSaver) questionMarkByEngine(qm *[]string) string {

	switch ds.kdb.Engine {
	case kaminodb.Mysql:
		return "?"
	case kaminodb.Postgres:
		return fmt.Sprintf("$%d", len(*qm)+1)
	}
	return ""
}

func (ds *DbSaver) getColNames(record common.Record) ([]string, []string, error) {

	var updateSet []string
	var questionmark []string
	rows, err := ds.db.QueryContext(ds.ctx, fmt.Sprintf("SELECT * from %s LIMIT 1", ds.table)) // We don't need data, we only needs the column names
	if err != nil {
		return nil, nil, err
	}

	columns, err := rows.ColumnTypes()
	if err != nil {
		return nil, nil, err
	}

	ds.wasEmpty = !rows.Next()
	if ds.mode == onlyIfEmpty && !ds.wasEmpty {
		log.Printf("Warning: the table %s of database %s is not empty, I will do nothing on it", ds.table, ds.database)
	}

	keyseen := false
	for _, col := range columns {
		_, ok := record[col.Name()]
		if !ok {
			log.Printf("Warning: the colum %s does not exist in source, for table %s of %s using table default value", col.Name(), ds.table, ds.database)
			continue
		}
		if strings.EqualFold(col.Name(), ds.key) {
			keyseen = true
			continue
		}
		ds.colNames = append(ds.colNames, col.Name())
		updateSet = append(updateSet, fmt.Sprintf("%s=%s", col.Name(), ds.questionMarkByEngine(&updateSet)))
		questionmark = append(questionmark, ds.questionMarkByEngine(&questionmark))
	}

	// By doing like this we ensure the primary key will be the last of column names and this array can be use for insert and update
	if ds.key != "" {
		ds.colNames = append(ds.colNames, ds.key)
		questionmark = append(questionmark, ds.questionMarkByEngine(&questionmark))

		if !keyseen {
			return nil, nil, fmt.Errorf("provided key %s is not a column of %s.%s ", ds.key, ds.database, ds.table)
		}
	}

	keyseen = false
	for colr := range record {
		if ds.key != "" && ds.key == colr {
			keyseen = true
		}
		seen := false
		for _, col := range columns {
			if col.Name() == colr {
				seen = true
			}
		}
		if !seen {
			log.Printf("Warning: the colum %s does not exist in destination table %s of %s", colr, ds.table, ds.database)
		}
	}

	if ds.key != "" && !keyseen {
		return nil, nil, fmt.Errorf("provided key %s is not available from filtered source for %s.%s ", ds.key, ds.database, ds.table)
	}

	return questionmark, updateSet, nil
}

// createStatement Query the destination table to determine the available colums, create the corresponding insert/update statement and save them in the dbSaver instance
func (ds *DbSaver) createStatement(record common.Record) error {
	questionmark, updateSet, err := ds.getColNames(record)
	if err != nil {
		return err
	}

	ds.insertString = fmt.Sprintf("INSERT INTO %s ( %s) VALUES ( %s )", ds.table, strings.Join(ds.colNames[:], ","), strings.Join(questionmark[:], ","))
	ds.updateString = fmt.Sprintf("UPDATE %s SET  %s WHERE %s = %s", ds.table, strings.Join(updateSet[:], ","), ds.key, ds.questionMarkByEngine(&updateSet))

	if ds.kdb.Transaction {
		ds.insertStmt, err = ds.tx.Prepare(ds.insertString)
	} else {
		ds.insertStmt, err = ds.db.Prepare(ds.insertString)
	}
	if err != nil {
		return err
	}
	if ds.mode == replace || ds.mode == update || ds.mode == exactCopy {
		if ds.kdb.Transaction {
			ds.updateStmt, err = ds.tx.Prepare(ds.updateString)
		} else {
			ds.updateStmt, err = ds.db.Prepare(ds.updateString)
		}
		if err != nil {
			return err
		}
	}

	return nil
}
