package mockdatasource

import (
	"database/sql"
)

//OpenDatabase open connection to the corresponding database
func (ds *MockDatasource) OpenDatabase(generateError bool, nodb bool) (*sql.DB, error) {
	if generateError {
		return nil, ds.ErrorOpenDb
	}
	return ds.MockedDb, nil
}
