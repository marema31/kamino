package mockprovider_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/marema31/kamino/mockdatasource"
	"github.com/marema31/kamino/mockprovider"
)

func TestOk(t *testing.T) {
	pf := &mockprovider.MockProvider{}

	saver, err := pf.NewSaver(context.Background(), &mockdatasource.MockDatasource{}, "", "", "")
	if err != nil {
		t.Fatalf("NewSaver should not return error and returned '%v'", err)
	}
	loader, err := pf.NewLoader(context.Background(), &mockdatasource.MockDatasource{}, "", "")
	if err != nil {
		t.Fatalf("NewLoader should not return error and returned '%v'", err)
	}

	loader.Content = []map[string]string{
		{"id": "1", "name": "Alice"},
		{"id": "2", "name": "Bob"},
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

	for index, row := range loader.Content {
		for k, v := range row {
			if saver.Content[index][k] != v {
				t.Errorf("The loader and saver have not the same content")
			}
		}
	}

	_, err = loader.Load()
	if err == nil {
		t.Fatalf("Load should return error")
	}

	loader.MockName = "myload"
	if loader.Name() != "myload" {
		t.Errorf("Loader name function does not return the correct name %s", loader.Name())
	}
	saver.MockName = "mysave"
	if saver.Name() != "mysave" {
		t.Errorf("Loader name function does not return the correct name %s", loader.Name())
	}

	err = saver.Close()
	if err != nil {
		t.Errorf("Saver close should not return error and returned '%v'", err)
	}

	err = loader.Close()
	if err != nil {
		t.Errorf("Saver close should not return error and returned '%v'", err)
	}

}

func TestError(t *testing.T) {
	pf := &mockprovider.MockProvider{}

	pf.ErrorLoader = fmt.Errorf("Fake error")
	saver, err := pf.NewSaver(context.Background(), &mockdatasource.MockDatasource{}, "", "", "")
	if err != nil {
		t.Fatalf("NewSaver should not return error and returned '%v'", err)
	}

	pf.ErrorSaver = fmt.Errorf("Fake error")
	_, err = pf.NewSaver(context.Background(), &mockdatasource.MockDatasource{}, "", "", "")
	if err == nil {
		t.Fatalf("NewSaver should return error")
	}

	pf.ErrorLoader = nil
	loader, err := pf.NewLoader(context.Background(), &mockdatasource.MockDatasource{}, "", "")
	if err != nil {
		t.Fatalf("NewLoader should not return error and returned '%v'", err)
	}
	pf.ErrorLoader = fmt.Errorf("Fake error")
	_, err = pf.NewLoader(context.Background(), &mockdatasource.MockDatasource{}, "", "")
	if err == nil {
		t.Fatalf("NewLoader should return error")
	}

	saver.ErrorReset = fmt.Errorf("Fake error")
	err = saver.Close()
	if err != nil {
		t.Fatalf("Saver close should not return error and returned '%v'", err)
	}
	saver.ErrorClose = fmt.Errorf("Fake error")

	err = saver.Close()
	if err == nil {
		t.Fatalf("Saver close should return error")
	}
	saver.ErrorReset = nil
	err = saver.Reset()
	if err != nil {
		t.Fatalf("Saver Resetshould not return error and returned '%v'", err)
	}
	saver.ErrorReset = fmt.Errorf("Fake error")

	err = saver.Reset()
	if err == nil {
		t.Fatalf("Saver Reset should return error")
	}

	err = loader.Close()
	if err != nil {
		t.Fatalf("Saver close should not return error and returned '%v'", err)
	}
	loader.ErrorClose = fmt.Errorf("Fake error")
	err = loader.Close()
	if err == nil {
		t.Fatalf("Loader close should return error")
	}

}
