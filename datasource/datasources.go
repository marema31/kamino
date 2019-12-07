package datasource

import (
	"time"

	"github.com/Sirupsen/logrus"
)

// Datasources is a collection of Datasource
type Datasources struct {
	//Dictionnary of datasource indexed by name
	datasources map[string]*Datasource
	// Datasource tag dictionnary for lookup
	tagToDatasource map[string][]string
	// Timeout of each database ping try
	conTimeout time.Duration
	// Number of retry of database ping
	conRetry int
}

// New returns a new Datasources object with elments initialized
func New(connectionTimeout time.Duration, connectionRetry int) *Datasources {
	var dss Datasources
	dss.datasources = make(map[string]*Datasource)
	dss.tagToDatasource = make(map[string][]string)
	dss.conTimeout = connectionTimeout
	dss.conRetry = connectionRetry
	return &dss
}

// Datasourcers interface to allow switching the way of storing the datasources
type Datasourcers interface {
	LoadAll(string, *logrus.Entry) error
	CloseAll(*logrus.Entry)
	Lookup(*logrus.Entry, []string, []Type, []Engine) []Datasourcer
}

// CloseAll close all filehandle and database connection still opened
func (dss *Datasources) CloseAll(log *logrus.Entry) {
	for _, ds := range dss.datasources {
		log.Debugf("Closing %s", ds.name)
		ds.file.Close()
		if ds.db != nil {
			ds.db.Close()
		}
		if ds.dbAdmin != nil {
			ds.dbAdmin.Close()
		}
		if ds.dbNoDb != nil {
			ds.dbNoDb.Close()
		}

	}
}
