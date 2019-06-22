package config

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/marema31/kamino/kaminodb"
)

//getDatabases Lookup the provided folders for database configuration files and return the corresponding map of configuration
func getDatabases(configPath string) (map[string]*kaminodb.KaminoDb, error) {
	dbs := make(map[string]*kaminodb.KaminoDb)

	dbfolder := filepath.Join(configPath, "databases")

	files, err := ioutil.ReadDir(dbfolder)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.Mode().IsRegular() {
			ext := filepath.Ext(file.Name())
			if ext == ".yml" || ext == ".yaml" || ext == ".json" || ext == ".toml" {
				name := strings.TrimSuffix(file.Name(), ext)
				db, err := kaminodb.New(dbfolder, name)
				if err != nil {
					return nil, err
				}
				dbs[name] = db
			}
		}
	}
	return dbs, nil
}
