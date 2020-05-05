package file

import (
	"io"
	"os"
)

// StdinReaderCloser type for Stdin with a non closing Close operation to avoid multiple close.
type StdinReaderCloser struct {
	io.Reader
	io.Closer
}

// NewStdinReaderCloser return an new instance of StdinReaderCloser.
func NewStdinReaderCloser() *StdinReaderCloser {
	var s StdinReaderCloser
	return &s
}

// Read read from Stdin.
func (s *StdinReaderCloser) Read(p []byte) (n int, err error) {
	r := os.Stdin
	return r.Read(p)
}

// Close fake the stream close.
func (s *StdinReaderCloser) Close() error {
	return nil
}

// StdoutWriterCloser type for Stdout with a non closing Close operation to avoid multiple close.
type StdoutWriterCloser struct {
	io.Writer
	io.Closer
}

// NewStdoutWriterCloser return an new instance of StdoutWriterCloser.
func NewStdoutWriterCloser() *StdoutWriterCloser {
	s := StdoutWriterCloser{}
	return &s
}

// Write write to Stdout.
func (s *StdoutWriterCloser) Write(p []byte) (n int, err error) {
	r := os.Stdout
	return r.Write(p)
}

// Close fake the stream close.
func (s *StdoutWriterCloser) Close() error {
	return nil
}
