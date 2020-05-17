package common

import (
	"context"

	"github.com/marema31/kamino/datasource"

	"github.com/Sirupsen/logrus"
)

// ToSkipDatabase run the query (likely a SELECT COUNT) on the datasource
// return true if the query return a non-zero value in the only column of the only row.
func ToSkipDatabase(ctx context.Context, log *logrus.Entry, ds datasource.Datasourcer, admin bool, nodb bool, queries []string) (bool, error) {
	var needskip int

	db, err := ds.OpenDatabase(log, admin, nodb)
	if err != nil {
		return false, err
	}

	for _, query := range queries {
		err = db.QueryRowContext(ctx, query).Scan(&needskip)
		if err != nil {
			log.Errorf("Query of skip phase failed : %v", err)
			return false, err
		}

		// If the query returned a 0 result we shoud do it, do not test other queries (AND) to allow test table exists and table contains something
		if needskip == 0 {
			return false, nil
		}
	}

	return true, nil
}
