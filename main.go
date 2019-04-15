package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/marema31/kamino/database"
)

var (
	ctx context.Context
	db  *sql.DB
)

func main() {
	host := "192.168.9.122"
	port := 3306
	user := "root"
	password := "123soleil"
	sourcedbname := "source1"
	destdbname := "copy2"
	table := "table1"

	ctx := context.Background()

	dbConnection := database.ConnectionInfo{
		Engine:   "mysql",
		Database: sourcedbname,
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
	}

	l, err := database.NewLoader(ctx, &dbConnection, table)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	dbConnection.Database = destdbname
	s, err := database.NewSaver(ctx, &dbConnection, table)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	for l.Next() {
		record, err := l.Load()
		if err != nil {
			log.Fatal(err)
		}

		//		fmt.Printf("%v\n", record["entier"])

		if err = s.Save(record); err != nil {
			log.Fatal(err)
		}

	}

	// Rows.Err will report the last error encountered by Rows.Scan.
	//	if err := rows.Err(); err != nil {
	//		log.Fatal(err)
	//	}
}
