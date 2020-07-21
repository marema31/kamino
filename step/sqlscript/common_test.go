package sqlscript_test

import (
	"context"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/mockdatasource"
	"github.com/spf13/viper"
)

func setupDo(path string, filename string) (context.Context, *logrus.Entry, datasource.Datasourcers, *viper.Viper, sqlmock.Sqlmock, error) {
	dss := mockdatasource.New()
	ds1 := mockdatasource.MockDatasource{Name: "ds1", Database: "db1", User: "user1", UserPw: "security4ever", Transaction: false, Tags: []string{"tag1"}}
	ds3 := mockdatasource.MockDatasource{Name: "ds3", Database: "db1", User: "user1", UserPw: "security4ever", Transaction: true, Tags: []string{"tag3"}}

	dss.Insert(true, []string{"tag1"}, []datasource.Type{datasource.Database}, []datasource.Engine{datasource.Mysql}, []*mockdatasource.MockDatasource{&ds1})
	dss.Insert(true, []string{"tag3"}, []datasource.Type{datasource.Database}, []datasource.Engine{datasource.Mysql}, []*mockdatasource.MockDatasource{&ds3})
	v := viper.New()
	v.SetConfigName(filename)
	v.AddConfigPath(path)
	v.ReadInConfig()
	ctx := context.Background()
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	logger.SetLevel(logrus.PanicLevel)
	db, mock, err := sqlmock.New()
	ds1.MockedDb = db
	ds3.MockedDb = db

	return ctx, log, dss, v, mock, err
}
