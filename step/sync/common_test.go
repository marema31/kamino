package sync_test

import (
	"context"
	"fmt"

	"github.com/marema31/kamino/mockprovider"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/mockdatasource"
)

func setupDo(path string, filename string) (context.Context, *logrus.Entry, datasource.Datasourcers, *viper.Viper, *mockprovider.MockProvider) {
	dss := mockdatasource.New()
	ds1 := mockdatasource.MockDatasource{Name: "ds1", Database: "db1", User: "user1", Tags: []string{"tag1", "tag2"}}
	ds2 := mockdatasource.MockDatasource{Name: "ds2", Database: "db2", User: "user2", Tags: []string{"tag2"}}
	ds3 := mockdatasource.MockDatasource{Name: "ds3", Database: "db3", User: "user3", Tags: []string{"tag3"}}
	ds4 := mockdatasource.MockDatasource{Name: "ds4", Database: "db4", User: "user4", Tags: []string{"tag3"}}
	ds5 := mockdatasource.MockDatasource{Name: "ds5", Database: "db4", User: "user4", Tags: []string{"tagerror"}, ErrorOpenDb: fmt.Errorf("fake error")}
	dscache := mockdatasource.MockDatasource{Name: "dscache", FilePath: "db4", Tags: []string{"tagcache"}}
	dscachenotexist := mockdatasource.MockDatasource{Name: "dscachenotexist", FilePath: "db4", Tags: []string{"tagcachenotexist"}, FileNotExists: true}
	dserrorfile := mockdatasource.MockDatasource{Name: "dserror", FilePath: "db4", Tags: []string{"tagcache"}, ErrorOpenFile: fmt.Errorf("fake error")}

	dss.Insert(true, []string{"tag1", "tag2"}, []datasource.Type{datasource.Database}, []datasource.Engine{datasource.Mysql}, []*mockdatasource.MockDatasource{&ds2})
	dss.Insert(true, []string{"tag3"}, []datasource.Type{datasource.Database}, []datasource.Engine{datasource.Mysql}, []*mockdatasource.MockDatasource{&ds3, &ds4})
	dss.Insert(true, []string{"tagsource"}, []datasource.Type{datasource.Database}, []datasource.Engine{datasource.Mysql}, []*mockdatasource.MockDatasource{&ds1})
	dss.Insert(true, []string{"tagcache"}, []datasource.Type{datasource.File}, []datasource.Engine{datasource.JSON}, []*mockdatasource.MockDatasource{&dscache})
	dss.Insert(true, []string{"tagcachenotexist"}, []datasource.Type{datasource.File}, []datasource.Engine{datasource.JSON}, []*mockdatasource.MockDatasource{&dscachenotexist})
	dss.Insert(true, []string{"tagerror"}, []datasource.Type{datasource.Database}, []datasource.Engine{datasource.Mysql}, []*mockdatasource.MockDatasource{&ds5})
	dss.Insert(true, []string{"tagerrorfile"}, []datasource.Type{datasource.File}, []datasource.Engine{datasource.JSON}, []*mockdatasource.MockDatasource{&dserrorfile})
	dss.Insert(true, []string{""}, []datasource.Type{datasource.Database}, []datasource.Engine{datasource.Mysql}, []*mockdatasource.MockDatasource{&ds2})
	v := viper.New()
	v.SetConfigName(filename)
	v.AddConfigPath(path)
	v.ReadInConfig()
	ctx := context.Background()
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	logger.SetLevel(logrus.PanicLevel)

	return ctx, log, dss, v, &mockprovider.MockProvider{}
}
