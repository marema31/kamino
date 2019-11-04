package datasource

import (
	"sort"
	"strings"
	"testing"

	"github.com/Sirupsen/logrus"
)

// We are providing the private datastore, we must be in same package

// setup fonction
func setupLookupDatastore() *Datasources {
	dss := Datasources{}

	dss.datasources = make(map[string]*Datasource)
	ds1 := Datasource{name: "ds1", dstype: Database, engine: Mysql}
	ds2 := Datasource{name: "ds2", dstype: Database, engine: Postgres}
	ds3 := Datasource{name: "ds3", dstype: File, engine: JSON}
	ds4 := Datasource{name: "ds4", dstype: File, engine: YAML}
	ds5 := Datasource{name: "ds5", dstype: File, engine: YAML}
	ds6 := Datasource{name: "ds6", dstype: File, engine: YAML}
	ds7 := Datasource{name: "ds7", dstype: File, engine: YAML}
	notag := Datasource{name: "notag", dstype: File, engine: YAML}

	dss.datasources["ds1"] = &ds1
	dss.datasources["ds2"] = &ds2
	dss.datasources["ds3"] = &ds3
	dss.datasources["ds4"] = &ds4
	dss.datasources["ds5"] = &ds5
	dss.datasources["ds6"] = &ds6
	dss.datasources["ds7"] = &ds7
	dss.datasources["notag"] = &notag

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
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")

	result := dss.Lookup(log, tags, dsTypes, engines)

	anames := []string{}
	for _, ds := range awaited {
		anames = append(anames, ds.name)
	}
	sort.Strings(anames)
	aw := strings.Join(anames, " ")

	rnames := []string{}
	for _, ds := range result {
		tv := ds.FillTmplValues()
		rnames = append(rnames, tv.Name)
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
