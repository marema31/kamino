package datasource

import (
	"database/sql"
	"fmt"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql" // Mysql library dynamically called by database/sql
	_ "github.com/lib/pq"              //Postgres library dynamically called by database/sql
	"github.com/spf13/viper"
)

// load a database type datasource from the viper configuration
func loadDatabaseDatasource(filename string, v *viper.Viper, ds *Datasource) error {
	ds.Database = v.GetString("database")
	if ds.Database == "" {
		return fmt.Errorf("the datasource %s does not provide the database name", ds.Name)
	}

	ds.Schema = v.GetString("schema")

	ds.Transaction = v.GetBool("transaction")

	host := v.GetString("host")
	if host == "" {
		host = "127.0.0.1"
	}
	port := v.GetString("port")

	user := v.GetString("user")
	admin := v.GetString("admin")
	userpw := v.GetString("password")
	adminpw := v.GetString("adminpassword")
	if adminpw == "" {
		adminpw = userpw
	}
	if userpw == "" {
		userpw = adminpw
	}

	switch ds.Engine {
	case Mysql:
		if user == "" {
			user = "root"
		}
		if admin == "" {
			admin = "root"
		}
		if port == "" {
			port = "3306"
		}

		ds.URL = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, userpw, host, port, ds.Database)
		ds.URLAdmin = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", admin, adminpw, host, port, ds.Database)
		ds.URLNoDb = fmt.Sprintf("%s:%s@tcp(%s:%s)", admin, adminpw, host, port)

	case Postgres:
		user := v.GetString("user")
		if user == "" {
			user = "postgres"
		}
		if admin == "" {
			admin = "postgres"
		}
		port := v.GetString("port")
		if port == "" {
			port = "5432"
		}
		//TODO: try without ssldisable or make it a optione od datasource
		//TODO: manage ds.Schema
		ds.URL = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, userpw, ds.Database)
		ds.URLAdmin = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, admin, adminpw, ds.Database)
		ds.URLNoDb = fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable", host, port, admin, adminpw)
	}
	return nil
}

//OpenDatabase open connection to the corresponding database
func (ds *Datasource) OpenDatabase(admin bool, nodb bool) (*sql.DB, error) {
	if ds.Type != Database {
		return nil, fmt.Errorf("The datasource %s is not a database cannot open it", ds.Name)
	}
	URL := ds.URL
	if admin {
		URL = ds.URLAdmin
	}
	if nodb {
		URL = ds.URLNoDb
	}

	var driver string
	switch ds.Engine {
	case Mysql:
		driver = "mysql"
	case Postgres:
		driver = "postgres"
	}

	db, err := sqlOpen(driver, URL)
	if err != nil {
		return nil, err
	}

	// Open does not really do a connection and therefore does not test for url is correct, ping will do
	err = db.Ping()
	if err != nil {
		return nil, err
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
