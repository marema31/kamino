package migration_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/marema31/kamino/step/migration"
)

func TestToSkipYesOk(t *testing.T) {
	ctx, log, dss, v, mock, err := setupDo("testdata/good/steps/", "migrationok")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}
	rows := sqlmock.NewRows([]string{"COUNT(id)"}).
		AddRow(10)

	mock.ExpectQuery("SELECT count\\(table_schema\\) FROM information_schema.tables WHERE table_catalog = db1 and table_schema = public").WillReturnRows(rows)

	_, steps, err := migration.Load(ctx, log, "testdata/good", "migrationok", 0, v, dss)
	if err != nil {
		t.Fatalf("Load should not returns an error, returned: %v", err)
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
	ctx, log, dss, v, mock, err := setupDo("testdata/good/steps/", "migrationok")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}
	rows := sqlmock.NewRows([]string{"COUNT(id)"}).
		AddRow(0)

	mock.ExpectQuery("SELECT count\\(table_schema\\) FROM information_schema.tables WHERE table_catalog = db1 and table_schema = public").WillReturnRows(rows)

	_, steps, err := migration.Load(ctx, log, "testdata/good", "migrationok", 0, v, dss)
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
	ctx, log, dss, v, mock, err := setupDo("testdata/good/steps/", "migrationok")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}

	mock.ExpectQuery("SELECT count\\(table_schema\\) FROM information_schema.tables WHERE table_catalog = db1 and table_schema = public").WillReturnError(fmt.Errorf("fake error"))

	_, steps, err := migration.Load(ctx, log, "testdata/good", "migrationok", 0, v, dss)
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
