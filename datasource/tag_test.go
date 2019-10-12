package datasource

import (
	"sort"
	"strings"
	"testing"
)

// We are providing the private datastore, we must be in same package

// setup fonction
func setupLookupDatastore() {
	datasources = map[string]*Datasource{
		"ds1":   &Datasource{Name: "ds1", Type: Database, Engine: Mysql},
		"ds2":   &Datasource{Name: "ds2", Type: Database, Engine: Postgres},
		"ds3":   &Datasource{Name: "ds3", Type: File, Engine: JSON},
		"ds4":   &Datasource{Name: "ds4", Type: File, Engine: YAML},
		"ds5":   &Datasource{Name: "ds5", Type: File, Engine: YAML},
		"ds6":   &Datasource{Name: "ds6", Type: File, Engine: YAML},
		"ds7":   &Datasource{Name: "ds7", Type: File, Engine: YAML},
		"notag": &Datasource{Name: "notag", Type: File, Engine: YAML},
	}

	tagToDatasource = map[string][]string{
		"tag1":           []string{"ds1", "ds2"},
		"tag2":           []string{"ds1", "ds3", "ds4"},
		"environment:us": []string{"ds5", "ds6"},
		"environment:fr": []string{"ds7", "ds2", "ds4"},
		"":               []string{"notag"},
	}
}

// teardown fonction
func teardownLookupDatastore() {
	datasources = nil
	tagToDatasource = nil
}

func helperTestlookupOneTag(t *testing.T, tag string, dsTypes []Type, engines []Engine, anames []string) {
	rnames := lookupOneTag(tag, dsTypes, engines)

	sort.Strings(anames)
	aw := strings.Join(anames, " ")

	sort.Strings(rnames)
	rw := strings.Join(rnames, " ")

	if rw != aw {
		t.Errorf("'%s, %v, %v' should returns [%s] but returned [%s]", tag, dsTypes, engines, aw, rw)
	}

}

func TestLookupOneTag(t *testing.T) {
	setupLookupDatastore()
	defer teardownLookupDatastore()

	helperTestlookupOneTag(
		t,
		"tag1",
		nil,
		nil,
		[]string{"ds1", "ds2"},
	)

	helperTestlookupOneTag(
		t,
		"environment:fr",
		nil,
		nil,
		[]string{"ds2", "ds4", "ds7"},
	)

	helperTestlookupOneTag(
		t,
		"",
		nil,
		nil,
		[]string{"notag"},
	)

	helperTestlookupOneTag(
		t,
		"tag1",
		[]Type{Database},
		nil,
		[]string{"ds1", "ds2"},
	)

	helperTestlookupOneTag(
		t,
		"",
		[]Type{File},
		nil,
		[]string{"notag"},
	)
}

func helperTestLookup(t *testing.T, tags []string, dsTypes []Type, engines []Engine, awaited []*Datasource) {
	result := Lookup(tags, dsTypes, engines)

	anames := []string{}
	for _, ds := range awaited {
		anames = append(anames, ds.Name)
	}
	sort.Strings(anames)
	aw := strings.Join(anames, " ")

	rnames := []string{}
	for _, ds := range result {
		rnames = append(rnames, ds.Name)
	}
	sort.Strings(rnames)
	rw := strings.Join(rnames, " ")

	if rw != aw {
		t.Errorf("'%v, %v, %v' should returns [%s] but returned [%s]", tags, dsTypes, engines, aw, rw)
	}

}

func TestLookup(t *testing.T) {
	setupLookupDatastore()
	defer teardownLookupDatastore()

	helperTestLookup(
		t,
		[]string{"tag1", "tag2"},
		nil,
		nil,
		[]*Datasource{
			datasources["ds1"],
			datasources["ds2"],
			datasources["ds3"],
			datasources["ds4"],
		},
	)

	helperTestLookup(
		t,
		[]string{"tag1.tag2"},
		nil,
		nil,
		[]*Datasource{
			datasources["ds1"],
		},
	)

	helperTestLookup(
		t,
		[]string{"tag1.environment:fr"},
		nil,
		nil,
		[]*Datasource{datasources["ds2"]},
	)

	helperTestLookup(
		t,
		[]string{"tag2.environment:us"},
		nil,
		nil,
		[]*Datasource{},
	)

	helperTestLookup(
		t,
		[]string{"tag1.tag2", "tag2.environment:fr"},
		nil,
		nil,
		[]*Datasource{
			datasources["ds1"],
			datasources["ds4"],
		},
	)

	helperTestLookup(
		t,
		[]string{"tag1.tag2", "tag2.environment:fr"},
		[]Type{Database},
		nil,
		[]*Datasource{
			datasources["ds1"],
		},
	)

	helperTestLookup(
		t,
		[]string{"tag1.tag2", "tag2.environment:fr"},
		[]Type{Database},
		[]Engine{Mysql},
		[]*Datasource{
			datasources["ds1"],
		},
	)

	helperTestLookup(
		t,
		[]string{"tag1.tag2", "tag2.environment:fr"},
		nil,
		[]Engine{JSON, YAML},
		[]*Datasource{
			datasources["ds4"],
		},
	)

	helperTestLookup(
		t,
		[]string{},
		[]Type{Database},
		nil,
		[]*Datasource{
			datasources["ds1"],
			datasources["ds2"],
		},
	)

	helperTestLookup(
		t,
		[]string{},
		[]Type{Database},
		[]Engine{Mysql},
		[]*Datasource{
			datasources["ds1"],
		},
	)

	helperTestLookup(
		t,
		[]string{},
		nil,
		[]Engine{JSON, YAML},
		[]*Datasource{
			datasources["ds3"],
			datasources["ds4"],
			datasources["ds5"],
			datasources["ds6"],
			datasources["ds7"],
			datasources["notag"],
		},
	)
}
