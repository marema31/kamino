package datasource

import (
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"strings"
	"sync"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Masterminds/sprig/v3"
	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql" // Mysql library dynamically called by database/sql
	_ "github.com/lib/pq"              //Postgres library dynamically called by database/sql
	"github.com/spf13/viper"
)

var openedDatabase = map[string]*sql.DB{}
var openedReadMutex = &sync.Mutex{}
var openMutex = &sync.Mutex{}

type tmplEnv struct {
	Environments map[string]string
}

func parseField(v *viper.Viper, data tmplEnv, field string, fieldDetailedName string) (string, error) {
	var buf bytes.Buffer

	fieldValue := v.GetString(field)

	tmpl, err := template.New(field).Funcs(sprig.FuncMap()).Parse(fieldValue)
	if err != nil {
		return "", fmt.Errorf("parsing %s provided: %w", fieldDetailedName, err)
	}

	if err = tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("expanding %s provided: %w", fieldDetailedName, err)
	}

	parsed := buf.String()

	return parsed, nil
}

// load a database type datasource from the viper configuration
//nolint: funlen
func loadDatabaseDatasource(log *logrus.Entry, filename string, v *viper.Viper, engine Engine, envVar map[string]string, connectionTimeout time.Duration, connectionRetry int) (Datasource, error) {
	log.Debugf("Loading %s file datasource", filename)

	var err error

	var ds Datasource
	ds.dstype = Database
	ds.engine = engine
	ds.name = filename

	data := tmplEnv{Environments: envVar}

	ds.database, err = parseField(v, data, "database", "database name")
	if err != nil {
		return Datasource{}, err
	}

	if ds.database == "" {
		return Datasource{}, fmt.Errorf("no database name provided: %w", errMissingParameter)
	}

	ds.tags = v.GetStringSlice("tags")
	if len(ds.tags) == 0 {
		ds.tags = []string{""}
	}

	ds.conTimeout = connectionTimeout
	ds.conRetry = connectionRetry

	ds.schema = v.GetString("schema")

	ds.transaction = v.GetBool("transaction")

	ds.host, err = parseField(v, data, "host", "host name")
	if err != nil {
		return Datasource{}, err
	}

	if ds.host == "" {
		ds.host = "127.0.0.1"
	}

	ds.port, err = parseField(v, data, "port", "port")
	if err != nil {
		return Datasource{}, err
	}

	ds.user, err = parseField(v, data, "user", "user name")
	if err != nil {
		return Datasource{}, err
	}

	ds.admin, err = parseField(v, data, "admin", "admin name")
	if err != nil {
		return Datasource{}, err
	}

	ds.userPw, err = parseField(v, data, "password", "user password")
	if err != nil {
		return Datasource{}, err
	}

	ds.adminPw, err = parseField(v, data, "adminpassword", "admin password")
	if err != nil {
		return Datasource{}, err
	}

	if ds.adminPw == "" {
		ds.adminPw = ds.userPw
	}

	if ds.userPw == "" {
		ds.userPw = ds.adminPw
	}

	dbOptions := v.GetStringSlice("options")

	switch ds.engine {
	case Mysql:
		if ds.user == "" {
			ds.user = "root"
		}

		if ds.admin == "" {
			ds.admin = "root"
		}

		if ds.port == "" {
			ds.port = "3306"
		}

		urlOptions := ""
		if len(dbOptions) > 0 {
			urlOptions = fmt.Sprintf("&%s", strings.Join(dbOptions, "&"))
		}
		//use parseTime=true to force date and time conversion
		ds.url = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true%s", ds.user, ds.userPw, ds.host, ds.port, ds.database, urlOptions)
		ds.urlAdmin = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true%s", ds.admin, ds.adminPw, ds.host, ds.port, ds.database, urlOptions)
		ds.urlNoDb = fmt.Sprintf("%s:%s@tcp(%s:%s)/mysql?parseTime=true%s", ds.admin, ds.adminPw, ds.host, ds.port, urlOptions)

	case Postgres:
		if ds.user == "" {
			ds.user = "postgres"
		}

		if ds.admin == "" {
			ds.admin = "postgres"
		}

		if ds.port == "" {
			ds.port = "5432"
		}

		urlOptions := ""
		if len(dbOptions) > 0 {
			urlOptions = strings.Join(dbOptions, " ")
		}

		ds.url = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s %s", ds.host, ds.port, ds.user, ds.userPw, ds.database, urlOptions)
		ds.urlAdmin = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s %s", ds.host, ds.port, ds.admin, ds.adminPw, ds.database, urlOptions)
		ds.urlNoDb = fmt.Sprintf("host=%s port=%s user=%s password=%s %s", ds.host, ds.port, ds.admin, ds.adminPw, urlOptions)
	}

	return ds, nil
}

//OpenDatabase open connection to the corresponding database.
func (ds *Datasource) OpenDatabase(log *logrus.Entry, admin bool, nodb bool) (*sql.DB, error) {
	logDb := log.WithField("engine", EngineToString(ds.engine))

	if ds.dstype != Database {
		logDb.Error("Can not open as a database")
		return nil, fmt.Errorf("can not open %s as a database: %w", ds.name, errWrongType)
	}

	var (
		URL string
		db  *sql.DB
	)

	switch {
	case nodb:
		URL = ds.urlNoDb
		db = ds.dbNoDb
	case admin:
		URL = ds.urlAdmin
		db = ds.dbAdmin
	default:
		URL = ds.url
		db = ds.db
	}

	if db != nil {
		//		logDb.Debug("The database is already opened, returning the current handler")
		return db, nil
	}

	var driver string

	switch ds.engine {
	case Mysql:
		log.Debugf("Opening Mysql database: %s", URL)

		driver = "mysql"
	case Postgres:
		log.Debugf("Opening Postgresql database: %v", URL)

		driver = "postgres"
	}

	db, err := sqlOpen(driver, URL)
	if err != nil {
		log.Error("Opening database failed")
		log.Error(err)

		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 5) //nolint:gomnd //for the moment no reason to make it parametrized
	db.SetMaxIdleConns(10)                 //nolint:gomnd //for the moment no reason to make it parametrized

	for databaseConnectionAttemptLoop := 0; databaseConnectionAttemptLoop < ds.conRetry; databaseConnectionAttemptLoop++ {
		// Open does not really do a connection and therefore does not test for url is correct, ping will do
		err = db.Ping()

		if err == nil {
			break // Here, if there is no error, it simply breaks out and does not retry again.
		}

		time.Sleep(ds.conTimeout)
	}

	err = db.Ping()
	if err != nil {
		log.Error("Ping database failed")
		log.Error(err)

		return nil, err
	}

	switch {
	case nodb:
		ds.dbNoDb = db
	case admin:
		ds.dbAdmin = db
	default:
		ds.db = db
	}

	return db, nil
}

//Only for unit testing of OpenDatabase function.
var mockingSQL = false

func sqlOpen(driver string, url string) (*sql.DB, error) {
	if mockingSQL {
		db, _, err := sqlmock.New()

		return db, err
	}

	var err error = nil
	// Most of the time the open will occurs on a already opened database
	openedReadMutex.Lock()
	db, ok := openedDatabase[url]
	openedReadMutex.Unlock()

	if ok {
		return db, nil
	}

	// Out of lock, we may have to open the database connection, avoid two concurrent opening
	openMutex.Lock()
	//Verify that in between the database as not been already opened
	openedReadMutex.Lock()
	db, ok = openedDatabase[url]
	openedReadMutex.Unlock()

	if !ok {
		//No we have to open it
		db, err = sql.Open(driver, url)
		if err == nil {
			openedReadMutex.Lock()
			openedDatabase[url] = db
			openedReadMutex.Unlock()
		}
	}
	openMutex.Unlock()

	return db, err
}
