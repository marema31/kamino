package csv

import (
	"context"
	"encoding/csv"
	"io"
	"sort"

	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/provider/types"
)

//KaminoCsvSaver specifc state for database Saver provider
type KaminoCsvSaver struct {
	ds       datasource.Datasourcer
	file     io.WriteCloser
	name     string
	writer   csv.Writer
	colNames []string
}

//NewSaver open the encoding process on provider file and return a Saver compatible object
func NewSaver(ctx context.Context, log *logrus.Entry, ds datasource.Datasourcer) (*KaminoCsvSaver, error) {
	logFile := log.WithField("datasource", ds.GetName())
	file, err := ds.OpenWriteFile(logFile)
	if err != nil {
		return nil, err
	}
	writer := csv.NewWriter(file)
	tv := ds.FillTmplValues()
	return &KaminoCsvSaver{file: file, ds: ds, name: tv.FilePath, writer: *writer, colNames: nil}, nil
}

//Save writes the record to the destination
func (cs *KaminoCsvSaver) Save(log *logrus.Entry, record types.Record) error {
	logFile := log.WithField("datasource", cs.ds.GetName())
	// Is this method is called for the first time
	//If yes fix the column order in csv file
	if cs.colNames == nil {
		//The order of columns could change between two executions of the test
		//TODO: add an options to datasource to provides the column order
		var keys []string
		for col := range record {
			keys = append(keys, col)
		}
		sort.Strings(keys)
		cs.colNames = keys
		err := cs.writer.Write(cs.colNames)
		if err != nil {
			logFile.Error("Writing file failed")
			logFile.Error(err)
			return nil
		}
	}

	row := make([]string, len(cs.colNames))

	for i, col := range cs.colNames {
		row[i] = record[col]
	}

	return cs.writer.Write(row)
}

//Close closes the destination
func (cs *KaminoCsvSaver) Close(log *logrus.Entry) error {
	logFile := log.WithField("datasource", cs.ds.GetName())
	cs.writer.Flush()
	return cs.ds.CloseFile(logFile)
}

//Name give the name of the destination
func (cs *KaminoCsvSaver) Name() string {
	return cs.name
}

//Reset reinitialize the destination (if possible)
func (cs *KaminoCsvSaver) Reset(log *logrus.Entry) error {
	logFile := log.WithField("datasource", cs.ds.GetName())
	return cs.ds.ResetFile(logFile)
}
