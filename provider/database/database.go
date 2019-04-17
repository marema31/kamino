package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // Mysql library dynamically called by database/sql
	_ "github.com/lib/pq"              //Postgres library dynamically called by database/sql
)

type kaminoDb struct {
	db       *sql.DB
	driver   string
	database string
	table    string
}

func newKaminoDb(config map[string]string) (*kaminoDb, error) {
	var url, driver string

	engine, ok := config["engine"]
	if !ok {
		return nil, fmt.Errorf("the configuration block does not provide the engine name")
	}
	database, ok := config["database"]
	if !ok {
		return nil, fmt.Errorf("the configuration block does not provide the database name")
	}
	table, ok := config["table"]
	if !ok {
		return nil, fmt.Errorf("the configuration block does not provide the table name")
	}

	password, ok := config["password"]
	if !ok {
		password = ""
	}
	host, ok := config["host"]
	if !ok {
		host = "127.0.0.1"
	}

	switch engine {
	case "mysql", "maria", "mariadb":
		user, ok := config["user"]
		if !ok {
			user = "root"
		}
		port, ok := config["port"]
		if !ok {
			port = "3306"
		}

		url = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, database)
		driver = "mysql"
	case "pgsql", "postgres":
		user, ok := config["user"]
		if !ok {
			user = "postgres"
		}
		port, ok := config["port"]
		if !ok {
			port = "3306"
		}
		url = fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, database)
		driver = "pq"
	default:
		return nil, fmt.Errorf("does not how to manage %s database engine", engine)
	}

	db, err := sql.Open(driver, url)
	if err != nil {
		return nil, err
	}

	// Open does not really open the connection and therefore does not test for url is correct, ping will do
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &kaminoDb{
		db:       db,
		driver:   driver,
		database: database,
		table:    table,
	}, nil
}
