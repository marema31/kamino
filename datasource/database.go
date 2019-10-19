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
func loadDatabaseDatasource(filename string, v *viper.Viper, engine Engine) (*Datasource, error) {
	var ds Datasource
	ds.Type = Database
	ds.Engine = engine
	ds.Name = filename
	ds.Database = v.GetString("database")
	if ds.Database == "" {
		return nil, fmt.Errorf("the datasource %s does not provide the database name", ds.Name)
	}

	ds.Tags = v.GetStringSlice("tags")
	if len(ds.Tags) == 0 {
		ds.Tags = []string{""}
	}

	ds.Schema = v.GetString("schema")

	ds.Transaction = v.GetBool("transaction")

	ds.Host = v.GetString("host")
	if ds.Host == "" {
		ds.Host = "127.0.0.1"
	}
	ds.Port = v.GetString("port")

	ds.User = v.GetString("user")
	ds.Admin = v.GetString("admin")
	ds.UserPw = v.GetString("password")
	ds.AdminPw = v.GetString("adminpassword")
	if ds.AdminPw == "" {
		ds.AdminPw = ds.UserPw
	}
	if ds.UserPw == "" {
		ds.UserPw = ds.AdminPw
	}

	switch ds.Engine {
	case Mysql:
		if ds.User == "" {
			ds.User = "root"
		}
		if ds.Admin == "" {
			ds.Admin = "root"
		}
		if ds.Port == "" {
			ds.Port = "3306"
		}

		ds.URL = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", ds.User, ds.UserPw, ds.Host, ds.Port, ds.Database)
		ds.URLAdmin = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", ds.Admin, ds.AdminPw, ds.Host, ds.Port, ds.Database)
		ds.URLNoDb = fmt.Sprintf("%s:%s@tcp(%s:%s)", ds.Admin, ds.AdminPw, ds.Host, ds.Port)

	case Postgres:
		ds.User = v.GetString("user")
		if ds.User == "" {
			ds.User = "postgres"
		}
		if ds.Admin == "" {
			ds.Admin = "postgres"
		}
		ds.Port = v.GetString("port")
		if ds.Port == "" {
			ds.Port = "5432"
		}
		//TODO: try without ssldisable or make it a option on datasource
		//TODO: manage ds.Schema
		ds.URL = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", ds.Host, ds.Port, ds.User, ds.UserPw, ds.Database)
		ds.URLAdmin = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", ds.Host, ds.Port, ds.Admin, ds.AdminPw, ds.Database)
		ds.URLNoDb = fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable", ds.Host, ds.Port, ds.Admin, ds.AdminPw)
	}
	return &ds, nil
}

//OpenDatabase open connection to the corresponding database
func (ds *Datasource) OpenDatabase(admin bool, nodb bool) (*sql.DB, error) {
	if ds.Type != Database {
		return nil, fmt.Errorf("the datasource %s is not a database cannot open it", ds.Name)
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
