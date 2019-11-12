package sqlscript_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/marema31/kamino/step/sqlscript"
)

func TestToSkipYesOk(t *testing.T) {
	ctx, log, dss, v, mock, err := setupDo("testdata/good/steps/", "sqlscriptok")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}
	rows := sqlmock.NewRows([]string{"COUNT(id)"}).
		AddRow(10)

	mock.ExpectQuery("SELECT \\* FROM USER WHERE user = 'user1';").WillReturnRows(rows)

	_, steps, err := sqlscript.Load(ctx, log, "testdata/good", "sqlscriptok", v, dss)
	if err != nil {
		t.Errorf("Load should not returns an error, returned: %v", err)
	}

	ok, err := steps[0].ToSkip(context.Background(), log)
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
	ctx, log, dss, v, mock, err := setupDo("testdata/good/steps/", "sqlscriptok")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}
	rows := sqlmock.NewRows([]string{"COUNT(id)"}).
		AddRow(0)

	mock.ExpectQuery("SELECT \\* FROM USER WHERE user = 'user1';").WillReturnRows(rows)
	if err != nil {
		t.Fatalf("NewSaver should not return error and returned '%v'", err)
	}

	_, steps, err := sqlscript.Load(ctx, log, "testdata/good", "sqlscriptok", v, dss)
	if err != nil {
		t.Errorf("Load should not returns an error, returned: %v", err)
	}

	ok, err := steps[0].ToSkip(context.Background(), log)
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
	ctx, log, dss, v, mock, err := setupDo("testdata/good/steps/", "sqlscriptok")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}

	mock.ExpectQuery("SELECT \\* FROM USER WHERE user = 'user1';").WillReturnError(fmt.Errorf("fake error"))
	if err != nil {
		t.Fatalf("NewSaver should not return error and returned '%v'", err)
	}

	_, steps, err := sqlscript.Load(ctx, log, "testdata/good", "sqlscriptok", v, dss)
	if err != nil {
		t.Errorf("Load should not returns an error, returned: %v", err)
	}

	_, err = steps[0].ToSkip(context.Background(), log)
	if err == nil {
		t.Errorf("ToSkip should return error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}
