package common

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
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
func OpenReader(config map[string]string) (io.ReadCloser, string, error) {
	file, okfile := config["file"]
	url, okurl := config["url"]
	inline, okinline := config["inline"]

	if !okfile && !okurl && !okinline {
		return nil, "", fmt.Errorf("the configuration block does not provide the filename or url or inline text")
	}

	var reader io.ReadCloser
	var err error

	if okfile && file == "-" {
		reader = NewStdinReaderCloser()
	} else if okfile {
		reader, err = os.Open(file)
		if err != nil {
			return nil, "", err
		}
	} else if okurl {
		resp, err := http.Get(url)
		if err != nil {
			return nil, "", err
		}

		reader = resp.Body
	} else if okinline {
		//string.NewReader returns a io.Reader, ioutil.NopCloser returns a io.ReadCloser with a Close implementation that do nothing
		reader = ioutil.NopCloser(strings.NewReader(inline))
	}

	gs, ok := config["gzip"]
	if ok {
		gb, err := strconv.ParseBool(gs)
		if err != nil {
			return nil, "", fmt.Errorf("the gzip element of configuration block must be true/false")
		}

		if gb {
			reader, err = gzip.NewReader(reader)
			if err != nil {
				return nil, "", err
			}

		}

	}

	return reader, file, nil
}
