package tmpl_test

import (
	"context"
	"testing"

	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/mockdatasource"
	"github.com/marema31/kamino/step/tmpl"
	"github.com/spf13/viper"
)

func setupLoad(path string, filename string) (context.Context, datasource.Datasourcers, *viper.Viper, error) {
	dss := mockdatasource.New()
	ds1 := mockdatasource.MockDatasource{Name: "ds1", Database: "db1", User: "user1", Tags: []string{"tag1a", "tag1b"}}
	ds2 := mockdatasource.MockDatasource{Name: "ds2", Database: "db2", User: "user2", Tags: []string{"tag2"}}
	ds3 := mockdatasource.MockDatasource{Name: "ds3", Database: "db3", User: "user3", Tags: []string{"tag2"}}

	dss.Insert([]string{"tag1", "tag2"}, []datasource.Type{datasource.Database}, []datasource.Engine{datasource.Mysql}, []*mockdatasource.MockDatasource{&ds1, &ds2, &ds3})
	v := viper.New()
	v.SetConfigName(filename) // The file will be named [filename].json, [filename].yaml or [filename.toml]
	v.AddConfigPath(path)
	err := v.ReadInConfig()
	ctx := context.Background()
	return ctx, dss, v, err
}

func TestTmplLoadOk(t *testing.T) {
	ctx, dss, v, err := setupLoad("testdata/good/steps/", "tmplok")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}

	priority, steps, err := tmpl.Load(ctx, "tmplok", v, dss)
	if err != nil {
		t.Errorf("Load should not returns an error, returned: %v", err)
	}

	if len(steps) != 3 {
		t.Fatalf("It should have been 3 steps created but it was created: %v", len(steps))
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

	if s.Name != "nametmplok" {
		t.Errorf("The name of the first step should be nametmplok, it was: %v", s.Name)
	}

	//Using black box strategy, we cannot test the others field members, they could be tested only via the Do test
}

func TestTmplLoadNoTag(t *testing.T) {
	ctx, dss, v, err := setupLoad("testdata/good/steps/", "notags")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}

	priority, steps, err := tmpl.Load(ctx, "notags", v, dss)
	if err != nil {
		t.Errorf("Load should not returns an error, returned: %v", err)
	}

	if len(steps) != 0 {
		t.Fatalf("It should have been 0 steps created but it was created: %v", steps)
	}

	if priority != 42 {
		t.Errorf("The priority should be 42, it was: %v", priority)
	}
}

func TestTmplLoadReplace(t *testing.T) {
	ctx, dss, v, err := setupLoad("testdata/good/steps/", "replace")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}

	priority, steps, err := tmpl.Load(ctx, "replace", v, dss)
	if err != nil {
		t.Errorf("Load should not returns an error, returned: %v", err)
	}

	if len(steps) != 0 {
		t.Fatalf("It should have been 0 steps created but it was created: %v", steps)
	}

	if priority != 42 {
		t.Errorf("The priority should be 42, it was: %v", priority)
	}
}
func TestTmplLoadNoMode(t *testing.T) {
	ctx, dss, v, err := setupLoad("testdata/good/steps/", "nomode")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}

	priority, steps, err := tmpl.Load(ctx, "nomode", v, dss)
	if err != nil {
		t.Errorf("Load should not returns an error, returned: %v", err)
	}

	if len(steps) != 0 {
		t.Fatalf("It should have been 0 steps created but it was created: %v", steps)
	}

	if priority != 42 {
		t.Errorf("The priority should be 42, it was: %v", priority)
	}
}

func TestTmplLoadNoTemplate(t *testing.T) {
	ctx, dss, v, err := setupLoad("testdata/fail/steps/", "notemplate")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}
	_, _, err = tmpl.Load(ctx, "notemplate", v, dss)
	if err == nil {
		t.Errorf("Load should returns an error")
	}
}

func TestTmplLoadNoTemplateFile(t *testing.T) {
	ctx, dss, v, err := setupLoad("testdata/fail/steps/", "notemplatefile")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}
	_, _, err = tmpl.Load(ctx, "notemplatefile", v, dss)
	if err == nil {
		t.Errorf("Load should returns an error")
	}
}

func TestTmplLoadNoDestination(t *testing.T) {
	ctx, dss, v, err := setupLoad("testdata/fail/steps/", "nodestination")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}
	_, _, err = tmpl.Load(ctx, "nodestination", v, dss)
	if err == nil {
		t.Errorf("Load should returns an error")
	}
}

func TestTmplLoadTemplateWrong(t *testing.T) {
	ctx, dss, v, err := setupLoad("testdata/fail/steps/", "wrongtemplate")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}
	_, _, err = tmpl.Load(ctx, "wrongtemplate", v, dss)
	if err == nil {
		t.Errorf("Load should returns an error")
	}
}

func TestTmplLoadDestinationTemplateWrong(t *testing.T) {
	ctx, dss, v, err := setupLoad("testdata/fail/steps/", "wrongdestination")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}
	_, _, err = tmpl.Load(ctx, "wrongdestination", v, dss)
	if err == nil {
		t.Errorf("Load should returns an error")
	}
}

func TestTmplLoadWrongEngine(t *testing.T) {
	ctx, dss, v, err := setupLoad("testdata/fail/steps/", "wrongengine")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}
	_, _, err = tmpl.Load(ctx, "wrongengine", v, dss)
	if err == nil {
		t.Errorf("Load should returns an error")
	}
}
