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
	_, err = database.NewSaver(context.Background(), log, &dest, "dtable", "mykey", "replace")
	if err == nil {
		t.Fatalf("NewSaver should return error")
	}
}

func TestCreateIdsListError(t *testing.T) {
	ddb, dmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	dmock.ExpectQuery("SELECT \\? from \\?").WithArgs("id", "dtable").WillReturnError(fmt.Errorf("fake error"))
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

	dmock.ExpectQuery("SELECT column_name AS name FROM information_schema.columns WHERE table_schema = 'blog' AND table_name ='dtable';").WillReturnError(fmt.Errorf("fake error"))

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

	rows = sqlmock.NewRows([]string{"name"}).
		AddRow("id").
		AddRow("title").
		AddRow("body").
		AddRow("summary")
	dmock.ExpectQuery("SELECT column_name AS name FROM information_schema.columns WHERE table_schema = 'blog' AND table_name ='dtable';").WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"count"}).
		AddRow(1)
	dmock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM dtable").WillReturnRows(rows)
	dmock.ExpectPrepare("INSERT INTO dtable \\( `title`,`body`,`id`\\) VALUES \\( \\?,\\?,\\? \\)")
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

	rows = sqlmock.NewRows([]string{"name"}).
		AddRow("id").
		AddRow("title")
	dmock.ExpectQuery("SELECT column_name AS name FROM information_schema.columns WHERE table_schema = 'blog' AND table_name ='dtable';").WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"count"}).
		AddRow(0)
	dmock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM dtable").WillReturnRows(rows)
	dmock.ExpectPrepare("INSERT INTO dtable \\( `title`,`id`\\) VALUES \\( \\?,\\? \\)")
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

	dmock.ExpectQuery("SELECT \\* from \\? LIMIT 1").WithArgs("dtable").WillReturnError(fmt.Errorf("fake error"))
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
	dmock.ExpectQuery("SELECT \\? from \\?").WithArgs("id", "dtable").WillReturnRows(rows)
	rows = sqlmock.NewRows([]string{"id", "title", "body"})
	dmock.ExpectQuery("SELECT \\* from \\? LIMIT 1").WithArgs("dtable").WillReturnRows(rows)
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
func TestPrepareInsertError(t *testing.T) {
	sdb, smock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "title", "body"}).
		AddRow(2, "post 2", "world")

	smock.ExpectQuery("SELECT (.+) from stable WHERE title like '%'").WillReturnRows(rows)
	source := mockdatasource.MockDatasource{MockedDb: sdb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog"}

	ddb, dmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows = sqlmock.NewRows([]string{"id", "title", "body"})
	dmock.ExpectQuery("SELECT \\* from \\? LIMIT 1").WithArgs("dtable").WillReturnRows(rows)
	dmock.ExpectPrepare("INSERT INTO dtable \\( `title`,`body`,`id`\\) VALUES \\( \\?,\\?,\\? \\)").WillReturnError(fmt.Errorf("fake error"))
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

	for loader.Next() {
		record, err := loader.Load(log)
		if err != nil {
			t.Fatalf("Load should not return error and returned '%v'", err)
		}

		if err = saver.Save(log, record); err == nil {
			t.Fatalf("Save should return error")
		}
	}
}

func TestPrepareUpdateError(t *testing.T) {
	sdb, smock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "title", "body"}).
		AddRow(2, "post 2", "world")

	smock.ExpectQuery("SELECT (.+) from stable WHERE title like '%'").WillReturnRows(rows)
	source := mockdatasource.MockDatasource{MockedDb: sdb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog"}

	ddb, dmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows = sqlmock.NewRows([]string{"id"}).
		AddRow(1).
		AddRow(2)
	dmock.ExpectQuery("SELECT \\? from \\?").WithArgs("id", "dtable").WillReturnRows(rows)
	rows = sqlmock.NewRows([]string{"id", "title", "body"})
	dmock.ExpectQuery("SELECT \\* from \\? LIMIT 1").WithArgs("dtable").WillReturnRows(rows)
	dmock.ExpectPrepare("INSERT INTO dtable \\( `title`,`body`,`id`\\) VALUES \\( \\?,\\?,\\? \\)")
	dmock.ExpectPrepare("UPDATE dtable SET  title=\\?,body=\\? WHERE id = \\?").WillReturnError(fmt.Errorf("fake error"))
	dest := mockdatasource.MockDatasource{MockedDb: ddb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")

	saver, err := database.NewSaver(context.Background(), log, &dest, "dtable", "id", "replace")
	if err != nil {
		t.Fatalf("NewSaver should not return error and returned '%v'", err)
	}
	loader, err := database.NewLoader(context.Background(), log, &source, "stable", "title like '%'")
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
}

func TestBeginTransactionError(t *testing.T) {
	sdb, smock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "title", "body"}).
		AddRow(2, "post 2", "world")

	smock.ExpectQuery("SELECT (.+) from stable").WillReturnRows(rows)
	smock.ExpectClose()
	source := mockdatasource.MockDatasource{MockedDb: sdb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog", Transaction: true}

	ddb, dmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	dmock.ExpectBegin().WillReturnError(fmt.Errorf("fake error"))
	rows = sqlmock.NewRows([]string{"id", "title", "body"})
	dmock.ExpectQuery("SELECT \\* from \\? LIMIT 1").WithArgs("dtable").WillReturnRows(rows)
	dmock.ExpectPrepare("INSERT INTO dtable \\( `title`,`body`,`id`\\) VALUES \\( \\?,\\?,\\? \\)")
	dmock.ExpectExec("INSERT INTO dtable").WithArgs("post 1", "hello", "1").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectCommit()
	dmock.ExpectClose()
	dmock.ExpectClose() // sqlmock and close/defer have erratic behavior
	dest := mockdatasource.MockDatasource{MockedDb: ddb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog", Transaction: true}
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

	for loader.Next() {
		record, err := loader.Load(log)
		if err != nil {
			t.Fatalf("Load should not return error and returned '%v'", err)
		}

		if err = saver.Save(log, record); err == nil {
			t.Fatalf("Save should return error")
		}
	}
}

func TestEndTransactionError(t *testing.T) {
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

	dmock.ExpectBegin()
	rows = sqlmock.NewRows([]string{"name"}).
		AddRow("id").
		AddRow("title").
		AddRow("body")
	dmock.ExpectQuery("SELECT column_name AS name FROM information_schema.columns WHERE table_schema = 'blog' AND table_name ='dtable';").WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"count"}).
		AddRow(0)
	dmock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM dtable").WillReturnRows(rows)
	dmock.ExpectPrepare("INSERT INTO dtable \\( `title`,`body`,`id`\\) VALUES \\( \\?,\\?,\\? \\)")
	dmock.ExpectExec("INSERT INTO dtable").WithArgs("post 1", "hello", "1").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectExec("INSERT INTO dtable").WithArgs("post 2", "world", "2").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectCommit().WillReturnError(fmt.Errorf("fake error"))
	dmock.ExpectCommit().WillReturnError(fmt.Errorf("fake error"))
	dest := mockdatasource.MockDatasource{MockedDb: ddb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog", Transaction: true}
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

func TestRollbackTransactionError(t *testing.T) {
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

	dmock.ExpectBegin()
	rows = sqlmock.NewRows([]string{"name"}).
		AddRow("id").
		AddRow("title").
		AddRow("body")
	dmock.ExpectQuery("SELECT column_name AS name FROM information_schema.columns WHERE table_schema = 'blog' AND table_name ='dtable';").WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"count"}).
		AddRow(0)
	dmock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM dtable").WillReturnRows(rows)
	dmock.ExpectPrepare("INSERT INTO dtable \\( `title`,`body`,`id`\\) VALUES \\( \\?,\\?,\\? \\)")
	dmock.ExpectExec("INSERT INTO dtable").WithArgs("post 1", "hello", "1").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectExec("INSERT INTO dtable").WithArgs("post 2", "world", "2").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectRollback().WillReturnError(fmt.Errorf("fake error"))
	dmock.ExpectClose()
	dmock.ExpectClose() // sqlmock and close/defer have erratic behavior
	dest := mockdatasource.MockDatasource{MockedDb: ddb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog", Transaction: true}
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

	for loader.Next() {
		record, err := loader.Load(log)
		if err != nil {
			t.Fatalf("Load should not return error and returned '%v'", err)
		}

		if err = saver.Save(log, record); err != nil {
			t.Fatalf("Save should not return error and returned '%v'", err)
		}
	}
	err = saver.Reset(log)
	if err == nil {
		t.Errorf("Saver close should return error")
	}

	err = loader.Close(log)
	if err != nil {
		t.Errorf("Loader close should not return error and returned '%v'", err)
	}
}

func TestTruncateError(t *testing.T) {
	sdb, smock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "title", "body"}).
		AddRow(2, "post 2", "world")

	smock.ExpectQuery("SELECT (.+) from stable WHERE title like '%'").WillReturnRows(rows)
	source := mockdatasource.MockDatasource{MockedDb: sdb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog"}

	ddb, dmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	dmock.ExpectBegin()
	rows = sqlmock.NewRows([]string{"name"}).
		AddRow("id").
		AddRow("title").
		AddRow("body")
	dmock.ExpectQuery("SELECT column_name AS name FROM information_schema.columns WHERE table_schema = 'blog' AND table_name ='dtable';").WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"count"}).
		AddRow(0)
	dmock.ExpectPrepare("INSERT INTO dtable \\( title,body,id\\) VALUES \\( \\?,\\?,\\? \\)")
	//	dmock.ExpectPrepare("UPDATE dtable SET  title=\\?,body=\\? WHERE id = \\?").WillReturnError(fmt.Errorf("fake error"))
	dmock.ExpectQuery("TRUNCATE TABLE dtable").WillReturnError(fmt.Errorf("fake error"))
	dest := mockdatasource.MockDatasource{MockedDb: ddb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")

	saver, err := database.NewSaver(context.Background(), log, &dest, "dtable", "id", "truncate")
	if err != nil {
		t.Fatalf("NewSaver should not return error and returned '%v'", err)
	}
	loader, err := database.NewLoader(context.Background(), log, &source, "stable", "title like '%'")
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
}

func TestColsListError(t *testing.T) {
	sdb, smock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	smock.ExpectQuery("SELECT (.+) from stable").WillReturnError(fmt.Errorf("fake error"))
	smock.ExpectClose()
	source := mockdatasource.MockDatasource{MockedDb: sdb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	logger.SetLevel(logrus.PanicLevel)

	ddb, dmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	dmock.ExpectQuery("SELECT \\* from \\? LIMIT 1").WithArgs("dtable").WillReturnError(fmt.Errorf("fake error"))
	dest := mockdatasource.MockDatasource{MockedDb: ddb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog"}

	_, err = database.NewSaver(context.Background(), log, &dest, "dtable", "id", "insert")
	if err != nil {
		t.Fatalf("NewSaver should not return error and returned '%v'", err)
	}
	_, err = database.NewLoader(context.Background(), log, &source, "stable", "")
	if err == nil {
		t.Fatalf("NewLoader should return error")
	}
}

func TestInsertError(t *testing.T) {
	sdb, smock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "title", "body"}).
		AddRow(2, "post 2", "world")

	smock.ExpectQuery("SELECT (.+) from stable WHERE title like '%'").WillReturnRows(rows)
	source := mockdatasource.MockDatasource{MockedDb: sdb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog"}

	ddb, dmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows = sqlmock.NewRows([]string{"id"}).
		AddRow(1).
		AddRow(2)
	dmock.ExpectQuery("SELECT \\? from \\?").WithArgs("id", "dtable").WillReturnRows(rows)
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
	dmock.ExpectExec("INSERT INTO dtable").WithArgs("post 2", "world", "2").WillReturnError(fmt.Errorf("fake error"))
	dest := mockdatasource.MockDatasource{MockedDb: ddb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")

	saver, err := database.NewSaver(context.Background(), log, &dest, "dtable", "id", "replace")
	if err != nil {
		t.Fatalf("NewSaver should not return error and returned '%v'", err)
	}
	loader, err := database.NewLoader(context.Background(), log, &source, "stable", "title like '%'")
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
}

func TestDeleteError(t *testing.T) {
	sdb, smock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "title", "body"}).
		AddRow(2, "post 2", "world")

	smock.ExpectQuery("SELECT (.+) from stable WHERE title like '%'").WillReturnRows(rows)
	source := mockdatasource.MockDatasource{MockedDb: sdb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog"}

	ddb, dmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows = sqlmock.NewRows([]string{"id"}).
		AddRow(1).
		AddRow(2)
	dmock.ExpectQuery("SELECT \\? from \\?").WithArgs("id", "dtable").WillReturnRows(rows)
	rows = sqlmock.NewRows([]string{"name"}).
		AddRow("id").
		AddRow("title").
		AddRow("body")
	dmock.ExpectQuery("SELECT column_name AS name FROM information_schema.columns WHERE table_schema = 'blog' AND table_name ='dtable';").WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"count"}).
		AddRow(0)
	dmock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM dtable").WillReturnRows(rows)
	dmock.ExpectPrepare("INSERT INTO dtable \\( `title`,`body`,`id`\\) VALUES \\( \\?,\\?,\\? \\)")
	dmock.ExpectPrepare("UPDATE dtable SET  title=\\?,body=\\? WHERE id = \\?")
	dmock.ExpectExec("UPDATE dtable SET title=\\?,body=\\? WHERE id = \\?").WithArgs("post 2", "world", "2").WillReturnResult(sqlmock.NewResult(1, 1))
	dmock.ExpectExec("DELETE from dtable WHERE id=1").WillReturnError(fmt.Errorf("fake error"))
	dest := mockdatasource.MockDatasource{MockedDb: ddb, Type: datasource.Database, Engine: datasource.Mysql, Database: "blog"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")

	saver, err := database.NewSaver(context.Background(), log, &dest, "dtable", "id", "exactCopy")
	if err != nil {
		t.Fatalf("NewSaver should not return error and returned '%v'", err)
	}
	loader, err := database.NewLoader(context.Background(), log, &source, "stable", "title like '%'")
	if err != nil {
		t.Fatalf("NewLoader should not return error and returned '%v'", err)
	}

	for loader.Next() {
		record, err := loader.Load(log)
		if err != nil {
			t.Fatalf("Load should not return error and returned '%v'", err)
		}

		if err = saver.Save(log, record); err != nil {
			t.Fatalf("Load should not return error and returned '%v'", err)
		}
	}
	err = saver.Close(log)
	if err == nil {
		t.Errorf("Saver close should return error")
	}
}
