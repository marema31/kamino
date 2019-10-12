package datasource

import "testing"

// We are using private function, we must be in same package

func TestLoadNoTag(t *testing.T) {
	if err := load("testdata/good/datasources", "notag"); err != nil {
		t.Errorf("Load should not returns an error, returned: %v", err)
	}
}
func TestLoadNoEngine(t *testing.T) {
	if err := load("testdata/fail/datasources", "noengine"); err == nil {
		t.Errorf("Load should returns an error")
	}
}

func TestLoadWrongEngine(t *testing.T) {
	if err := load("testdata/fail/datasources", "wrongengine"); err == nil {
		t.Errorf("Load should returns an error")
	}
}

func TestLoadAllGood(t *testing.T) {
	if err := LoadAll("testdata/good"); err != nil {
		t.Errorf("Load should not returns an error, was : %v", err)
	}

}

func TestLoadAllWrong(t *testing.T) {
	if err := LoadAll("testdata/fail"); err == nil {
		t.Errorf("Load should returns an error")
	}
}
