package datasource

import (
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
)

// We are using private function, we must be in same package
func setupLoadTest() *Datasources {
	return New(time.Microsecond*2, 2)
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
	logger := logrus.New()
	log := logger.WithField("a", 1)
	err := dss.LoadAll("testdata/good", log)
	if err != nil {
		t.Errorf("Load should not returns an error, was : %v", err)
	}

}

func TestLoadAllWrong(t *testing.T) {
	dss := setupLoadTest()
	logger := logrus.New()
	log := logger.WithField("a", 1)
	err := dss.LoadAll("testdata/fail", log)
	if err == nil {
		t.Errorf("Load should returns an error")
	}
}

func TestLoadAllUnknown(t *testing.T) {
	dss := setupLoadTest()
	logger := logrus.New()
	log := logger.WithField("a", 1)
	err := dss.LoadAll("testdata/unknown", log)
	if err == nil {
		t.Errorf("Load should returns an error")
	}
}
func TestLoadAllEmpty(t *testing.T) {
	dss := setupLoadTest()
	logger := logrus.New()
	log := logger.WithField("a", 1)
	err := dss.LoadAll("testdata/empty", log)
	if err == nil {
		t.Errorf("Load should returns an error")
	}
}

func TestGetNamedTags(t *testing.T) {
	ds := Datasource{}

	ds.tags = []string{"az1", "environment:production", "instance:fr"}
	value := ds.GetNamedTag("environment")
	if value != "production" {
		t.Errorf("Should return 'production' and returned '%s'", value)
	}

	value = ds.GetNamedTag("country")
	if value != "" {
		t.Errorf("Should return '' and returned '%s'", value)
	}
}

func TestGetName(t *testing.T) {
	ds := Datasource{}

	ds.name = "production"
	value := ds.GetName()
	if value != "production" {
		t.Errorf("Should return 'production' and returned '%s'", value)
	}
}

func TestGetEngine(t *testing.T) {
	ds := Datasource{}

	ds.engine = JSON
	value := ds.GetEngine()
	if value != JSON {
		t.Errorf("Should return 'production' and returned '%s'", EngineToString(value))
	}
}
