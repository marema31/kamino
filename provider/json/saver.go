package json

import (
	"context"
	"encoding/json"
	"io"

	"github.com/marema31/kamino/provider/common"
)

//KaminoJsonSaver specifc state for database Saver provider
type KaminoJSONSaver struct {
	file    io.WriteCloser
	content []map[string]string
}

//NewSaver open the encoding process on provider file and return a Saver compatible object
func NewSaver(ctx context.Context, config map[string]string, file io.WriteCloser) (*KaminoJSONSaver, error) {
	content := make([]map[string]string, 0)
	return &KaminoJSONSaver{file: file, content: content}, nil
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
