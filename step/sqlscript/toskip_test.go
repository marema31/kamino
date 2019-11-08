package sqlscript_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/mockdatasource"
	"github.com/marema31/kamino/step/sqlscript"
	"github.com/spf13/viper"
)

func setupDo(path string, filename string) (context.Context, *logrus.Entry, datasource.Datasourcers, *viper.Viper, sqlmock.Sqlmock, error) {
	dss := mockdatasource.New()
	ds1 := mockdatasource.MockDatasource{Name: "ds1", Database: "db1", User: "user1", Tags: []string{"tag1"}}

	dss.Insert([]string{"tag1", "tag2"}, []datasource.Type{datasource.Database}, []datasource.Engine{datasource.Mysql}, []*mockdatasource.MockDatasource{&ds1})
	v := viper.New()
	v.SetConfigName(filename) // The file will be named [filename].json, [filename].yaml or [filename.toml]
	v.AddConfigPath(path)
	v.ReadInConfig()
	ctx := context.Background()
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	logger.SetLevel(logrus.PanicLevel)
	db, mock, err := sqlmock.New()
	ds1.MockedDb = db

	return ctx, log, dss, v, mock, err
}
func TestToSkipYesOk(t *testing.T) {
	ctx, log, dss, v, mock, err := setupDo("testdata/good/steps/", "sqlscriptok")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}
	rows := sqlmock.NewRows([]string{"COUNT(id)"}).
		AddRow(10)

	mock.ExpectQuery("SELECT \\* FROM USER WHERE user = 'user1';").WillReturnRows(rows)
	if err != nil {
		t.Fatalf("NewSaver should not return error and returned '%v'", err)
	}

	_, steps, err := sqlscript.Load(ctx, log, "testdata/good", "sqlscriptok", v, dss)
	if err != nil {
		t.Errorf("Load should not returns an error, returned: %v", err)
	}

	ok, err := steps[0].ToSkip(context.Background(), log)
	if err != nil {
		t.Errorf("ToSkip should not return error, returned: %v", err)
	}

	if !ok {
		t.Error("ToSkip should return true")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}
