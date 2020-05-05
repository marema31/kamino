package yaml

import (
	"context"
	"io"

	"gopkg.in/yaml.v2"

	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/provider/types"
)

//KaminoYAMLSaver specifc state for database Saver provider.
type KaminoYAMLSaver struct {
	ds      datasource.Datasourcer
	file    io.WriteCloser
	name    string
	content []map[string]string
}

//NewSaver open the encoding process on provider file and return a Saver compatible object.
func NewSaver(ctx context.Context, log *logrus.Entry, ds datasource.Datasourcer) (*KaminoYAMLSaver, error) {
	logFile := log.WithField("datasource", ds.GetName())

	file, err := ds.OpenWriteFile(logFile)
	if err != nil {
		return nil, err
	}

	content := make([]map[string]string, 0)
	tv := ds.FillTmplValues()

	return &KaminoYAMLSaver{file: file, ds: ds, name: tv.FilePath, content: content}, nil
}

//Save writes the record to the destination.
func (ys *KaminoYAMLSaver) Save(log *logrus.Entry, record types.Record) error {
	ys.content = append(ys.content, record)
	return nil
}

//Close closes the destination.
func (ys *KaminoYAMLSaver) Close(log *logrus.Entry) error {
	logFile := log.WithField("datasource", ys.ds.GetName())

	yamlStr, err := yaml.Marshal(ys.content)
	if err != nil {
		logFile.Error("Converting to YAML failed")
		logFile.Error(err)

		return err
	}

	_, err = ys.file.Write(yamlStr)
	if err != nil {
		logFile.Error("Writing file failed")
		logFile.Error(err)

		return err
	}

	return ys.ds.CloseFile(logFile)
}

//Name give the name of the destination.
func (ys *KaminoYAMLSaver) Name() string {
	return ys.name
}

//Reset reinitialize the destination (if possible).
func (ys *KaminoYAMLSaver) Reset(log *logrus.Entry) error {
	logFile := log.WithField("datasource", ys.ds.GetName())
	return ys.ds.ResetFile(logFile)
}
