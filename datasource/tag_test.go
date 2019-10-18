package datasource

import (
	"sort"
	"strings"
	"testing"
)

// We are providing the private datastore, we must be in same package

// setup fonction
func setupLookupDatastore() *Datasources {
	dss := Datasources{}

	dss.datasources = map[string]*Datasource{
		"ds1":   {Name: "ds1", Type: Database, Engine: Mysql},
		"ds2":   {Name: "ds2", Type: Database, Engine: Postgres},
		"ds3":   {Name: "ds3", Type: File, Engine: JSON},
		"ds4":   {Name: "ds4", Type: File, Engine: YAML},
		"ds5":   {Name: "ds5", Type: File, Engine: YAML},
		"ds6":   {Name: "ds6", Type: File, Engine: YAML},
		"ds7":   {Name: "ds7", Type: File, Engine: YAML},
		"notag": {Name: "notag", Type: File, Engine: YAML},
	}

	dss.tagToDatasource = map[string][]string{
		"tag1":           {"ds1", "ds2"},
		"tag2":           {"ds1", "ds3", "ds4"},
		"environment:us": {"ds5", "ds6"},
		"environment:fr": {"ds7", "ds2", "ds4"},
		"":               {"notag"},
	}
	return &dss
}

// teardown fonction
func teardownLookupDatastore(dss *Datasources) {
	dss.datasources = nil
	dss.tagToDatasource = nil
}

func helperTestlookupOneTag(t *testing.T, dss *Datasources, tag string, dsTypes []Type, engines []Engine, anames []string) {
	rnames := dss.lookupOneTag(tag, dsTypes, engines)

	sort.Strings(anames)
	aw := strings.Join(anames, " ")

	sort.Strings(rnames)
	rw := strings.Join(rnames, " ")

	if rw != aw {
		t.Errorf("'%s, %v, %v' should returns [%s] but returned [%s]", tag, dsTypes, engines, aw, rw)
	}

}

func TestLookupOneTag(t *testing.T) {
	dss := setupLookupDatastore()
	defer teardownLookupDatastore(dss)

	helperTestlookupOneTag(
		t,
		dss,
		"tag1",
		nil,
		nil,
		[]string{"ds1", "ds2"},
	)

	helperTestlookupOneTag(
		t,
		dss,
		"environment:fr",
		nil,
		nil,
		[]string{"ds2", "ds4", "ds7"},
	)

	helperTestlookupOneTag(
		t,
		dss,
		"",
		nil,
		nil,
		[]string{"notag"},
	)

	helperTestlookupOneTag(
		t,
		dss,
		"tag1",
		[]Type{Database},
		nil,
		[]string{"ds1", "ds2"},
	)

	helperTestlookupOneTag(
		t,
		dss,
		"",
		[]Type{File},
		nil,
		[]string{"notag"},
	)
}

func helperTestLookup(t *testing.T, dss *Datasources, tags []string, dsTypes []Type, engines []Engine, awaited []*Datasource) {
	result := dss.Lookup(tags, dsTypes, engines)

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
	dss := setupLookupDatastore()
	defer teardownLookupDatastore(dss)

	helperTestLookup(
		t,
		dss,
		[]string{"tag1", "tag2"},
		nil,
		nil,
		[]*Datasource{
			dss.datasources["ds1"],
			dss.datasources["ds2"],
			dss.datasources["ds3"],
			dss.datasources["ds4"],
		},
	)

	helperTestLookup(
		t,
		dss,
		[]string{"tag1.tag2"},
		nil,
		nil,
		[]*Datasource{
			dss.datasources["ds1"],
		},
	)

	helperTestLookup(
		t,
		dss,
		[]string{"tag1.environment:fr"},
		nil,
		nil,
		[]*Datasource{dss.datasources["ds2"]},
	)

	helperTestLookup(
		t,
		dss,
		[]string{"tag2.environment:us"},
		nil,
		nil,
		[]*Datasource{},
	)

	helperTestLookup(
		t,
		dss,
		[]string{"tag1.tag2", "tag2.environment:fr"},
		nil,
		nil,
		[]*Datasource{
			dss.datasources["ds1"],
			dss.datasources["ds4"],
		},
	)

	helperTestLookup(
		t,
		dss,
		[]string{"tag1.tag2", "tag2.environment:fr"},
		[]Type{Database},
		nil,
		[]*Datasource{
			dss.datasources["ds1"],
		},
	)

	helperTestLookup(
		t,
		dss,
		[]string{"tag1.tag2", "tag2.environment:fr"},
		[]Type{Database},
		[]Engine{Mysql},
		[]*Datasource{
			dss.datasources["ds1"],
		},
	)

	helperTestLookup(
		t,
		dss,
		[]string{"tag1.tag2", "tag2.environment:fr"},
		nil,
		[]Engine{JSON, YAML},
		[]*Datasource{
			dss.datasources["ds4"],
		},
	)

	helperTestLookup(
		t,
		dss,
		[]string{},
		[]Type{Database},
		nil,
		[]*Datasource{
			dss.datasources["ds1"],
			dss.datasources["ds2"],
		},
	)

	helperTestLookup(
		t,
		dss,
		[]string{},
		[]Type{Database},
		[]Engine{Mysql},
		[]*Datasource{
			dss.datasources["ds1"],
		},
	)

	helperTestLookup(
		t,
		dss,
		[]string{},
		nil,
		[]Engine{JSON, YAML},
		[]*Datasource{
			dss.datasources["ds3"],
			dss.datasources["ds4"],
			dss.datasources["ds5"],
			dss.datasources["ds6"],
			dss.datasources["ds7"],
			dss.datasources["notag"],
		},
	)
}
