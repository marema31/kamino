package datasource

import "testing"

// We are using private function, we must be in same package
func setupLoadTest() *Datasources {
	return New()
}

func TestLoadNoTag(t *testing.T) {
	dss := setupLoadTest()
	_, err := dss.load("testdata/good", "notag")
	if err != nil {
		t.Errorf("Load should not returns an error, returned: %v", err)
	}
}
func TestLoadNoEngine(t *testing.T) {
	dss := setupLoadTest()
	_, err := dss.load("testdata/fail", "noengine")
	if err == nil {
		t.Errorf("Load should returns an error")
	}
}

func TestLoadWrongEngine(t *testing.T) {
	dss := setupLoadTest()
	_, err := dss.load("testdata/fail", "wrongengine")
	if err == nil {
		t.Errorf("Load should returns an error")
	}
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
	err := dss.LoadAll("testdata/fail")
	if err == nil {
		t.Errorf("Load should returns an error")
	}
}
