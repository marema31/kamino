package file

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
)

// File managed file operation for kamino
type File struct {
	FilePath    string
	tmpFilePath string
	Gzip        bool
	Zip         bool
	fileHandle  io.Closer
	filewriter  bool
	URL         string
	Inline      string
	ZippedExt   string
}

//OpenReadFile open and return a io.ReadCloser used by datasource, providers and template destination
func (fi *File) OpenReadFile(log *logrus.Entry) (io.ReadCloser, error) {
	logFile := log.WithField("file", fi.FilePath)
	var reader io.ReadCloser
	var err error
	fi.tmpFilePath = ""

	logFile.Debugf("Opening %s in read mode", fi.FilePath)
	if fi.FilePath == "-" {
		reader = NewStdinReaderCloser()
	} else if fi.FilePath != "" {
		if fi.Zip {
			logFile.Debug("Opening as zip")
			archive, err := zip.OpenReader(fi.FilePath)
			if err != nil {
				logFile.Errorf("Opening zip file %s failed", fi.FilePath)
				logFile.Error(err)
				return nil, err
			}
			return archive.File[0].Open()
		}
		reader, err = os.Open(fi.FilePath)
		if err != nil {
			logFile.Errorf("Opening file %s failed", fi.FilePath)
			logFile.Error(err)
			return nil, err
		}
	} else if fi.URL != "" {
		resp, err := http.Get(fi.URL)
		if err != nil {
			logFile.Errorf("Opening URL %s failed", fi.URL)
			logFile.Error(err)
			return nil, err
		}

		reader = resp.Body
	} else if fi.Inline != "" {
		//string.NewReader returns a io.Reader, ioutil.NopCloser returns a io.ReadCloser with a Close implementation that do nothing
		reader = ioutil.NopCloser(strings.NewReader(fi.Inline))
	}

	if reader == nil {
		logFile.Error("Do not know which file to open for reading")
		return nil, fmt.Errorf("do not know which file to open for reading")
	}

	if fi.Gzip {
		logFile.Debug("Opening as gzip")
		reader, err = gzip.NewReader(reader)
		if err != nil {
			logFile.Error("Opening gzip file")
			logFile.Error(err)
			return nil, err
		}
	}

	fi.filewriter = false
	fi.fileHandle = reader
	return reader, nil
}

//OpenWriteFile open and return a io.WriteCloser corresponding to the datasource to be used by providers
func (fi *File) OpenWriteFile(log *logrus.Entry) (io.WriteCloser, error) {
	logFile := log.WithField("file", fi.FilePath)
	var writer io.WriteCloser

	logFile.Debugf("Opening %s in write mode", fi.FilePath)
	if fi.FilePath == "-" {
		writer = NewStdoutWriterCloser()

	} else {
		dir, pattern := filepath.Split(fi.FilePath)
		cache, err := ioutil.TempFile(dir, pattern+".")
		if err != nil {
			logFile.Error("Opening temporary file failed")
			logFile.Error(err)
			return nil, err
		}
		writer = cache
		fi.tmpFilePath = cache.Name()
		logFile.Debugf("Opened temporary file %s", fi.tmpFilePath)
	}

	if fi.Gzip {
		logFile.Debug("Opening as gzip")
		writer = gzip.NewWriter(writer)
	}

	fi.filewriter = true
	fi.fileHandle = writer
	return writer, nil
}

//ResetFile close the file and remove the temporary file
func (fi *File) ResetFile(log *logrus.Entry) error {
	logFile := log.WithField("file", fi.FilePath)
	logFile.Debugf("Resetting file by removing the temporary file %s", fi.tmpFilePath)

	if fi.fileHandle == nil {
		logFile.Debug("Skipping already closed")
		return nil
	}
	fi.Close()

	if fi.filewriter && fi.tmpFilePath != "" {
		err := os.Remove(fi.tmpFilePath)
		if err != nil {
			logFile.Error("Resetting file failed")
			logFile.Error(err)
			return err
		}
	}
	fi.tmpFilePath = ""
	return nil
}

//CloseFile close the file and rename the temporary file to real name (if exists)
func (fi *File) CloseFile(log *logrus.Entry) error {
	logFile := log.WithField("file", fi.FilePath)
	if fi.fileHandle == nil {
		logFile.Debug("Skipping already closed")
		return nil
	}
	fi.Close()

	if !fi.filewriter || fi.tmpFilePath == "" || fi.FilePath == "-" {
		logFile.Debugf("Closing %s", fi.FilePath)
		return nil // For file opened for read or stdin/stdout nothing more to do
	}
	logFile.Debugf("Closing temporary file %s", fi.tmpFilePath)

	logFile.Debugf("Removing destination file %s", fi.FilePath)
	if _, err := os.Stat(fi.FilePath); !os.IsNotExist(err) {
		err = os.Remove(fi.FilePath)
		if err != nil {
			logFile.Error("Removing file failed")
			logFile.Error(err)
			return err
		}
	}

	if fi.Zip {
		logFile.Debug("Creating zip archive")
		archivew, err := os.Create(fi.FilePath)
		if err != nil {
			logFile.Error("Creating zip file failed")
			logFile.Error(err)
			return err
		}

		archive := zip.NewWriter(archivew)

		reader, err := os.Open(fi.tmpFilePath)
		if err != nil {
			logFile.Errorf("Opening temporary file %s failed", fi.tmpFilePath)
			logFile.Error(err)
			return err
		}

		name := filepath.Base(strings.TrimSuffix(fi.FilePath, "zip"))
		writer, err := archive.Create(name + fi.ZippedExt)
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

		return os.Remove(fi.tmpFilePath)
	}

	logFile.Debugf("Renaming the temporary file %s to destination file %s", fi.tmpFilePath, fi.FilePath)
	err := os.Rename(fi.tmpFilePath, fi.FilePath)
	if err != nil {
		logFile.Errorf("Renaming temporary file %s to destination file %s", fi.tmpFilePath, fi.FilePath)
		logFile.Error(err)
		return err
	}
	fi.tmpFilePath = ""
	return nil
}

// Close the file if it is still opened
func (fi *File) Close() {
	if fi.fileHandle != nil {
		fi.fileHandle.Close()
		fi.fileHandle = nil
	}
}

// Stat returns os.FileInfo on the file of the datasource
func (fi *File) Stat() (os.FileInfo, error) {
	return os.Stat(fi.FilePath)
}
