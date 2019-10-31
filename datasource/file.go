package datasource

import (
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// load a dile type datasource from the viper configuration
func loadFileDatasource(recipePath string, filename string, v *viper.Viper, engine Engine) (Datasource, error) {
	var ds Datasource
	ds.dstype = File
	ds.engine = engine
	ds.name = filename
	ds.inline = v.GetString("inline")
	ds.filePath = v.GetString("file")
	if ds.filePath != "" && ds.filePath != "-" && ds.filePath[0] != '/' {
		ds.filePath = filepath.Join(recipePath, ds.filePath)
	}

	ds.url = v.GetString("URL")
	if ds.filePath == "" && ds.url == "" && ds.inline == "" {
		return Datasource{}, fmt.Errorf("the datasource %s does not provide the file path or URL", ds.name)
	}

	ds.tags = v.GetStringSlice("tags")
	if len(ds.tags) == 0 {
		ds.tags = []string{""}
	}

	ds.zip = v.GetBool("zip")
	ds.gzip = v.GetBool("gzip")
	return ds, nil
}

//OpenReadFile open and return a io.ReadCloser corresponding to the datasource to be used by providers
func (ds *Datasource) OpenReadFile() (io.ReadCloser, error) {
	var reader io.ReadCloser
	var err error
	ds.tmpFilePath = ""

	if ds.filePath == "-" {
		reader = NewStdinReaderCloser()
	} else if ds.filePath != "" {
		if ds.zip {
			archive, err := zip.OpenReader(ds.filePath)
			if err != nil {
				return nil, err
			}
			return archive.File[0].Open()
		}
		reader, err = os.Open(ds.filePath)
		if err != nil {
			return nil, err
		}
	} else if ds.url != "" {
		resp, err := http.Get(ds.url)
		if err != nil {
			return nil, err
		}

		reader = resp.Body
	} else if ds.inline != "" {
		//string.NewReader returns a io.Reader, ioutil.NopCloser returns a io.ReadCloser with a Close implementation that do nothing
		reader = ioutil.NopCloser(strings.NewReader(ds.inline))
	}

	if reader == nil {
		return nil, fmt.Errorf("do not known what file to open for datasource: %s", ds.name)
	}

	if ds.gzip {
		reader, err = gzip.NewReader(reader)
		if err != nil {
			return nil, err
		}
	}

	ds.filewriter = false
	ds.fileHandle = reader
	return reader, nil
}

//OpenWriteFile open and return a io.WriteCloser corresponding to the datasource to be used by providers
func (ds *Datasource) OpenWriteFile() (io.WriteCloser, error) {
	var writer io.WriteCloser

	if ds.filePath == "-" {
		writer = NewStdoutWriterCloser()

	} else {
		dir, pattern := filepath.Split(ds.filePath)
		cache, err := ioutil.TempFile(dir, pattern+".")
		if err != nil {
			return nil, err
		}
		writer = cache
		ds.tmpFilePath = cache.Name()
	}

	if writer == nil {
		return nil, fmt.Errorf("do not known what file to open for datasource: %s", ds.name)
	}

	if ds.gzip {
		writer = gzip.NewWriter(writer)
	}

	ds.filewriter = true
	ds.fileHandle = writer
	return writer, nil
}

//ResetFile close the file and remove the temporary file
func (ds *Datasource) ResetFile() error {
	ds.fileHandle.Close()
	if ds.filewriter && ds.tmpFilePath != "" {
		err := os.Remove(ds.tmpFilePath)
		if err != nil {
			return err
		}
	}
	ds.tmpFilePath = ""
	return nil
}

//CloseFile close the file and rename the temporary file to real name (if exists)
func (ds *Datasource) CloseFile() error {

	ds.fileHandle.Close()

	if !ds.filewriter || ds.tmpFilePath == "" || ds.filePath == "-" {
		return nil // For file opened for read or stdin/stdout nothing more to do
	}

	if _, err := os.Stat(ds.filePath); !os.IsNotExist(err) {
		err = os.Remove(ds.filePath)
		if err != nil {
			return err
		}

	}

	if ds.zip {
		archivew, err := os.Create(ds.filePath)
		if err != nil {
			return err
		}

		archive := zip.NewWriter(archivew)

		reader, err := os.Open(ds.tmpFilePath)
		if err != nil {
			return err
		}

		name := filepath.Base(strings.TrimSuffix(ds.filePath, "zip"))
		writer, err := archive.Create(name + EngineToString(ds.engine))
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, reader)
		if err != nil {
			return err
		}

		err = archive.Close()
		if err != nil {
			return err
		}

		return os.Remove(ds.tmpFilePath)
	}

	//rename the tempfile for cache to its real name
	err := os.Rename(ds.tmpFilePath, ds.filePath)
	if err != nil {
		return err
	}
	ds.tmpFilePath = ""
	return nil
}
