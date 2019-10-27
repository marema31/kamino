package json

import (
	"context"
	"encoding/json"
	"io"

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
func NewSaver(ctx context.Context, ds datasource.Datasourcer) (*KaminoJSONSaver, error) {
	file, err := ds.OpenWriteFile()
	if err != nil {
		return nil, err
	}
	content := make([]map[string]string, 0)

	tv := ds.FillTmplValues()
	return &KaminoJSONSaver{file: file, ds: ds, name: tv.FilePath, content: content}, nil
}

//Save writes the record to the destination
func (js *KaminoJSONSaver) Save(record types.Record) error {
	js.content = append(js.content, record)
	return nil
}

//Close closes the destination
func (js *KaminoJSONSaver) Close() error {
	jsonStr, err := json.MarshalIndent(js.content, "", "    ")
	if err != nil {
		return err
	}
	_, err = js.file.Write(jsonStr)
	if err != nil {
		return err
	}
	return js.ds.CloseFile()
}

//Name give the name of the destination
func (js *KaminoJSONSaver) Name() string {
	return js.name
}

//Reset reinitialize the destination (if possible)
func (js *KaminoJSONSaver) Reset() error {
	return js.ds.ResetFile()
}
