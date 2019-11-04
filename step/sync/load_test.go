package sync_test

import (
	"context"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/mockdatasource"
	"github.com/marema31/kamino/mockprovider"
	"github.com/marema31/kamino/step/sync"
	"github.com/spf13/viper"
)

func setupLoad(path string, filename string) (context.Context, *logrus.Entry, datasource.Datasourcers, *viper.Viper, *mockprovider.MockProvider, error) {
	dss := mockdatasource.New()
	ds1 := mockdatasource.MockDatasource{Name: "ds1", Database: "db1", User: "user1", Tags: []string{"tag1", "tag2"}}
	ds2 := mockdatasource.MockDatasource{Name: "ds2", Database: "db2", User: "user2", Tags: []string{"tag2"}}
	ds3 := mockdatasource.MockDatasource{Name: "ds3", Database: "db3", User: "user3", Tags: []string{"tag3"}}
	ds4 := mockdatasource.MockDatasource{Name: "ds4", Database: "db4", User: "user4", Tags: []string{"tag3"}}

	dss.Insert([]string{"tag1", "tag2"}, []datasource.Type{datasource.Database}, []datasource.Engine{datasource.Mysql}, []*mockdatasource.MockDatasource{&ds2})
	dss.Insert([]string{"tag3"}, []datasource.Type{datasource.Database}, []datasource.Engine{datasource.Mysql}, []*mockdatasource.MockDatasource{&ds3, &ds4})
	dss.Insert([]string{"tagsource"}, []datasource.Type{datasource.Database}, []datasource.Engine{datasource.Mysql}, []*mockdatasource.MockDatasource{&ds1})
	dss.Insert([]string{"tagcache"}, []datasource.Type{datasource.Database}, []datasource.Engine{datasource.Mysql}, []*mockdatasource.MockDatasource{&ds2})
	v := viper.New()
	v.SetConfigName(filename) // The file will be named [filename].json, [filename].yaml or [filename.toml]
	v.AddConfigPath(path)
	err := v.ReadInConfig()
	ctx := context.Background()
	prov := &mockprovider.MockProvider{}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	logger.SetLevel(logrus.PanicLevel) //TODO: Modify this to ease issue 1 resolution

	return ctx, log, dss, v, prov, err
}

func TestSyncLoadOk(t *testing.T) {
	ctx, log, dss, v, prov, err := setupLoad("testdata/good/steps/", "syncok")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}

	priority, steps, err := sync.Load(ctx, log, "syncok", v, dss, prov)
	if err != nil {
		t.Errorf("Load should not returns an error, returned: %v", err)
	}

	if len(steps) != 1 {
		t.Fatalf("It should have been 1 steps created but it was created: %v", len(steps))
	}

	if priority != 42 {
		t.Errorf("The priority should be 42, it was: %v", priority)
	}

	step := steps[0]

	//For test purpose we must see what is inside the step and for this convert the interface to the presumed type
	s, ok := step.(*sync.Step)

	if !ok {
		t.Fatalf("The first step should be a sync step")
	}

	if s.Name != "namesyncok" {
		t.Errorf("The name of the first step should be namesyncok, it was: %v", s.Name)
	}

	//Using black box strategy, we cannot test the others field members, they could be tested only via the Do test
}

func TestSyncNoSource(t *testing.T) {
	ctx, log, dss, v, prov, err := setupLoad("testdata/fail/steps/", "nosource")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}

	_, _, err = sync.Load(ctx, log, "nosource", v, dss, prov)
	if err == nil {
		t.Errorf("Load should returns an error")
	}

}

func TestSyncNoDestination(t *testing.T) {
	ctx, log, dss, v, prov, err := setupLoad("testdata/fail/steps/", "nodestination")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}

	_, _, err = sync.Load(ctx, log, "nodestination", v, dss, prov)
	if err == nil {
		t.Errorf("Load should returns an error")
	}

}

func TestSyncNoSourceDatasource(t *testing.T) {
	ctx, log, dss, v, prov, err := setupLoad("testdata/fail/steps/", "nosourceds")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}

	_, _, err = sync.Load(ctx, log, "nosourceds", v, dss, prov)
	if err == nil {
		t.Errorf("Load should returns an error")
	}

}

func TestSyncNoDestinationDatasource(t *testing.T) {
	ctx, log, dss, v, prov, err := setupLoad("testdata/fail/steps/", "nodestinationds")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}

	_, _, err = sync.Load(ctx, log, "nodestinationds", v, dss, prov)
	if err == nil {
		t.Errorf("Load should returns an error")
	}

}

func TestSyncTooManySourceDatasource(t *testing.T) {
	ctx, log, dss, v, prov, err := setupLoad("testdata/fail/steps/", "toomanysourceds")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}

	_, _, err = sync.Load(ctx, log, "toomanysourceds", v, dss, prov)
	if err == nil {
		t.Errorf("Load should returns an error")
	}

}
