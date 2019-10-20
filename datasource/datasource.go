//Package datasource manage the list of datasources and their
// selection from tags
package datasource

import (
	"fmt"
	"io"
	"strings"
)

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

//TODO: comment the fields

// Datasource is handle to the corresponding datasource (either file/database)
type Datasource struct {
	Name        string
	Type        Type
	Database    string
	Engine      Engine
	Inline      string
	Host        string
	Port        string
	User        string
	UserPw      string
	Admin       string
	AdminPw     string
	URL         string //TODO: not exported ?
	URLAdmin    string //TODO: not exported ?
	URLNoDb     string //TODO: not exported ?
	Transaction bool
	Schema      string
	FilePath    string
	tmpFilePath string
	Gzip        bool // TODO: not exported ?
	Zip         bool // TODO: not exported ?
	fileHandle  io.Closer
	filewriter  bool
	Tags        []string
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

//GetEngine return a string containing the engine name
func (ds *Datasource) GetEngine() string {
	switch ds.Engine {
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

//GetNamedTag return the value of the tag with the provided name or "" if not exists
func (ds *Datasource) GetNamedTag(name string) string {
	for _, tag := range ds.Tags {
		if strings.Contains(tag, name+":") {
			return tag[len(name)+1:]
		}
	}
	return ""
}
