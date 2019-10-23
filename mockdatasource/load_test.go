package mockdatasource_test

import (
	"testing"

	"github.com/marema31/kamino/mockdatasource"
)

// We are using private function, we must be in same package
func setupLoadTest() *mockdatasource.MockDatasources {
	return mockdatasource.New()
}

func TestLoadAllGood(t *testing.T) {
	dss := setupLoadTest()
	err := dss.LoadAll("testdata/good")
	if err != nil {
		t.Errorf("Load should not returns an error, was : %v", err)
	}

}

func TestLoadAllWrong(t *testing.T) {
	dss := setupLoadTest()
	err := dss.LoadAll("")
	if err == nil {
		t.Errorf("Load should returns an error")
	}
}
