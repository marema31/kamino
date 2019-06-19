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
	name    string
	tmpName string
	content []map[string]string
}

//NewSaver open the encoding process on provider file and return a Saver compatible object
func NewSaver(ctx context.Context, config map[string]string, name string, tmpName string, file io.WriteCloser) (*KaminoYAMLSaver, error) {
	content := make([]map[string]string, 0)
	return &KaminoYAMLSaver{file: file, name: name, tmpName: tmpName, content: content}, nil
}

//Save writes the record to the destination
func (ys *KaminoYAMLSaver) Save(record common.Record) error {
	ys.content = append(ys.content, record)
	return nil
}

//Close closes the destination
func (ys *KaminoYAMLSaver) Close() error {
	yamlStr, err := yaml.Marshal(ys.content)
	if err != nil {
		return err
	}
	_, err = ys.file.Write(yamlStr)
	if err != nil {
		return err
	}
	return common.CloseWriter(ys.file, ys.tmpName, ys.name)
}

//Name give the name of the destination
func (ys *KaminoYAMLSaver) Name() string {
	return ys.name
}

//Reset reinitialize the destination (if possible)
func (ys *KaminoYAMLSaver) Reset() error {
	return common.ResetWriter(ys.file, ys.tmpName, ys.name)
}
