package common

import (
	"context"

	"github.com/marema31/kamino/datasource"

	"github.com/Sirupsen/logrus"
)

// ToSkipDatabase run the query (likely a SELECT COUNT) on the datasource
// return true if the query return a non-zero value in the only column of the only row.
func ToSkipDatabase(ctx context.Context, log *logrus.Entry, ds datasource.Datasourcer, admin bool, nodb bool, query string) (bool, error) {
	var needskip int

	db, err := ds.OpenDatabase(log, admin, nodb)
	if err != nil {
		return false, err
	}

	err = db.QueryRowContext(ctx, query).Scan(&needskip)
	if err != nil {
		log.Error("Query of skip phase failed")
		log.Error(err)

		return false, err
	}

	return needskip != 0, nil
}
