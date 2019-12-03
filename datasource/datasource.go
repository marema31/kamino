//Package datasource manage the list of datasources and their
// selection from tags
package datasource

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/marema31/kamino/file"
)

//Datasourcer interface for allowing mocking of Datasource object
type Datasourcer interface {
	FillTmplValues() TmplValues
	OpenDatabase(*logrus.Entry, bool, bool) (*sql.DB, error)
	OpenReadFile(*logrus.Entry) (io.ReadCloser, error)
	OpenWriteFile(*logrus.Entry) (io.WriteCloser, error)
	ResetFile(*logrus.Entry) error
	CloseFile(*logrus.Entry) error
	GetName() string
	GetEngine() Engine
	GetType() Type
	IsTransaction() bool
	Stat() (os.FileInfo, error)
}

//Engine constants for file/database engine
type Engine int

const (
	// Mysql / MariaDB database engine
	Mysql Engine = iota
	// Postgres database engine
	Postgres Engine = iota
	// JSON file engine
	JSON Engine = iota
	// YAML file engine
	YAML Engine = iota
	// CSV file engine
	CSV Engine = iota
)

//Type discriminate the type of datasource
type Type int

const (
	// Database (mariadb, mysql, postgres, ...)
	Database Type = iota
	// File (JSON,YAML,CSV, ...)
	File Type = iota
)

// Datasource is handle to the corresponding datasource (either file/database)
type Datasource struct {
	name        string
	dstype      Type
	database    string
	engine      Engine
	host        string
	port        string
	user        string
	userPw      string
	admin       string
	adminPw     string
	url         string
	urlAdmin    string
	urlNoDb     string
	conTimeout  time.Duration
	conRetry    int
	db          *sql.DB
	dbAdmin     *sql.DB
	dbNoDb      *sql.DB
	transaction bool
	schema      string
	file        file.File
	tags        []string
}

//StringToType convert string to corresponding Type typed value
func StringToType(dsType string) (Type, error) {
	switch strings.ToLower(dsType) {
	case "db", "database", "databases":
		return Database, nil
	case "file", "files":
		return File, nil
	}
	return File, fmt.Errorf("does not how to manage %s datasource type", dsType)
}

//StringsToTypes convert string slice to an slice of corresponding Type typed values
func StringsToTypes(dsTypes []string) ([]Type, error) {
	typeSlice := make([]Type, 0)
	for _, dsType := range dsTypes {
		t, err := StringToType(dsType)
		if err != nil {
			return nil, err
		}
		typeSlice = append(typeSlice, t)
	}
	return typeSlice, nil
}

//StringToEngine convert string to corresponding Engine typed value
func StringToEngine(engine string) (Engine, error) {
	switch strings.ToLower(engine) {
	case "mysql", "maria", "mariadb":
		return Mysql, nil
	case "pgsql", "postgres", "postgresql":
		return Postgres, nil
	case "json":
		return JSON, nil
	case "yaml":
		return YAML, nil
	case "csv":
		return CSV, nil
	}
	return CSV, fmt.Errorf("does not how to manage %s datasource engine", engine)
}

//StringsToEngines convert string slice to an slice of corresponding Engine typed values
func StringsToEngines(engines []string) ([]Engine, error) {
	engineSlice := make([]Engine, 0)
	for _, engine := range engines {
		e, err := StringToEngine(engine)
		if err != nil {
			return nil, err
		}
		engineSlice = append(engineSlice, e)
	}
	return engineSlice, nil
}

//EngineToString return do the conversion
func EngineToString(engine Engine) string {
	switch engine {
	case Mysql:
		return "mysql"
	case Postgres:
		return "postgresql"
	case JSON:
		return "json"
	case YAML:
		return "yaml"
	case CSV:
		return "csv"
	}
	return "Unknown" // We will never arrive here
}

//TypeToString return do the conversion
func TypeToString(dsType Type) string {
	switch dsType {
	case Database:
		return "database"
	case File:
		return "file"
	}
	return "Unknown" // We will never arrive here
}

//GetEngine return the engine enum value
func (ds *Datasource) GetEngine() Engine {
	return ds.engine
}

//GetType return the type enum value
func (ds *Datasource) GetType() Type {
	return ds.dstype
}

//IsTransaction return true if the datasource has transaction
func (ds *Datasource) IsTransaction() bool {
	return ds.transaction
}

//GetNamedTag return the value of the tag with the provided name or "" if not exists
func (ds *Datasource) GetNamedTag(name string) string {
	for _, tag := range ds.tags {
		if strings.Contains(tag, name+":") {
			return tag[len(name)+1:]
		}
	}
	return ""
}

//GetName return the name of the datasource
func (ds *Datasource) GetName() string {
	return ds.name
}
