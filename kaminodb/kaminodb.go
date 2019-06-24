package kaminodb

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // Mysql library dynamically called by database/sql
	_ "github.com/lib/pq"              //Postgres library dynamically called by database/sql
	"github.com/spf13/viper"
)

//DbEngine constants for database engine
type DbEngine int

const (
	// Mysql / MariaDB database engine
	Mysql DbEngine = iota
	// Postgres database engine
	Postgres DbEngine = iota
)

//KaminoDb struct provide
type KaminoDb struct {
	Driver      string //TODO: useful ?
	Database    string
	Engine      DbEngine
	URL         string // TODO: useful ?
	Transaction bool
	Schema      string
}

// New initialize a KaminoDb from a config file
func New(path string, filename string) (*KaminoDb, error) {
	v := viper.New()

	v.SetConfigName(filename) // The file will be named [filename].json, [filename].yaml or [filename.toml]
	v.AddConfigPath(path)
	err := v.ReadInConfig()

	if err != nil {
		return nil, err
	}

	var kd KaminoDb

	kd.Database = v.GetString("database")
	if kd.Database == "" {
		return nil, fmt.Errorf("the configuration %s does not provide the database name", filename)
	}

	password := v.GetString("password")

	kd.Schema = v.GetString("schema")

	kd.Transaction = v.GetBool("transaction")

	host := v.GetString("host")
	if host == "" {
		host = "127.0.0.1"
	}

	engine := v.GetString("engine")
	if engine == "" {
		return nil, fmt.Errorf("the configuration %s does not provide the engine name", filename)
	}

	switch engine {
	case "mysql", "maria", "mariadb":
		kd.Engine = Mysql
		user := v.GetString("user")
		if user == "" {
			user = "root"
		}
		port := v.GetString("port")
		if port == "" {
			port = "3306"
		}

		kd.URL = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, kd.Database)
		kd.Driver = "mysql"

	case "pgsql", "postgres":
		kd.Engine = Postgres
		user := v.GetString("user")
		if user == "" {
			user = "postgres"
		}
		port := v.GetString("port")
		if port == "" {
			port = "5432"
		}
		kd.URL = fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, kd.Database)
		kd.Driver = "postgres"
	default:
		return nil, fmt.Errorf("does not how to manage %s database engine", engine)
	}
	return &kd, nil
}

//Open open connection to the corresponding database
func (kd *KaminoDb) Open() (*sql.DB, error) {
	db, err := sql.Open(kd.Driver, kd.URL)
	if err != nil {
		return nil, err
	}

	// Open does not really open the connection and therefore does not test for url is correct, ping will do
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
