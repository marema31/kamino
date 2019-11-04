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

func TestNoTableError(t *testing.T) {
	sdb, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	ddb, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	source := mockdatasource.MockDatasource{Type: datasource.Database, Engine: datasource.Mysql, Database: "source", MockedDb: sdb}
	dest := mockdatasource.MockDatasource{Type: datasource.Database, Engine: datasource.Mysql, Database: "dest", MockedDb: ddb}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	logger.SetLevel(logrus.PanicLevel)

	_, err = database.NewSaver(context.Background(), log, &dest, "", "id", "replace")
	if err == nil {
		t.Fatalf("NewSaver should return error")
	}

	_, err = database.NewLoader(context.Background(), log, &source, "", "")
	if err == nil {
		t.Fatalf("NewLoader should return error")
	}

}

func TestNoKeyError(t *testing.T) {
	ddb, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	dest := mockdatasource.MockDatasource{Type: datasource.Database, Engine: datasource.Mysql, Database: "dest", MockedDb: ddb}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	logger.SetLevel(logrus.PanicLevel)

	_, err = database.NewSaver(context.Background(), log, &dest, "dtable", "", "exactCopy")
	if err == nil {
		t.Fatalf("NewSaver should return error")
	}
}

func TestCreateIdsListError(t *testing.T) {
	ddb, dmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	dmock.ExpectQuery("SELECT id from dtable").WillReturnError(fmt.Errorf("fake error"))
	dest := mockdatasource.MockDatasource{MockedDb: ddb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog", Transaction: true}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	logger.SetLevel(logrus.PanicLevel)

	_, err = database.NewSaver(context.Background(), log, &dest, "dtable", "id", "exactCopy")
	if err == nil {
		t.Fatalf("NewSaver should return error")
	}
	if err := dmock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations on Saver: %s", err)
	}

}

func TestGetColNamesError(t *testing.T) {
	sdb, smock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "title", "body"}).
		AddRow(1, "post 1", "hello").
		AddRow(2, "post 2", "world")

	smock.ExpectQuery("SELECT (.+) from stable").WillReturnRows(rows)
	smock.ExpectClose()
	source := mockdatasource.MockDatasource{MockedDb: sdb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	logger.SetLevel(logrus.PanicLevel)

	ddb, dmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	dmock.ExpectQuery("SELECT \\* from dtable LIMIT 1").WillReturnError(fmt.Errorf("fake error"))

	dest := mockdatasource.MockDatasource{MockedDb: ddb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog"}

	saver, err := database.NewSaver(context.Background(), log, &dest, "dtable", "", "insert")
	if err != nil {
		t.Fatalf("NewSaver should not return error and returned '%v'", err)
	}
	loader, err := database.NewLoader(context.Background(), log, &source, "stable", "")
	if err != nil {
		t.Fatalf("NewLoader should not return error and returned '%v'", err)
	}

	for loader.Next() {
		record, err := loader.Load(log)
		if err != nil {
			t.Fatalf("Load should not return error and returned '%v'", err)
		}

		if err = saver.Save(log, record); err == nil {
			t.Fatalf("Save should return error")
		}
	}
	if err := dmock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations on Saver: %s", err)
	}

}

func TestDestHasMoreColsOk(t *testing.T) {
	sdb, smock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "title", "body"}).
		AddRow(1, "post 1", "hello").
		AddRow(2, "post 2", "world")

	smock.ExpectQuery("SELECT (.+) from stable").WillReturnRows(rows)
	smock.ExpectClose()
	source := mockdatasource.MockDatasource{MockedDb: sdb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	logger.SetLevel(logrus.PanicLevel)

	ddb, dmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows = sqlmock.NewRows([]string{"id", "title", "body", "summary"})
	dmock.ExpectQuery("SELECT \\* from dtable LIMIT 1").WillReturnRows(rows)
	dmock.ExpectPrepare("INSERT INTO dtable \\( title,body,id\\) VALUES \\( \\?,\\?,\\? \\)")
	dmock.ExpectExec("INSERT INTO dtable").WithArgs("post 1", "hello", "1").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectExec("INSERT INTO dtable").WithArgs("post 2", "world", "2").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectClose()
	dest := mockdatasource.MockDatasource{MockedDb: ddb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog"}

	saver, err := database.NewSaver(context.Background(), log, &dest, "dtable", "id", "insert")
	if err != nil {
		t.Fatalf("NewSaver should not return error and returned '%v'", err)
	}
	loader, err := database.NewLoader(context.Background(), log, &source, "stable", "")
	if err != nil {
		t.Fatalf("NewLoader should not return error and returned '%v'", err)
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
}

func TestDestHasLessColsOk(t *testing.T) {
	sdb, smock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "title", "body"}).
		AddRow(1, "post 1", "hello").
		AddRow(2, "post 2", "world")

	smock.ExpectQuery("SELECT (.+) from stable").WillReturnRows(rows)
	smock.ExpectClose()
	source := mockdatasource.MockDatasource{MockedDb: sdb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	logger.SetLevel(logrus.PanicLevel)

	ddb, dmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows = sqlmock.NewRows([]string{"id", "title"})
	dmock.ExpectQuery("SELECT \\* from dtable LIMIT 1").WillReturnRows(rows)
	dmock.ExpectPrepare("INSERT INTO dtable \\( title,id\\) VALUES \\( \\?,\\? \\)")
	dmock.ExpectExec("INSERT INTO dtable").WithArgs("post 1", "1").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectExec("INSERT INTO dtable").WithArgs("post 2", "2").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectClose()
	dest := mockdatasource.MockDatasource{MockedDb: ddb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog"}

	saver, err := database.NewSaver(context.Background(), log, &dest, "dtable", "id", "insert")
	if err != nil {
		t.Fatalf("NewSaver should not return error and returned '%v'", err)
	}
	loader, err := database.NewLoader(context.Background(), log, &source, "stable", "")
	if err != nil {
		t.Fatalf("NewLoader should not return error and returned '%v'", err)
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
}
func TestKeyNotExistsDestError(t *testing.T) {
	ddb, dmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	dmock.ExpectQuery("SELECT id from dtable").WillReturnError(fmt.Errorf("fake error"))
	rows := sqlmock.NewRows([]string{"title", "body"})
	dmock.ExpectQuery("SELECT \\* from dtable LIMIT 1").WillReturnRows(rows)
	dmock.ExpectPrepare("INSERT INTO dtable \\( title,body,id\\) VALUES \\( \\?,\\?,\\? \\)")
	dest := mockdatasource.MockDatasource{MockedDb: ddb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	logger.SetLevel(logrus.PanicLevel)

	_, err = database.NewSaver(context.Background(), log, &dest, "dtable", "id", "update")
	if err == nil {
		t.Fatalf("NewSaver should return error")
	}
}

func TestKeyNotExistsSourceError(t *testing.T) {
	sdb, smock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"title", "body"}).
		AddRow("post 1", "hello").
		AddRow("post 2", "world")

	smock.ExpectQuery("SELECT (.+) from stable").WillReturnRows(rows)
	smock.ExpectClose()
	source := mockdatasource.MockDatasource{MockedDb: sdb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog"}

	ddb, dmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows = sqlmock.NewRows([]string{"id"}).
		AddRow(3).
		AddRow(2)
	dmock.ExpectQuery("SELECT id from dtable").WillReturnRows(rows)
	rows = sqlmock.NewRows([]string{"id", "title", "body"})
	dmock.ExpectQuery("SELECT \\* from dtable LIMIT 1").WillReturnRows(rows)
	dmock.ExpectPrepare("INSERT INTO dtable \\( title,body,id\\) VALUES \\( \\?,\\?,\\? \\)")
	dmock.ExpectExec("INSERT INTO dtable").WithArgs("post 1", "hello", "1").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectExec("INSERT INTO dtable").WithArgs("post 2", "world", "2").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectClose()
	dest := mockdatasource.MockDatasource{MockedDb: ddb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	logger.SetLevel(logrus.PanicLevel)

	saver, err := database.NewSaver(context.Background(), log, &dest, "dtable", "id", "update")
	if err != nil {
		t.Fatalf("NewSaver should not return error and returned '%v'", err)
	}
	loader, err := database.NewLoader(context.Background(), log, &source, "stable", "")
	if err != nil {
		t.Fatalf("NewLoader should not return error and returned '%v'", err)
	}

	loader.Next()
	record, err := loader.Load(log)
	if err != nil {
		t.Fatalf("Load should not return error and returned '%v'", err)
	}

	err = saver.Save(log, record)
	if err == nil {
		t.Fatalf("Save should return error")
	}
}
