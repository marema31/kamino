package main

import (
	"context"
	"fmt"
	"log"

	"github.com/marema31/kamino/provider"
)

func main() {

	ctx := context.Background()

	dbConnection := map[string]string{
		"type":     "database",
		"engine":   "mysql",
		"database": "source1",
		"host":     "192.168.9.122",
		"port":     "3306",
		"user":     "root",
		"password": "123soleil",
		"table":    "table1",
	}
	//source, err := provider.NewLoader(ctx, dbConnection)

	/*	csvInput := map[string]string{
			"type": "csv",
			"file": "/tmp/kaminoIn.csv",
		}
		source, err := provider.NewLoader(ctx, csvInput)
	*/
	jsonInput := map[string]string{
		"type": "json",
		"file": "/tmp/kaminoIn.json",
	}

	source, err := provider.NewLoader(ctx, jsonInput)
	if err != nil {
		log.Fatal(err)
	}
	defer source.Close()

	var destinations []provider.Saver

	for i := 1; 3 > i; i++ {
		dbConnection["database"] = fmt.Sprintf("copy%d", i)
		d, err := provider.NewSaver(ctx, dbConnection)
		if err != nil {
			log.Fatal(err)
		}
		defer d.Close()

		destinations = append(destinations, d)
	}

	csvOutput := map[string]string{
		"type": "csv",
		"file": "/tmp/kaminoOut.csv",
	}
	d, err := provider.NewSaver(ctx, csvOutput)
	if err != nil {
		log.Fatal(err)
	}
	defer d.Close()

	jsonOutput := map[string]string{
		"type": "json",
		"file": "/tmp/kaminoOut.json",
	}

	d, err = provider.NewSaver(ctx, jsonOutput)
	if err != nil {
		log.Fatal(err)
	}
	defer d.Close()

	destinations = append(destinations, d)

	stdOutput := map[string]string{
		"type": "json",
		"std":  "true",
	}
	d, err = provider.NewSaver(ctx, stdOutput)
	if err != nil {
		log.Fatal(err)
	}
	defer d.Close()

	destinations = append(destinations, d)

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

}
