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
func LoadAll(configPath string) error {
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
				err := load(dsfolder, name)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func load(path string, filename string) error {
	v := viper.New()

	v.SetConfigName(filename) // The file will be named [filename].json, [filename].yaml or [filename.toml]
	v.AddConfigPath(path)
	err := v.ReadInConfig()

	if err != nil {
		return err
	}

	var ds Datasource
	ds.Name = filename
	engine := strings.ToLower(v.GetString("engine"))
	if engine == "" {
		return fmt.Errorf("the datasource %s does not provide the engine name", ds.Name)
	}

	datasources[ds.Name] = &ds
	ds.Tags = v.GetStringSlice("tags")
	if len(ds.Tags) == 0 {
		ds.Tags = []string{""}
	}
	insertTag(ds.Tags, ds.Name)

	switch engine {
	case "mysql", "maria", "mariadb":
		ds.Type = Database
		ds.Engine = Mysql
		return loadDatabaseDatasource(filename, v, &ds)
	case "pgsql", "postgres":
		ds.Type = Database
		ds.Engine = Postgres
		return loadDatabaseDatasource(filename, v, &ds)
	case "json":
		ds.Type = File
		ds.Engine = JSON
		return loadFileDatasource(filename, v, &ds)
	case "yaml":
		ds.Type = File
		ds.Engine = YAML
		return loadFileDatasource(filename, v, &ds)
	case "csv":
		ds.Type = File
		ds.Engine = CSV
		return loadFileDatasource(filename, v, &ds)
	default:
		return fmt.Errorf("does not how to manage %s datasource engine", engine)
	}
}
