package datasource

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

//LoadAll Lookup the provided folder for datasource configuration files
func (dss *Datasources) LoadAll(configPath string) error {
	dsfolder := filepath.Join(configPath, "datasources")

	files, err := ioutil.ReadDir(dsfolder)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.Mode().IsRegular() {
			ext := filepath.Ext(file.Name())
			if ext == ".yml" || ext == ".yaml" || ext == ".json" || ext == ".toml" {
				name := strings.TrimSuffix(file.Name(), ext)
				ds, err := dss.load(dsfolder, name)
				if err != nil {
					return err
				}
				dss.datasources[name] = ds

			}
		}
	}
	// Insert the datasource name in all entry of the dictionnary
	// that correspond to one tag of the tag list
	for _, ds := range dss.datasources {
		for _, tag := range ds.Tags {
			if _, ok := dss.tagToDatasource[tag]; ok {
				dss.tagToDatasource[tag] = append(dss.tagToDatasource[tag], ds.Name)
			} else {
				dl := make([]string, 0, 1)
				dl = append(dl, ds.Name)
				dss.tagToDatasource[tag] = dl
			}
		}
	}
	return nil
}

func (dss *Datasources) load(path string, filename string) (*Datasource, error) {
	v := viper.New()

	v.SetConfigName(filename) // The file will be named [filename].json, [filename].yaml or [filename.toml]
	v.AddConfigPath(path)
	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	engine := strings.ToLower(v.GetString("engine"))
	if engine == "" {
		return nil, fmt.Errorf("the datasource %s does not provide the engine name", filename)
	}

	switch engine {
	case "mysql", "maria", "mariadb":
		return loadDatabaseDatasource(filename, v, Mysql)
	case "pgsql", "postgres":
		return loadDatabaseDatasource(filename, v, Postgres)
	case "json":
		return loadFileDatasource(filename, v, JSON)
	case "yaml":
		return loadFileDatasource(filename, v, YAML)
	case "csv":
		return loadFileDatasource(filename, v, CSV)
	default:
		return nil, fmt.Errorf("does not how to manage %s datasource engine", engine)
	}
}
