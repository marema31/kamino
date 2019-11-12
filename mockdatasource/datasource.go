//Package mockdatasource provides mock objects for datasource package
package mockdatasource

import (
	"bytes"
	"database/sql"
	"io"
	"strings"

	"github.com/marema31/kamino/datasource"
)

// MockDatasource is fake datasource object for test purpose
type MockDatasource struct {
	Name          string
	Type          datasource.Type
	Database      string
	Engine        datasource.Engine
	Inline        string
	Host          string
	Port          string
	User          string
	UserPw        string
	Admin         string
	AdminPw       string
	URL           string
	URLAdmin      string
	URLNoDb       string
	Transaction   bool
	Schema        string
	FilePath      string
	TmpFilePath   string
	Gzip          bool
	Zip           bool
	FileHandle    io.Closer
	Filewriter    bool
	Tags          []string
	ErrorOpenDb   error
	ErrorOpenFile error
	ErrorReset    error
	ErrorClose    error
	MockedDb      *sql.DB
	WriteBuf      bytes.Buffer
}

//GetEngine return the engine enum value
func (ds *MockDatasource) GetEngine() datasource.Engine {
	return ds.Engine
}

//IsTransaction return true if the datasource has transaction
func (ds *MockDatasource) IsTransaction() bool {
	return ds.Transaction
}

//GetNamedTag return the value of the tag with the provided name or "" if not exists
func (ds *MockDatasource) GetNamedTag(name string) string {
	for _, tag := range ds.Tags {
		if strings.Contains(tag, name+":") {
			return tag[len(name)+1:]
		}
	}
	return ""
}

//GetName return the name of the datasource
func (ds *MockDatasource) GetName() string {
	return ds.Name
}
