package json

import (
	"context"
	"encoding/json"
	"io"

	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/provider/types"
)

//KaminoJSONSaver specifc state for database Saver provider
type KaminoJSONSaver struct {
	ds      datasource.Datasourcer
	file    io.WriteCloser
	name    string
	content []map[string]string
}

//NewSaver open the encoding process on provider file and return a Saver compatible object
func NewSaver(ctx context.Context, log *logrus.Entry, ds datasource.Datasourcer) (*KaminoJSONSaver, error) {
	logFile := log.WithField("datasource", ds.GetName())
	file, err := ds.OpenWriteFile(logFile)
	if err != nil {
		return nil, err
	}
	content := make([]map[string]string, 0)

	tv := ds.FillTmplValues()
	return &KaminoJSONSaver{file: file, ds: ds, name: tv.FilePath, content: content}, nil
}

//Save writes the record to the destination
func (js *KaminoJSONSaver) Save(log *logrus.Entry, record types.Record) error {
	js.content = append(js.content, record)
	return nil
}

//Close closes the destination
func (js *KaminoJSONSaver) Close(log *logrus.Entry) error {
	logFile := log.WithField("datasource", js.ds.GetName())
	jsonStr, err := json.MarshalIndent(js.content, "", "    ")
	if err != nil {
		logFile.Error("Converting to JSON failed")
		logFile.Error(err)
		return err
	}
	_, err = js.file.Write(jsonStr)
	if err != nil {
		logFile.Error("Writing file failed")
		logFile.Error(err)
		return err
	}
	return js.ds.CloseFile(logFile)
}

//Name give the name of the destination
func (js *KaminoJSONSaver) Name() string {
	return js.name
}

//Reset reinitialize the destination (if possible)
func (js *KaminoJSONSaver) Reset(log *logrus.Entry) error {
	logFile := log.WithField("datasource", js.ds.GetName())
	return js.ds.ResetFile(logFile)
}
