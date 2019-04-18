package csv

import (
	"context"
	"encoding/csv"
	"io"

	"github.com/marema31/kamino/provider/common"
)

//KaminoCsvSaver specifc state for database Saver provider
type KaminoCsvSaver struct {
	file     io.WriteCloser
	writer   csv.Writer
	colNames []string
}

//NewSaver open the encoding process on provider file and return a Saver compatible object
func NewSaver(ctx context.Context, config map[string]string, file io.WriteCloser) (*KaminoCsvSaver, error) {
	writer := csv.NewWriter(file)
	return &KaminoCsvSaver{file: file, writer: *writer, colNames: nil}, nil
}

//Save writes the record to the destination
func (cs *KaminoCsvSaver) Save(record common.Record) error {
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
func (cs *KaminoCsvSaver) Close() {
	cs.writer.Flush()
	cs.file.Close()
}
