package datasource

import (
	"os"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
)

// Datasources is a collection of Datasource.
type Datasources struct {
	//Dictionnary of datasource indexed by name
	datasources map[string]*Datasource
	// Datasource tag dictionary for lookup
	tagToDatasource map[string][]string
	// Timeout of each database ping try
	conTimeout time.Duration
	// Number of retry of database ping
	conRetry int
	// Environment variables for templating
	envVar map[string]string
}

// New returns a new Datasources object with elments initialized.
func New(connectionTimeout time.Duration, connectionRetry int) *Datasources {
	var dss Datasources
	dss.datasources = make(map[string]*Datasource)
	dss.tagToDatasource = make(map[string][]string)
	dss.conTimeout = connectionTimeout
	dss.conRetry = connectionRetry

	dss.envVar = make(map[string]string)

	for _, v := range os.Environ() {
		splitV := strings.Split(v, "=")
		dss.envVar[splitV[0]] = splitV[1]
	}

	return &dss
}

// Datasourcers interface to allow switching the way of storing the datasources.
type Datasourcers interface {
	LoadAll(string, *logrus.Entry) error
	CloseAll(*logrus.Entry)
	Lookup(*logrus.Entry, []string, []string, []Type, []Engine) ([]Datasourcer, []Datasourcer, error)
}

// CloseAll close all filehandle and database connection still opened.
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
