package mockdatasource

import (
	"github.com/marema31/kamino/datasource"
)

// FillTmplValues return a struct for template operation with value corresponding to the provided datasource.
func (ds *MockDatasource) FillTmplValues() datasource.TmplValues {
	var tv datasource.TmplValues
	tv.Name = ds.Name
	tv.Transaction = ds.Transaction
	tv.Database = ds.Database
	tv.User = ds.User
	tv.Password = ds.UserPw
	tv.Schema = ds.Schema
	tv.Host = ds.Host
	tv.Port = ds.Port
	tv.Tags = ds.Tags
	tv.Type = datasource.TypeToString(ds.Type)
	tv.Engine = datasource.EngineToString(ds.Engine)
	tv.FilePath = ds.FilePath

	return tv
}
