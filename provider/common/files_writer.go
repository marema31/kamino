package common

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strconv"
)

// StdoutWriterCloser type for Stdout with a non closing Close operation to avoid multiple close
type StdoutWriterCloser struct {
	io.Writer
	io.Closer
}

// NewStdoutWriterCloser return an new instance of StdoutWriterCloser
func NewStdoutWriterCloser() *StdoutWriterCloser {
	var s StdoutWriterCloser
	return &s
}

// Write write to Stdout
func (s *StdoutWriterCloser) Write(p []byte) (n int, err error) {
	r := os.Stdout
	return r.Write(p)
}

// Close fake the stream close
func (s *StdoutWriterCloser) Close() error {
	return nil
}

//OpenWriter analyze the config block and return the corresponding io.WriteCloser to be used by other providers
func OpenWriter(config map[string]string) (io.WriteCloser, error) {
	file, okfile := config["file"]
	if !okfile {
		return nil, fmt.Errorf("the configuration block does not provide the filename")
	}

	var writer io.WriteCloser
	var err error

	if okfile && file == "-" {
		writer = NewStdoutWriterCloser()
	} else if okfile {
		writer, err = os.Create(file)
		if err != nil {
			return nil, err
		}
	}

	gs, ok := config["gzip"]
	if ok {
		gb, err := strconv.ParseBool(gs)
		if err != nil {
			return nil, fmt.Errorf("the gzip element of configuration block must be true/false")
		}

		if gb {
			writer = gzip.NewWriter(writer)
		}

	}

	return writer, nil
}
