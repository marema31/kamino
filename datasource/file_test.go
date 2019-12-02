package datasource

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
)

// We are using private function, we must be in same package
func setupFileTest() *Datasources {
	return &Datasources{datasources: make(map[string]*Datasource)}
}
func TestLoadCsvEngine(t *testing.T) {
	dss := setupFileTest()
	ds, err := dss.load("testdata/good", "csv")
	if err != nil {
		t.Errorf("Load returns an error %v", err)
	}

	if ds.dstype != File {
		t.Errorf("Should be recognized as file datasource")
	}

	if ds.engine != CSV {
		t.Errorf("Should be recognized as CSV datasource but was recognized as '%s'", EngineToString(ds.GetEngine()))
	}
	if ds.filePath != "testdata/good/tmp/file.csv" {
		t.Errorf("The file path is '%s'", ds.filePath)
	}

	if ds.zip {
		t.Errorf("Should not be zipped")
	}

	if ds.gzip {
		t.Errorf("Should not be gzipped")
	}

	if len(ds.tags) != 0 && ds.tags[0] != "tagcsv" {
		t.Errorf("The tag should be found")
	}
}

func TestLoadZipCsvEngine(t *testing.T) {
	dss := setupFileTest()
	ds, err := dss.load("testdata/good", "zipcsv")
	if err != nil {
		t.Errorf("Load returns an error %v", err)
	}

	if ds.dstype != File {
		t.Errorf("Should be recognized as file datasource")
	}

	if ds.engine != CSV {
		t.Errorf("Should be recognized as CSV datasource but was recognized as '%s'", EngineToString(ds.GetEngine()))
	}
	if ds.filePath != "testdata/good/tmp/file.zip" {
		t.Errorf("The file path is '%s'", ds.filePath)
	}

	if !ds.zip {
		t.Errorf("Should be zipped")
	}

	if ds.gzip {
		t.Errorf("Should not be gzipped")
	}

	if len(ds.tags) != 0 && ds.tags[0] != "tagzipcsv" {
		t.Errorf("The tag should be found")
	}
}

func TestLoadGZipCsvEngine(t *testing.T) {
	dss := setupFileTest()
	ds, err := dss.load("testdata/good", "gzipcsv")
	if err != nil {
		t.Errorf("Load returns an error %v", err)
	}

	if ds.dstype != File {
		t.Errorf("Should be recognized as file datasource")
	}

	if ds.engine != CSV {
		t.Errorf("Should be recognized as CSV datasource but was recognized as '%s'", EngineToString(ds.GetEngine()))
	}
	if ds.filePath != "testdata/good/tmp/file.csv.gz" {
		t.Errorf("The file path is '%s'", ds.filePath)
	}

	if ds.zip {
		t.Errorf("Should not be zipped")
	}

	if !ds.gzip {
		t.Errorf("Should be gzipped")
	}

	if len(ds.tags) != 0 && ds.tags[0] != "taggzipcsv" {
		t.Errorf("The tag should be found")
	}
}

func TestLoadYamlEngine(t *testing.T) {
	dss := setupFileTest()
	ds, err := dss.load("testdata/good", "yaml")
	if err != nil {
		t.Errorf("Load returns an error %v", err)
	}

	if ds.dstype != File {
		t.Errorf("Should be recognized as file datasource")
	}

	if ds.engine != YAML {
		t.Errorf("Should be recognized as CSV datasource but was recognized as '%s'", EngineToString(ds.GetEngine()))
	}
	if ds.filePath != "testdata/good/tmp/file.yaml" {
		t.Errorf("The file path is '%s'", ds.filePath)
	}

	if ds.zip {
		t.Errorf("Should not be zipped")
	}

	if ds.gzip {
		t.Errorf("Should not be gzipped")
	}

	if len(ds.tags) != 0 && ds.tags[0] != "tagyaml" {
		t.Errorf("The tag should be found")
	}
}

func TestLoadJsonEngine(t *testing.T) {
	dss := setupFileTest()
	ds, err := dss.load("testdata/good", "json")
	if err != nil {
		t.Errorf("Load returns an error %v", err)
	}

	if ds.dstype != File {
		t.Errorf("Should be recognized as file datasource")
	}

	if ds.engine != JSON {
		t.Errorf("Should be recognized as JSON datasource but was recognized as '%s'", EngineToString(ds.GetEngine()))
	}
	if ds.filePath != "testdata/good/tmp/file.json" {
		t.Errorf("The file path is '%s'", ds.filePath)
	}

	if ds.zip {
		t.Errorf("Should not be zipped")
	}

	if ds.gzip {
		t.Errorf("Should not be gzipped")
	}

	if len(ds.tags) != 0 && ds.tags[0] != "tagjson" {
		t.Errorf("The tag should be found")
	}
}

func TestLoadStdio(t *testing.T) {
	dss := setupFileTest()
	ds, err := dss.load("testdata/good", "stdio")
	if err != nil {
		t.Errorf("Load returns an error %v", err)
	}

	if ds.filePath != "-" {
		t.Errorf("The file path is '%s'", ds.filePath)
	}
}

func TestLoadURL(t *testing.T) {
	dss := setupFileTest()
	ds, err := dss.load("testdata/good", "url")
	if err != nil {
		t.Errorf("Load returns an error %v", err)
	}

	if ds.url != "http://127.0.0.1/file.json" {
		t.Errorf("The URL is '%s'", ds.url)
	}
}

func TestLoadInline(t *testing.T) {
	dss := setupFileTest()
	ds, err := dss.load("testdata/good", "inline")
	if err != nil {
		t.Errorf("Load returns an error %v", err)
	}

	if ds.inline != "[{\"id\":1,\"value\":\"test\"}]" {
		t.Errorf("Inline is '%s'", ds.inline)
	}
}

func TestLoadNoPath(t *testing.T) {
	dss := setupFileTest()
	_, err := dss.load("testdata/fail", "nopath")
	if err == nil {
		t.Errorf("Load should returns an error")
	}

}

func TestOpenStdio(t *testing.T) {
	ds := Datasource{dstype: File, zip: false, gzip: false, filePath: "-"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")

	_, err := ds.OpenWriteFile(log)
	if err != nil {
		t.Fatalf("Should not return error and returned '%v'", err)
	}
	err = ds.CloseFile(log)
	if err != nil {
		t.Errorf("Should not return error and returned '%v'", err)
	}

	_, err = ds.OpenReadFile(log)
	if err != nil {
		t.Fatalf("Should not return error and returned '%v'", err)
	}
	err = ds.CloseFile(log)
	if err != nil {
		t.Errorf("Should not return error and returned '%v'", err)
	}

}

func TestOpenInline(t *testing.T) {
	ds := Datasource{dstype: File, zip: false, gzip: false, inline: "testinline"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")

	reader, err := ds.OpenReadFile(log)
	if err != nil {
		t.Fatalf("OpenReadFile Should not return error and returned '%v'", err)
	}

	test, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll Should not return error and returned '%v'", err)
	}
	if string(test) != "testinline" {
		t.Errorf("The content of inline is not the one we waits for :%v", test)
	}
	err = ds.CloseFile(log)
	if err != nil {
		t.Errorf("Close Should not return error and returned '%v'", err)
	}
}

func TestOpenFile(t *testing.T) {
	if _, err := os.Stat("testdata/tmp"); os.IsNotExist(err) {
		os.Mkdir("testdata/tmp", 0777)
	}

	test := []byte{1, 2, 3}

	ds := Datasource{dstype: File, zip: false, gzip: false, filePath: "testdata/tmp/testfile"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")

	_, err := ds.OpenWriteFile(log)
	if err != nil {
		t.Errorf("OpenWriteFile Should not return error and returned '%v'", err)
	}
	//	writer.Write(test)

	err = ds.CloseFile(log)
	if err != nil {
		t.Errorf("CloseFile Should not return error and returned '%v'", err)
	}

	fi, err := ds.Stat()
	if err != nil {
		t.Errorf("Stat Should not return error and returned '%v'", err)
	}
	if !fi.Mode().IsRegular() {
		t.Fatalf("Should be a file")
	}

	reader, err := ds.OpenReadFile(log)
	if err != nil {
		t.Fatalf("OpenReadFile Should not return error and returned '%v'", err)
	}
	reader.Read(test)
	if test[2] != 3 {
		t.Errorf("The content of file is not the one we waits for :%v", test)
	}
	err = ds.CloseFile(log)
	if err != nil {
		t.Errorf("Close Should not return error and returned '%v'", err)
	}
	err = ds.CloseFile(log)
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

	ds := Datasource{dstype: File, zip: true, gzip: false, filePath: "testdata/tmp/testfile.zip"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")

	writer, err := ds.OpenWriteFile(log)
	if err != nil {
		t.Fatalf("OpenWriteFile Should not return error and returned '%v'", err)
	}
	writer.Write(test)

	err = ds.CloseFile(log)
	if err != nil {
		t.Errorf("CloseFile Should not return error and returned '%v'", err)
	}

	reader, err := ds.OpenReadFile(log)
	if err != nil {
		t.Fatalf("OpenReadFile Should not return error and returned '%v'", err)
	}
	reader.Read(test)
	if test[2] != 3 {
		t.Errorf("The content of file is not the one we waits for :%v", test)
	}
	err = ds.CloseFile(log)
	if err != nil {
		t.Errorf("Close Should not return error and returned '%v'", err)
	}

	os.Remove("testdata/tmp/testfile.zip")
}

func TestOpenGzipFile(t *testing.T) {
	if _, err := os.Stat("testdata/tmp"); os.IsNotExist(err) {
		os.Mkdir("testdata/tmp", 0777)
	}

	test := []byte{1, 2, 3}

	ds := Datasource{dstype: File, zip: false, gzip: true, filePath: "testdata/tmp/testfile.gz"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")

	writer, err := ds.OpenWriteFile(log)
	if err != nil {
		t.Fatalf("OpenWriteFile Should not return error and returned '%v'", err)
	}
	writer.Write(test)

	err = ds.CloseFile(log)
	if err != nil {
		t.Errorf("CloseFile Should not return error and returned '%v'", err)
	}

	reader, err := ds.OpenReadFile(log)
	if err != nil {
		t.Fatalf("OpenReadFile Should not return error and returned '%v'", err)
	}
	reader.Read(test)
	if test[2] != 3 {
		t.Errorf("The content of file is not the one we waits for :%v", test)
	}
	err = ds.CloseFile(log)
	if err != nil {
		t.Errorf("Close Should not return error and returned '%v'", err)
	}

	os.Remove("testdata/tmp/testfile.gz")
}

func TestReadWrongZip(t *testing.T) {
	if _, err := os.Stat("testdata/tmp"); os.IsNotExist(err) {
		os.Mkdir("testdata/tmp", 0777)
	}

	test := []byte{1, 2, 3}

	ds := Datasource{dstype: File, zip: false, gzip: false, filePath: "testdata/tmp/testfile"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")

	writer, err := ds.OpenWriteFile(log)
	if err != nil {
		t.Fatalf("OpenWriteFile Should not return error and returned '%v'", err)
	}
	writer.Write(test)

	err = ds.CloseFile(log)
	if err != nil {
		t.Errorf("CloseFile Should not return error and returned '%v'", err)
	}

	ds.zip = true
	_, err = ds.OpenReadFile(log)
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

	ds := Datasource{dstype: File, zip: false, gzip: false, filePath: "testdata/tmp/testfile"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")

	writer, err := ds.OpenWriteFile(log)
	if err != nil {
		t.Fatalf("OpenWriteFile Should not return error and returned '%v'", err)
	}
	writer.Write(test)

	err = ds.CloseFile(log)
	if err != nil {
		t.Errorf("CloseFile Should not return error and returned '%v'", err)
	}

	ds.gzip = true
	_, err = ds.OpenReadFile(log)
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

	ds := Datasource{dstype: File, zip: false, gzip: false, filePath: "testdata/tmp/testfile"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")

	writer, err := ds.OpenWriteFile(log)
	if err != nil {
		t.Fatalf("OpenWriteFile Should not return error and returned '%v'", err)
	}
	writer.Write(test)

	err = ds.ResetFile(log)
	if err != nil {
		t.Errorf("ResetFile Should not return error and returned '%v'", err)
	}

	if _, err := os.Stat(ds.tmpFilePath); os.IsExist(err) {
		os.Remove(ds.tmpFilePath)
		t.Errorf("ResetFile Should have removed the temporary file'")
	}

	err = ds.ResetFile(log)
	if err != nil {
		t.Errorf("ResetFile Should not return error and returned '%v'", err)
	}
}

func TestOpenFileError(t *testing.T) {
	if _, err := os.Stat("testdata/tmp"); os.IsNotExist(err) {
		os.Mkdir("testdata/tmp", 0777)
	}

	ds := Datasource{dstype: File, zip: false, gzip: false, filePath: "testdata/tmp/nofile"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	logger.SetLevel(logrus.PanicLevel)

	_, err := ds.OpenReadFile(log)
	if err == nil {
		t.Fatalf("OpenReadFile Should return an error")
	}
}

func TestOpenUrlError(t *testing.T) {

	ds := Datasource{dstype: File, zip: false, gzip: false, url: "http://1.2.3.4.5"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	logger.SetLevel(logrus.PanicLevel)

	_, err := ds.OpenReadFile(log)
	if err == nil {
		t.Fatalf("OpenReadFile Should return an error")
	}
}

func TestOpenNoFileNoUrlError(t *testing.T) {

	ds := Datasource{dstype: File, zip: false, gzip: false}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	logger.SetLevel(logrus.PanicLevel)

	_, err := ds.OpenReadFile(log)
	if err == nil {
		t.Fatalf("OpenReadFile Should return an error")
	}
}

func TestOpenTmpFileError(t *testing.T) {
	if _, err := os.Stat("testdata/tmp"); os.IsNotExist(err) {
		os.Mkdir("testdata/tmp", 0777)
	}

	ds := Datasource{dstype: File, zip: false, gzip: false, filePath: "testdata/tmp/nodir/nofile"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	logger.SetLevel(logrus.PanicLevel)

	_, err := ds.OpenWriteFile(log)
	if err == nil {
		t.Fatalf("OpenWriteFile Should return an error")
	}
}

func TestResetTmpFileError(t *testing.T) {
	if _, err := os.Stat("testdata/tmp"); os.IsNotExist(err) {
		os.Mkdir("testdata/tmp", 0777)
	}

	ds := Datasource{dstype: File, zip: false, gzip: false, tmpFilePath: "testdata/tmp/nodir/nofile"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	logger.SetLevel(logrus.PanicLevel)

	tmpFile, err := ioutil.TempFile("testdata/tmp", "reset.")
	if err != nil {
		t.Fatalf("Should not return error and returned '%v'", err)
	}
	ds.fileHandle = tmpFile
	ds.filewriter = true
	err = ds.ResetFile(log)
	if err == nil {
		t.Fatalf("ResetFile Should return an error")
	}
}

func TestCloseFileError(t *testing.T) {
	if _, err := os.Stat("testdata/tmp"); os.IsNotExist(err) {
		os.Mkdir("testdata/tmp", 0777)
	}

	ds := Datasource{dstype: File, zip: false, gzip: false, tmpFilePath: "testdata/tmp/nodir/nofile"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	logger.SetLevel(logrus.PanicLevel)

	tmpFile, err := ioutil.TempFile("testdata/tmp", "reset.")
	if err != nil {
		t.Fatalf("Should not return error and returned '%v'", err)
	}
	ds.fileHandle = tmpFile
	ds.filewriter = true
	err = ds.CloseFile(log)
	if err == nil {
		t.Fatalf("ResetFile Should return an error")
	}
}

func TestCloseFileZipError(t *testing.T) {
	if _, err := os.Stat("testdata/tmp"); os.IsNotExist(err) {
		os.Mkdir("testdata/tmp", 0777)
	}

	ds := Datasource{dstype: File, zip: true, gzip: false, tmpFilePath: "testdata/tmp/nodir/nofile"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	logger.SetLevel(logrus.PanicLevel)

	tmpFile, err := ioutil.TempFile("testdata/tmp", "reset.")
	if err != nil {
		t.Fatalf("Should not return error and returned '%v'", err)
	}
	ds.fileHandle = tmpFile
	ds.filewriter = true
	err = ds.CloseFile(log)
	if err == nil {
		t.Fatalf("ResetFile Should return an error")
	}
}

func TestCloseFileZipNoDataError(t *testing.T) {
	if _, err := os.Stat("testdata/tmp"); os.IsNotExist(err) {
		os.Mkdir("testdata/tmp", 0777)
	}

	ds := Datasource{dstype: File, zip: true, gzip: false, filePath: "testdata/tmp/nodata.zip", tmpFilePath: "testdata/tmp/nodir/nofile"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	logger.SetLevel(logrus.PanicLevel)

	tmpFile, err := ioutil.TempFile("testdata/tmp", "reset.")
	if err != nil {
		t.Fatalf("Should not return error and returned '%v'", err)
	}
	ds.fileHandle = tmpFile
	ds.filewriter = true
	err = ds.CloseFile(log)
	if err == nil {
		t.Fatalf("ResetFile Should return an error")
	}
}

func TestFileCloseAll(t *testing.T) {
	dss := setupFileTest()
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	logger.SetLevel(logrus.PanicLevel)
	ds, err := dss.load("testdata/good", "csv")
	if err != nil {
		t.Errorf("Load returns an error %v", err)
	}
	dss.datasources["test"] = &ds
	if _, err := os.Stat("testdata/good/tmp"); os.IsNotExist(err) {
		os.Mkdir("testdata/good/tmp", 0777)
	}
	_, err = ds.OpenWriteFile(log)
	if err != nil {
		t.Fatalf("OpenWriteFile Should not return error and returned '%v'", err)
	}

	dss.CloseAll(log)
}

func TestStat(t *testing.T) {
	if _, err := os.Stat("testdata/tmp"); os.IsNotExist(err) {
		os.Mkdir("testdata/tmp", 0777)
	}

	ds := Datasource{dstype: File, zip: true, gzip: false, filePath: "testdata/tmp", tmpFilePath: ""}

	if _, err := ds.Stat(); os.IsNotExist(err) {
		t.Errorf("Stat Should have seen the file'")
	}
}
