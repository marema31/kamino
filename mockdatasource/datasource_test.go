package mockdatasource_test

import (
	"testing"

	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/mockdatasource"
)

func TestGetEngineOK(t *testing.T) {
	ds := mockdatasource.MockDatasource{}

	set := []datasource.Engine{
		datasource.Mysql,
		datasource.Postgres,
		datasource.YAML,
		datasource.JSON,
		datasource.CSV,
	}

	for _, v := range set {
		ds.Engine = v
		engine := ds.GetEngine()
		if engine != v {
			t.Errorf("GetEngine should return '%d' and returned '%d'", v, engine)
		}

	}
}

func TestTypeOK(t *testing.T) {
	ds := mockdatasource.MockDatasource{}

	set := []datasource.Type{
		datasource.Database,
		datasource.File,
	}

	for _, v := range set {
		ds.Type = v
		Type := ds.GetType()
		if Type != v {
			t.Errorf("GetType should return '%d' and returned '%d'", v, Type)
		}

	}
}

func TestGetNamedTags(t *testing.T) {
	ds := mockdatasource.MockDatasource{}

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
