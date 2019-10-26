package csv

import (
	"context"
	"encoding/csv"
	"io"

	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/provider/types"
)

//KaminoCsvLoader specifc state for database Saver provider
type KaminoCsvLoader struct {
	ds           datasource.Datasourcer
	file         io.ReadCloser
	reader       csv.Reader
	name         string
	colNames     []string
	currentRow   []string
	currentError error
}

//NewLoader open the encoding process on provider file, read the header from the first line and return a Loader compatible object
func NewLoader(ctx context.Context, ds datasource.Datasourcer) (*KaminoCsvLoader, error) {
	file, err := ds.OpenReadFile()
	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(file)
	//Read the header to dertermine the column order
	row, err := reader.Read()
	if err != nil {
		return nil, err
	}
	tv := ds.FillTmplValues()

	return &KaminoCsvLoader{ds: ds, file: file, name: tv.FilePath, reader: *reader, colNames: row, currentRow: nil, currentError: nil}, nil
}

//Next moves to next record and return false if there is no more records
func (cl *KaminoCsvLoader) Next() bool {
	row, err := cl.reader.Read()
	if err == io.EOF {
		cl.currentRow = nil
		return false
	} else if err != nil {
		// To conserve the interface, we can return the error here but in Load call
		cl.currentRow = nil
		cl.currentError = err
		return true
	}

	cl.currentRow = row
	return true
}

//Load reads the next record and return it
func (cl *KaminoCsvLoader) Load() (types.Record, error) {
	if cl.currentError != nil {
		return nil, cl.currentError
	}

	record := make(types.Record, len(cl.colNames))
	for i, col := range cl.colNames {
		record[col] = string(cl.currentRow[i])
	}

	return record, nil

}

//Close closes the datasource
func (cl *KaminoCsvLoader) Close() error {
	return cl.ds.CloseFile()
}

//Name give the name of the destination
func (cl *KaminoCsvLoader) Name() string {
	return cl.name
}
