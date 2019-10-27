package database_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/mockdatasource"
	"github.com/marema31/kamino/provider/database"
)

func TestMySqlOk(t *testing.T) {
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
	dmock.ExpectClose()
	dest := mockdatasource.MockDatasource{MockedDb: ddb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog"}

	saver, err := database.NewSaver(context.Background(), &dest, "dtable", "id", "insert")
	if err != nil {
		t.Fatalf("NewSaver should not return error and returned '%v'", err)
	}
	loader, err := database.NewLoader(context.Background(), &source, "stable", "title like '%'")
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
		record, err := loader.Load()
		if err != nil {
			t.Fatalf("Load should not return error and returned '%v'", err)
		}

		if err = saver.Save(record); err != nil {
			t.Fatalf("Save should not return error and returned '%v'", err)
		}
	}

	_, err = loader.Load()
	if err == nil {
		t.Errorf("Load should return error ")
	}

	err = saver.Close()
	if err != nil {
		t.Errorf("Saver close should not return error and returned '%v'", err)
	}

	err = loader.Close()
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

func TestMySqlTransactionOk(t *testing.T) {
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
	dmock.ExpectCommit()
	dmock.ExpectClose()
	dest := mockdatasource.MockDatasource{MockedDb: ddb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog", Transaction: true}

	saver, err := database.NewSaver(context.Background(), &dest, "dtable", "id", "insert")
	if err != nil {
		t.Fatalf("NewSaver should not return error and returned '%v'", err)
	}
	loader, err := database.NewLoader(context.Background(), &source, "stable", "")
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
		record, err := loader.Load()
		if err != nil {
			t.Fatalf("Load should not return error and returned '%v'", err)
		}

		if err = saver.Save(record); err != nil {
			t.Fatalf("Save should not return error and returned '%v'", err)
		}
	}

	_, err = loader.Load()
	if err == nil {
		t.Errorf("Load should return error ")
	}

	err = saver.Close()
	if err != nil {
		t.Errorf("Saver close should not return error and returned '%v'", err)
	}

	err = loader.Close()
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