//Package mockdatasource provides mock objects for datasource package
package mockdatasource

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"io"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/datasource"
)

// MockDatasource is fake datasource object for test purpose.
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
	FileNotExists bool
	TableExists   bool
	TableEmpty    bool
	MockedDb      *sql.DB
	WriteBuf      bytes.Buffer
}

//GetEngine return the engine enum value.
func (ds *MockDatasource) GetEngine() datasource.Engine {
	return ds.Engine
}

//GetType return the type enum value.
func (ds *MockDatasource) GetType() datasource.Type {
	return ds.Type
}

//IsTransaction return true if the datasource has transaction.
func (ds *MockDatasource) IsTransaction() bool {
	return ds.Transaction
}

//GetNamedTag return the value of the tag with the provided name or "" if not exists.
func (ds *MockDatasource) GetNamedTag(name string) string {
	for _, tag := range ds.Tags {
		if strings.Contains(tag, name+":") {
			return tag[len(name)+1:]
		}
	}

	return ""
}

//GetName return the name of the datasource.
func (ds *MockDatasource) GetName() string {
	return ds.Name
}

//GetHash returns uniq hash for the datasource final destination (more than one datasource could have the same hash by example same database engine).
func (ds *MockDatasource) GetHash(log *logrus.Entry, admin bool, nodb bool) string {
	var toHash string

	if ds.Type == datasource.File {
		toHash = ds.FilePath
	} else {
		switch {
		case nodb:
			toHash = ds.URLNoDb
		case admin:
			toHash = ds.URLAdmin
		default:
			toHash = ds.URL
		}
	}

	hashed := sha256.Sum256([]byte(toHash))

	return fmt.Sprintf("%x", hashed)
}
