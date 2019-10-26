package json_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/mockdatasource"
	"github.com/marema31/kamino/provider/json"
)

func TestOk(t *testing.T) {
	source := mockdatasource.MockDatasource{Type: datasource.File, Engine: datasource.JSON, Zip: false, Gzip: false, FilePath: "sourcefile"}
	dest := mockdatasource.MockDatasource{Type: datasource.File, Engine: datasource.JSON, Zip: false, Gzip: false, FilePath: "destfile"}

	testString := []byte("[\n    {\n        \"id\": \"1\",\n        \"name\": \"Alice\"\n    },\n    {\n        \"id\": \"2\",\n        \"name\": \"Bob\"\n    }\n]")
	_, err := source.WriteBuf.Write(testString)
	if err != nil {
		t.Fatalf("Writing to the mocked source file should not return error and returned '%v'", err)
	}

	saver, err := json.NewSaver(context.Background(), &dest)
	if err != nil {
		t.Fatalf("NewSaver should not return error and returned '%v'", err)
	}
	loader, err := json.NewLoader(context.Background(), &source)
	if err != nil {
		t.Fatalf("NewLoader should not return error and returned '%v'", err)
	}

	lname := loader.Name()
	if lname != "sourcefile" {
		t.Errorf("Loader name function does not return the correct name %s", lname)
	}
	sname := saver.Name()
	if sname != "destfile" {
		t.Errorf("Saver name function does not return the correct name %s", sname)
	}

	for loader.Next() {
		record, err := loader.Load()
		if err != nil {
			t.Fatalf("Load should not return error and returned '%v'", err)
		}

		if err = saver.Save(record); err != nil {
			t.Fatalf("Save should not return error and returned '%v'", err)
		}
	}

	_, err = loader.Load()
	if err == nil {
		t.Errorf("Load should return error ")
	}

	err = saver.Close()
	if err != nil {
		t.Errorf("Saver close should not return error and returned '%v'", err)
	}

	err = loader.Close()
	if err != nil {
		t.Errorf("Saver close should not return error and returned '%v'", err)
	}

	readString := make([]byte, len(testString))
	_, err = dest.WriteBuf.Read(readString)
	if err != nil {
		t.Fatalf("Reading the mocked dest file should not return error and returned '%v'", err)
	}

	if !bytes.Equal(testString, readString) {
		t.Errorf("The read string is not equal to written one: '%s' != '%s'  ", string(testString), string(readString))
	}
}

func TestOpenError(t *testing.T) {
	source := mockdatasource.MockDatasource{ErrorOpenFile: fmt.Errorf("fake error"), Type: datasource.File, Engine: datasource.JSON, Zip: false, Gzip: false, FilePath: "sourcefile"}
	dest := mockdatasource.MockDatasource{ErrorOpenFile: fmt.Errorf("fake error"), Type: datasource.File, Engine: datasource.JSON, Zip: false, Gzip: false, FilePath: "destfile"}

	_, err := json.NewSaver(context.Background(), &dest)
	if err == nil {
		t.Fatalf("NewSaver should return error")
	}

	_, err = json.NewLoader(context.Background(), &source)
	if err == nil {
		t.Fatalf("NewLoader should return error")
	}

}

func TestCloseError(t *testing.T) {
	source := mockdatasource.MockDatasource{ErrorClose: fmt.Errorf("fake error"), Type: datasource.File, Engine: datasource.JSON, Zip: false, Gzip: false, FilePath: "sourcefile"}
	dest := mockdatasource.MockDatasource{ErrorClose: fmt.Errorf("fake error"), Type: datasource.File, Engine: datasource.JSON, Zip: false, Gzip: false, FilePath: "destfile"}

	testString := []byte("[{\"id\":\"1\",\"name\":\"Alice\"},{\"id\":\"2\",\"name\":\"Bob\"}]")
	_, err := source.WriteBuf.Write(testString)
	if err != nil {
		t.Fatalf("Writing to the mocked source file should not return error and returned '%v'", err)
	}

	saver, err := json.NewSaver(context.Background(), &dest)
	if err != nil {
		t.Fatalf("NewSaver should not return error and returned '%v'", err)
	}
	loader, err := json.NewLoader(context.Background(), &source)
	if err != nil {
		t.Fatalf("NewLoader should not return error and returned '%v'", err)
	}

	err = saver.Close()
	if err == nil {
		t.Fatalf("Saver close should return error")
	}

	err = loader.Close()
	if err == nil {
		t.Fatalf("Loader close should return error")
	}
}

func TestResetOK(t *testing.T) {
	dest := mockdatasource.MockDatasource{Type: datasource.File, Engine: datasource.JSON, Zip: false, Gzip: false, FilePath: "destfile"}

	saver, err := json.NewSaver(context.Background(), &dest)
	if err != nil {
		t.Fatalf("NewSaver should not return error and returned '%v'", err)
	}

	err = saver.Reset()
	if err != nil {
		t.Fatalf("Saver Resetshould not return error and returned '%v'", err)
	}

}

func TestResetError(t *testing.T) {
	dest := mockdatasource.MockDatasource{ErrorReset: fmt.Errorf("fake error"), Type: datasource.File, Engine: datasource.JSON, Zip: false, Gzip: false, FilePath: "destfile"}

	saver, err := json.NewSaver(context.Background(), &dest)
	if err != nil {
		t.Fatalf("NewSaver should not return error and returned '%v'", err)
	}

	err = saver.Reset()
	if err == nil {
		t.Fatalf("Saver Reset should return error")
	}
}

func TestWrongFormat(t *testing.T) {
	source := mockdatasource.MockDatasource{ErrorOpenFile: fmt.Errorf("fake error"), Type: datasource.File, Engine: datasource.JSON, Zip: false, Gzip: false, FilePath: "sourcefile"}
	testString := []byte("{\"id\":\"1\",\"name\":\"Alice\"},{\"id\":\"2\",\"name\":\"Bob\"}]")
	_, err := source.WriteBuf.Write(testString)
	if err != nil {
		t.Fatalf("Writing to the mocked source file should not return error and returned '%v'", err)
	}

	_, err = json.NewLoader(context.Background(), &source)
	if err == nil {
		t.Fatalf("NewLoader should return error")
	}

}
