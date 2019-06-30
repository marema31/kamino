package json

import (
	"context"
	"encoding/json"
	"io"

	"github.com/marema31/kamino/provider/common"
)

//KaminoJSONSaver specifc state for database Saver provider
type KaminoJSONSaver struct {
	file    io.WriteCloser
	name    string
	tmpName string
	content []map[string]string
}

//NewSaver open the encoding process on provider file and return a Saver compatible object
func NewSaver(ctx context.Context, config map[string]string, name string, tmpName string, file io.WriteCloser) (*KaminoJSONSaver, error) {
	content := make([]map[string]string, 0)
	return &KaminoJSONSaver{file: file, name: name, tmpName: tmpName, content: content}, nil
}

//Save writes the record to the destination
func (js *KaminoJSONSaver) Save(record common.Record) error {
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
	return common.CloseWriter(js.file, js.tmpName, js.name)
}

//Name give the name of the destination
func (js *KaminoJSONSaver) Name() string {
	return js.name
}

//Reset reinitialize the destination (if possible)
func (js *KaminoJSONSaver) Reset() error {
	return common.ResetWriter(js.file, js.tmpName, js.name)
}
