package common

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
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
func OpenWriter(config map[string]string) (io.WriteCloser, string, string, error) {
	file, okfile := config["file"]
	if !okfile {
		return nil, "", "", fmt.Errorf("the configuration block does not provide the filename")
	}

	var writer io.WriteCloser
	var tmpFile = file

	if okfile && file == "-" {
		writer = NewStdoutWriterCloser()

	} else if okfile {
		dir, pattern := filepath.Split(file)
		cache, err := ioutil.TempFile(dir, pattern)
		if err != nil {
			return nil, "", "", err
		}
		writer = cache
		tmpFile = cache.Name()
	}

	gs, ok := config["gzip"]
	if ok {
		gb, err := strconv.ParseBool(gs)
		if err != nil {
			return nil, "", "", fmt.Errorf("the gzip element of configuration block must be true/false")
		}

		if gb {
			writer = gzip.NewWriter(writer)
		}

	}

	return writer, file, tmpFile, nil
}

//ResetWriter close the writer and remove the temporary file since the synchronization was not OK
func ResetWriter(writer io.Closer, tmpFile string, file string) error {
	writer.Close()

	err := os.Remove(tmpFile)
	if err != nil {
		return err
	}
	return nil
}

//CloseWriter close the writer and rename the temporary file to real name since synchronization was OK
func CloseWriter(writer io.Closer, tmpFile string, file string) error {

	writer.Close()

	if file == "-" {
		return nil // For Stdout nothing more to do
	}
	//rename the tempfile for cache to its real name
	if _, err := os.Stat(file); !os.IsNotExist(err) {
		err = os.Remove(file)
		if err != nil {
			return err
		}

	}
	err := os.Rename(tmpFile, file)
	if err != nil {
		return err
	}
	return nil
}
