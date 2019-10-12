package common

//TODO: This file is replaced by Datasource.OpenReadFile
import (
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/marema31/kamino/config"
)

// StdinReaderCloser type for Stdin with a non closing Close operation to avoid multiple close
type StdinReaderCloser struct {
	io.Reader
	io.Closer
}

// NewStdinReaderCloser return an new instance of StdinReaderCloser
func NewStdinReaderCloser() *StdinReaderCloser {
	var s StdinReaderCloser
	return &s
}

// Read read from Stdin
func (s *StdinReaderCloser) Read(p []byte) (n int, err error) {
	r := os.Stdin
	return r.Read(p)
}

// Close fake the stream close
func (s *StdinReaderCloser) Close() error {
	return nil
}

//OpenReader analyze the config block and return the corresponding io.ReadCloser to be used by other providers
func OpenReader(config config.SourceConfig) (io.ReadCloser, error) {
	if config.File == "" && config.URL == "" && config.Inline == "" {
		return nil, fmt.Errorf("the configuration block does not provide the filename or url or inline text")
	}

	var reader io.ReadCloser
	var err error

	if config.File == "-" {
		reader = NewStdinReaderCloser()
	} else if config.File != "" {
		if config.Zip {
			archive, err := zip.OpenReader(config.File)
			if err != nil {
				return nil, err
			}
			return archive.File[0].Open()
		}
		reader, err = os.Open(config.File)
		if err != nil {
			return nil, err
		}
	} else if config.URL != "" {
		resp, err := http.Get(config.URL)
		if err != nil {
			return nil, err
		}

		reader = resp.Body
	} else if config.Inline != "" {
		//string.NewReader returns a io.Reader, ioutil.NopCloser returns a io.ReadCloser with a Close implementation that do nothing
		reader = ioutil.NopCloser(strings.NewReader(config.Inline))
	}

	if config.Gzip {
		reader, err = gzip.NewReader(reader)
		if err != nil {
			return nil, err
		}
	}

	return reader, nil
}
