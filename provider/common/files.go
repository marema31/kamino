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

//OpenReader analyze the config block and return the corresponding io.ReadCloser to be used by other providers
func OpenReader(config map[string]string) (io.ReadCloser, error) {
	_, okstd := config["std"]
	file, okfile := config["file"]
	url, okurl := config["url"]
	inline, okinline := config["inline"]

	if !okfile && !okurl && !okinline && !okstd {
		return nil, fmt.Errorf("the configuration block does not provide the filename or url or inline text")
	}

	var reader io.ReadCloser
	var err error

	if okstd {
		reader = os.Stdin
	} else if okfile {
		reader, err = os.Open(file)
		if err != nil {
			return nil, err
		}
	} else if okurl {
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
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
			return nil, fmt.Errorf("the gzip element of configuration block must be true/false")
		}

		if gb {
			reader, err = gzip.NewReader(reader)
			if err != nil {
				return nil, err
			}

		}

	}

	return reader, nil
}

//OpenWriter analyze the config block and return the corresponding io.WriteCloser to be used by other providers
func OpenWriter(config map[string]string) (io.WriteCloser, error) {
	_, okstd := config["std"]
	file, okfile := config["file"]
	if !okfile && !okstd {
		return nil, fmt.Errorf("the configuration block does not provide the filename or url or inline text")
	}

	var writer io.WriteCloser
	var err error

	if okstd {
		writer = os.Stdout
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
