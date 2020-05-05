package mockdatasource

import (
	"database/sql"

	"github.com/Sirupsen/logrus"
)

//OpenDatabase open connection to the corresponding database.
func (ds *MockDatasource) OpenDatabase(log *logrus.Entry, admin bool, nodb bool) (*sql.DB, error) {
	if ds.ErrorOpenDb != nil {
		return nil, ds.ErrorOpenDb
	}

	return ds.MockedDb, nil
}
