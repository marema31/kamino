package mockdatasource

import (
	"context"
	"database/sql"

	"github.com/Sirupsen/logrus"
)

//IsTableExists return true if the table exists.
func (ds *MockDatasource) IsTableExists(ctx context.Context, log *logrus.Entry, table string) (bool, error) {
	if ds.ErrorOpenDb != nil {
		return false, ds.ErrorOpenDb
	}

	return ds.TableExists, nil
}

//IsTableEmpty return true if the table empty.
func (ds *MockDatasource) IsTableEmpty(ctx context.Context, log *logrus.Entry, table string) (bool, error) {
	if ds.ErrorOpenDb != nil {
		return false, ds.ErrorOpenDb
	}

	return ds.TableEmpty, nil
}

//OpenDatabase open connection to the corresponding database.
func (ds *MockDatasource) OpenDatabase(log *logrus.Entry, admin bool, nodb bool) (*sql.DB, error) {
	if ds.ErrorOpenDb != nil {
		return nil, ds.ErrorOpenDb
	}

	return ds.MockedDb, nil
}

//CloseDatabase close connection to the corresponding database only if no more used.
func (ds *MockDatasource) CloseDatabase(log *logrus.Entry, admin bool, nodb bool) error {
	return ds.ErrorClose
}
