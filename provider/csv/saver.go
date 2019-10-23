package csv

import (
	"context"
	"encoding/csv"
	"io"

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
func NewSaver(ctx context.Context, ds datasource.Datasourcer) (*KaminoCsvSaver, error) {
	file, err := ds.OpenWriteFile()
	if err != nil {
		return nil, err
	}
	writer := csv.NewWriter(file)
	tv := ds.FillTmplValues()
	return &KaminoCsvSaver{file: file, ds: ds, name: tv.FilePath, writer: *writer, colNames: nil}, nil
}

//Save writes the record to the destination
func (cs *KaminoCsvSaver) Save(record types.Record) error {
	// Is this method is called for the first time
	//If yes fix the column order in csv file
	if cs.colNames == nil {
		for col := range record {
			cs.colNames = append(cs.colNames, col)
		}
		err := cs.writer.Write(cs.colNames)
		if err != nil {
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
func (cs *KaminoCsvSaver) Close() error {
	cs.writer.Flush()
	return cs.ds.CloseFile()
}

//Name give the name of the destination
func (cs *KaminoCsvSaver) Name() string {
	return cs.name
}

//Reset reinitialize the destination (if possible)
func (cs *KaminoCsvSaver) Reset() error {
	return cs.ds.ResetFile()
}
