package json

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/marema31/kamino/provider/common"
)

//KaminoJSONSaver specifc state for database Saver provider
type KaminoJSONSaver struct {
	file    io.WriteCloser
	name    string
	content []map[string]string
}

//NewSaver open the encoding process on provider file and return a Saver compatible object
func NewSaver(ctx context.Context, config map[string]string, name string, file io.WriteCloser) (*KaminoJSONSaver, error) {
	content := make([]map[string]string, 0)
	fmt.Println(name)
	return &KaminoJSONSaver{file: file, name: name, content: content}, nil
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
	js.file.Write(jsonStr)
	js.file.Close()
	return nil
}

//Name give the name of the destination
func (js *KaminoJSONSaver) Name() string {
	return js.name
}

//Reset reinitialize the destination (if possible)
func (js *KaminoJSONSaver) Reset() error {
	return nil
}
