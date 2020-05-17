package common_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/mockdatasource"
	"github.com/marema31/kamino/step/common"
)

func TestToSkipYesOk(t *testing.T) {
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	ds := &mockdatasource.MockDatasource{Name: "ds1", Database: "db1", User: "user1"}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	ds.MockedDb = db
	rows := sqlmock.NewRows([]string{"COUNT(id)"}).
		AddRow(10)

	mock.ExpectQuery("SELECT COUNT\\(id\\) from dtable WHERE title like '%'").WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"COUNT(id)"}).
		AddRow(1)
	mock.ExpectQuery("SELECT COUNT\\(id\\) from stable WHERE title like '%'").WillReturnRows(rows)

	ok, err := common.ToSkipDatabase(context.Background(), log, ds, false, false, []string{"SELECT COUNT(id) from dtable WHERE title like '%'", "SELECT COUNT(id) from stable WHERE title like '%'"})
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

func TestToSkipNoOk(t *testing.T) {
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	ds := &mockdatasource.MockDatasource{Name: "ds1", Database: "db1", User: "user1"}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	ds.MockedDb = db
	rows := sqlmock.NewRows([]string{"COUNT(id)"}).
		AddRow(0)

	mock.ExpectQuery("SELECT COUNT\\(id\\) from stable WHERE title like '%'").WillReturnRows(rows)
	if err != nil {
		t.Fatalf("NewSaver should not return error and returned '%v'", err)
	}

	ok, err := common.ToSkipDatabase(context.Background(), log, ds, false, false, []string{"SELECT COUNT(id) from stable WHERE title like '%'", "SELECT COUNT(id) from dtable WHERE title like '%'"})
	if err != nil {
		t.Errorf("ToSkip should not return error, returned: %v", err)
	}

	if ok {
		t.Error("ToSkip should return false")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func TestToSkipError(t *testing.T) {
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	ds := &mockdatasource.MockDatasource{Name: "ds1", Database: "db1", User: "user1"}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	ds.MockedDb = db
	mock.ExpectQuery("SELECT COUNT\\(id\\) from stable WHERE title like '%'").WillReturnError(fmt.Errorf("fake error"))
	if err != nil {
		t.Fatalf("NewSaver should not return error and returned '%v'", err)
	}

	_, err = common.ToSkipDatabase(context.Background(), log, ds, false, false, []string{"SELECT COUNT(id) from stable WHERE title like '%'"})
	if err == nil {
		t.Errorf("ToSkip should return error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func TestToSkipOpenError(t *testing.T) {
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	ds := &mockdatasource.MockDatasource{Name: "ds1", Database: "db1", User: "user1"}

	ds.ErrorOpenDb = fmt.Errorf("Fake OpenDatabase error")

	_, err := common.ToSkipDatabase(context.Background(), log, ds, false, false, []string{"SELECT COUNT(id) from stable WHERE title like '%'"})
	if err == nil {
		t.Errorf("ToSkip should return error")
	}
}
