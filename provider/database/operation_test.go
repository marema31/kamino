package database_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/mockdatasource"
	"github.com/marema31/kamino/provider/database"
)

func TestOpenError(t *testing.T) {
	sdb, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	ddb, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	source := mockdatasource.MockDatasource{ErrorOpenDb: fmt.Errorf("fake error"), Type: datasource.Database, Engine: datasource.Mysql, Database: "source", MockedDb: sdb}
	dest := mockdatasource.MockDatasource{ErrorOpenDb: fmt.Errorf("fake error"), Type: datasource.Database, Engine: datasource.Mysql, Database: "dest", MockedDb: ddb}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")

	_, err = database.NewSaver(context.Background(), log, &dest, "dtable", "id", "replace")
	if err == nil {
		t.Fatalf("NewSaver should return error")
	}

	_, err = database.NewLoader(context.Background(), log, &source, "stable", "")
	if err == nil {
		t.Fatalf("NewLoader should return error")
	}

}

func TestResetOk(t *testing.T) {
	sdb, smock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "title", "body"}).
		AddRow(1, "post 1", "hello").
		AddRow(2, "post 2", "world")

	smock.ExpectQuery("SELECT (.+) from stable WHERE title like '%'").WillReturnRows(rows)
	smock.ExpectClose()
	source := mockdatasource.MockDatasource{MockedDb: sdb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog"}

	ddb, dmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows = sqlmock.NewRows([]string{"id", "title", "body"})
	dmock.ExpectQuery("SELECT \\* from dtable LIMIT 1").WillReturnRows(rows)
	dmock.ExpectPrepare("INSERT INTO dtable \\( title,body,id\\) VALUES \\( \\?,\\?,\\? \\)")
	dmock.ExpectExec("INSERT INTO dtable").WithArgs("post 1", "hello", "1").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectExec("INSERT INTO dtable").WithArgs("post 2", "world", "2").WillReturnResult(sqlmock.NewResult(1, 1))
	dest := mockdatasource.MockDatasource{MockedDb: ddb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")

	saver, err := database.NewSaver(context.Background(), log, &dest, "dtable", "id", "insert")
	if err != nil {
		t.Fatalf("NewSaver should not return error and returned '%v'", err)
	}
	loader, err := database.NewLoader(context.Background(), log, &source, "stable", "title like '%'")
	if err != nil {
		t.Fatalf("NewLoader should not return error and returned '%v'", err)
	}

	lname := loader.Name()
	if lname != "blog_stable" {
		t.Errorf("Loader name function does not return the correct name %s", lname)
	}
	sname := saver.Name()
	if sname != "blog_dtable" {
		t.Errorf("Saver name function does not return the correct name %s", sname)
	}

	for loader.Next() {
		record, err := loader.Load(log)
		if err != nil {
			t.Fatalf("Load should not return error and returned '%v'", err)
		}

		if err = saver.Save(log, record); err != nil {
			t.Fatalf("Save should not return error and returned '%v'", err)
		}
	}

	_, err = loader.Load(log)
	if err == nil {
		t.Errorf("Load should return error ")
	}

	err = saver.Reset(log)
	if err != nil {
		t.Errorf("Saver close should not return error and returned '%v'", err)
	}

	err = loader.Close(log)
	if err != nil {
		t.Errorf("Loader close should not return error and returned '%v'", err)
	}

	if err := smock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations on Loader: %s", err)
	}
	if err := dmock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations on Saver: %s", err)
	}

}

func TestResetTransactionOk(t *testing.T) {
	sdb, smock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "title", "body"}).
		AddRow(1, "post 1", "hello").
		AddRow(2, "post 2", "world")

	smock.ExpectQuery("SELECT (.+) from stable").WillReturnRows(rows)
	smock.ExpectClose()
	source := mockdatasource.MockDatasource{MockedDb: sdb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog", Transaction: true}

	ddb, dmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows = sqlmock.NewRows([]string{"id", "title", "body"})
	dmock.ExpectBegin()
	dmock.ExpectQuery("SELECT \\* from dtable LIMIT 1").WillReturnRows(rows)
	dmock.ExpectPrepare("INSERT INTO dtable \\( title,body,id\\) VALUES \\( \\?,\\?,\\? \\)")
	dmock.ExpectExec("INSERT INTO dtable").WithArgs("post 1", "hello", "1").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectExec("INSERT INTO dtable").WithArgs("post 2", "world", "2").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectRollback()
	dest := mockdatasource.MockDatasource{MockedDb: ddb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog", Transaction: true}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")

	saver, err := database.NewSaver(context.Background(), log, &dest, "dtable", "id", "insert")
	if err != nil {
		t.Fatalf("NewSaver should not return error and returned '%v'", err)
	}
	loader, err := database.NewLoader(context.Background(), log, &source, "stable", "")
	if err != nil {
		t.Fatalf("NewLoader should not return error and returned '%v'", err)
	}

	lname := loader.Name()
	if lname != "blog_stable" {
		t.Errorf("Loader name function does not return the correct name %s", lname)
	}
	sname := saver.Name()
	if sname != "blog_dtable" {
		t.Errorf("Saver name function does not return the correct name %s", sname)
	}

	for loader.Next() {
		record, err := loader.Load(log)
		if err != nil {
			t.Fatalf("Load should not return error and returned '%v'", err)
		}

		if err = saver.Save(log, record); err != nil {
			t.Fatalf("Save should not return error and returned '%v'", err)
		}
	}

	_, err = loader.Load(log)
	if err == nil {
		t.Errorf("Load should return error ")
	}

	err = saver.Reset(log)
	if err != nil {
		t.Errorf("Saver close should not return error and returned '%v'", err)
	}

	err = loader.Close(log)
	if err != nil {
		t.Errorf("Loader close should not return error and returned '%v'", err)
	}

	if err := smock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations on Loader: %s", err)
	}
	if err := dmock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations on Saver: %s", err)
	}

}
