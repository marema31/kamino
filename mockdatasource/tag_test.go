package mockdatasource_test

import (
	"sort"
	"strings"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/mockdatasource"
)

// We are providing the private datastore, we must be in same package

// setup fonction
func setupLookupDatastore() *mockdatasource.MockDatasources {
	dss := mockdatasource.New()
	ds1 := mockdatasource.MockDatasource{Name: "ds1"}
	ds2 := mockdatasource.MockDatasource{Name: "ds2"}
	ds3 := mockdatasource.MockDatasource{Name: "ds3"}
	ds4 := mockdatasource.MockDatasource{Name: "ds4"}
	ds5 := mockdatasource.MockDatasource{Name: "ds5"}
	ds6 := mockdatasource.MockDatasource{Name: "ds6"}
	notag := mockdatasource.MockDatasource{Name: "notag"}

	dss.Insert(true, []string{"tag1", "tag2"}, []datasource.Type{}, []datasource.Engine{}, []*mockdatasource.MockDatasource{&ds1, &ds2, &ds3, &ds4})
	dss.Insert(true, []string{"tag1.environment:fr"}, []datasource.Type{datasource.Database}, []datasource.Engine{datasource.Mysql}, []*mockdatasource.MockDatasource{&ds1, &ds5, &ds6})
	dss.Insert(true, []string{""}, []datasource.Type{datasource.Database}, []datasource.Engine{datasource.Mysql}, []*mockdatasource.MockDatasource{&notag})
	dss.Insert(true, []string{"empty"}, []datasource.Type{datasource.Database}, []datasource.Engine{datasource.Mysql}, []*mockdatasource.MockDatasource{})

	return dss
}

// teardown fonction
func teardownLookupDatastore(dss *mockdatasource.MockDatasources) {
}

func helperTestLookup(t *testing.T, dss *mockdatasource.MockDatasources, tags []string, dsTypes []datasource.Type, engines []datasource.Engine, awaited []string) {
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")

	result, _, _ := dss.Lookup(log, tags, nil, dsTypes, engines)

	sort.Strings(awaited)
	aw := strings.Join(awaited, " ")

	rnames := []string{}
	for _, ds := range result {
		tv := ds.FillTmplValues()
		if tv.Name != ds.GetName() {
			t.Errorf("FillTmplValue and GetName should give the same name but %s != %s", tv.Name, ds.GetName())
		}
		rnames = append(rnames, ds.GetName())
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
		[]string{"ds1", "ds2", "ds3", "ds4"},
	)

	helperTestLookup(
		t,
		dss,
		[]string{"tag1.tag2"},
		nil,
		nil,
		[]string{},
	)

	helperTestLookup(
		t,
		dss,
		[]string{"tag1.environment:fr"},
		[]datasource.Type{datasource.Database},
		[]datasource.Engine{datasource.Mysql},
		[]string{"ds1", "ds5", "ds6"},
	)

	helperTestLookup(
		t,
		dss,
		[]string{"empty"},
		nil,
		nil,
		[]string{""},
	)

	helperTestLookup(
		t,
		dss,
		[]string{},
		[]datasource.Type{datasource.Database},
		[]datasource.Engine{datasource.Mysql},
		[]string{"notag"},
	)

}
