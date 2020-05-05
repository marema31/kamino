package datasource

import (
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
)

// We are using private function, we must be in same package
func setupFileTest() (*Datasources, *logrus.Entry) {
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")

	return &Datasources{datasources: make(map[string]*Datasource)}, log
}
func TestLoadCsvEngine(t *testing.T) {
	dss, log := setupFileTest()
	ds, err := dss.load(log, "testdata/good", "datasources", "csv")
	if err != nil {
		t.Errorf("Load returns an error %v", err)
	}

	if ds.dstype != File {
		t.Errorf("Should be recognized as file datasource")
	}

	if ds.engine != CSV {
		t.Errorf("Should be recognized as CSV datasource but was recognized as '%s'", EngineToString(ds.GetEngine()))
	}
	if ds.file.FilePath != "testdata/good/tmp/file.csv" {
		t.Errorf("The file path is '%s'", ds.file.FilePath)
	}

	if ds.file.Zip {
		t.Errorf("Should not be zipped")
	}

	if ds.file.Gzip {
		t.Errorf("Should not be Gzipped")
	}

	if len(ds.tags) != 0 && ds.tags[0] != "tagcsv" {
		t.Errorf("The tag should be found")
	}
}

func TestLoadZipCsvEngine(t *testing.T) {
	dss, log := setupFileTest()
	ds, err := dss.load(log, "testdata/good", "datasources", "zipcsv")
	if err != nil {
		t.Errorf("Load returns an error %v", err)
	}

	if ds.dstype != File {
		t.Errorf("Should be recognized as file datasource")
	}

	if ds.engine != CSV {
		t.Errorf("Should be recognized as CSV datasource but was recognized as '%s'", EngineToString(ds.GetEngine()))
	}
	if ds.file.FilePath != "testdata/good/tmp/file.zip" {
		t.Errorf("The file path is '%s'", ds.file.FilePath)
	}

	if !ds.file.Zip {
		t.Errorf("Should be zipped")
	}

	if ds.file.Gzip {
		t.Errorf("Should not be Gzipped")
	}

	if len(ds.tags) != 0 && ds.tags[0] != "tagzipcsv" {
		t.Errorf("The tag should be found")
	}
}

func TestLoadGzipCsvEngine(t *testing.T) {
	dss, log := setupFileTest()
	ds, err := dss.load(log, "testdata/good", "datasources", "gzipcsv")
	if err != nil {
		t.Errorf("Load returns an error %v", err)
	}

	if ds.dstype != File {
		t.Errorf("Should be recognized as file datasource")
	}

	if ds.engine != CSV {
		t.Errorf("Should be recognized as CSV datasource but was recognized as '%s'", EngineToString(ds.GetEngine()))
	}
	if ds.file.FilePath != "testdata/good/tmp/file.csv.gz" {
		t.Errorf("The file path is '%s'", ds.file.FilePath)
	}

	if ds.file.Zip {
		t.Errorf("Should not be zipped")
	}

	if !ds.file.Gzip {
		t.Errorf("Should be Gzipped")
	}

	if len(ds.tags) != 0 && ds.tags[0] != "taggzipcsv" {
		t.Errorf("The tag should be found")
	}
}

func TestLoadYamlEngine(t *testing.T) {
	dss, log := setupFileTest()
	ds, err := dss.load(log, "testdata/good", "datasources", "yaml")
	if err != nil {
		t.Errorf("Load returns an error %v", err)
	}

	if ds.dstype != File {
		t.Errorf("Should be recognized as file datasource")
	}

	if ds.engine != YAML {
		t.Errorf("Should be recognized as CSV datasource but was recognized as '%s'", EngineToString(ds.GetEngine()))
	}
	if ds.file.FilePath != "testdata/good/tmp/file.yaml" {
		t.Errorf("The file path is '%s'", ds.file.FilePath)
	}

	if ds.file.Zip {
		t.Errorf("Should not be zipped")
	}

	if ds.file.Gzip {
		t.Errorf("Should not be Gzipped")
	}

	if len(ds.tags) != 0 && ds.tags[0] != "tagyaml" {
		t.Errorf("The tag should be found")
	}
}

func TestLoadJsonEngine(t *testing.T) {
	dss, log := setupFileTest()
	ds, err := dss.load(log, "testdata/good", "datasources", "json")
	if err != nil {
		t.Errorf("Load returns an error %v", err)
	}

	if ds.dstype != File {
		t.Errorf("Should be recognized as file datasource")
	}

	if ds.engine != JSON {
		t.Errorf("Should be recognized as JSON datasource but was recognized as '%s'", EngineToString(ds.GetEngine()))
	}
	if ds.file.FilePath != "testdata/good/tmp/file.json" {
		t.Errorf("The file path is '%s'", ds.file.FilePath)
	}

	if ds.file.Zip {
		t.Errorf("Should not be zipped")
	}

	if ds.file.Gzip {
		t.Errorf("Should not be Gzipped")
	}

	if len(ds.tags) != 0 && ds.tags[0] != "tagjson" {
		t.Errorf("The tag should be found")
	}
}

func TestLoadStdio(t *testing.T) {
	dss, log := setupFileTest()
	ds, err := dss.load(log, "testdata/good", "datasources", "stdio")
	if err != nil {
		t.Errorf("Load returns an error %v", err)
	}

	if ds.file.FilePath != "-" {
		t.Errorf("The file path is '%s'", ds.file.FilePath)
	}
}

func TestLoadURL(t *testing.T) {
	dss, log := setupFileTest()
	ds, err := dss.load(log, "testdata/good", "datasources", "url")
	if err != nil {
		t.Errorf("Load returns an error %v", err)
	}

	if ds.file.URL != "http://127.0.0.1/file.json" {
		t.Errorf("The URL is '%s'", ds.file.URL)
	}
}

func TestLoadInline(t *testing.T) {
	dss, log := setupFileTest()
	ds, err := dss.load(log, "testdata/good", "datasources", "inline")
	if err != nil {
		t.Errorf("Load returns an error %v", err)
	}

	if ds.file.Inline != "[{\"id\":1,\"value\":\"test\"}]" {
		t.Errorf("Inline is '%s'", ds.file.Inline)
	}
}

func TestLoadNoPath(t *testing.T) {
	dss, log := setupFileTest()
	_, err := dss.load(log, "testdata/fail", "datasources", "nopath")
	if err == nil {
		t.Errorf("Load should returns an error")
	}

}

func TestOpenStdio(t *testing.T) {
	ds := Datasource{dstype: File}
	ds.file.Zip = false
	ds.file.Gzip = false
	ds.file.FilePath = "-"
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

func TestFileCloseAll(t *testing.T) {
	dss, log := setupFileTest()
	logger := logrus.New()
	log = logger.WithField("appname", "kamino")
	logger.SetLevel(logrus.PanicLevel)
	ds, err := dss.load(log, "testdata/good", "datasources", "csv")
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

func TestResetFile(t *testing.T) {
	if _, err := os.Stat("testdata/tmp"); os.IsNotExist(err) {
		os.Mkdir("testdata/tmp", 0777)
	}

	test := []byte{1, 2, 3}

	ds := Datasource{dstype: File}
	ds.file.Zip = false
	ds.file.Gzip = false
	ds.file.FilePath = "testdata/tmp/testfile"
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

	err = ds.ResetFile(log)
	if err != nil {
		t.Errorf("ResetFile Should not return error and returned '%v'", err)
	}
}

func TestStat(t *testing.T) {
	if _, err := os.Stat("testdata/tmp"); os.IsNotExist(err) {
		os.Mkdir("testdata/tmp", 0777)
	}

	ds := Datasource{dstype: File}
	ds.file.Zip = false
	ds.file.Gzip = false
	ds.file.FilePath = "testdata/tmp"

	if _, err := ds.Stat(); os.IsNotExist(err) {
		t.Errorf("Stat Should have seen the file'")
	}
}
