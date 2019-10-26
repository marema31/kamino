package mockprovider

import (
	"github.com/marema31/kamino/provider/types"
)

//MockSaver specifc state for database Saver provider
type MockSaver struct {
	MockName   string
	Content    []map[string]string
	ErrorClose error
	ErrorReset error
}

//Save writes the record to the destination
func (ms *MockSaver) Save(record types.Record) error {
	ms.Content = append(ms.Content, record)
	return nil
}

//Close closes the destination
func (ms *MockSaver) Close() error {
	return ms.ErrorClose
}

//Name give the name of the destination
func (ms *MockSaver) Name() string {
	return ms.MockName
}

//Reset reinitialize the destination (if possible)
func (ms *MockSaver) Reset() error {
	return ms.ErrorReset
}
