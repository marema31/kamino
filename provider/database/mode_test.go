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

func TestOnlyIfEmptyOk(t *testing.T) {
	sdb, smock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "title", "body"}).
		AddRow(1, "post 1", "hello").
		AddRow(2, "post 2", "world")

	smock.ExpectQuery("SELECT (.+) from stable").WillReturnRows(rows)
	source := mockdatasource.MockDatasource{MockedDb: sdb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog"}

	ddbfull, dmockfull, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows = sqlmock.NewRows([]string{"name"}).
		AddRow("id").
		AddRow("title").
		AddRow("body")
	dmockfull.ExpectQuery("SELECT column_name AS name FROM information_schema.columns WHERE table_schema = 'blog' AND table_name ='dtable';").WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"count"}).
		AddRow(3)
	dmockfull.ExpectQuery("SELECT COUNT\\(\\*\\) FROM dtable").WillReturnRows(rows)
	dmockfull.ExpectPrepare("INSERT INTO dtable \\( title,body,id\\) VALUES \\( \\?,\\?,\\? \\)")
	destfull := mockdatasource.MockDatasource{MockedDb: ddbfull, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	logger.SetLevel(logrus.PanicLevel)

	saverfull, err := database.NewSaver(context.Background(), log, &destfull, "dtable", "id", "onlyIfEmpty")
	if err != nil {
		t.Fatalf("NewSaver should not return error and returned '%v'", err)
	}
	ddbempty, dmockempty, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows = sqlmock.NewRows([]string{"name"}).
		AddRow("id").
		AddRow("title").
		AddRow("body")
	dmockempty.ExpectQuery("SELECT column_name AS name FROM information_schema.columns WHERE table_schema = 'blog' AND table_name ='dtable';").WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"count"}).
		AddRow(0)
	dmockempty.ExpectQuery("SELECT COUNT\\(\\*\\) FROM dtable").WillReturnRows(rows)
	dmockempty.ExpectPrepare("INSERT INTO dtable \\( title,body,id\\) VALUES \\( \\?,\\?,\\? \\)")
	dmockempty.ExpectExec("INSERT INTO dtable").WithArgs("post 1", "hello", "1").WillReturnResult(sqlmock.NewResult(1, 1))
	dmockempty.ExpectExec("INSERT INTO dtable").WithArgs("post 2", "world", "2").WillReturnResult(sqlmock.NewResult(1, 1))
	destempty := mockdatasource.MockDatasource{MockedDb: ddbempty, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog"}

	saverempty, err := database.NewSaver(context.Background(), log, &destempty, "dtable", "id", "onlyIfEmpty")
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

		if err = saverempty.Save(log, record); err != nil {
			t.Fatalf("Save should not return error and returned '%v'", err)
		}
		if err = saverfull.Save(log, record); err != nil {
			t.Fatalf("Save should not return error and returned '%v'", err)
		}
	}

	err = saverfull.Close(log)
	if err != nil {
		t.Errorf("Saver close should not return error and returned '%v'", err)
	}

	err = saverempty.Close(log)
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
	if err := dmockempty.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations on Saver: %s", err)
	}
	if err := dmockfull.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations on Saver: %s", err)
	}

}

func TestInsertOk(t *testing.T) {
	sdb, smock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "title", "body"}).
		AddRow(1, "post 1", "hello").
		AddRow(2, "post 2", "world")

	smock.ExpectQuery("SELECT (.+) from stable").WillReturnRows(rows)
	source := mockdatasource.MockDatasource{MockedDb: sdb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog"}

	ddb, dmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows = sqlmock.NewRows([]string{"name"}).
		AddRow("id").
		AddRow("title").
		AddRow("body")
	dmock.ExpectQuery("SELECT column_name AS name FROM information_schema.columns WHERE table_schema = 'blog' AND table_name ='dtable';").WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"count"}).
		AddRow(0)
	dmock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM dtable").WillReturnRows(rows)
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

func TestTruncateOk(t *testing.T) {
	sdb, smock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "title", "body"}).
		AddRow(1, "post 1", "hello").
		AddRow(2, "post 2", "world")

	smock.ExpectQuery("SELECT (.+) from stable").WillReturnRows(rows)
	source := mockdatasource.MockDatasource{MockedDb: sdb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog"}

	ddb, dmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows = sqlmock.NewRows([]string{"id", "title", "body"})
	rows = sqlmock.NewRows([]string{"name"}).
		AddRow("id").
		AddRow("title").
		AddRow("body")
	dmock.ExpectQuery("SELECT column_name AS name FROM information_schema.columns WHERE table_schema = 'blog' AND table_name ='dtable';").WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"count"}).
		AddRow(0)
	dmock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM dtable").WillReturnRows(rows)
	dmock.ExpectPrepare("INSERT INTO dtable \\( title,body,id\\) VALUES \\( \\?,\\?,\\? \\)")
	dmock.ExpectExec("RUNCATE TABLE dtable").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectExec("INSERT INTO dtable").WithArgs("post 1", "hello", "1").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectExec("INSERT INTO dtable").WithArgs("post 2", "world", "2").WillReturnResult(sqlmock.NewResult(1, 1))
	dest := mockdatasource.MockDatasource{MockedDb: ddb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")

	saver, err := database.NewSaver(context.Background(), log, &dest, "dtable", "id", "truncate")
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

func TestTruncateTransactionOk(t *testing.T) {
	sdb, smock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "title", "body"}).
		AddRow(1, "post 1", "hello").
		AddRow(2, "post 2", "world")

	smock.ExpectQuery("SELECT (.+) from stable").WillReturnRows(rows)
	source := mockdatasource.MockDatasource{MockedDb: sdb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog"}

	ddb, dmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows = sqlmock.NewRows([]string{"id", "title", "body"})
	dmock.ExpectBegin()
	rows = sqlmock.NewRows([]string{"name"}).
		AddRow("id").
		AddRow("title").
		AddRow("body")
	dmock.ExpectQuery("SELECT column_name AS name FROM information_schema.columns WHERE table_schema = 'blog' AND table_name ='dtable';").WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"count"}).
		AddRow(0)
	dmock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM dtable").WillReturnRows(rows)
	dmock.ExpectPrepare("INSERT INTO dtable \\( title,body,id\\) VALUES \\( \\?,\\?,\\? \\)")
	dmock.ExpectExec("TRUNCATE TABLE dtable").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectExec("INSERT INTO dtable").WithArgs("post 1", "hello", "1").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectExec("INSERT INTO dtable").WithArgs("post 2", "world", "2").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectCommit()
	dest := mockdatasource.MockDatasource{MockedDb: ddb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog", Transaction: true}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")

	saver, err := database.NewSaver(context.Background(), log, &dest, "dtable", "id", "truncate")
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

func TestUpdateOk(t *testing.T) {
	sdb, smock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "title", "body"}).
		AddRow(1, "post 1", "hello").
		AddRow(2, "post 2", "world")

	smock.ExpectQuery("SELECT (.+) from stable").WillReturnRows(rows)
	source := mockdatasource.MockDatasource{MockedDb: sdb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog"}

	ddb, dmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows = sqlmock.NewRows([]string{"id"}).
		AddRow(3).
		AddRow(2)
	dmock.ExpectQuery("SELECT \\? from \\?").WithArgs("id", "dtable").WillReturnRows(rows)
	rows = sqlmock.NewRows([]string{"id", "title", "body"}).
		AddRow(3, "post 3", "good").
		AddRow(2, "post 2 bis", "planet")
	rows = sqlmock.NewRows([]string{"name"}).
		AddRow("id").
		AddRow("title").
		AddRow("body")
	dmock.ExpectQuery("SELECT column_name AS name FROM information_schema.columns WHERE table_schema = 'blog' AND table_name ='dtable';").WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"count"}).
		AddRow(0)
	dmock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM dtable").WillReturnRows(rows)
	dmock.ExpectPrepare("INSERT INTO dtable \\( title,body,id\\) VALUES \\( \\?,\\?,\\? \\)")
	dmock.ExpectPrepare("UPDATE dtable SET  title=\\?,body=\\? WHERE id = \\?")
	dmock.ExpectExec("UPDATE dtable SET  title=\\?,body=\\? WHERE id = \\?").WithArgs("post 1", "hello", "1").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectExec("UPDATE dtable SET  title=\\?,body=\\? WHERE id = \\?").WithArgs("post 2", "world", "2").WillReturnResult(sqlmock.NewResult(1, 1))
	dest := mockdatasource.MockDatasource{MockedDb: ddb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")

	saver, err := database.NewSaver(context.Background(), log, &dest, "dtable", "id", "update")
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

func TestReplaceOk(t *testing.T) {
	sdb, smock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "title", "body"}).
		AddRow(1, "post 1", "hello").
		AddRow(2, "post 2", "world")

	smock.ExpectQuery("SELECT (.+) from stable").WillReturnRows(rows)
	source := mockdatasource.MockDatasource{MockedDb: sdb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog"}

	ddb, dmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	rows = sqlmock.NewRows([]string{"id"}).
		AddRow(3).
		AddRow(2)
	dmock.ExpectQuery("SELECT \\? from \\?").WithArgs("id", "dtable").WillReturnRows(rows)
	rows = sqlmock.NewRows([]string{"id", "title", "body"}).
		AddRow(3, "post 3", "good").
		AddRow(2, "post 2 bis", "planet")
	rows = sqlmock.NewRows([]string{"name"}).
		AddRow("id").
		AddRow("title").
		AddRow("body")
	dmock.ExpectQuery("SELECT column_name AS name FROM information_schema.columns WHERE table_schema = 'blog' AND table_name ='dtable';").WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"count"}).
		AddRow(0)
	dmock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM dtable").WillReturnRows(rows)
	dmock.ExpectPrepare("INSERT INTO dtable \\( title,body,id\\) VALUES \\( \\?,\\?,\\? \\)")
	dmock.ExpectPrepare("UPDATE dtable SET  title=\\?,body=\\? WHERE id = \\?")
	dmock.ExpectExec("INSERT INTO dtable").WithArgs("post 1", "hello", "1").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectExec("UPDATE dtable SET  title=\\?,body=\\? WHERE id = \\?").WithArgs("post 2", "world", "2").WillReturnResult(sqlmock.NewResult(1, 1))
	dest := mockdatasource.MockDatasource{MockedDb: ddb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")

	saver, err := database.NewSaver(context.Background(), log, &dest, "dtable", "id", "replace")
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

func TestExactCopyOk(t *testing.T) {
	sdb, smock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "title", "body"}).
		AddRow(1, "post 1", "hello").
		AddRow(2, "post 2", "world")

	smock.ExpectQuery("SELECT (.+) from stable").WillReturnRows(rows)
	source := mockdatasource.MockDatasource{MockedDb: sdb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog"}

	ddb, dmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows = sqlmock.NewRows([]string{"id"}).
		AddRow(3).
		AddRow(2)
	dmock.ExpectQuery("SELECT \\? from \\?").WithArgs("id", "dtable").WillReturnRows(rows)
	rows = sqlmock.NewRows([]string{"id", "title", "body"}).
		AddRow(3, "post 3", "good").
		AddRow(2, "post 2 bis", "planet")
	rows = sqlmock.NewRows([]string{"name"}).
		AddRow("id").
		AddRow("title").
		AddRow("body")
	dmock.ExpectQuery("SELECT column_name AS name FROM information_schema.columns WHERE table_schema = 'blog' AND table_name ='dtable';").WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"count"}).
		AddRow(0)
	dmock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM dtable").WillReturnRows(rows)
	dmock.ExpectPrepare("INSERT INTO dtable \\( title,body,id\\) VALUES \\( \\?,\\?,\\? \\)")
	dmock.ExpectPrepare("UPDATE dtable SET  title=\\?,body=\\? WHERE id = \\?")
	dmock.ExpectExec("INSERT INTO dtable").WithArgs("post 1", "hello", "1").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectExec("UPDATE dtable SET  title=\\?,body=\\? WHERE id = \\?").WithArgs("post 2", "world", "2").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectExec("DELETE from \\? WHERE \\?=\\?").WithArgs("dtable", "id", "3").WillReturnResult(sqlmock.NewResult(1, 1))
	dest := mockdatasource.MockDatasource{MockedDb: ddb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")

	saver, err := database.NewSaver(context.Background(), log, &dest, "dtable", "id", "exactCopy")
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

func TestExactCopyTransactionOk(t *testing.T) {
	sdb, smock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "title", "body"}).
		AddRow(1, "post 1", "hello").
		AddRow(2, "post 2", "world")

	smock.ExpectQuery("SELECT (.+) from stable").WillReturnRows(rows)
	source := mockdatasource.MockDatasource{MockedDb: sdb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog"}

	ddb, dmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows = sqlmock.NewRows([]string{"id"}).
		AddRow(3).
		AddRow(2)
	dmock.ExpectQuery("SELECT \\? from \\?").WithArgs("id", "dtable").WillReturnRows(rows)
	rows = sqlmock.NewRows([]string{"id", "title", "body"}).
		AddRow(3, "post 3", "good").
		AddRow(2, "post 2 bis", "planet")
	dmock.ExpectBegin()
	rows = sqlmock.NewRows([]string{"name"}).
		AddRow("id").
		AddRow("title").
		AddRow("body")
	dmock.ExpectQuery("SELECT column_name AS name FROM information_schema.columns WHERE table_schema = 'blog' AND table_name ='dtable';").WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"count"}).
		AddRow(0)
	dmock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM dtable").WillReturnRows(rows)
	dmock.ExpectPrepare("INSERT INTO dtable \\( title,body,id\\) VALUES \\( \\?,\\?,\\? \\)")
	dmock.ExpectPrepare("UPDATE dtable SET  title=\\?,body=\\? WHERE id = \\?")
	dmock.ExpectExec("INSERT INTO dtable").WithArgs("post 1", "hello", "1").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectExec("UPDATE dtable SET  title=\\?,body=\\? WHERE id = \\?").WithArgs("post 2", "world", "2").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectExec("DELETE from \\? WHERE \\?=\\?").WithArgs("dtable", "id", "3").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectCommit()
	dest := mockdatasource.MockDatasource{MockedDb: ddb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog", Transaction: true}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")

	saver, err := database.NewSaver(context.Background(), log, &dest, "dtable", "id", "exactCopy")
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
