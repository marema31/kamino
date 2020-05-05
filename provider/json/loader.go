package json

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/provider/common"
	"github.com/marema31/kamino/provider/types"
)

//KaminoJSONLoader specifc state for database Saver provider.
type KaminoJSONLoader struct {
	ds         datasource.Datasourcer
	file       io.ReadCloser
	name       string
	content    []map[string]string
	currentRow int
}

//NewLoader open the encoding process on provider file, read the whole file, unMarshal it and return a Loader compatible object.
func NewLoader(ctx context.Context, log *logrus.Entry, ds datasource.Datasourcer) (*KaminoJSONLoader, error) {
	logFile := log.WithField("datasource", ds.GetName())

	file, err := ds.OpenReadFile(logFile)
	if err != nil {
		return nil, err
	}

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		logFile.Error("Reading the JSON file failed")
		logFile.Error(err)

		return nil, err
	}

	tv := ds.FillTmplValues()
	k := KaminoJSONLoader{ds: ds, file: file, name: tv.FilePath, content: nil, currentRow: 0}
	k.content = make([]map[string]string, 0)

	err = json.Unmarshal(byteValue, &k.content)
	if err != nil {
		logFile.Error("Parsing the JSON file failed")
		logFile.Error(err)

		return nil, err
	}

	return &k, nil
}

//Next moves to next record and return false if there is no more records.
func (jl *KaminoJSONLoader) Next() bool {
	return len(jl.content) > jl.currentRow
}

//Load reads the next record and return it.
func (jl *KaminoJSONLoader) Load(log *logrus.Entry) (types.Record, error) {
	logFile := log.WithField("datasource", jl.ds.GetName())

	if jl.currentRow >= len(jl.content) {
		logFile.Error("no more data to read")
		return nil, fmt.Errorf("no more data to read: %w", common.ErrEOF)
	}

	record := jl.content[jl.currentRow]
	jl.currentRow++

	return record, nil
}

//Close closes the datasource.
func (jl *KaminoJSONLoader) Close(log *logrus.Entry) error {
	logFile := log.WithField("datasource", jl.ds.GetName())
	return jl.ds.CloseFile(logFile)
}

//Name give the name of the destination.
func (jl *KaminoJSONLoader) Name() string {
	return jl.name
}
