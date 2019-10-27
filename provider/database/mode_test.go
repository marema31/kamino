package database_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
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
	smock.ExpectClose()
	source := mockdatasource.MockDatasource{MockedDb: sdb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog"}

	ddbfull, dmockfull, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows = sqlmock.NewRows([]string{"id", "title", "body"}).
		AddRow(3, "post 3", "good").
		AddRow(2, "post 2", "world")
	dmockfull.ExpectQuery("SELECT \\* from dtable LIMIT 1").WillReturnRows(rows)
	dmockfull.ExpectPrepare("INSERT INTO dtable \\( title,body,id\\) VALUES \\( \\?,\\?,\\? \\)")
	dmockfull.ExpectClose()
	destfull := mockdatasource.MockDatasource{MockedDb: ddbfull, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog"}

	saverfull, err := database.NewSaver(context.Background(), &destfull, "dtable", "id", "onlyIfEmpty")
	if err != nil {
		t.Fatalf("NewSaver should not return error and returned '%v'", err)
	}
	ddbempty, dmockempty, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows = sqlmock.NewRows([]string{"id", "title", "body"})
	dmockempty.ExpectQuery("SELECT \\* from dtable LIMIT 1").WillReturnRows(rows)
	dmockempty.ExpectPrepare("INSERT INTO dtable \\( title,body,id\\) VALUES \\( \\?,\\?,\\? \\)")
	dmockempty.ExpectExec("INSERT INTO dtable").WithArgs("post 1", "hello", "1").WillReturnResult(sqlmock.NewResult(1, 1))
	dmockempty.ExpectExec("INSERT INTO dtable").WithArgs("post 2", "world", "2").WillReturnResult(sqlmock.NewResult(1, 1))
	dmockempty.ExpectClose()
	destempty := mockdatasource.MockDatasource{MockedDb: ddbempty, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog"}

	saverempty, err := database.NewSaver(context.Background(), &destempty, "dtable", "id", "onlyIfEmpty")
	if err != nil {
		t.Fatalf("NewSaver should not return error and returned '%v'", err)
	}
	loader, err := database.NewLoader(context.Background(), &source, "stable", "")
	if err != nil {
		t.Fatalf("NewLoader should not return error and returned '%v'", err)
	}

	for loader.Next() {
		record, err := loader.Load()
		if err != nil {
			t.Fatalf("Load should not return error and returned '%v'", err)
		}

		if err = saverempty.Save(record); err != nil {
			t.Fatalf("Save should not return error and returned '%v'", err)
		}
		if err = saverfull.Save(record); err != nil {
			t.Fatalf("Save should not return error and returned '%v'", err)
		}
	}

	err = saverfull.Close()
	if err != nil {
		t.Errorf("Saver close should not return error and returned '%v'", err)
	}

	err = saverempty.Close()
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
	loader, err := database.NewLoader(context.Background(), &source, "stable", "")
	if err != nil {
		t.Fatalf("NewLoader should not return error and returned '%v'", err)
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

func TestTruncateOk(t *testing.T) {
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

	ddb, dmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows = sqlmock.NewRows([]string{"id", "title", "body"})
	dmock.ExpectQuery("SELECT \\* from dtable LIMIT 1").WillReturnRows(rows)
	dmock.ExpectPrepare("INSERT INTO dtable \\( title,body,id\\) VALUES \\( \\?,\\?,\\? \\)")
	dmock.ExpectExec("RUNCATE TABLE dtable").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectExec("INSERT INTO dtable").WithArgs("post 1", "hello", "1").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectExec("INSERT INTO dtable").WithArgs("post 2", "world", "2").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectClose()
	dest := mockdatasource.MockDatasource{MockedDb: ddb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog"}

	saver, err := database.NewSaver(context.Background(), &dest, "dtable", "id", "truncate")
	if err != nil {
		t.Fatalf("NewSaver should not return error and returned '%v'", err)
	}
	loader, err := database.NewLoader(context.Background(), &source, "stable", "")
	if err != nil {
		t.Fatalf("NewLoader should not return error and returned '%v'", err)
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

func TestTruncateTransactionOk(t *testing.T) {
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

	ddb, dmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows = sqlmock.NewRows([]string{"id", "title", "body"})
	dmock.ExpectBegin()
	dmock.ExpectQuery("SELECT \\* from dtable LIMIT 1").WillReturnRows(rows)
	dmock.ExpectPrepare("INSERT INTO dtable \\( title,body,id\\) VALUES \\( \\?,\\?,\\? \\)")
	dmock.ExpectExec("RUNCATE TABLE dtable").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectExec("INSERT INTO dtable").WithArgs("post 1", "hello", "1").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectExec("INSERT INTO dtable").WithArgs("post 2", "world", "2").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectCommit()
	dmock.ExpectClose()
	dest := mockdatasource.MockDatasource{MockedDb: ddb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog", Transaction: true}

	saver, err := database.NewSaver(context.Background(), &dest, "dtable", "id", "truncate")
	if err != nil {
		t.Fatalf("NewSaver should not return error and returned '%v'", err)
	}
	loader, err := database.NewLoader(context.Background(), &source, "stable", "")
	if err != nil {
		t.Fatalf("NewLoader should not return error and returned '%v'", err)
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

func TestUpdateOk(t *testing.T) {
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

	ddb, dmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows = sqlmock.NewRows([]string{"id"}).
		AddRow(3).
		AddRow(2)
	dmock.ExpectQuery("SELECT id from dtable").WillReturnRows(rows)
	rows = sqlmock.NewRows([]string{"id", "title", "body"}).
		AddRow(3, "post 3", "good").
		AddRow(2, "post 2 bis", "planet")
	dmock.ExpectQuery("SELECT \\* from dtable LIMIT 1").WillReturnRows(rows)
	dmock.ExpectPrepare("INSERT INTO dtable \\( title,body,id\\) VALUES \\( \\?,\\?,\\? \\)")
	dmock.ExpectPrepare("UPDATE dtable SET  title=\\?,body=\\? WHERE id = \\?")
	dmock.ExpectExec("UPDATE dtable SET  title=\\?,body=\\? WHERE id = \\?").WithArgs("post 1", "hello", "1").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectExec("UPDATE dtable SET  title=\\?,body=\\? WHERE id = \\?").WithArgs("post 2", "world", "2").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectClose()
	dest := mockdatasource.MockDatasource{MockedDb: ddb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog"}

	saver, err := database.NewSaver(context.Background(), &dest, "dtable", "id", "update")
	if err != nil {
		t.Fatalf("NewSaver should not return error and returned '%v'", err)
	}
	loader, err := database.NewLoader(context.Background(), &source, "stable", "")
	if err != nil {
		t.Fatalf("NewLoader should not return error and returned '%v'", err)
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

func TestReplaceOk(t *testing.T) {
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

	ddb, dmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	rows = sqlmock.NewRows([]string{"id"}).
		AddRow(3).
		AddRow(2)
	dmock.ExpectQuery("SELECT id from dtable").WillReturnRows(rows)
	rows = sqlmock.NewRows([]string{"id", "title", "body"}).
		AddRow(3, "post 3", "good").
		AddRow(2, "post 2 bis", "planet")
	dmock.ExpectQuery("SELECT \\* from dtable LIMIT 1").WillReturnRows(rows)
	dmock.ExpectPrepare("INSERT INTO dtable \\( title,body,id\\) VALUES \\( \\?,\\?,\\? \\)")
	dmock.ExpectPrepare("UPDATE dtable SET  title=\\?,body=\\? WHERE id = \\?")
	dmock.ExpectExec("INSERT INTO dtable").WithArgs("post 1", "hello", "1").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectExec("UPDATE dtable SET  title=\\?,body=\\? WHERE id = \\?").WithArgs("post 2", "world", "2").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectClose()
	dest := mockdatasource.MockDatasource{MockedDb: ddb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog"}

	saver, err := database.NewSaver(context.Background(), &dest, "dtable", "id", "replace")
	if err != nil {
		t.Fatalf("NewSaver should not return error and returned '%v'", err)
	}
	loader, err := database.NewLoader(context.Background(), &source, "stable", "")
	if err != nil {
		t.Fatalf("NewLoader should not return error and returned '%v'", err)
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

func TestExactCopyOk(t *testing.T) {
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

	ddb, dmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows = sqlmock.NewRows([]string{"id"}).
		AddRow(3).
		AddRow(2).
		AddRow(4)
	dmock.ExpectQuery("SELECT id from dtable").WillReturnRows(rows)
	rows = sqlmock.NewRows([]string{"id", "title", "body"}).
		AddRow(3, "post 3", "good").
		AddRow(4, "post 4", "bye").
		AddRow(2, "post 2 bis", "planet")
	dmock.ExpectQuery("SELECT \\* from dtable LIMIT 1").WillReturnRows(rows)
	dmock.ExpectPrepare("INSERT INTO dtable \\( title,body,id\\) VALUES \\( \\?,\\?,\\? \\)")
	dmock.ExpectPrepare("UPDATE dtable SET  title=\\?,body=\\? WHERE id = \\?")
	dmock.ExpectExec("INSERT INTO dtable").WithArgs("post 1", "hello", "1").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectExec("UPDATE dtable SET  title=\\?,body=\\? WHERE id = \\?").WithArgs("post 2", "world", "2").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectExec("DELETE from dtable WHERE id=.+").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectExec("DELETE from dtable WHERE id=.+").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectClose()
	dest := mockdatasource.MockDatasource{MockedDb: ddb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog"}

	saver, err := database.NewSaver(context.Background(), &dest, "dtable", "id", "exactCopy")
	if err != nil {
		t.Fatalf("NewSaver should not return error and returned '%v'", err)
	}
	loader, err := database.NewLoader(context.Background(), &source, "stable", "")
	if err != nil {
		t.Fatalf("NewLoader should not return error and returned '%v'", err)
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

func TestExactCopyTransactionOk(t *testing.T) {
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

	ddb, dmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows = sqlmock.NewRows([]string{"id"}).
		AddRow(3).
		AddRow(2).
		AddRow(4)
	dmock.ExpectQuery("SELECT id from dtable").WillReturnRows(rows)
	rows = sqlmock.NewRows([]string{"id", "title", "body"}).
		AddRow(3, "post 3", "good").
		AddRow(4, "post 4", "bye").
		AddRow(2, "post 2 bis", "planet")
	dmock.ExpectBegin()
	dmock.ExpectQuery("SELECT \\* from dtable LIMIT 1").WillReturnRows(rows)
	dmock.ExpectPrepare("INSERT INTO dtable \\( title,body,id\\) VALUES \\( \\?,\\?,\\? \\)")
	dmock.ExpectPrepare("UPDATE dtable SET  title=\\?,body=\\? WHERE id = \\?")
	dmock.ExpectExec("INSERT INTO dtable").WithArgs("post 1", "hello", "1").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectExec("UPDATE dtable SET  title=\\?,body=\\? WHERE id = \\?").WithArgs("post 2", "world", "2").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectExec("DELETE from dtable WHERE id=.+").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectExec("DELETE from dtable WHERE id=.+").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectCommit()
	dmock.ExpectClose()
	dest := mockdatasource.MockDatasource{MockedDb: ddb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog", Transaction: true}

	saver, err := database.NewSaver(context.Background(), &dest, "dtable", "id", "exactCopy")
	if err != nil {
		t.Fatalf("NewSaver should not return error and returned '%v'", err)
	}
	loader, err := database.NewLoader(context.Background(), &source, "stable", "")
	if err != nil {
		t.Fatalf("NewLoader should not return error and returned '%v'", err)
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
