package mockdatasource_test

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/mockdatasource"
)

func TestDatabaseOpenMysql(t *testing.T) {
	ds := mockdatasource.MockDatasource{URL: "bob:123soleil@tcp(localhost:1234)/dbmc", URLAdmin: "URLAdmin", URLNoDb: "URLNoDb"}
	db, _, err := sqlmock.New()
	if err != nil {
		t.Errorf("sqlmock.New should not returns an error, was: %v", err)
	}
	ds.MockedDb = db
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")

	if _, err = ds.OpenDatabase(log, false, false); err != nil {
		t.Errorf("OpenDatabase should not returns an error, was: %v", err)
	}

	ds.ErrorOpenDb = fmt.Errorf("Fake OpenDatabase error")

	if _, err = ds.OpenDatabase(log, true, false); err == nil {
		t.Errorf("OpenDatabase should returns an error")
	}
	if ds.IsTransaction() {
		t.Errorf("The datasource should not have transaction")
	}

}
