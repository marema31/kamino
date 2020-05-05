package datasource

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

//LoadAll Lookup the provided folder for datasource configuration files.
func (dss *Datasources) LoadAll(recipePath string, log *logrus.Entry) error {
	firstError := dss.findRecipes(recipePath, "datasources", log)

	if len(dss.datasources) == 0 {
		log.Errorf("no datasources configuration files found in %s/datasources", recipePath)
		return fmt.Errorf("no datasources configuration files found in %s/datasources: %w", recipePath, errNoConfiguration)
	}

	// Insert the datasource name in all entry of the dictionary
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

	return firstError
}

func (dss *Datasources) findRecipes(recipePath string, subpath string, log *logrus.Entry) error {
	var firstError error = nil

	dsfolder := filepath.Join(recipePath, subpath)

	files, err := ioutil.ReadDir(dsfolder)
	if err != nil {
		log.Errorf("Can not list datasources configuration folder: %v", err)
		return err
	}

	for _, file := range files {
		ext := filepath.Ext(file.Name())

		if file.Mode().IsDir() {
			err = dss.findRecipes(recipePath, filepath.Join(subpath, file.Name()), log)
			if err != nil && firstError == nil {
				firstError = err
			}
		} else if file.Mode().IsRegular() && ext == ".yml" || ext == ".yaml" || ext == ".json" || ext == ".toml" {
			name := strings.TrimSuffix(file.Name(), ext)
			logDatasource := log.WithField("datasource", name)
			logDatasource.Debug("Parsing datasource configuration")

			ds, err := dss.load(recipePath, subpath, name)
			if err != nil {
				logDatasource.Errorf("Unable to parse configuration: %v", err)
				if firstError == nil {
					firstError = err
				}
			}

			dss.datasources[name] = &ds
		}
	}

	return firstError
}

func (dss *Datasources) load(recipePath string, subpath string, filename string) (Datasource, error) {
	v := viper.New()

	v.SetConfigName(filename) // The file will be named [filename].json, [filename].yaml or [filename.toml]
	v.AddConfigPath(filepath.Join(recipePath, subpath))

	err := v.ReadInConfig()
	if err != nil {
		return Datasource{}, err
	}

	engine := strings.ToLower(v.GetString("engine"))
	if engine == "" {
		return Datasource{}, fmt.Errorf("the datasource %s does not provide the engine name: %w", filename, errMissingParameter)
	}

	e, err := StringToEngine(engine)
	if err != nil {
		return Datasource{}, err
	}

	switch e {
	case Mysql, Postgres:
		return loadDatabaseDatasource(filename, v, e, dss.envVar, dss.conTimeout, dss.conRetry)
	case JSON, YAML, CSV:
		return loadFileDatasource(recipePath, filename, v, e, dss.envVar)
	}
	//Should never come here, error will be raised by StringToEngine
	return Datasource{}, fmt.Errorf("does not how to manage %s datasource engine: %w", engine, errWrongParameterValue)
}
