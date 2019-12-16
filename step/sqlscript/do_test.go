package sqlscript_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/marema31/kamino/step/sqlscript"
)

func TestDoOk(t *testing.T) {
	ctx, log, dss, v, mock, err := setupDo("testdata/good/steps/", "sqlscriptok")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}

	mock.ExpectExec("CREATE DATABASE mydb;").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("CREATE USER 'user1' IDENTIFIED BY 'security4ever';").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("CREATE FUNCTION inc\\(val integer\\) RETURNS integer AS \\$\\$ BEGIN RETURN val \\+ 1; END; \\$\\$ LANGUAGE PLPGSQL;").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("REVOKE ALL ON DATABASE db1 FROM PUBLIC;").WillReturnResult(sqlmock.NewResult(1, 1))

	_, steps, err := sqlscript.Load(ctx, log, "testdata/good", "sqlscriptok", 0, v, dss, false, false)
	if err != nil {
		t.Fatalf("Load should not returns an error, returned: %v", err)
	}

	err = steps[0].Init(ctx, log)
	if err != nil {
		t.Fatalf("Init should not returns an error, returned: %v", err)
	}

	err = steps[0].Do(context.Background(), log)
	if err != nil {
		t.Errorf("Do should not return error, returned: %v", err)
	}

	steps[0].Finish(log)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDoTransactionOk(t *testing.T) {
	ctx, log, dss, v, mock, err := setupDo("testdata/good/steps/", "transaction")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}

	mock.ExpectBegin()
	mock.ExpectExec("CREATE DATABASE mydb;").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("CREATE USER 'user1' IDENTIFIED BY 'security4ever';").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("CREATE FUNCTION inc\\(val integer\\) RETURNS integer AS \\$\\$ BEGIN RETURN val \\+ 1; END; \\$\\$ LANGUAGE PLPGSQL;").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("REVOKE ALL ON DATABASE db1 FROM PUBLIC;").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	_, steps, err := sqlscript.Load(ctx, log, "testdata/good", "transaction", 0, v, dss, false, false)
	if err != nil {
		t.Fatalf("Load should not returns an error, returned: %v", err)
	}

	err = steps[0].Init(ctx, log)
	if err != nil {
		t.Fatalf("Init should not returns an error, returned: %v", err)
	}

	err = steps[0].Do(context.Background(), log)
	if err != nil {
		t.Errorf("Do should not return error, returned: %v", err)
	}

	steps[0].Finish(log)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDoTransactionCancelOk(t *testing.T) {
	ctx, log, dss, v, mock, err := setupDo("testdata/good/steps/", "transaction")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}

	mock.ExpectBegin()
	mock.ExpectExec("CREATE DATABASE mydb;").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("CREATE USER 'user1' IDENTIFIED BY 'security4ever';").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("CREATE FUNCTION inc\\(val integer\\) RETURNS integer AS \\$\\$ BEGIN RETURN val \\+ 1; END; \\$\\$ LANGUAGE PLPGSQL;").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("REVOKE ALL ON DATABASE db1 FROM PUBLIC;").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectRollback()
	_, steps, err := sqlscript.Load(ctx, log, "testdata/good", "transaction", 0, v, dss, false, false)
	if err != nil {
		t.Fatalf("Load should not returns an error, returned: %v", err)
	}

	err = steps[0].Init(ctx, log)
	if err != nil {
		t.Fatalf("Init should not returns an error, returned: %v", err)
	}

	err = steps[0].Do(context.Background(), log)
	if err != nil {
		t.Errorf("Do should not return error, returned: %v", err)
	}

	steps[0].Cancel(log)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDoDryRun(t *testing.T) {
	ctx, log, dss, v, mock, err := setupDo("testdata/good/steps/", "sqlscriptok")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}

	_, steps, err := sqlscript.Load(ctx, log, "testdata/good", "sqlscriptok", 0, v, dss, false, true)
	if err != nil {
		t.Fatalf("Load should not returns an error, returned: %v", err)
	}

	err = steps[0].Init(ctx, log)
	if err != nil {
		t.Fatalf("Init should not returns an error, returned: %v", err)
	}

	err = steps[0].Do(context.Background(), log)
	if err != nil {
		t.Errorf("Do should not return error, returned: %v", err)
	}

	steps[0].Finish(log)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
