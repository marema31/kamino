package datasource_test

import (
	"testing"

	"github.com/marema31/kamino/datasource"
)

func TestGetEngineOK(t *testing.T) {
	ds := datasource.Datasource{}

	set := map[datasource.Engine]string{
		datasource.Mysql:    "mysql",
		datasource.Postgres: "postgresql",
		datasource.YAML:     "yaml",
		datasource.JSON:     "json",
		datasource.CSV:      "csv",
	}

	for k, v := range set {
		ds.Engine = k
		engine := ds.GetEngine()
		if engine != v {
			t.Errorf("GetEngine should return '%s' and returned '%s'", v, engine)
		}

	}
}

func TestGetEngineFail(t *testing.T) {
	ds := datasource.Datasource{}

	ds.Engine = datasource.Engine(15)
	engine := ds.GetEngine()
	if engine != "Unknown" {
		t.Errorf("GetEngine should return 'Unknown' and returned '%s'", engine)
	}

}

func TestGetNamedTags(t *testing.T) {
	ds := datasource.Datasource{}

	ds.Tags = []string{"az1", "environment:production", "instance:fr"}
	value := ds.GetNamedTag("environment")
	if value != "production" {
		t.Errorf("Should return 'production' and returned '%s'", value)
	}

	value = ds.GetNamedTag("country")
	if value != "" {
		t.Errorf("Should return '' and returned '%s'", value)
	}
}
