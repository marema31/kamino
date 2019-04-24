package yaml

import (
	"context"
	"io"

	"gopkg.in/yaml.v2"

	"github.com/marema31/kamino/provider/common"
)

//KaminoYAMLSaver specifc state for database Saver provider
type KaminoYAMLSaver struct {
	file    io.WriteCloser
	content []map[string]string
}

//NewSaver open the encoding process on provider file and return a Saver compatible object
func NewSaver(ctx context.Context, config map[string]string, file io.WriteCloser) (*KaminoYAMLSaver, error) {
	content := make([]map[string]string, 0)
	return &KaminoYAMLSaver{file: file, content: content}, nil
}

//Save writes the record to the destination
func (js *KaminoYAMLSaver) Save(record common.Record) error {
	js.content = append(js.content, record)
	return nil
}

//Close closes the destination
func (js *KaminoYAMLSaver) Close() error {
	yamlStr, err := yaml.Marshal(js.content)
	if err != nil {
		return err
	}
	js.file.Write(yamlStr)
	js.file.Close()
	return nil
}
