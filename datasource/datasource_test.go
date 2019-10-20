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

func TestStringToEngineOK(t *testing.T) {

	set := map[datasource.Engine]string{
		datasource.Mysql:    "mysql",
		datasource.Postgres: "postgresql",
		datasource.YAML:     "yaml",
		datasource.JSON:     "json",
		datasource.CSV:      "csv",
	}

	for k, v := range set {
		engine, err := datasource.StringToEngine(v)
		if err != nil {
			t.Fatalf("StringToEngine should not return error, returned: %v", err)
		}
		if engine != k {
			t.Errorf("StringToEngine should return '%d' and returned '%d'", k, engine)
		}

	}
}

func TestStringToEngineFail(t *testing.T) {
	engine, err := datasource.StringToEngine("unknown")

	if err == nil {
		t.Fatalf("StringToEngine should not return error, returned: %d", engine)
	}
}

func TestStringsToEnginesOk(t *testing.T) {
	engines, err := datasource.StringsToEngines([]string{"mariadb", "postgres"})
	if err != nil {
		t.Fatalf("StringsToEngines should not return error, returned: %v", err)
	}

	if len(engines) != 2 {
		t.Fatalf("StringsToEngines should return a slice of two elements, but it returned %d elements", len(engines))
	}

	if engines[0] != datasource.Mysql || engines[1] != datasource.Postgres {
		t.Fatalf("StringsToEngines should return %v, returned %v", engines, []datasource.Engine{datasource.Mysql, datasource.Postgres})
	}
}

func TestStringsToEnginesFail(t *testing.T) {
	_, err := datasource.StringsToEngines([]string{"mariadb", "unknown"})
	if err == nil {
		t.Fatalf("StringsToEngines should return error")
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
