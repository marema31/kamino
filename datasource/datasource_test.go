package datasource_test

import (
	"testing"

	"github.com/marema31/kamino/datasource"
)

func TestGetEngineOK(t *testing.T) {
	set := map[datasource.Engine]string{
		datasource.Mysql:    "mysql",
		datasource.Postgres: "postgresql",
		datasource.YAML:     "yaml",
		datasource.JSON:     "json",
		datasource.CSV:      "csv",
	}

	for k, v := range set {

		enginestr := datasource.EngineToString(k)
		if enginestr != v {
			t.Errorf("EngineToString should return '%s' and returned '%s'", v, enginestr)
		}
		engine, _ := datasource.StringToEngine(v)
		if engine != k {
			t.Errorf("GetEngine should return '%d' and returned '%d'", k, engine)
		}
	}
}

func TestGetEngineFail(t *testing.T) {
	engine := datasource.EngineToString(15)
	if engine != "Unknown" {
		t.Errorf("GetEngine should return 'Unknown' and returned '%s'", engine)
	}

}

func TestStringToTypeOK(t *testing.T) {

	set := map[datasource.Type]string{
		datasource.Database: "database",
		datasource.File:     "file",
	}

	for k, v := range set {
		dstype, err := datasource.StringToType(v)
		if err != nil {
			t.Fatalf("StringToEngine should not return error, returned: %v", err)
		}
		if dstype != k {
			t.Errorf("StringToEngine should return '%d' and returned '%d'", k, dstype)
		}
		dstypestr := datasource.TypeToString(k)
		if dstypestr != v {
			t.Errorf("StringToEngine should return '%s' and returned '%s'", dstypestr, v)
		}

	}
}

func TestStringToTypeFail(t *testing.T) {
	dstype, err := datasource.StringToType("unknown")

	if err == nil {
		t.Fatalf("StringToEngine should not return error, returned: %d", dstype)
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

/*
func TestGetNamedTags(t *testing.T) {
	ds := datasource.Datasource{}

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
*/
