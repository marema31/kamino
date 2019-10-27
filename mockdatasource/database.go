package mockdatasource

import (
	"database/sql"
)

//OpenDatabase open connection to the corresponding database
func (ds *MockDatasource) OpenDatabase(admin bool, nodb bool) (*sql.DB, error) {
	if ds.ErrorOpenDb != nil {
		return nil, ds.ErrorOpenDb
	}
	return ds.MockedDb, nil
}
