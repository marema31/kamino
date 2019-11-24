package datasource

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql" // Mysql library dynamically called by database/sql
	_ "github.com/lib/pq"              //Postgres library dynamically called by database/sql
	"github.com/spf13/viper"
)

// load a database type datasource from the viper configuration
func loadDatabaseDatasource(filename string, v *viper.Viper, engine Engine, connectionTimeout time.Duration, connectionRetry int) (Datasource, error) {
	var ds Datasource
	ds.dstype = Database
	ds.engine = engine
	ds.name = filename
	ds.database = v.GetString("database")
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

	ds.host = v.GetString("host")
	if ds.host == "" {
		ds.host = "127.0.0.1"
	}
	ds.port = v.GetString("port")

	ds.user = v.GetString("user")
	ds.admin = v.GetString("admin")
	ds.userPw = v.GetString("password")
	ds.adminPw = v.GetString("adminpassword")
	if ds.adminPw == "" {
		ds.adminPw = ds.userPw
	}
	if ds.userPw == "" {
		ds.userPw = ds.adminPw
	}

	//TODO: Allow the user to define options for the DSN (like MySql allowOldPasswords=1 use for some very old DB)
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

		//use parseTime=true to force date and time conversion
		//TODO: may disturb the sync mechanism
		ds.url = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", ds.user, ds.userPw, ds.host, ds.port, ds.database)
		ds.urlAdmin = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", ds.admin, ds.adminPw, ds.host, ds.port, ds.database)
		ds.urlNoDb = fmt.Sprintf("%s:%s@tcp(%s:%s)/mysql?parseTime=true", ds.admin, ds.adminPw, ds.host, ds.port)

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
		//TODO: try without ssldisable or make it a option on datasource (for testing see: https://gist.github.com/mrw34/c97bb03ea1054afb551886ffc8b63c3b)
		//TODO: manage ds.Schema
		ds.url = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", ds.host, ds.port, ds.user, ds.userPw, ds.database)
		ds.urlAdmin = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", ds.host, ds.port, ds.admin, ds.adminPw, ds.database)
		ds.urlNoDb = fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable", ds.host, ds.port, ds.admin, ds.adminPw)
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
	var URL string
	var db *sql.DB
	if nodb {
		//		logDb.Debug("Openning connection to database engine in Admin")
		URL = ds.urlNoDb
		db = ds.dbNoDb
	} else if admin {
		//		logDb.Debugf("Openning database %s in Admin", ds.database)
		URL = ds.urlAdmin
		db = ds.dbAdmin
	} else {
		//		logDb.Debugf("Openning database %s in User", ds.database)
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
		driver = "mysql"
		log.Debug("Opening Mysql database")
		log.Debug(URL)
	case Postgres:
		driver = "postgres"
		log.Debug("Opening Postgresql database")
		log.Debug(URL)
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
	if nodb {
		ds.dbNoDb = db
	} else if admin {
		ds.dbAdmin = db
	} else {
		ds.db = db
	}
	return db, nil
}

//Only for unit testing of OpenDatabase function
var mockingSQL = false

func sqlOpen(driver string, URL string) (*sql.DB, error) {
	if !mockingSQL {
		return sql.Open(driver, URL)
	}
	db, _, err := sqlmock.New()
	return db, err
}
