package file

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
)

func TestOpenStdio(t *testing.T) {
	fi := File{Zip: false, Gzip: false, FilePath: "-"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")

	_, err := fi.OpenWriteFile(log)
	if err != nil {
		t.Fatalf("Should not return error and returned '%v'", err)
	}
	err = fi.CloseFile(log)
	if err != nil {
		t.Errorf("Should not return error and returned '%v'", err)
	}

	_, err = fi.OpenReadFile(log)
	if err != nil {
		t.Fatalf("Should not return error and returned '%v'", err)
	}
	err = fi.CloseFile(log)
	if err != nil {
		t.Errorf("Should not return error and returned '%v'", err)
	}
}

func TestOpenInline(t *testing.T) {
	fi := File{Zip: false, Gzip: false, Inline: "testInline"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")

	reader, err := fi.OpenReadFile(log)
	if err != nil {
		t.Fatalf("OpenReadFile Should not return error and returned '%v'", err)
	}

	test, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll Should not return error and returned '%v'", err)
	}
	if string(test) != "testInline" {
		t.Errorf("The content of Inline is not the one we waits for :%v", test)
	}
	err = fi.CloseFile(log)
	if err != nil {
		t.Errorf("Close Should not return error and returned '%v'", err)
	}
}

func TestOpenFile(t *testing.T) {
	if _, err := os.Stat("testdata/tmp"); os.IsNotExist(err) {
		os.Mkdir("testdata/tmp", 0777)
	}

	test := []byte{1, 2, 3}

	fi := File{Zip: false, Gzip: false, FilePath: "testdata/tmp/testfile"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")

	_, err := fi.OpenWriteFile(log)
	if err != nil {
		t.Errorf("OpenWriteFile Should not return error and returned '%v'", err)
	}
	//	writer.Write(test)

	err = fi.CloseFile(log)
	if err != nil {
		t.Errorf("CloseFile Should not return error and returned '%v'", err)
	}

	fileInfo, err := fi.Stat()
	if err != nil {
		t.Errorf("Stat Should not return error and returned '%v'", err)
	}
	if !fileInfo.Mode().IsRegular() {
		t.Fatalf("Should be a file")
	}

	reader, err := fi.OpenReadFile(log)
	if err != nil {
		t.Fatalf("OpenReadFile Should not return error and returned '%v'", err)
	}
	reader.Read(test)
	if test[2] != 3 {
		t.Errorf("The content of file is not the one we waits for :%v", test)
	}
	err = fi.CloseFile(log)
	if err != nil {
		t.Errorf("Close Should not return error and returned '%v'", err)
	}
	err = fi.CloseFile(log)
	if err != nil {
		t.Errorf("Close Should not return error and returned '%v'", err)
	}

	os.Remove("testdata/tmp/testfile")
}

func TestOpenZipFile(t *testing.T) {
	if _, err := os.Stat("testdata/tmp"); os.IsNotExist(err) {
		os.Mkdir("testdata/tmp", 0777)
	}

	test := []byte{1, 2, 3}

	fi := File{Zip: true, Gzip: false, FilePath: "testdata/tmp/testzip"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")

	writer, err := fi.OpenWriteFile(log)
	if err != nil {
		t.Fatalf("OpenWriteFile Should not return error and returned '%v'", err)
	}
	writer.Write(test)

	err = fi.CloseFile(log)
	if err != nil {
		t.Errorf("CloseFile Should not return error and returned '%v'", err)
	}

	reader, err := fi.OpenReadFile(log)
	if err != nil {
		t.Fatalf("OpenReadFile Should not return error and returned '%v'", err)
	}
	reader.Read(test)
	if test[2] != 3 {
		t.Errorf("The content of file is not the one we waits for :%v", test)
	}
	err = fi.CloseFile(log)
	if err != nil {
		t.Errorf("Close Should not return error and returned '%v'", err)
	}

	os.Remove("testdata/tmp/testzip")
}

func TestOpenGzipFile(t *testing.T) {
	if _, err := os.Stat("testdata/tmp"); os.IsNotExist(err) {
		os.Mkdir("testdata/tmp", 0777)
	}

	test := []byte{1, 2, 3}

	fi := File{Zip: false, Gzip: true, FilePath: "testdata/tmp/testgz"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")

	writer, err := fi.OpenWriteFile(log)
	if err != nil {
		t.Fatalf("OpenWriteFile Should not return error and returned '%v'", err)
	}
	writer.Write(test)

	err = fi.CloseFile(log)
	if err != nil {
		t.Errorf("CloseFile Should not return error and returned '%v'", err)
	}

	reader, err := fi.OpenReadFile(log)
	if err != nil {
		t.Fatalf("OpenReadFile Should not return error and returned '%v'", err)
	}
	reader.Read(test)
	if test[2] != 3 {
		t.Errorf("The content of file is not the one we waits for :%v", test)
	}
	err = fi.CloseFile(log)
	if err != nil {
		t.Errorf("Close Should not return error and returned '%v'", err)
	}

	os.Remove("testdata/tmp/testgz")
}

func TestReadWronGzip(t *testing.T) {
	if _, err := os.Stat("testdata/tmp"); os.IsNotExist(err) {
		os.Mkdir("testdata/tmp", 0777)
	}

	test := []byte{1, 2, 3}

	fi := File{Zip: false, Gzip: false, FilePath: "testdata/tmp/testfile"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")

	writer, err := fi.OpenWriteFile(log)
	if err != nil {
		t.Fatalf("OpenWriteFile Should not return error and returned '%v'", err)
	}
	writer.Write(test)

	err = fi.CloseFile(log)
	if err != nil {
		t.Errorf("CloseFile Should not return error and returned '%v'", err)
	}

	fi.Zip = true
	_, err = fi.OpenReadFile(log)
	if err == nil {
		t.Errorf("OpenReadFile Should return error")
	}

	os.Remove("testdata/tmp/testfile")
}

func TestReadWrongGzip(t *testing.T) {
	if _, err := os.Stat("testdata/tmp"); os.IsNotExist(err) {
		os.Mkdir("testdata/tmp", 0777)
	}

	test := []byte{1, 2, 3}

	fi := File{Zip: false, Gzip: false, FilePath: "testdata/tmp/testfile"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")

	writer, err := fi.OpenWriteFile(log)
	if err != nil {
		t.Fatalf("OpenWriteFile Should not return error and returned '%v'", err)
	}
	writer.Write(test)

	err = fi.CloseFile(log)
	if err != nil {
		t.Errorf("CloseFile Should not return error and returned '%v'", err)
	}

	fi.Gzip = true
	_, err = fi.OpenReadFile(log)
	if err == nil {
		t.Errorf("OpenReadFile Should return error")
	}

	os.Remove("testdata/tmp/testfile")
}

func TestResetFile(t *testing.T) {
	if _, err := os.Stat("testdata/tmp"); os.IsNotExist(err) {
		os.Mkdir("testdata/tmp", 0777)
	}

	test := []byte{1, 2, 3}

	fi := File{Zip: false, Gzip: false, FilePath: "testdata/tmp/testfile"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")

	writer, err := fi.OpenWriteFile(log)
	if err != nil {
		t.Fatalf("OpenWriteFile Should not return error and returned '%v'", err)
	}
	writer.Write(test)

	err = fi.ResetFile(log)
	if err != nil {
		t.Errorf("ResetFile Should not return error and returned '%v'", err)
	}

	if _, err := os.Stat(fi.tmpFilePath); os.IsExist(err) {
		os.Remove(fi.tmpFilePath)
		t.Errorf("ResetFile Should have removed the temporary file'")
	}

	err = fi.ResetFile(log)
	if err != nil {
		t.Errorf("ResetFile Should not return error and returned '%v'", err)
	}
}

func TestOpenFileError(t *testing.T) {
	if _, err := os.Stat("testdata/tmp"); os.IsNotExist(err) {
		os.Mkdir("testdata/tmp", 0777)
	}

	fi := File{Zip: false, Gzip: false, FilePath: "testdata/tmp/nofile"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	logger.SetLevel(logrus.PanicLevel)

	_, err := fi.OpenReadFile(log)
	if err == nil {
		t.Fatalf("OpenReadFile Should return an error")
	}
}

func TestOpenUrlError(t *testing.T) {

	fi := File{Zip: false, Gzip: false, URL: "http://1.2.3.4.5"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	logger.SetLevel(logrus.PanicLevel)

	_, err := fi.OpenReadFile(log)
	if err == nil {
		t.Fatalf("OpenReadFile Should return an error")
	}
}

func TestOpenNoFileNoUrlError(t *testing.T) {

	fi := File{Zip: false, Gzip: false}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	logger.SetLevel(logrus.PanicLevel)

	_, err := fi.OpenReadFile(log)
	if err == nil {
		t.Fatalf("OpenReadFile Should return an error")
	}
}

func TestOpenTmpFileError(t *testing.T) {
	if _, err := os.Stat("testdata/tmp"); os.IsNotExist(err) {
		os.Mkdir("testdata/tmp", 0777)
	}

	fi := File{Zip: false, Gzip: false, FilePath: "testdata/tmp/nodir/nofile"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	logger.SetLevel(logrus.PanicLevel)

	_, err := fi.OpenWriteFile(log)
	if err == nil {
		t.Fatalf("OpenWriteFile Should return an error")
	}
}

func TestResetTmpFileError(t *testing.T) {
	if _, err := os.Stat("testdata/tmp"); os.IsNotExist(err) {
		os.Mkdir("testdata/tmp", 0777)
	}

	fi := File{Zip: false, Gzip: false, tmpFilePath: "testdata/tmp/nodir/nofile"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	logger.SetLevel(logrus.PanicLevel)

	tmpFile, err := ioutil.TempFile("testdata/tmp", "reset.")
	if err != nil {
		t.Fatalf("Should not return error and returned '%v'", err)
	}
	fi.fileHandle = tmpFile
	fi.filewriter = true
	err = fi.ResetFile(log)
	if err == nil {
		t.Fatalf("ResetFile Should return an error")
	}
}

func TestCloseFileError(t *testing.T) {
	if _, err := os.Stat("testdata/tmp"); os.IsNotExist(err) {
		os.Mkdir("testdata/tmp", 0777)
	}

	fi := File{Zip: false, Gzip: false, tmpFilePath: "testdata/tmp/nodir/nofile"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	logger.SetLevel(logrus.PanicLevel)

	tmpFile, err := ioutil.TempFile("testdata/tmp", "reset.")
	if err != nil {
		t.Fatalf("Should not return error and returned '%v'", err)
	}
	fi.fileHandle = tmpFile
	fi.filewriter = true
	err = fi.CloseFile(log)
	if err == nil {
		t.Fatalf("ResetFile Should return an error")
	}
}

func TestCloseFileZipError(t *testing.T) {
	if _, err := os.Stat("testdata/tmp"); os.IsNotExist(err) {
		os.Mkdir("testdata/tmp", 0777)
	}

	fi := File{Zip: true, Gzip: false, tmpFilePath: "testdata/tmp/nodir/nofile"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	logger.SetLevel(logrus.PanicLevel)

	tmpFile, err := ioutil.TempFile("testdata/tmp", "reset.")
	if err != nil {
		t.Fatalf("Should not return error and returned '%v'", err)
	}
	fi.fileHandle = tmpFile
	fi.filewriter = true
	err = fi.CloseFile(log)
	if err == nil {
		t.Fatalf("ResetFile Should return an error")
	}
}

func TestCloseFileZipNoDataError(t *testing.T) {
	if _, err := os.Stat("testdata/tmp"); os.IsNotExist(err) {
		os.Mkdir("testdata/tmp", 0777)
	}

	fi := File{Zip: true, Gzip: false, FilePath: "testdata/tmp/nodata.zip", tmpFilePath: "testdata/tmp/nodir/nofile"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	logger.SetLevel(logrus.PanicLevel)

	tmpFile, err := ioutil.TempFile("testdata/tmp", "reset.")
	if err != nil {
		t.Fatalf("Should not return error and returned '%v'", err)
	}
	fi.fileHandle = tmpFile
	fi.filewriter = true
	err = fi.CloseFile(log)
	if err == nil {
		t.Fatalf("ResetFile Should return an error")
	}
}

func TestStat(t *testing.T) {
	if _, err := os.Stat("testdata/tmp"); os.IsNotExist(err) {
		os.Mkdir("testdata/tmp", 0777)
	}

	fi := File{Zip: true, Gzip: false, FilePath: "testdata/tmp", tmpFilePath: ""}

	if _, err := fi.Stat(); os.IsNotExist(err) {
		t.Errorf("Stat Should have seen the file'")
	}
}
