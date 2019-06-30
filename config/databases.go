package config

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/marema31/kamino/kaminodb"
)

//getDatabases Lookup the provided folders for database configuration files and return the corresponding map of configuration
func getDatabases(configPath string) (map[string]map[string]map[string]*kaminodb.KaminoDb, error) {
	dbs := make(map[string]map[string]map[string]*kaminodb.KaminoDb)

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
				if _, ok := dbs[db.Name]; !ok {
					dbs[db.Name] = make(map[string]map[string]*kaminodb.KaminoDb)
				}
				if _, ok := dbs[db.Name][db.Environment]; !ok {
					dbs[db.Name][db.Environment] = make(map[string]*kaminodb.KaminoDb)
				}
				dbs[db.Name][db.Environment][db.Instance] = db
			}
		}
	}
	return dbs, nil
}
