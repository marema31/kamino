package step_test

import (
	"context"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/mockdatasource"
	"github.com/marema31/kamino/mockprovider"
	"github.com/marema31/kamino/provider"
	"github.com/marema31/kamino/step"
	"github.com/marema31/kamino/step/migration"
	"github.com/marema31/kamino/step/shell"
	"github.com/marema31/kamino/step/sqlscript"
	"github.com/marema31/kamino/step/sync"
	"github.com/marema31/kamino/step/tmpl"
)

func setupLoad() (context.Context, *logrus.Entry, datasource.Datasourcers, provider.Provider, step.Creater) {
	dss := mockdatasource.New()
	ds1 := mockdatasource.MockDatasource{Name: "ds1", Database: "db1", User: "user1", Tags: []string{"tag1a", "tag1b"}}
	ds2 := mockdatasource.MockDatasource{Name: "ds2", Database: "db2", User: "user2", Tags: []string{"tag2"}}
	ds3 := mockdatasource.MockDatasource{Name: "ds3", Database: "db3", User: "user3", Schema: "az", Tags: []string{"tag2"}}
	ds4 := mockdatasource.MockDatasource{Name: "ds4", Database: "db4", User: "user4", Tags: []string{"tag3"}}

	dss.Insert([]string{"tag1", "tag2"}, []datasource.Type{datasource.Database}, []datasource.Engine{datasource.Mysql}, []*mockdatasource.MockDatasource{&ds1, &ds2, &ds3})
	dss.Insert([]string{"tag1"}, []datasource.Type{datasource.Database}, []datasource.Engine{datasource.Mysql}, []*mockdatasource.MockDatasource{&ds1, &ds2, &ds3})
	dss.Insert([]string{"tag3"}, []datasource.Type{datasource.Database}, []datasource.Engine{datasource.Mysql}, []*mockdatasource.MockDatasource{&ds3, &ds4})
	dss.Insert([]string{"tagsource"}, []datasource.Type{datasource.Database}, []datasource.Engine{datasource.Mysql}, []*mockdatasource.MockDatasource{&ds1})
	dss.Insert([]string{"tagcache"}, []datasource.Type{datasource.Database}, []datasource.Engine{datasource.Mysql}, []*mockdatasource.MockDatasource{&ds2})
	prov := &mockprovider.MockProvider{}
	ctx := context.Background()
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	logger.SetLevel(logrus.PanicLevel)
	sf := &step.Factory{}
	return ctx, log, dss, prov, sf
}

func TestStepLoadError(t *testing.T) {
	ctx, log, dss, prov, sf := setupLoad()

	_, _, _, err := sf.Load(ctx, log, "testdata/fail", "nofolder", dss, prov, []string{}, []string{}, []string{}, true, false)
	if err == nil {
		t.Errorf("Load should returns an error")
	}

	_, _, _, err = sf.Load(ctx, log, "testdata/fail", "dummy", dss, prov, []string{}, []string{}, []string{}, true, false)
	if err == nil {
		t.Errorf("Load should returns an error")
	}

	_, _, _, err = sf.Load(ctx, log, "testdata/fail", "notype", dss, prov, []string{}, []string{}, []string{}, true, false)
	if err == nil {
		t.Errorf("Load should returns an error")
	}

	_, _, _, err = sf.Load(ctx, log, "testdata/fail", "unknown", dss, prov, []string{}, []string{}, []string{}, true, false)
	if err == nil {
		t.Errorf("Load should returns an error")
	}
}
func TestStepLoadMigrationOk(t *testing.T) {
	ctx, log, dss, prov, sf := setupLoad()

	priority, _, steps, err := sf.Load(ctx, log, "testdata/good", "migrationok", dss, prov, []string{}, []string{}, []string{}, true, false)
	if err != nil {
		t.Errorf("Load should not returns an error, returned: %v", err)
	}

	if len(steps) != 3 {
		t.Fatalf("It should have been 3 steps created but it was created: %d", len(steps))
	}

	if priority != 42 {
		t.Errorf("The priority should be 42, it was: %v", priority)
	}

	step := steps[0]

	//For test purpose we must see what is inside the step and for this convert the interface to the presumed type
	s, ok := step.(*migration.Step)

	if !ok {
		t.Fatalf("The first step should be a migration step")
	}

	if s.Name != "namemigrationok:0" {
		t.Errorf("The name of the first step should be namemigrationok:0, it was: %v", s.Name)
	}
}

func TestStepLoadShellOk(t *testing.T) {
	ctx, log, dss, prov, sf := setupLoad()

	priority, _, steps, err := sf.Load(ctx, log, "testdata/good", "shellok", dss, prov, []string{}, []string{}, []string{}, true, false)
	if err != nil {
		t.Errorf("Load should not returns an error, returned: %v", err)
	}

	if len(steps) != 3 {
		t.Fatalf("It should have been 3 steps created but it was created: %d", len(steps))
	}

	if priority != 42 {
		t.Errorf("The priority should be 42, it was: %v", priority)
	}

	step := steps[0]

	//For test purpose we must see what is inside the step and for this convert the interface to the presumed type
	s, ok := step.(*shell.Step)

	if !ok {
		t.Fatalf("The first step should be a shell step")
	}

	if s.Name != "nameshellok:0" {
		t.Errorf("The name of the first step should be nameshellok:0, it was: %v", s.Name)
	}
}
func TestStepLoadSqlscriptOk(t *testing.T) {
	ctx, log, dss, prov, sf := setupLoad()

	priority, _, steps, err := sf.Load(ctx, log, "testdata/good", "sqlscriptok", dss, prov, []string{}, []string{}, []string{}, true, false)
	if err != nil {
		t.Errorf("Load should not returns an error, returned: %v", err)
	}

	if len(steps) != 3 {
		t.Fatalf("It should have been 3 steps created but it was created: %d", len(steps))
	}

	if priority != 42 {
		t.Errorf("The priority should be 42, it was: %v", priority)
	}

	step := steps[0]

	//For test purpose we must see what is inside the step and for this convert the interface to the presumed type
	s, ok := step.(*sqlscript.Step)

	if !ok {
		t.Fatalf("The first step should be a sqlscript step")
	}

	if s.Name != "namesqlscriptok:0" {
		t.Errorf("The name of the first step should be namesqlscriptok:0, it was: %v", s.Name)
	}
}

func TestStepLoadSyncOk(t *testing.T) {
	ctx, log, dss, prov, sf := setupLoad()

	priority, _, steps, err := sf.Load(ctx, log, "testdata/good", "syncok", dss, prov, []string{}, []string{}, []string{}, true, false)
	if err != nil {
		t.Errorf("Load should not returns an error, returned: %v", err)
	}

	if len(steps) != 1 {
		t.Fatalf("It should have been 1 steps created but it was created: %d", len(steps))
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

	if s.Name != "namesyncok:0" {
		t.Errorf("The name of the first step should be namesyncok:0, it was: %v", s.Name)
	}
}

func TestStepLoadTmplOk(t *testing.T) {
	ctx, log, dss, prov, sf := setupLoad()

	priority, _, steps, err := sf.Load(ctx, log, "testdata/good", "tmplok", dss, prov, []string{}, []string{}, []string{}, true, false)
	if err != nil {
		t.Errorf("Load should not returns an error, returned: %v", err)
	}

	if len(steps) != 3 {
		t.Fatalf("It should have been 3 steps created but it was created: %d", len(steps))
	}

	if priority != 42 {
		t.Errorf("The priority should be 42, it was: %v", priority)
	}

	step := steps[0]

	//For test purpose we must see what is inside the step and for this convert the interface to the presumed type
	s, ok := step.(*tmpl.Step)

	if !ok {
		t.Fatalf("The first step should be a tmpl step")
	}

	if s.Name != "nametmplok:0" {
		t.Errorf("The name of the first step should be nametmplok:0, it was: %v", s.Name)
	}
}

func TestStepNameFound(t *testing.T) {
	ctx, log, dss, prov, sf := setupLoad()

	priority, _, steps, err := sf.Load(ctx, log, "testdata/good", "tmplok", dss, prov, []string{}, []string{"nametmplok"}, []string{}, true, false)
	if err != nil {
		t.Errorf("Load should not returns an error, returned: %v", err)
	}

	if len(steps) != 3 {
		t.Fatalf("It should have been 3 steps created but it was created: %d", len(steps))
	}

	if priority != 42 {
		t.Errorf("The priority should be 42, it was: %v", priority)
	}
}

func TestGlobbedStepNameFound(t *testing.T) {
	ctx, log, dss, prov, sf := setupLoad()

	priority, _, steps, err := sf.Load(ctx, log, "testdata/good", "tmplok", dss, prov, []string{}, []string{"nametmpl*"}, []string{}, true, false)
	if err != nil {
		t.Errorf("Load should not returns an error, returned: %v", err)
	}

	if len(steps) != 3 {
		t.Fatalf("It should have been 3 steps created but it was created: %d", len(steps))
	}

	if priority != 42 {
		t.Errorf("The priority should be 42, it was: %v", priority)
	}
}

func TestStepNameNotFound(t *testing.T) {
	ctx, log, dss, prov, sf := setupLoad()

	priority, _, steps, err := sf.Load(ctx, log, "testdata/good", "tmplok", dss, prov, []string{}, []string{"unknown"}, []string{}, true, false)
	if err != nil {
		t.Errorf("Load should not returns an error, returned: %v", err)
	}

	if len(steps) != 0 {
		t.Fatalf("It should have been 0 steps created but it was created: %d", len(steps))
	}

	if priority != 0 {
		t.Errorf("The priority should be 0, it was: %v", priority)
	}
}

func TestStepTypeFound(t *testing.T) {
	ctx, log, dss, prov, sf := setupLoad()

	priority, _, steps, err := sf.Load(ctx, log, "testdata/good", "tmplok", dss, prov, []string{}, []string{}, []string{"template"}, true, false)
	if err != nil {
		t.Errorf("Load should not returns an error, returned: %v", err)
	}

	if len(steps) != 3 {
		t.Fatalf("It should have been 3 steps created but it was created: %d", len(steps))
	}

	if priority != 42 {
		t.Errorf("The priority should be 42, it was: %v", priority)
	}
}

func TestStepTypeNotFound(t *testing.T) {
	ctx, log, dss, prov, sf := setupLoad()

	priority, _, steps, err := sf.Load(ctx, log, "testdata/good", "tmplok", dss, prov, []string{}, []string{}, []string{"shell"}, true, false)
	if err != nil {
		t.Errorf("Load should not returns an error, returned: %v", err)
	}

	if len(steps) != 0 {
		t.Fatalf("It should have been 0 steps created but it was created: %d", len(steps))
	}

	if priority != 0 {
		t.Errorf("The priority should be 0, it was: %v", priority)
	}

	_, _, _, err = sf.Load(ctx, log, "testdata/good", "tmplok", dss, prov, []string{}, []string{}, []string{"unknown"}, true, false)
	if err == nil {
		t.Errorf("Load should returns an error")
	}
}
