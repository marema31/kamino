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

	if ds == nil {
		t.Fatalf("Should have been inserted in datasources with the name csv")
	}

	if ds.Type != File {
		t.Errorf("Should be recognized as file datasource")
	}

	if ds.Engine != CSV {
		t.Errorf("Should be recognized as CSV datasource but was recognized as '%s'", ds.GetEngine())
	}
	if ds.FilePath != "tmp/file.csv" {
		t.Errorf("The file path is '%s'", ds.FilePath)
	}

	if ds.Zip {
		t.Errorf("Should not be zipped")
	}

	if ds.Gzip {
		t.Errorf("Should not be gzipped")
	}

	if len(ds.Tags) != 0 && ds.Tags[0] != "tagcsv" {
		t.Errorf("The tag should be found")
	}
}

func TestLoadZipCsvEngine(t *testing.T) {
	dss := setupFileTest()
	ds, err := dss.load("testdata/good/datasources", "zipcsv")
	if err != nil {
		t.Errorf("Load returns an error %v", err)
	}

	if ds == nil {
		t.Fatalf("Should have been inserted in datasources with the name zipcsv")
	}

	if ds.Type != File {
		t.Errorf("Should be recognized as file datasource")
	}

	if ds.Engine != CSV {
		t.Errorf("Should be recognized as CSV datasource but was recognized as '%s'", ds.GetEngine())
	}
	if ds.FilePath != "tmp/file.zip" {
		t.Errorf("The file path is '%s'", ds.FilePath)
	}

	if !ds.Zip {
		t.Errorf("Should be zipped")
	}

	if ds.Gzip {
		t.Errorf("Should not be gzipped")
	}

	if len(ds.Tags) != 0 && ds.Tags[0] != "tagzipcsv" {
		t.Errorf("The tag should be found")
	}
}

func TestLoadGZipCsvEngine(t *testing.T) {
	dss := setupFileTest()
	ds, err := dss.load("testdata/good/datasources", "gzipcsv")
	if err != nil {
		t.Errorf("Load returns an error %v", err)
	}

	if ds == nil {
		t.Fatalf("Should have been inserted in datasources with the name gzipcsv")
	}

	if ds.Type != File {
		t.Errorf("Should be recognized as file datasource")
	}

	if ds.Engine != CSV {
		t.Errorf("Should be recognized as CSV datasource but was recognized as '%s'", ds.GetEngine())
	}
	if ds.FilePath != "tmp/file.csv.gz" {
		t.Errorf("The file path is '%s'", ds.FilePath)
	}

	if ds.Zip {
		t.Errorf("Should not be zipped")
	}

	if !ds.Gzip {
		t.Errorf("Should be gzipped")
	}

	if len(ds.Tags) != 0 && ds.Tags[0] != "taggzipcsv" {
		t.Errorf("The tag should be found")
	}
}

func TestLoadYamlEngine(t *testing.T) {
	dss := setupFileTest()
	ds, err := dss.load("testdata/good/datasources", "yaml")
	if err != nil {
		t.Errorf("Load returns an error %v", err)
	}

	if ds == nil {
		t.Fatalf("Should have been inserted in datasources with the name yaml")
	}

	if ds.Type != File {
		t.Errorf("Should be recognized as file datasource")
	}

	if ds.Engine != YAML {
		t.Errorf("Should be recognized as CSV datasource but was recognized as '%s'", ds.GetEngine())
	}
	if ds.FilePath != "tmp/file.yaml" {
		t.Errorf("The file path is '%s'", ds.FilePath)
	}

	if ds.Zip {
		t.Errorf("Should not be zipped")
	}

	if ds.Gzip {
		t.Errorf("Should not be gzipped")
	}

	if len(ds.Tags) != 0 && ds.Tags[0] != "tagyaml" {
		t.Errorf("The tag should be found")
	}
}

func TestLoadJsonEngine(t *testing.T) {
	dss := setupFileTest()
	ds, err := dss.load("testdata/good/datasources", "json")
	if err != nil {
		t.Errorf("Load returns an error %v", err)
	}

	if ds == nil {
		t.Fatalf("Should have been inserted in datasources with the name json")
	}

	if ds.Type != File {
		t.Errorf("Should be recognized as file datasource")
	}

	if ds.Engine != JSON {
		t.Errorf("Should be recognized as JSON datasource but was recognized as '%s'", ds.GetEngine())
	}
	if ds.FilePath != "tmp/file.json" {
		t.Errorf("The file path is '%s'", ds.FilePath)
	}

	if ds.Zip {
		t.Errorf("Should not be zipped")
	}

	if ds.Gzip {
		t.Errorf("Should not be gzipped")
	}

	if len(ds.Tags) != 0 && ds.Tags[0] != "tagjson" {
		t.Errorf("The tag should be found")
	}
}

func TestLoadStdio(t *testing.T) {
	dss := setupFileTest()
	ds, err := dss.load("testdata/good/datasources", "stdio")
	if err != nil {
		t.Errorf("Load returns an error %v", err)
	}

	if ds == nil {
		t.Fatalf("Should have been inserted in datasources with the name stdio")
	}

	if ds.FilePath != "-" {
		t.Errorf("The file path is '%s'", ds.FilePath)
	}
}

func TestLoadURL(t *testing.T) {
	dss := setupFileTest()
	ds, err := dss.load("testdata/good/datasources", "url")
	if err != nil {
		t.Errorf("Load returns an error %v", err)
	}

	if ds == nil {
		t.Fatalf("Should have been inserted in datasources with the name url")
	}

	if ds.URL != "http://127.0.0.1/file.json" {
		t.Errorf("The URL is '%s'", ds.URL)
	}
}

func TestLoadInline(t *testing.T) {
	dss := setupFileTest()
	ds, err := dss.load("testdata/good/datasources", "inline")
	if err != nil {
		t.Errorf("Load returns an error %v", err)
	}

	if ds == nil {
		t.Fatalf("Should have been inserted in datasources with the name inline")
	}

	if ds.Inline != "[{\"id\":1,\"value\":\"test\"}]" {
		t.Errorf("Inline is '%s'", ds.Inline)
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
	ds := Datasource{Type: File, Zip: false, Gzip: false, FilePath: "-"}

	_, err := ds.OpenWriteFile()
	if err != nil {
		t.Fatalf("Should not return error and returned '%v'", err)
	}
	ds.CloseFile()
	if err != nil {
		t.Errorf("Should not return error and returned '%v'", err)
	}

	_, err = ds.OpenReadFile()
	if err != nil {
		t.Fatalf("Should not return error and returned '%v'", err)
	}
	ds.CloseFile()
	if err != nil {
		t.Errorf("Should not return error and returned '%v'", err)
	}

}

func TestOpenInline(t *testing.T) {
	ds := Datasource{Type: File, Zip: false, Gzip: false, Inline: "testinline"}

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
	ds.CloseFile()
	if err != nil {
		t.Errorf("Close Should not return error and returned '%v'", err)
	}
}

func TestOpenFile(t *testing.T) {
	if _, err := os.Stat("testdata/tmp"); os.IsNotExist(err) {
		os.Mkdir("testdata/tmp", 0777)
	}

	test := []byte{1, 2, 3}

	ds := Datasource{Type: File, Zip: false, Gzip: false, FilePath: "testdata/tmp/testfile"}

	writer, err := ds.OpenWriteFile()
	if err != nil {
		t.Fatalf("OpenWriteFile Should not return error and returned '%v'", err)
	}
	writer.Write(test)

	ds.CloseFile()
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
	ds.CloseFile()
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

	ds := Datasource{Type: File, Zip: true, Gzip: false, FilePath: "testdata/tmp/testfile.zip"}

	writer, err := ds.OpenWriteFile()
	if err != nil {
		t.Fatalf("OpenWriteFile Should not return error and returned '%v'", err)
	}
	writer.Write(test)

	ds.CloseFile()
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
	ds.CloseFile()
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

	ds := Datasource{Type: File, Zip: false, Gzip: true, FilePath: "testdata/tmp/testfile.gz"}

	writer, err := ds.OpenWriteFile()
	if err != nil {
		t.Fatalf("OpenWriteFile Should not return error and returned '%v'", err)
	}
	writer.Write(test)

	ds.CloseFile()
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
	ds.CloseFile()
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

	ds := Datasource{Type: File, Zip: false, Gzip: false, FilePath: "testdata/tmp/testfile"}

	writer, err := ds.OpenWriteFile()
	if err != nil {
		t.Fatalf("OpenWriteFile Should not return error and returned '%v'", err)
	}
	writer.Write(test)

	ds.CloseFile()
	if err != nil {
		t.Errorf("CloseFile Should not return error and returned '%v'", err)
	}

	ds.Zip = true
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

	ds := Datasource{Type: File, Zip: false, Gzip: false, FilePath: "testdata/tmp/testfile"}

	writer, err := ds.OpenWriteFile()
	if err != nil {
		t.Fatalf("OpenWriteFile Should not return error and returned '%v'", err)
	}
	writer.Write(test)

	ds.CloseFile()
	if err != nil {
		t.Errorf("CloseFile Should not return error and returned '%v'", err)
	}

	ds.Gzip = true
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

	ds := Datasource{Type: File, Zip: false, Gzip: false, FilePath: "testdata/tmp/testfile"}

	writer, err := ds.OpenWriteFile()
	if err != nil {
		t.Fatalf("OpenWriteFile Should not return error and returned '%v'", err)
	}
	writer.Write(test)

	ds.ResetFile()
	if err != nil {
		t.Errorf("ResetFile Should not return error and returned '%v'", err)
	}

	if _, err := os.Stat(ds.tmpFilePath); os.IsExist(err) {
		os.Remove(ds.tmpFilePath)
		t.Errorf("ResetFile Should have removed the temporary file'")
	}

}
