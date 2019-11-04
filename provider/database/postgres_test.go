package database_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/mockdatasource"
	"github.com/marema31/kamino/provider/database"
)

func TestPostgresOk(t *testing.T) {
	sdb, smock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "title", "body"}).
		AddRow(1, "post 1", "hello").
		AddRow(2, "post 2", "world")

	smock.ExpectQuery("SELECT (.+) from stable").WillReturnRows(rows)
	smock.ExpectClose()
	source := mockdatasource.MockDatasource{MockedDb: sdb, Type: datasource.Database, Engine: datasource.Postgres, Database: "blog"}

	ddb, dmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows = sqlmock.NewRows([]string{"id", "title", "body"})
	dmock.ExpectQuery("SELECT \\* from dtable LIMIT 1").WillReturnRows(rows)
	dmock.ExpectPrepare("INSERT INTO dtable \\( title,body,id\\) VALUES \\( \\$1,\\$2,\\$3 \\)")
	dmock.ExpectExec("INSERT INTO dtable").WithArgs("post 1", "hello", "1").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectExec("INSERT INTO dtable").WithArgs("post 2", "world", "2").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectClose()
	dest := mockdatasource.MockDatasource{MockedDb: ddb, Type: datasource.Database, Engine: datasource.Postgres, Database: "blog"}
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

	err = saver.Close(log)
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

func TestPostgresTransactionOk(t *testing.T) {
	sdb, smock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "title", "body"}).
		AddRow(1, "post 1", "hello").
		AddRow(2, "post 2", "world")

	smock.ExpectQuery("SELECT (.+) from stable").WillReturnRows(rows)
	smock.ExpectClose()
	source := mockdatasource.MockDatasource{MockedDb: sdb, Type: datasource.Database, Engine: datasource.Postgres, Database: "blog", Transaction: true}

	ddb, dmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows = sqlmock.NewRows([]string{"id", "title", "body"})
	dmock.ExpectBegin()
	dmock.ExpectQuery("SELECT \\* from dtable LIMIT 1").WillReturnRows(rows)
	dmock.ExpectPrepare("INSERT INTO dtable \\( title,body,id\\) VALUES \\( \\$1,\\$2,\\$3 \\)")
	dmock.ExpectExec("INSERT INTO dtable").WithArgs("post 1", "hello", "1").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectExec("INSERT INTO dtable").WithArgs("post 2", "world", "2").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectCommit()
	dmock.ExpectClose()
	dmock.ExpectClose() // sqlmock and close/defer have erratic behavior
	dest := mockdatasource.MockDatasource{MockedDb: ddb, Type: datasource.Database, Engine: datasource.Postgres, Database: "blog", Transaction: true}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	logger.SetLevel(logrus.PanicLevel)

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

	err = saver.Close(log)
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

func TestPostgresSchemaOk(t *testing.T) {
	sdb, smock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "title", "body"}).
		AddRow(1, "post 1", "hello").
		AddRow(2, "post 2", "world")

	smock.ExpectQuery("SELECT (.+) from greatbob.stable").WillReturnRows(rows)
	smock.ExpectClose()
	source := mockdatasource.MockDatasource{MockedDb: sdb, Type: datasource.Database, Engine: datasource.Postgres, Database: "blog", Schema: "greatbob"}

	ddb, dmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows = sqlmock.NewRows([]string{"id", "title", "body"})
	dmock.ExpectQuery("SELECT \\* from greatbob.dtable LIMIT 1").WillReturnRows(rows)
	dmock.ExpectPrepare("INSERT INTO greatbob.dtable \\( title,body,id\\) VALUES \\( \\$1,\\$2,\\$3 \\)")
	dmock.ExpectExec("INSERT INTO greatbob.dtable").WithArgs("post 1", "hello", "1").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectExec("INSERT INTO greatbob.dtable").WithArgs("post 2", "world", "2").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectClose()
	dest := mockdatasource.MockDatasource{MockedDb: ddb, Type: datasource.Database, Engine: datasource.Postgres, Database: "blog", Schema: "greatbob"}
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
	if lname != "blog_greatbob.stable" {
		t.Errorf("Loader name function does not return the correct name %s", lname)
	}
	sname := saver.Name()
	if sname != "blog_greatbob.dtable" {
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

	err = saver.Close(log)
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
