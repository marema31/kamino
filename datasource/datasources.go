package datasource

import "github.com/Sirupsen/logrus"

// Datasources is a collection of Datasource
type Datasources struct {
	//Dictionnary of datasource indexed by name
	datasources map[string]*Datasource
	// Datasource tag dictionnary for lookup
	tagToDatasource map[string][]string
}

// New returns a new Datasources object with elments initialized
func New() *Datasources {
	var dss Datasources
	dss.datasources = make(map[string]*Datasource)
	dss.tagToDatasource = make(map[string][]string)
	return &dss
}

// Datasourcers interface to allow switching the way of storing the datasources
type Datasourcers interface {
	LoadAll(string, *logrus.Entry) error
	Lookup(*logrus.Entry, []string, []Type, []Engine) []Datasourcer
}
