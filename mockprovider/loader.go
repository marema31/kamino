package mockprovider

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/provider/common"
	"github.com/marema31/kamino/provider/types"
)

//MockLoader specifc state for database Saver provider.
type MockLoader struct {
	MockName   string
	Content    []map[string]string
	CurrentRow int
	ErrorClose error
	ErrorLoad  error
}

//Next moves to next record and return false if there is no more records.
func (ml *MockLoader) Next() bool {
	return len(ml.Content) > ml.CurrentRow
}

//Load reads the next record and return it.
func (ml *MockLoader) Load(log *logrus.Entry) (types.Record, error) {
	if ml.ErrorLoad != nil {
		return nil, ml.ErrorLoad
	}

	if ml.CurrentRow >= len(ml.Content) {
		return nil, fmt.Errorf("no more data to read: %w", common.ErrEOF)
	}

	record := ml.Content[ml.CurrentRow]
	ml.CurrentRow++

	return record, nil
}

//Close closes the datasource.
func (ml *MockLoader) Close(log *logrus.Entry) error {
	return ml.ErrorClose
}

//Name give the name of the destination.
func (ml *MockLoader) Name() string {
	return ml.MockName
}
