package datasource

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

//LoadAll Lookup the provided folder for datasource configuration files
func (dss *Datasources) LoadAll(recipePath string, log *logrus.Entry) error {
	dsfolder := filepath.Join(recipePath, "datasources")

	files, err := ioutil.ReadDir(dsfolder)
	if err != nil {
		log.Errorf("Can not list datasources configuration folder: %v", err)
		return err
	}

	for _, file := range files {
		if file.Mode().IsRegular() {
			ext := filepath.Ext(file.Name())
			if ext == ".yml" || ext == ".yaml" || ext == ".json" || ext == ".toml" {
				name := strings.TrimSuffix(file.Name(), ext)
				logDatasource := log.WithField("datasource", name)
				logDatasource.Debug("Parsing datasource configuration")
				ds, err := dss.load(recipePath, name)
				if err != nil {
					logDatasource.Errorf("Unable to parse configuration: %v", err)
					return err
				}
				dss.datasources[name] = &ds

			}
		}
	}

	if len(dss.datasources) == 0 {
		log.Errorf("no datasources configuration files found in %s", dsfolder)
		return fmt.Errorf("no datasources configuration files found in %s", dsfolder)
	}
	// Insert the datasource name in all entry of the dictionnary
	// that correspond to one tag of the tag list
	for _, ds := range dss.datasources {
		for _, tag := range ds.tags {
			if _, ok := dss.tagToDatasource[tag]; ok {
				dss.tagToDatasource[tag] = append(dss.tagToDatasource[tag], ds.name)
			} else {
				dl := make([]string, 0, 1)
				dl = append(dl, ds.name)
				dss.tagToDatasource[tag] = dl
			}
		}
	}
	return nil
}

func (dss *Datasources) load(recipePath string, filename string) (Datasource, error) {
	v := viper.New()

	v.SetConfigName(filename) // The file will be named [filename].json, [filename].yaml or [filename.toml]
	dsfolder := filepath.Join(recipePath, "datasources")
	v.AddConfigPath(dsfolder)
	err := v.ReadInConfig()
	if err != nil {
		return Datasource{}, err
	}

	engine := strings.ToLower(v.GetString("engine"))
	if engine == "" {
		return Datasource{}, fmt.Errorf("the datasource %s does not provide the engine name", filename)
	}

	e, err := StringToEngine(engine)
	if err != nil {
		return Datasource{}, err
	}

	switch e {
	case Mysql, Postgres:
		return loadDatabaseDatasource(filename, v, e, dss.conTimeout, dss.conRetry)
	case JSON, YAML, CSV:
		return loadFileDatasource(recipePath, filename, v, e)
	}
	return Datasource{}, fmt.Errorf("does not how to manage %s datasource engine", engine)
}
