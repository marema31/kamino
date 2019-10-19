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

	_ "github.com/go-sql-driver/mysql" // Mysql library dynamically called by database/sql
	_ "github.com/lib/pq"              //Postgres library dynamically called by database/sql
	"github.com/spf13/viper"
)

// load a dile type datasource from the viper configuration
func loadFileDatasource(filename string, v *viper.Viper, engine Engine) (*Datasource, error) {
	var ds Datasource
	ds.Type = File
	ds.Engine = engine
	ds.Name = filename
	ds.Inline = v.GetString("inline")
	ds.FilePath = v.GetString("file")
	ds.URL = v.GetString("URL")
	if ds.FilePath == "" && ds.URL == "" && ds.Inline == "" {
		return nil, fmt.Errorf("the datasource %s does not provide the file path or URL", ds.Name)
	}

	ds.Tags = v.GetStringSlice("tags")
	if len(ds.Tags) == 0 {
		ds.Tags = []string{""}
	}

	ds.Zip = v.GetBool("zip")
	ds.Gzip = v.GetBool("gzip")
	return &ds, nil
}

//OpenReadFile open and return a io.ReadCloser corresponding to the datasource to be used by providers
func (ds *Datasource) OpenReadFile() (io.ReadCloser, error) {
	var reader io.ReadCloser
	var err error

	if ds.FilePath == "-" {
		reader = NewStdinReaderCloser()
	} else if ds.FilePath != "" {
		if ds.Zip {
			archive, err := zip.OpenReader(ds.FilePath)
			if err != nil {
				return nil, err
			}
			return archive.File[0].Open()
		}
		reader, err = os.Open(ds.FilePath)
		if err != nil {
			return nil, err
		}
	} else if ds.URL != "" {
		resp, err := http.Get(ds.URL)
		if err != nil {
			return nil, err
		}

		reader = resp.Body
	} else if ds.Inline != "" {
		//string.NewReader returns a io.Reader, ioutil.NopCloser returns a io.ReadCloser with a Close implementation that do nothing
		reader = ioutil.NopCloser(strings.NewReader(ds.Inline))
	}

	if ds.Gzip {
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

	if ds.FilePath == "-" {
		writer = NewStdoutWriterCloser()

	} else {
		dir, pattern := filepath.Split(ds.FilePath)
		cache, err := ioutil.TempFile(dir, pattern+".")
		if err != nil {
			return nil, err
		}
		writer = cache
		ds.tmpFilePath = cache.Name()
	}

	if ds.Gzip {
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
	return nil
}

//CloseFile close the file and rename the temporary file to real name (if exists)
func (ds *Datasource) CloseFile() error {

	ds.fileHandle.Close()

	if !ds.filewriter || ds.tmpFilePath == "" || ds.FilePath == "-" {
		return nil // For file opened for read or stdin/stdout nothing more to do
	}

	if _, err := os.Stat(ds.FilePath); !os.IsNotExist(err) {
		err = os.Remove(ds.FilePath)
		if err != nil {
			return err
		}

	}

	if ds.Zip {
		archivew, err := os.Create(ds.FilePath)
		if err != nil {
			return err
		}

		archive := zip.NewWriter(archivew)

		reader, err := os.Open(ds.tmpFilePath)
		if err != nil {
			return err
		}

		name := filepath.Base(strings.TrimSuffix(ds.FilePath, "zip"))
		writer, err := archive.Create(name + ds.GetEngine())
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
	err := os.Rename(ds.tmpFilePath, ds.FilePath)
	if err != nil {
		return err
	}
	return nil
}
