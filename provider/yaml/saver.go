package yaml

import (
	"context"
	"io"

	"gopkg.in/yaml.v2"

	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/provider/common"
)

//KaminoYAMLSaver specifc state for database Saver provider
type KaminoYAMLSaver struct {
	ds      *datasource.Datasource
	file    io.WriteCloser
	name    string
	content []map[string]string
}

//NewSaver open the encoding process on provider file and return a Saver compatible object
func NewSaver(ctx context.Context, ds *datasource.Datasource) (*KaminoYAMLSaver, error) {
	file, err := ds.OpenWriteFile()
	if err != nil {
		return nil, err
	}
	content := make([]map[string]string, 0)
	return &KaminoYAMLSaver{file: file, ds: ds, name: ds.FilePath, content: content}, nil
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
	return ys.ds.CloseFile()
}

//Name give the name of the destination
func (ys *KaminoYAMLSaver) Name() string {
	return ys.name
}

//Reset reinitialize the destination (if possible)
func (ys *KaminoYAMLSaver) Reset() error {
	return ys.ds.ResetFile()
}
