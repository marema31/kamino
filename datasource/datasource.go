//Package datasource manage the list of datasources and their
// selction from tags
package datasource

import (
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

//Dictionnary of datasource indexed by name
var datasources = make(map[string]*Datasource)

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
		if strings.Index(tag, name+":") != -1 {
			return tag[len(name)+1 : len(tag)]
		}
	}
	return ""
}
