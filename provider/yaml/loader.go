package yaml

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/provider/types"
)

//KaminoYAMLLoader specifc state for database Saver provider
type KaminoYAMLLoader struct {
	ds         datasource.Datasourcer
	file       io.ReadCloser
	name       string
	content    []map[string]string
	currentRow int
}

//NewLoader open the encoding process on provider file, read the whole file, unMarshal it and return a Loader compatible object
func NewLoader(ctx context.Context, log *logrus.Entry, ds datasource.Datasourcer) (*KaminoYAMLLoader, error) {
	logFile := log.WithField("datasource", ds.GetName())
	file, err := ds.OpenReadFile(logFile)
	if err != nil {
		return nil, err
	}

	logFile.Debug("Reading the YAML file to be able to parse it")
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		logFile.Error("Reading the YAML file failed")
		logFile.Error(err)
		return nil, err
	}

	tv := ds.FillTmplValues()
	k := KaminoYAMLLoader{ds: ds, file: file, name: tv.FilePath, content: nil, currentRow: 0}
	k.content = make([]map[string]string, 0)

	err = yaml.Unmarshal([]byte(byteValue), &k.content)
	if err != nil {
		logFile.Error("Parsing the YAML file failed")
		logFile.Error(err)
		return nil, err
	}

	return &k, nil
}

//Next moves to next record and return false if there is no more records
func (yl *KaminoYAMLLoader) Next() bool {

	return len(yl.content) > yl.currentRow
}

//Load reads the next record and return it
func (yl *KaminoYAMLLoader) Load(log *logrus.Entry) (types.Record, error) {
	logFile := log.WithField("datasource", yl.ds.GetName())
	if yl.currentRow >= len(yl.content) {
		logFile.Error("no more data to read")
		return nil, fmt.Errorf("no more data to read")
	}

	record := yl.content[yl.currentRow]
	yl.currentRow++
	return record, nil

}

//Close closes the datasource
func (yl *KaminoYAMLLoader) Close(log *logrus.Entry) error {
	logFile := log.WithField("datasource", yl.ds.GetName())
	return yl.ds.CloseFile(logFile)
}

//Name give the name of the destination
func (yl *KaminoYAMLLoader) Name() string {
	return yl.name
}
