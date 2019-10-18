package database

/*TODO: Uncomment
//parseConfig parse the config to extract the mode and the primary key and save them in the dbSaver instance
func (ds *DbSaver) parseConfig(saverConfig config.DestinationConfig) error {
	ds.key = saverConfig.Key

	ds.mode = exactCopy
	modestr := saverConfig.Mode
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

	if ds.key == "" && (ds.mode == update || ds.mode == replace) {
		return fmt.Errorf("mode for %s.%s is %s and no key is provided", ds.database, ds.table, modestr)
	}
	return nil
}

//createIdsList store in the instance the list of all values of column described in 'key' configuration entry
func (ds *DbSaver) createIdsList() error {
	rows, err := ds.db.QueryContext(ds.ctx, fmt.Sprintf("SELECT %s from %s", ds.key, ds.table)) // We don't need data, we only needs the column names
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
*/
