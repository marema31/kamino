package csv

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/provider/common"
	"github.com/marema31/kamino/provider/types"
)

//KaminoCsvLoader specifc state for database Saver provider.
type KaminoCsvLoader struct {
	ds           datasource.Datasourcer
	file         io.ReadCloser
	reader       csv.Reader
	name         string
	colNames     []string
	currentRow   []string
	currentError error
}

//NewLoader open the encoding process on provider file, read the header from the first line and return a Loader compatible object.
func NewLoader(ctx context.Context, log *logrus.Entry, ds datasource.Datasourcer) (*KaminoCsvLoader, error) {
	logFile := log.WithField("datasource", ds.GetName())

	file, err := ds.OpenReadFile(logFile)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(file)

	logFile.Debug("Reading the header to determine the column order")

	row, err := reader.Read()
	if err != nil {
		logFile.Error("Reading CSV header failed")
		logFile.Error(err)

		return nil, err
	}

	colNames := make([]string, 0, len(row))

	for _, col := range row {
		colNames = append(colNames, strings.TrimSpace(col))
	}

	tv := ds.FillTmplValues()

	return &KaminoCsvLoader{ds: ds, file: file, name: tv.FilePath, reader: *reader, colNames: colNames, currentRow: nil, currentError: nil}, nil
}

//Next moves to next record and return false if there is no more records.
func (cl *KaminoCsvLoader) Next() bool {
	row, err := cl.reader.Read()
	if err == io.EOF {
		cl.currentRow = nil
		return false
	} else if err != nil {
		// To conserve the interface, we can not return the error here but in Load call
		cl.currentRow = nil
		cl.currentError = err

		return true
	}

	cl.currentRow = row

	return true
}

//Load reads the next record and return it.
func (cl *KaminoCsvLoader) Load(log *logrus.Entry) (types.Record, error) {
	logFile := log.WithField("datasource", cl.ds.GetName())

	if cl.currentError != nil {
		logFile.Error("Reading CSV next line failed")
		logFile.Error(cl.currentError)

		return nil, cl.currentError
	}

	if cl.currentRow == nil {
		logFile.Error("EOF reached")
		return nil, fmt.Errorf("no more data: %w", common.ErrEOF)
	}

	record := make(types.Record, len(cl.colNames))

	for i, col := range cl.colNames {
		record[col] = cl.currentRow[i]
	}

	return record, nil
}

//Close closes the datasource.
func (cl *KaminoCsvLoader) Close(log *logrus.Entry) error {
	logFile := log.WithField("datasource", cl.ds.GetName())
	return cl.ds.CloseFile(logFile)
}

//Name give the name of the destination.
func (cl *KaminoCsvLoader) Name() string {
	return cl.name
}
