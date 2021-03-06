package mockdatasource_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/mockdatasource"
)

func TestOpenFile(t *testing.T) {
	ds := mockdatasource.MockDatasource{Type: datasource.File, Zip: false, Gzip: false, FilePath: "testdata/tmp/testfile"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")

	testString := []byte("test,string,for,writing\n  test_string: - for writing\n{'test':[ 'string','for','writing']}")

	writer, err := ds.OpenWriteFile(log)
	if err != nil {
		t.Fatalf("OpenWriteFile Should not return error and returned '%v'", err)
	}

	_, err = writer.Write(testString)
	if err != nil {
		t.Fatalf("Write Should not return error and returned '%v'", err)
	}

	err = ds.CloseFile(log)
	if err != nil {
		t.Errorf("CloseFile Should not return error and returned '%v'", err)
	}

	fi, err := ds.Stat()
	if err != nil {
		t.Errorf("Stat Should not return error and returned '%v'", err)
	}
	if !fi.Mode().IsDir() { //Since we use in-memory in the mock, we mock this with the current directory
		t.Fatalf("Should be a directory")
	}

	reader, err := ds.OpenReadFile(log)
	if err != nil {
		t.Fatalf("OpenReadFile Should not return error and returned '%v'", err)
	}

	readString := make([]byte, len(testString))
	_, err = reader.Read(readString)
	if err != nil {
		t.Fatalf("Read Should not return error and returned '%v'", err)
	}

	if !bytes.Equal(testString, readString) {
		t.Errorf("The read string is not equal to written one: '%s' != '%s'  ", string(testString), string(readString))
	}
	err = ds.CloseFile(log)
	if err != nil {
		t.Errorf("Close Should not return error and returned '%v'", err)
	}
}
func TestResetFile(t *testing.T) {
	ds := mockdatasource.MockDatasource{Type: datasource.File, Zip: false, Gzip: false, FilePath: "testdata/tmp/testfile"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")

	_, err := ds.OpenWriteFile(log)
	if err != nil {
		t.Fatalf("OpenWriteFile Should not return error and returned '%v'", err)
	}

	err = ds.ResetFile(log)
	if err != nil {
		t.Errorf("ResetFile Should not return error and returned '%v'", err)
	}
}

func TestErrorFile(t *testing.T) {
	ds := mockdatasource.MockDatasource{Type: datasource.File, Zip: false, Gzip: false, FilePath: "testdata/tmp/testfile"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")

	ds.ErrorClose = fmt.Errorf("Fake error")
	ds.ErrorReset = fmt.Errorf("Fake error")
	_, err := ds.OpenReadFile(log)
	if err != nil {
		t.Fatalf("OpenReadFile Should not return error and returned '%v'", err)
	}

	_, err = ds.OpenWriteFile(log)
	if err != nil {
		t.Fatalf("OpenWriteFile Should not return error and returned '%v'", err)
	}

	err = ds.CloseFile(log)
	if err == nil {
		t.Errorf("CloseFile Should return error")
	}

	err = ds.ResetFile(log)
	if err == nil {
		t.Errorf("ResetFile Should return error")
	}

	ds.ErrorOpenFile = fmt.Errorf("Fake error")
	_, err = ds.OpenWriteFile(log)
	if err == nil {
		t.Fatalf("OpenWriteFile Should return error ")
	}
	_, err = ds.OpenReadFile(log)
	if err == nil {
		t.Fatalf("OpenReadFile Should return error ")
	}

	ds.ErrorClose = nil
	err = ds.CloseFile(log)
	if err != nil {
		t.Errorf("Close Should not return error and returned '%v'", err)
	}

	ds.ErrorClose = fmt.Errorf("Fake error")
	ds.ErrorReset = nil
	err = ds.ResetFile(log)
	if err != nil {
		t.Errorf("ResetFile Should not return error and returned '%v'", err)
	}
}

func TestStat(t *testing.T) {
	ds := mockdatasource.MockDatasource{Type: datasource.File, Zip: false, Gzip: false, FilePath: "testdata/tmp/testfile"}

	fi, err := ds.Stat()
	if err != nil {
		t.Errorf("Stat Should not return error and returned '%v'", err)
	}
	if !fi.Mode().IsDir() { //Since we use in-memory in the mock, we mock this with the current directory
		t.Fatalf("Should be a directory")
	}

	ds.FileNotExists = true
	_, err = ds.Stat()
	if err == nil {
		t.Errorf("Stat Should return error")
	}
	if os.IsExist(err) { //Since we use in-memory in the mock, we mock this with inexistant directory
		t.Fatalf("Should not exists")
	}

}
