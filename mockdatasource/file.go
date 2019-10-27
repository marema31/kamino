package mockdatasource

import (
	"io"
	"io/ioutil"
)

type nopWriteCloser struct {
	io.Writer
}

func (w *nopWriteCloser) Close() error {
	return nil
}

// NopWriteCloser returns a nopWriteCloser.
func NopWriteCloser(w io.Writer) io.WriteCloser {
	return &nopWriteCloser{w}
}

//OpenReadFile open and return a io.ReadCloser on the WriteBuf element
func (ds *MockDatasource) OpenReadFile() (io.ReadCloser, error) {
	//string.NewReader returns a io.Reader, ioutil.NopCloser returns a io.ReadCloser with a Close implementation that do nothing
	reader := ioutil.NopCloser(&ds.WriteBuf)

	ds.Filewriter = false
	ds.FileHandle = reader
	return reader, ds.ErrorOpenFile
}

//OpenWriteFile open and return a io.WriteCloser on the WriteBuf element
func (ds *MockDatasource) OpenWriteFile() (io.WriteCloser, error) {
	writer := NopWriteCloser(&ds.WriteBuf)
	ds.Filewriter = true
	ds.FileHandle = writer
	return writer, ds.ErrorOpenFile
}

//ResetFile close the file and remove the temporary file
func (ds *MockDatasource) ResetFile() error {
	return ds.ErrorReset
}

//CloseFile close the file and rename the temporary file to real name (if exists)
func (ds *MockDatasource) CloseFile() error {
	return ds.ErrorClose
}
