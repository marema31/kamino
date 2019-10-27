package datasource

import (
	"io/ioutil"
	"os"
	"testing"
)

// We are using private function, we must be in same package
func setupFileTest() *Datasources {
	return &Datasources{}
}
func TestLoadCsvEngine(t *testing.T) {
	dss := setupFileTest()
	ds, err := dss.load("testdata/good/datasources", "csv")
	if err != nil {
		t.Errorf("Load returns an error %v", err)
	}

	if ds.dstype != File {
		t.Errorf("Should be recognized as file datasource")
	}

	if ds.engine != CSV {
		t.Errorf("Should be recognized as CSV datasource but was recognized as '%s'", EngineToString(ds.GetEngine()))
	}
	if ds.filePath != "tmp/file.csv" {
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
	ds, err := dss.load("testdata/good/datasources", "zipcsv")
	if err != nil {
		t.Errorf("Load returns an error %v", err)
	}

	if ds.dstype != File {
		t.Errorf("Should be recognized as file datasource")
	}

	if ds.engine != CSV {
		t.Errorf("Should be recognized as CSV datasource but was recognized as '%s'", EngineToString(ds.GetEngine()))
	}
	if ds.filePath != "tmp/file.zip" {
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
	ds, err := dss.load("testdata/good/datasources", "gzipcsv")
	if err != nil {
		t.Errorf("Load returns an error %v", err)
	}

	if ds.dstype != File {
		t.Errorf("Should be recognized as file datasource")
	}

	if ds.engine != CSV {
		t.Errorf("Should be recognized as CSV datasource but was recognized as '%s'", EngineToString(ds.GetEngine()))
	}
	if ds.filePath != "tmp/file.csv.gz" {
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
	ds, err := dss.load("testdata/good/datasources", "yaml")
	if err != nil {
		t.Errorf("Load returns an error %v", err)
	}

	if ds.dstype != File {
		t.Errorf("Should be recognized as file datasource")
	}

	if ds.engine != YAML {
		t.Errorf("Should be recognized as CSV datasource but was recognized as '%s'", EngineToString(ds.GetEngine()))
	}
	if ds.filePath != "tmp/file.yaml" {
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
	ds, err := dss.load("testdata/good/datasources", "json")
	if err != nil {
		t.Errorf("Load returns an error %v", err)
	}

	if ds.dstype != File {
		t.Errorf("Should be recognized as file datasource")
	}

	if ds.engine != JSON {
		t.Errorf("Should be recognized as JSON datasource but was recognized as '%s'", EngineToString(ds.GetEngine()))
	}
	if ds.filePath != "tmp/file.json" {
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
	ds, err := dss.load("testdata/good/datasources", "stdio")
	if err != nil {
		t.Errorf("Load returns an error %v", err)
	}

	if ds.filePath != "-" {
		t.Errorf("The file path is '%s'", ds.filePath)
	}
}

func TestLoadURL(t *testing.T) {
	dss := setupFileTest()
	ds, err := dss.load("testdata/good/datasources", "url")
	if err != nil {
		t.Errorf("Load returns an error %v", err)
	}

	if ds.url != "http://127.0.0.1/file.json" {
		t.Errorf("The URL is '%s'", ds.url)
	}
}

func TestLoadInline(t *testing.T) {
	dss := setupFileTest()
	ds, err := dss.load("testdata/good/datasources", "inline")
	if err != nil {
		t.Errorf("Load returns an error %v", err)
	}

	if ds.inline != "[{\"id\":1,\"value\":\"test\"}]" {
		t.Errorf("Inline is '%s'", ds.inline)
	}
}

func TestLoadNoPath(t *testing.T) {
	dss := setupFileTest()
	_, err := dss.load("testdata/fail/datasources", "nopath")
	if err == nil {
		t.Errorf("Load should returns an error")
	}

}

func TestOpenStdio(t *testing.T) {
	ds := Datasource{dstype: File, zip: false, gzip: false, filePath: "-"}

	_, err := ds.OpenWriteFile()
	if err != nil {
		t.Fatalf("Should not return error and returned '%v'", err)
	}
	err = ds.CloseFile()
	if err != nil {
		t.Errorf("Should not return error and returned '%v'", err)
	}

	_, err = ds.OpenReadFile()
	if err != nil {
		t.Fatalf("Should not return error and returned '%v'", err)
	}
	err = ds.CloseFile()
	if err != nil {
		t.Errorf("Should not return error and returned '%v'", err)
	}

}

func TestOpenInline(t *testing.T) {
	ds := Datasource{dstype: File, zip: false, gzip: false, inline: "testinline"}

	reader, err := ds.OpenReadFile()
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
	err = ds.CloseFile()
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

	_, err := ds.OpenWriteFile()
	if err != nil {
		t.Errorf("OpenWriteFile Should not return error and returned '%v'", err)
	}
	//	writer.Write(test)

	err = ds.CloseFile()
	if err != nil {
		t.Errorf("CloseFile Should not return error and returned '%v'", err)
	}

	reader, err := ds.OpenReadFile()
	if err != nil {
		t.Fatalf("OpenReadFile Should not return error and returned '%v'", err)
	}
	reader.Read(test)
	if test[2] != 3 {
		t.Errorf("The content of file is not the one we waits for :%v", test)
	}
	err = ds.CloseFile()
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

	writer, err := ds.OpenWriteFile()
	if err != nil {
		t.Fatalf("OpenWriteFile Should not return error and returned '%v'", err)
	}
	writer.Write(test)

	err = ds.CloseFile()
	if err != nil {
		t.Errorf("CloseFile Should not return error and returned '%v'", err)
	}

	reader, err := ds.OpenReadFile()
	if err != nil {
		t.Fatalf("OpenReadFile Should not return error and returned '%v'", err)
	}
	reader.Read(test)
	if test[2] != 3 {
		t.Errorf("The content of file is not the one we waits for :%v", test)
	}
	err = ds.CloseFile()
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

	writer, err := ds.OpenWriteFile()
	if err != nil {
		t.Fatalf("OpenWriteFile Should not return error and returned '%v'", err)
	}
	writer.Write(test)

	err = ds.CloseFile()
	if err != nil {
		t.Errorf("CloseFile Should not return error and returned '%v'", err)
	}

	reader, err := ds.OpenReadFile()
	if err != nil {
		t.Fatalf("OpenReadFile Should not return error and returned '%v'", err)
	}
	reader.Read(test)
	if test[2] != 3 {
		t.Errorf("The content of file is not the one we waits for :%v", test)
	}
	err = ds.CloseFile()
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

	writer, err := ds.OpenWriteFile()
	if err != nil {
		t.Fatalf("OpenWriteFile Should not return error and returned '%v'", err)
	}
	writer.Write(test)

	err = ds.CloseFile()
	if err != nil {
		t.Errorf("CloseFile Should not return error and returned '%v'", err)
	}

	ds.zip = true
	_, err = ds.OpenReadFile()
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

	writer, err := ds.OpenWriteFile()
	if err != nil {
		t.Fatalf("OpenWriteFile Should not return error and returned '%v'", err)
	}
	writer.Write(test)

	err = ds.CloseFile()
	if err != nil {
		t.Errorf("CloseFile Should not return error and returned '%v'", err)
	}

	ds.gzip = true
	_, err = ds.OpenReadFile()
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

	writer, err := ds.OpenWriteFile()
	if err != nil {
		t.Fatalf("OpenWriteFile Should not return error and returned '%v'", err)
	}
	writer.Write(test)

	err = ds.ResetFile()
	if err != nil {
		t.Errorf("ResetFile Should not return error and returned '%v'", err)
	}

	if _, err := os.Stat(ds.tmpFilePath); os.IsExist(err) {
		os.Remove(ds.tmpFilePath)
		t.Errorf("ResetFile Should have removed the temporary file'")
	}

}
