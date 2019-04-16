package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/marema31/kamino/provider"

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

	source, err := database.NewLoader(ctx, &dbConnection, table)
	if err != nil {
		log.Fatal(err)
	}
	defer source.Close()

	var destinations []provider.Saver

	for i := 1; 3 > i; i++ {
		dbConnection.Database = fmt.Sprintf("copy%d", i)
		d, err := database.NewSaver(ctx, &dbConnection, table)
		if err != nil {
			log.Fatal(err)
		}
		defer d.Close()

		destinations = append(destinations, d)
	}

	for source.Next() {
		record, err := source.Load()
		if err != nil {
			log.Fatal(err)
		}

		//		fmt.Printf("%v\n", record["entier"])

		for _, d := range destinations {
			if err = d.Save(record); err != nil {
				log.Fatal(err)
			}
		}

	}

	// Rows.Err will report the last error encountered by Rows.Scan.
	//	if err := rows.Err(); err != nil {
	//		log.Fatal(err)
	//	}
}
