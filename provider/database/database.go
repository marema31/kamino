package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

type ConnectionInfo struct {
	Engine   string
	Database string
	Host     string
	Port     int
	User     string
	Password string
}

type kaminoDb struct {
	db       *sql.DB
	driver   string
	database string
	table    string
}

func New(c *ConnectionInfo, table string) (*kaminoDb, error) {
	var url, driver string

	switch c.Engine {
	case "mysql", "maria", "mariadb":
		url = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", c.User, c.Password, c.Host, c.Port, c.Database)
		driver = "mysql"
	case "pgsql", "postgres":
		url = fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Password, c.Database)
		driver = "pq"
	default:
		return nil, fmt.Errorf("Does not how to manage %s database engine", c.Engine)
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
		database: c.Database,
		table:    table,
	}, nil
}
