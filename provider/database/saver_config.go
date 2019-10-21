package database

import (
	"fmt"
	"strings"
)

func stringToMode(modestr string) dbSaverMode {

	switch strings.ToLower(modestr) {
	case "onlyifempty":
		return onlyIfEmpty
	case "insert":
		return insert
	case "update":
		return update
	case "replace":
		return replace
	case "copy":
		return exactCopy
	case "truncate":
		return truncate
	default:
		return exactCopy
	}
}

//createIdsList store in the instance the list of all values of column described in 'key' configuration entry
func (saver *DbSaver) createIdsList() error {
	rows, err := saver.db.QueryContext(saver.ctx, fmt.Sprintf("SELECT %s from %s", saver.key, saver.table))
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return err
		}
		saver.ids[id] = false
	}
	return nil
}
