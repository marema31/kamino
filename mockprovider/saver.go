package mockprovider

import (
	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/provider/types"
)

//MockSaver specifc state for database Saver provider.
type MockSaver struct {
	MockName   string
	Content    []map[string]string
	ErrorClose error
	ErrorReset error
	ErrorSave  error
}

//Save writes the record to the destination.
func (ms *MockSaver) Save(log *logrus.Entry, record types.Record) error {
	if ms.ErrorSave != nil {
		return ms.ErrorSave
	}

	ms.Content = append(ms.Content, record)

	return nil
}

//Close closes the destination.
func (ms *MockSaver) Close(log *logrus.Entry) error {
	return ms.ErrorClose
}

//Name give the name of the destination.
func (ms *MockSaver) Name() string {
	return ms.MockName
}

//Reset reinitialize the destination (if possible).
func (ms *MockSaver) Reset(log *logrus.Entry) error {
	return ms.ErrorReset
}
