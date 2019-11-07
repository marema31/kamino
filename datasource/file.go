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

	"github.com/Sirupsen/logrus"
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
		return Datasource{}, fmt.Errorf("no file path or URL provided")
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
func (ds *Datasource) OpenReadFile(log *logrus.Entry) (io.ReadCloser, error) {
	logFile := log.WithField("engine", EngineToString(ds.engine))
	var reader io.ReadCloser
	var err error
	ds.tmpFilePath = ""

	logFile.Debugf("Opening %s in read mode", ds.filePath)
	if ds.filePath == "-" {
		reader = NewStdinReaderCloser()
	} else if ds.filePath != "" {
		if ds.zip {
			logFile.Debug("Opening as zip")
			archive, err := zip.OpenReader(ds.filePath)
			if err != nil {
				logFile.Errorf("Opening zip file %s failed", ds.filePath)
				logFile.Error(err)
				return nil, err
			}
			return archive.File[0].Open()
		}
		reader, err = os.Open(ds.filePath)
		if err != nil {
			logFile.Errorf("Opening file %s failed", ds.filePath)
			logFile.Error(err)
			return nil, err
		}
	} else if ds.url != "" {
		resp, err := http.Get(ds.url)
		if err != nil {
			logFile.Errorf("Opening URL %s failed", ds.url)
			logFile.Error(err)
			return nil, err
		}

		reader = resp.Body
	} else if ds.inline != "" {
		//string.NewReader returns a io.Reader, ioutil.NopCloser returns a io.ReadCloser with a Close implementation that do nothing
		reader = ioutil.NopCloser(strings.NewReader(ds.inline))
	}

	if reader == nil {
		logFile.Error("Do not know which file to open for reading")
		return nil, fmt.Errorf("do not known what file to open for datasource: %s", ds.name)
	}

	if ds.gzip {
		logFile.Debug("Opening as gzip")
		reader, err = gzip.NewReader(reader)
		if err != nil {
			logFile.Error("Opening gzip file")
			logFile.Error(err)
			return nil, err
		}
	}

	ds.filewriter = false
	ds.fileHandle = reader
	return reader, nil
}

//OpenWriteFile open and return a io.WriteCloser corresponding to the datasource to be used by providers
func (ds *Datasource) OpenWriteFile(log *logrus.Entry) (io.WriteCloser, error) {
	logFile := log.WithField("engine", EngineToString(ds.engine))
	var writer io.WriteCloser

	logFile.Debugf("Opening %s in write mode", ds.filePath)
	if ds.filePath == "-" {
		writer = NewStdoutWriterCloser()

	} else {
		dir, pattern := filepath.Split(ds.filePath)
		cache, err := ioutil.TempFile(dir, pattern+".")
		if err != nil {
			logFile.Error("Opening temporary file failed")
			logFile.Error(err)
			return nil, err
		}
		writer = cache
		ds.tmpFilePath = cache.Name()
		logFile.Debugf("Opened temporary file %s", ds.tmpFilePath)
	}

	if ds.gzip {
		logFile.Debug("Opening as gzip")
		writer = gzip.NewWriter(writer)
	}

	ds.filewriter = true
	ds.fileHandle = writer
	return writer, nil
}

//ResetFile close the file and remove the temporary file
func (ds *Datasource) ResetFile(log *logrus.Entry) error {
	logFile := log.WithField("engine", EngineToString(ds.engine))
	logFile.Debugf("Resetting file by removing the temporary file %s", ds.tmpFilePath)
	ds.fileHandle.Close()
	if ds.filewriter && ds.tmpFilePath != "" {
		err := os.Remove(ds.tmpFilePath)
		if err != nil {
			logFile.Error("Resetting file failed")
			logFile.Error(err)
			return err
		}
	}
	ds.tmpFilePath = ""
	return nil
}

//CloseFile close the file and rename the temporary file to real name (if exists)
func (ds *Datasource) CloseFile(log *logrus.Entry) error {
	logFile := log.WithField("engine", EngineToString(ds.engine))

	ds.fileHandle.Close()

	if !ds.filewriter || ds.tmpFilePath == "" || ds.filePath == "-" {
		logFile.Debugf("Closing %s", ds.filePath)
		return nil // For file opened for read or stdin/stdout nothing more to do
	}
	logFile.Debugf("Closing temporary file %s", ds.tmpFilePath)

	logFile.Debugf("Removing destination file %s", ds.filePath)
	if _, err := os.Stat(ds.filePath); !os.IsNotExist(err) {
		err = os.Remove(ds.filePath)
		if err != nil {
			logFile.Error("Removing file failed")
			logFile.Error(err)
			return err
		}
	}

	if ds.zip {
		logFile.Debug("Creating zip archive")
		archivew, err := os.Create(ds.filePath)
		if err != nil {
			logFile.Error("Creating zip file failed")
			logFile.Error(err)
			return err
		}

		archive := zip.NewWriter(archivew)

		reader, err := os.Open(ds.tmpFilePath)
		if err != nil {
			logFile.Errorf("Opening temporary file %s failed", ds.tmpFilePath)
			logFile.Error(err)
			return err
		}

		name := filepath.Base(strings.TrimSuffix(ds.filePath, "zip"))
		writer, err := archive.Create(name + EngineToString(ds.engine))
		if err != nil {
			logFile.Error("Creating archive entry zip file failed")
			logFile.Error(err)
			return err
		}

		_, err = io.Copy(writer, reader)
		if err != nil {
			logFile.Error("Compression failed")
			logFile.Error(err)
			return err
		}

		err = archive.Close()
		if err != nil {
			logFile.Error("Closing zip file failed")
			logFile.Error(err)
			return err
		}

		return os.Remove(ds.tmpFilePath)
	}

	logFile.Debugf("Renaming the temporary file %s to destination file %s", ds.tmpFilePath, ds.filePath)
	err := os.Rename(ds.tmpFilePath, ds.filePath)
	if err != nil {
		logFile.Errorf("Renaming temporary file %s to destination file %s", ds.tmpFilePath, ds.filePath)
		logFile.Error(err)
		return err
	}
	ds.tmpFilePath = ""
	return nil
}
