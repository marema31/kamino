package datasource

import (
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Masterminds/sprig/v3"
	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql" // Mysql library dynamically called by database/sql
	_ "github.com/lib/pq"              //Postgres library dynamically called by database/sql
	"github.com/spf13/viper"
)

// load a database type datasource from the viper configuration
//nolint: funlen
func loadDatabaseDatasource(filename string, v *viper.Viper, engine Engine, envVar map[string]string, connectionTimeout time.Duration, connectionRetry int) (Datasource, error) {
	var ds Datasource
	ds.dstype = Database
	ds.engine = engine
	ds.name = filename

	type tmplEnv struct {
		Environments map[string]string
	}

	data := tmplEnv{Environments: envVar}

	databaseTmpl, err := template.New("database").Funcs(sprig.FuncMap()).Parse(v.GetString("database"))
	if err != nil {
		return Datasource{}, fmt.Errorf("parsing database name provided")
	}

	var database bytes.Buffer
	if err = databaseTmpl.Execute(&database, data); err != nil {
		return Datasource{}, fmt.Errorf("expanding database name provided")
	}

	ds.database = database.String()
	if ds.database == "" {
		return Datasource{}, fmt.Errorf("no database name provided")
	}

	ds.tags = v.GetStringSlice("tags")
	if len(ds.tags) == 0 {
		ds.tags = []string{""}
	}

	ds.conTimeout = connectionTimeout
	ds.conRetry = connectionRetry

	ds.schema = v.GetString("schema")

	ds.transaction = v.GetBool("transaction")

	hostTmpl, err := template.New("host").Funcs(sprig.FuncMap()).Parse(v.GetString("host"))
	if err != nil {
		return Datasource{}, fmt.Errorf("parsing host name provided")
	}

	var host bytes.Buffer
	if err = hostTmpl.Execute(&host, data); err != nil {
		return Datasource{}, fmt.Errorf("expanding host name provided")
	}

	ds.host = host.String()
	if ds.host == "" {
		ds.host = "127.0.0.1"
	}

	portTmpl, err := template.New("port").Funcs(sprig.FuncMap()).Parse(v.GetString("port"))
	if err != nil {
		return Datasource{}, fmt.Errorf("parsing port provided")
	}

	var port bytes.Buffer
	if err = portTmpl.Execute(&port, data); err != nil {
		return Datasource{}, fmt.Errorf("expanding port provided")
	}

	ds.port = port.String()

	userTmpl, err := template.New("user").Funcs(sprig.FuncMap()).Parse(v.GetString("user"))
	if err != nil {
		return Datasource{}, fmt.Errorf("parsing user name provided")
	}

	var user bytes.Buffer
	if err = userTmpl.Execute(&user, data); err != nil {
		return Datasource{}, fmt.Errorf("expanding user name provided")
	}

	ds.user = user.String()

	adminTmpl, err := template.New("admin").Funcs(sprig.FuncMap()).Parse(v.GetString("admin"))
	if err != nil {
		return Datasource{}, fmt.Errorf("parsing admin name provided")
	}

	var admin bytes.Buffer
	if err = adminTmpl.Execute(&admin, data); err != nil {
		return Datasource{}, fmt.Errorf("expanding admin name provided")
	}

	ds.admin = admin.String()

	userPwTmpl, err := template.New("userPw").Funcs(sprig.FuncMap()).Parse(v.GetString("password"))
	if err != nil {
		return Datasource{}, fmt.Errorf("parsing user password provided")
	}

	var userPw bytes.Buffer
	if err = userPwTmpl.Execute(&userPw, data); err != nil {
		return Datasource{}, fmt.Errorf("expanding user password provided")
	}

	ds.userPw = userPw.String()

	adminPwTmpl, err := template.New("adminPw").Funcs(sprig.FuncMap()).Parse(v.GetString("adminpassword"))
	if err != nil {
		return Datasource{}, fmt.Errorf("parsing admin password provided")
	}

	var adminPw bytes.Buffer
	if err = adminPwTmpl.Execute(&adminPw, data); err != nil {
		return Datasource{}, fmt.Errorf("expanding admin password provided")
	}

	ds.adminPw = adminPw.String()

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
		ds.user = v.GetString("user")
		if ds.user == "" {
			ds.user = "postgres"
		}

		if ds.admin == "" {
			ds.admin = "postgres"
		}

		ds.port = v.GetString("port")
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

//OpenDatabase open connection to the corresponding database
func (ds *Datasource) OpenDatabase(log *logrus.Entry, admin bool, nodb bool) (*sql.DB, error) {
	logDb := log.WithField("engine", EngineToString(ds.engine))

	if ds.dstype != Database {
		logDb.Error("Can not open as a database")
		return nil, fmt.Errorf("can not open %s as a database", ds.name)
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
		log.Debug("Opening Mysql database")
		log.Debug(URL)

		driver = "mysql"
	case Postgres:
		log.Debug("Opening Postgresql database")
		log.Debug(URL)

		driver = "postgres"
	}

	db, err := sqlOpen(driver, URL)
	if err != nil {
		log.Error("Opening database failed")
		log.Error(err)

		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 5)
	db.SetMaxIdleConns(10)

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

//Only for unit testing of OpenDatabase function
var mockingSQL = false

func sqlOpen(driver string, url string) (*sql.DB, error) {
	if !mockingSQL {
		return sql.Open(driver, url)
	}

	db, _, err := sqlmock.New()

	return db, err
}
