package migration_test

import (
	"context"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/mockdatasource"
	"github.com/marema31/kamino/step/migration"
	"github.com/spf13/viper"
)

func setupLoad(path string, filename string) (context.Context, *logrus.Entry, datasource.Datasourcers, *viper.Viper, error) {
	dss := mockdatasource.New()
	ds1 := mockdatasource.MockDatasource{Name: "ds1", Database: "db1", User: "user1", Tags: []string{"tag1a", "tag1b"}}
	ds2 := mockdatasource.MockDatasource{Name: "ds2", Database: "db2", User: "user2", Tags: []string{"tag2"}}
	ds3 := mockdatasource.MockDatasource{Name: "ds3", Database: "db3", User: "user3", Schema: "az", Tags: []string{"tag2"}}

	dss.Insert([]string{"tag1", "tag2"}, []datasource.Type{datasource.Database}, []datasource.Engine{datasource.Mysql}, []*mockdatasource.MockDatasource{&ds1, &ds2, &ds3})
	dss.Insert([]string{"tag1"}, []datasource.Type{datasource.Database}, []datasource.Engine{datasource.Mysql}, []*mockdatasource.MockDatasource{&ds1, &ds2, &ds3})
	v := viper.New()
	v.SetConfigName(filename) // The file will be named [filename].json, [filename].yaml or [filename.toml]
	v.AddConfigPath(path)
	err := v.ReadInConfig()
	ctx := context.Background()
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	logger.SetLevel(logrus.PanicLevel)

	return ctx, log, dss, v, err
}

func TestMigrationLoadOk(t *testing.T) {
	ctx, log, dss, v, err := setupLoad("testdata/good/steps/", "migrationok")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}

	priority, steps, err := migration.Load(ctx, log, "testdata/good", "namemigrationok", 0, v, dss, false, false)
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

	//Using black box strategy, we cannot test the others field members, they could be tested only via the Do test
}

func TestMigrationLoadNoTag(t *testing.T) {
	ctx, log, dss, v, err := setupLoad("testdata/good/steps/", "notags")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}

	priority, steps, err := migration.Load(ctx, log, "testdata/good", "notags", 0, v, dss, false, false)
	if err != nil {
		t.Errorf("Load should not returns an error, returned: %v", err)
	}

	if len(steps) != 0 {
		t.Fatalf("It should have been 0 steps created but it was created: %v", len(steps))
	}

	if priority != 42 {
		t.Errorf("The priority should be 42, it was: %v", priority)
	}
}

func TestMigrationLoadNoFolder(t *testing.T) {
	ctx, log, dss, v, err := setupLoad("testdata/fail/steps/", "nofolder")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}
	_, _, err = migration.Load(ctx, log, "testdata/fail", "nofolder", 0, v, dss, false, false)
	if err == nil {
		t.Errorf("Load should returns an error")
	}
}

func TestMigrationLoadDFolderTemplateWrong(t *testing.T) {
	ctx, log, dss, v, err := setupLoad("testdata/fail/steps/", "wrongfolder")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}
	_, _, err = migration.Load(ctx, log, "testdata/fail", "wrongfolder", 0, v, dss, false, false)
	if err == nil {
		t.Errorf("Load should returns an error")
	}
}

func TestMigrationLoadWrongEngine(t *testing.T) {
	ctx, log, dss, v, err := setupLoad("testdata/fail/steps/", "wrongengine")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}
	_, _, err = migration.Load(ctx, log, "testdata/fail", "wrongengine", 0, v, dss, false, false)
	if err == nil {
		t.Errorf("Load should returns an error")
	}
}

func TestMigrationLoadNoQuery(t *testing.T) {
	ctx, log, dss, v, err := setupLoad("testdata/fail/steps/", "noquery")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}
	_, _, err = migration.Load(ctx, log, "testdata/fail", "noquery", 0, v, dss, false, false)
	if err == nil {
		t.Errorf("Load should returns an error")
	}
}
func TestMigrationLoadWrongQuery(t *testing.T) {
	ctx, log, dss, v, err := setupLoad("testdata/fail/steps/", "wrongquery")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}
	_, _, err = migration.Load(ctx, log, "testdata/fail", "wrongquery", 0, v, dss, false, false)
	if err == nil {
		t.Errorf("Load should returns an error")
	}
}

func TestMigrationPostLoad(t *testing.T) {
	ctx, log, dss, v, err := setupLoad("testdata/good/steps/", "migrationok")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}

	_, steps, err := migration.Load(ctx, log, "testdata/good", "namemigrationok", 0, v, dss, false, false)
	if err != nil {
		t.Errorf("Load should not returns an error, returned: %v", err)
	}

	superseed := make(map[string]string)
	err = steps[0].PostLoad(log, superseed)
	if err != nil {
		t.Errorf("PostLoad should not returns an error, returned: %v", err)
	}

	superseed["unknown.limit"] = "novalue"
	err = steps[0].PostLoad(log, superseed)
	if err != nil {
		t.Errorf("PostLoad should not returns an error, returned: %v", err)
	}

	superseed["migration.limit"] = "1"
	superseed["migration.noAdmin"] = "true"
	superseed["migration.noUser"] = "false"
	superseed["migration.dir"] = "up"
	err = steps[0].PostLoad(log, superseed)
	if err != nil {
		t.Errorf("PostLoad should not returns an error, returned: %v", err)
	}
	superseed["migration.dir"] = "down"
	err = steps[0].PostLoad(log, superseed)
	if err != nil {
		t.Errorf("PostLoad should not returns an error, returned: %v", err)
	}
	superseed["migration.dir"] = "status"
	err = steps[0].PostLoad(log, superseed)
	if err != nil {
		t.Errorf("PostLoad should not returns an error, returned: %v", err)
	}
	superseed["migration.dir"] = "dummy"
	err = steps[0].PostLoad(log, superseed)
	if err == nil {
		t.Errorf("PostLoad should returns an error")
	}
}
