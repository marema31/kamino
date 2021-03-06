package mockprovider_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/mockdatasource"
	"github.com/marema31/kamino/mockprovider"
)

func TestOk(t *testing.T) {
	pf := &mockprovider.MockProvider{}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")

	saver, err := pf.NewSaver(context.Background(), log, &mockdatasource.MockDatasource{}, "", "", "")
	if err != nil {
		t.Fatalf("NewSaver should not return error and returned '%v'", err)
	}
	mockedSaver := pf.Savers[0]

	loader, err := pf.NewLoader(context.Background(), log, &mockdatasource.MockDatasource{}, "", "")
	if err != nil {
		t.Fatalf("NewLoader should not return error and returned '%v'", err)
	}

	mockedLoader := pf.Loader
	mockedLoader.Content = []map[string]string{
		{"id": "1", "name": "Alice"},
		{"id": "2", "name": "Bob"},
	}
	mockedLoader.MockName = "myload"

	mockedLoader.ErrorLoad = fmt.Errorf("fake error")
	_, err = loader.Load(log)
	if err == nil {
		t.Fatalf("Load should return error")
	}

	mockedLoader.ErrorLoad = nil
	for loader.Next() {
		record, err := loader.Load(log)
		if err != nil {
			t.Fatalf("Load should not return error and returned '%v'", err)
		}

		if err = saver.Save(log, record); err != nil {
			t.Fatalf("Save should not return error and returned '%v'", err)
		}
	}

	for index, row := range mockedLoader.Content {
		for k, v := range row {
			if mockedSaver.Content[index][k] != v {
				t.Errorf("The loader and saver have not the same content")
			}
		}
	}

	_, err = loader.Load(log)
	if err == nil {
		t.Fatalf("Load should return error")
	}

	mockedSaver.ErrorSave = fmt.Errorf("fake error")
	err = saver.Save(log, nil)
	if err == nil {
		t.Fatalf("Save should return error")
	}

	if loader.Name() != "myload" {
		t.Errorf("Loader name function does not return the correct name %s", loader.Name())
	}
	mockedSaver.MockName = "mysave"
	if saver.Name() != "mysave" {
		t.Errorf("Loader name function does not return the correct name %s", loader.Name())
	}

	err = saver.Close(log)
	if err != nil {
		t.Errorf("Saver close should not return error and returned '%v'", err)
	}

	err = loader.Close(log)
	if err != nil {
		t.Errorf("Saver close should not return error and returned '%v'", err)
	}

}

func TestError(t *testing.T) {
	pf := &mockprovider.MockProvider{}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")

	pf.ErrorLoader = fmt.Errorf("Fake error")
	saver, err := pf.NewSaver(context.Background(), log, &mockdatasource.MockDatasource{}, "", "", "")
	if err != nil {
		t.Fatalf("NewSaver should not return error and returned '%v'", err)
	}

	pf.ErrorSaver = fmt.Errorf("Fake error")
	pf.SaverToFail = 1
	_, err = pf.NewSaver(context.Background(), log, &mockdatasource.MockDatasource{}, "", "", "")
	if err == nil {
		t.Fatalf("NewSaver should return error")
	}
	mockedSaver := pf.Savers[0]

	pf.ErrorLoader = nil
	loader, err := pf.NewLoader(context.Background(), log, &mockdatasource.MockDatasource{}, "", "")
	if err != nil {
		t.Fatalf("NewLoader should not return error and returned '%v'", err)
	}
	mockedLoader := pf.Loader

	pf.LoaderToFail = 1
	pf.ErrorLoader = fmt.Errorf("Fake error")
	_, err = pf.NewLoader(context.Background(), log, &mockdatasource.MockDatasource{}, "", "")
	if err == nil {
		t.Fatalf("NewLoader should return error")
	}

	mockedSaver.ErrorReset = fmt.Errorf("Fake error")
	err = saver.Close(log)
	if err != nil {
		t.Fatalf("Saver close should not return error and returned '%v'", err)
	}
	mockedSaver.ErrorClose = fmt.Errorf("Fake error")

	err = saver.Close(log)
	if err == nil {
		t.Fatalf("Saver close should return error")
	}
	mockedSaver.ErrorReset = nil
	err = saver.Reset(log)
	if err != nil {
		t.Fatalf("Saver Resetshould not return error and returned '%v'", err)
	}
	mockedSaver.ErrorReset = fmt.Errorf("Fake error")

	err = saver.Reset(log)
	if err == nil {
		t.Fatalf("Saver Reset should return error")
	}

	mockedLoader.ErrorClose = fmt.Errorf("fake error")
	err = loader.Close(log)
	if err == nil {
		t.Fatalf("Loader close should return error")
	}

}
