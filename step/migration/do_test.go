package migration_test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/marema31/kamino/step/migration"
)

func TestDoUpOk(t *testing.T) {
	ctx, log, dss, v, mock, err := setupDo("testdata/good/steps/", "migrationok")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}
	rows := sqlmock.NewRows([]string{"NOW()"}).
		AddRow(time.Now())

	mock.ExpectQuery("SELECT NOW()").WillReturnRows(rows)
	mock.ExpectExec("create table if not exists `kamino_admin_migrations` \\(`id` varchar\\(255\\) not null primary key, `applied_at` datetime\\) engine=InnoDB charset=UTF8;").WillReturnResult(sqlmock.NewResult(1, 1))
	rows = sqlmock.NewRows([]string{"id", "applied_at"})
	mock.ExpectQuery("SELECT \\* FROM `kamino_admin_migrations`").WillReturnRows(rows)
	mock.ExpectBegin()
	mock.ExpectExec("CREATE TABLE admin \\( id INT, name VARCHAR\\(255\\) \\);").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("insert into `kamino_admin_migrations` \\(`id`,`applied_at`\\)").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	rows = sqlmock.NewRows([]string{"NOW()"}).
		AddRow(time.Now())
	mock.ExpectQuery("SELECT NOW()").WillReturnRows(rows)
	mock.ExpectExec("create table if not exists `kamino_admin_migrations` \\(`id` varchar\\(255\\) not null primary key, `applied_at` datetime\\) engine=InnoDB charset=UTF8;").WillReturnResult(sqlmock.NewResult(1, 1))
	rows = sqlmock.NewRows([]string{"id", "applied_at"}).
		AddRow("v0001.sql", time.Now())
	mock.ExpectQuery("SELECT \\* FROM `kamino_admin_migrations`").WillReturnRows(rows)
	rows = sqlmock.NewRows([]string{"NOW()"}).
		AddRow(time.Now())
	mock.ExpectQuery("SELECT NOW()").WillReturnRows(rows)
	mock.ExpectExec("create table if not exists `kamino_user_migrations` \\(`id` varchar\\(255\\) not null primary key, `applied_at` datetime\\) engine=InnoDB charset=UTF8;").WillReturnResult(sqlmock.NewResult(1, 1))
	rows = sqlmock.NewRows([]string{"id", "applied_at"})
	mock.ExpectQuery("SELECT \\* FROM `kamino_user_migrations`").WillReturnRows(rows)
	mock.ExpectBegin()
	mock.ExpectExec("CREATE TABLE user \\( id INT, name VARCHAR\\(255\\) \\);").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("insert into `kamino_user_migrations` \\(`id`,`applied_at`\\)").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	rows = sqlmock.NewRows([]string{"NOW()"}).
		AddRow(time.Now())
	mock.ExpectQuery("SELECT NOW()").WillReturnRows(rows)
	mock.ExpectExec("create table if not exists `kamino_user_migrations` \\(`id` varchar\\(255\\) not null primary key, `applied_at` datetime\\) engine=InnoDB charset=UTF8;").WillReturnResult(sqlmock.NewResult(1, 1))
	rows = sqlmock.NewRows([]string{"id", "applied_at"}).
		AddRow("v0001.sql", time.Now())
	mock.ExpectQuery("SELECT \\* FROM `kamino_user_migrations`").WillReturnRows(rows)

	_, steps, err := migration.Load(ctx, log, "testdata/good", "migrationok", 0, v, dss, false, false, nil)
	if err != nil {
		t.Fatalf("Load should not returns an error, returned: %v", err)
	}

	superseed := make(map[string]string)
	superseed["migration.limit"] = "2"
	steps[0].PostLoad(log, superseed)
	err = steps[0].Init(ctx, log)
	if err != nil {
		t.Fatalf("Init should not returns an error, returned: %v", err)
	}

	err = steps[0].Do(context.Background(), log)
	if err != nil {
		t.Errorf("Do should not return error, returned: %v", err)
	}

	steps[0].Cancel(log)
	steps[0].Finish(log)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDoUpAdminError(t *testing.T) {
	ctx, log, dss, v, mock, err := setupDo("testdata/good/steps/", "migrationok")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}
	rows := sqlmock.NewRows([]string{"NOW()"}).
		AddRow(time.Now())

	mock.ExpectQuery("SELECT NOW()").WillReturnRows(rows)
	mock.ExpectExec("create table if not exists `kamino_admin_migrations` \\(`id` varchar\\(255\\) not null primary key, `applied_at` datetime\\) engine=InnoDB charset=UTF8;").WillReturnResult(sqlmock.NewResult(1, 1))
	rows = sqlmock.NewRows([]string{"id", "applied_at"})
	mock.ExpectQuery("SELECT \\* FROM `kamino_admin_migrations`").WillReturnRows(rows)
	mock.ExpectBegin()
	mock.ExpectExec("CREATE TABLE admin \\( id INT, name VARCHAR\\(255\\) \\);").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("insert into `kamino_admin_migrations` \\(`id`,`applied_at`\\)").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	rows = sqlmock.NewRows([]string{"NOW()"}).
		AddRow(time.Now())
	mock.ExpectQuery("SELECT NOW()").WillReturnRows(rows)
	mock.ExpectExec("create table if not exists `kamino_admin_migrations` \\(`id` varchar\\(255\\) not null primary key, `applied_at` datetime\\) engine=InnoDB charset=UTF8;").WillReturnResult(sqlmock.NewResult(1, 1))
	rows = sqlmock.NewRows([]string{"id", "applied_at"})
	mock.ExpectQuery("SELECT \\* FROM `kamino_admin_migrations`").WillReturnRows(rows)

	_, steps, err := migration.Load(ctx, log, "testdata/good", "migrationok", 0, v, dss, false, false, nil)
	if err != nil {
		t.Fatalf("Load should not returns an error, returned: %v", err)
	}

	superseed := make(map[string]string)
	superseed["migration.limit"] = "2"
	steps[0].PostLoad(log, superseed)
	err = steps[0].Init(ctx, log)
	if err != nil {
		t.Fatalf("Init should not returns an error, returned: %v", err)
	}

	err = steps[0].Do(context.Background(), log)
	if err == nil {
		t.Errorf("Do should return error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDoUpLimitedOk(t *testing.T) {
	ctx, log, dss, v, mock, err := setupDo("testdata/good/steps/", "migrationok")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}
	rows := sqlmock.NewRows([]string{"NOW()"}).
		AddRow(time.Now())

	mock.ExpectQuery("SELECT NOW()").WillReturnRows(rows)
	mock.ExpectExec("create table if not exists `kamino_admin_migrations` \\(`id` varchar\\(255\\) not null primary key, `applied_at` datetime\\) engine=InnoDB charset=UTF8;").WillReturnResult(sqlmock.NewResult(1, 1))
	rows = sqlmock.NewRows([]string{"id", "applied_at"})
	mock.ExpectQuery("SELECT \\* FROM `kamino_admin_migrations`").WillReturnRows(rows)
	mock.ExpectBegin()
	mock.ExpectExec("CREATE TABLE admin \\( id INT, name VARCHAR\\(255\\) \\);").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("insert into `kamino_admin_migrations` \\(`id`,`applied_at`\\)").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	rows = sqlmock.NewRows([]string{"NOW()"}).
		AddRow(time.Now())
	mock.ExpectQuery("SELECT NOW()").WillReturnRows(rows)
	mock.ExpectExec("create table if not exists `kamino_admin_migrations` \\(`id` varchar\\(255\\) not null primary key, `applied_at` datetime\\) engine=InnoDB charset=UTF8;").WillReturnResult(sqlmock.NewResult(1, 1))
	rows = sqlmock.NewRows([]string{"id", "applied_at"}).
		AddRow("v0001.sql", time.Now())
	mock.ExpectQuery("SELECT \\* FROM `kamino_admin_migrations`").WillReturnRows(rows)

	_, steps, err := migration.Load(ctx, log, "testdata/good", "migrationok", 0, v, dss, false, false, nil)
	if err != nil {
		t.Fatalf("Load should not returns an error, returned: %v", err)
	}

	superseed := make(map[string]string)
	superseed["migration.limit"] = "1"
	steps[0].PostLoad(log, superseed)
	err = steps[0].Init(ctx, log)
	if err != nil {
		t.Fatalf("Init should not returns an error, returned: %v", err)
	}

	err = steps[0].Do(context.Background(), log)
	if err != nil {
		t.Errorf("Do should not return error, returned: %v", err)
	}

	steps[0].Cancel(log)
	steps[0].Finish(log)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDoDownOk(t *testing.T) {
	ctx, log, dss, v, mock, err := setupDo("testdata/good/steps/", "migrationok")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}
	rows := sqlmock.NewRows([]string{"NOW()"}).
		AddRow(time.Now())
	mock.ExpectQuery("SELECT NOW()").WillReturnRows(rows)
	mock.ExpectExec("create table if not exists `kamino_user_migrations` \\(`id` varchar\\(255\\) not null primary key, `applied_at` datetime\\) engine=InnoDB charset=UTF8;").WillReturnResult(sqlmock.NewResult(1, 1))
	rows = sqlmock.NewRows([]string{"id", "applied_at"}).
		AddRow("v0001.sql", time.Now())
	mock.ExpectQuery("SELECT \\* FROM `kamino_user_migrations`").WillReturnRows(rows)
	mock.ExpectBegin()
	mock.ExpectExec("DROP TABLE user;").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("delete from `kamino_user_migrations` where `id`=?").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	rows = sqlmock.NewRows([]string{"NOW()"}).
		AddRow(time.Now())
	mock.ExpectQuery("SELECT NOW()").WillReturnRows(rows)
	mock.ExpectExec("create table if not exists `kamino_user_migrations` \\(`id` varchar\\(255\\) not null primary key, `applied_at` datetime\\) engine=InnoDB charset=UTF8;").WillReturnResult(sqlmock.NewResult(1, 1))
	rows = sqlmock.NewRows([]string{"id", "applied_at"})
	mock.ExpectQuery("SELECT \\* FROM `kamino_user_migrations`").WillReturnRows(rows)
	rows = sqlmock.NewRows([]string{"NOW()"}).
		AddRow(time.Now())
	mock.ExpectQuery("SELECT NOW()").WillReturnRows(rows)
	mock.ExpectExec("create table if not exists `kamino_admin_migrations` \\(`id` varchar\\(255\\) not null primary key, `applied_at` datetime\\) engine=InnoDB charset=UTF8;").WillReturnResult(sqlmock.NewResult(1, 1))
	rows = sqlmock.NewRows([]string{"id", "applied_at"}).
		AddRow("v0001.sql", time.Now())
	mock.ExpectQuery("SELECT \\* FROM `kamino_admin_migrations`").WillReturnRows(rows)
	mock.ExpectBegin()
	mock.ExpectExec("DROP TABLE admin;").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("delete from `kamino_admin_migrations` where `id`=?").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	rows = sqlmock.NewRows([]string{"NOW()"}).
		AddRow(time.Now())
	mock.ExpectQuery("SELECT NOW()").WillReturnRows(rows)
	mock.ExpectExec("create table if not exists `kamino_admin_migrations` \\(`id` varchar\\(255\\) not null primary key, `applied_at` datetime\\) engine=InnoDB charset=UTF8;").WillReturnResult(sqlmock.NewResult(1, 1))
	rows = sqlmock.NewRows([]string{"id", "applied_at"})
	mock.ExpectQuery("SELECT \\* FROM `kamino_admin_migrations`").WillReturnRows(rows)

	_, steps, err := migration.Load(ctx, log, "testdata/good", "migrationok", 0, v, dss, false, false, nil)
	if err != nil {
		t.Fatalf("Load should not returns an error, returned: %v", err)
	}

	superseed := make(map[string]string)
	superseed["migration.dir"] = "down"
	steps[0].PostLoad(log, superseed)
	err = steps[0].Init(ctx, log)
	if err != nil {
		t.Fatalf("Init should not returns an error, returned: %v", err)
	}

	err = steps[0].Do(context.Background(), log)
	if err != nil {
		t.Errorf("Do should not return error, returned: %v", err)
	}

	steps[0].Cancel(log)
	steps[0].Finish(log)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDoDownAdminError(t *testing.T) {
	ctx, log, dss, v, mock, err := setupDo("testdata/good/steps/", "migrationok")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}
	rows := sqlmock.NewRows([]string{"NOW()"}).
		AddRow(time.Now())
	mock.ExpectQuery("SELECT NOW()").WillReturnRows(rows)
	mock.ExpectExec("create table if not exists `kamino_user_migrations` \\(`id` varchar\\(255\\) not null primary key, `applied_at` datetime\\) engine=InnoDB charset=UTF8;").WillReturnResult(sqlmock.NewResult(1, 1))
	rows = sqlmock.NewRows([]string{"id", "applied_at"}).
		AddRow("v0001.sql", time.Now())
	mock.ExpectQuery("SELECT \\* FROM `kamino_user_migrations`").WillReturnRows(rows)
	mock.ExpectBegin()
	mock.ExpectExec("DROP TABLE user;").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("delete from `kamino_user_migrations` where `id`=?").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	rows = sqlmock.NewRows([]string{"NOW()"}).
		AddRow(time.Now())
	mock.ExpectQuery("SELECT NOW()").WillReturnRows(rows)
	mock.ExpectExec("create table if not exists `kamino_user_migrations` \\(`id` varchar\\(255\\) not null primary key, `applied_at` datetime\\) engine=InnoDB charset=UTF8;").WillReturnResult(sqlmock.NewResult(1, 1))
	rows = sqlmock.NewRows([]string{"id", "applied_at"})
	mock.ExpectQuery("SELECT \\* FROM `kamino_user_migrations`").WillReturnRows(rows)
	rows = sqlmock.NewRows([]string{"NOW()"}).
		AddRow(time.Now())
	mock.ExpectQuery("SELECT NOW()").WillReturnRows(rows)
	mock.ExpectExec("create table if not exists `kamino_admin_migrations` \\(`id` varchar\\(255\\) not null primary key, `applied_at` datetime\\) engine=InnoDB charset=UTF8;").WillReturnResult(sqlmock.NewResult(1, 1))
	rows = sqlmock.NewRows([]string{"id", "applied_at"}).
		AddRow("v0001.sql", time.Now())
	mock.ExpectQuery("SELECT \\* FROM `kamino_admin_migrations`").WillReturnRows(rows)
	mock.ExpectBegin()
	mock.ExpectExec("DROP TABLE admin;").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("delete from `kamino_admin_migrations` where `id`=?").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	rows = sqlmock.NewRows([]string{"NOW()"}).
		AddRow(time.Now())
	mock.ExpectQuery("SELECT NOW()").WillReturnRows(rows)
	mock.ExpectExec("create table if not exists `kamino_admin_migrations` \\(`id` varchar\\(255\\) not null primary key, `applied_at` datetime\\) engine=InnoDB charset=UTF8;").WillReturnResult(sqlmock.NewResult(1, 1))
	rows = sqlmock.NewRows([]string{"id", "applied_at"}).
		AddRow("v0001.sql", time.Now())
	mock.ExpectQuery("SELECT \\* FROM `kamino_admin_migrations`").WillReturnRows(rows)

	_, steps, err := migration.Load(ctx, log, "testdata/good", "migrationok", 0, v, dss, false, false, nil)
	if err != nil {
		t.Fatalf("Load should not returns an error, returned: %v", err)
	}

	superseed := make(map[string]string)
	superseed["migration.dir"] = "down"
	steps[0].PostLoad(log, superseed)
	err = steps[0].Init(ctx, log)
	if err != nil {
		t.Fatalf("Init should not returns an error, returned: %v", err)
	}

	err = steps[0].Do(context.Background(), log)
	if err == nil {
		t.Errorf("Do should return error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDoDownUserError(t *testing.T) {
	ctx, log, dss, v, mock, err := setupDo("testdata/good/steps/", "migrationok")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}
	rows := sqlmock.NewRows([]string{"NOW()"}).
		AddRow(time.Now())
	mock.ExpectQuery("SELECT NOW()").WillReturnRows(rows)
	mock.ExpectExec("create table if not exists `kamino_user_migrations` \\(`id` varchar\\(255\\) not null primary key, `applied_at` datetime\\) engine=InnoDB charset=UTF8;").WillReturnResult(sqlmock.NewResult(1, 1))
	rows = sqlmock.NewRows([]string{"id", "applied_at"}).
		AddRow("v0001.sql", time.Now())
	mock.ExpectQuery("SELECT \\* FROM `kamino_user_migrations`").WillReturnRows(rows)
	mock.ExpectBegin()
	mock.ExpectExec("DROP TABLE user;").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("delete from `kamino_user_migrations` where `id`=?").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	rows = sqlmock.NewRows([]string{"NOW()"}).
		AddRow(time.Now())
	mock.ExpectQuery("SELECT NOW()").WillReturnRows(rows)
	mock.ExpectExec("create table if not exists `kamino_user_migrations` \\(`id` varchar\\(255\\) not null primary key, `applied_at` datetime\\) engine=InnoDB charset=UTF8;").WillReturnResult(sqlmock.NewResult(1, 1))
	rows = sqlmock.NewRows([]string{"id", "applied_at"}).
		AddRow("v0001.sql", time.Now())
	mock.ExpectQuery("SELECT \\* FROM `kamino_user_migrations`").WillReturnRows(rows)

	_, steps, err := migration.Load(ctx, log, "testdata/good", "migrationok", 0, v, dss, false, false, nil)
	if err != nil {
		t.Fatalf("Load should not returns an error, returned: %v", err)
	}

	superseed := make(map[string]string)
	superseed["migration.dir"] = "down"
	steps[0].PostLoad(log, superseed)
	err = steps[0].Init(ctx, log)
	if err != nil {
		t.Fatalf("Init should not returns an error, returned: %v", err)
	}

	err = steps[0].Do(context.Background(), log)
	if err == nil {
		t.Errorf("Do should return error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
func TestDoDownLimitedOk(t *testing.T) {
	ctx, log, dss, v, mock, err := setupDo("testdata/good/steps/", "migrationok")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}
	rows := sqlmock.NewRows([]string{"NOW()"}).
		AddRow(time.Now())
	mock.ExpectQuery("SELECT NOW()").WillReturnRows(rows)
	mock.ExpectExec("create table if not exists `kamino_user_migrations` \\(`id` varchar\\(255\\) not null primary key, `applied_at` datetime\\) engine=InnoDB charset=UTF8;").WillReturnResult(sqlmock.NewResult(1, 1))
	rows = sqlmock.NewRows([]string{"id", "applied_at"}).
		AddRow("v0001.sql", time.Now())
	mock.ExpectQuery("SELECT \\* FROM `kamino_user_migrations`").WillReturnRows(rows)
	mock.ExpectBegin()
	mock.ExpectExec("DROP TABLE user;").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("delete from `kamino_user_migrations` where `id`=?").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	rows = sqlmock.NewRows([]string{"NOW()"}).
		AddRow(time.Now())
	mock.ExpectQuery("SELECT NOW()").WillReturnRows(rows)
	mock.ExpectExec("create table if not exists `kamino_user_migrations` \\(`id` varchar\\(255\\) not null primary key, `applied_at` datetime\\) engine=InnoDB charset=UTF8;").WillReturnResult(sqlmock.NewResult(1, 1))
	rows = sqlmock.NewRows([]string{"id", "applied_at"})
	mock.ExpectQuery("SELECT \\* FROM `kamino_user_migrations`").WillReturnRows(rows)

	_, steps, err := migration.Load(ctx, log, "testdata/good", "migrationok", 0, v, dss, false, false, nil)
	if err != nil {
		t.Fatalf("Load should not returns an error, returned: %v", err)
	}

	superseed := make(map[string]string)
	superseed["migration.dir"] = "down"
	superseed["migration.limit"] = "1"
	steps[0].PostLoad(log, superseed)
	err = steps[0].Init(ctx, log)
	if err != nil {
		t.Fatalf("Init should not returns an error, returned: %v", err)
	}

	err = steps[0].Do(context.Background(), log)
	if err != nil {
		t.Errorf("Do should not returns an error, returned: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDoDownExecError(t *testing.T) {
	ctx, log, dss, v, mock, err := setupDo("testdata/good/steps/", "migrationok")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}
	rows := sqlmock.NewRows([]string{"NOW()"}).
		AddRow(time.Now())

	mock.ExpectQuery("SELECT NOW()").WillReturnRows(rows)
	mock.ExpectExec("create table if not exists `kamino_user_migrations` \\(`id` varchar\\(255\\) not null primary key, `applied_at` datetime\\) engine=InnoDB charset=UTF8;").WillReturnResult(sqlmock.NewResult(1, 1))
	rows = sqlmock.NewRows([]string{"id", "applied_at"}).
		AddRow("v0099.sql", time.Now())
	mock.ExpectQuery("SELECT \\* FROM `kamino_user_migrations`").WillReturnRows(rows)

	_, steps, err := migration.Load(ctx, log, "testdata/good", "migrationok", 0, v, dss, false, false, nil)
	if err != nil {
		t.Fatalf("Load should not returns an error, returned: %v", err)
	}

	superseed := make(map[string]string)
	superseed["migration.dir"] = "down"
	superseed["migration.limit"] = "1"
	steps[0].PostLoad(log, superseed)
	err = steps[0].Init(ctx, log)
	if err != nil {
		t.Fatalf("Init should not returns an error, returned: %v", err)
	}

	err = steps[0].Do(context.Background(), log)
	if err == nil {
		t.Errorf("Do should return error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDoPrintOk(t *testing.T) {
	ctx, log, dss, v, mock, err := setupDo("testdata/good/steps/", "migrationok")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}
	rows := sqlmock.NewRows([]string{"NOW()"}).
		AddRow(time.Now())

	mock.ExpectQuery("SELECT NOW()").WillReturnRows(rows)
	mock.ExpectExec("create table if not exists `kamino_admin_migrations` \\(`id` varchar\\(255\\) not null primary key, `applied_at` datetime\\) engine=InnoDB charset=UTF8;").WillReturnResult(sqlmock.NewResult(1, 1))
	rows = sqlmock.NewRows([]string{"id", "applied_at"}).
		AddRow("v0001.sql", time.Now())
	mock.ExpectQuery("SELECT \\* FROM `kamino_admin_migrations` ORDER BY id ASC").WillReturnRows(rows)
	rows = sqlmock.NewRows([]string{"NOW()"}).
		AddRow(time.Now())
	mock.ExpectQuery("SELECT NOW()").WillReturnRows(rows)
	mock.ExpectExec("create table if not exists `kamino_user_migrations` \\(`id` varchar\\(255\\) not null primary key, `applied_at` datetime\\) engine=InnoDB charset=UTF8;").WillReturnResult(sqlmock.NewResult(1, 1))
	rows = sqlmock.NewRows([]string{"id", "applied_at"}).
		AddRow("v0001.sql", time.Now())
	mock.ExpectQuery("SELECT \\* FROM `kamino_user_migrations`").WillReturnRows(rows)
	_, steps, err := migration.Load(ctx, log, "testdata/good", "migrationok", 0, v, dss, false, false, nil)
	if err != nil {
		t.Fatalf("Load should not returns an error, returned: %v", err)
	}

	superseed := make(map[string]string)
	superseed["migration.dir"] = "status"
	steps[0].PostLoad(log, superseed)
	err = steps[0].Init(ctx, log)
	if err != nil {
		t.Fatalf("Init should not returns an error, returned: %v", err)
	}

	err = steps[0].Do(context.Background(), log)
	if err != nil {
		t.Errorf("Do should not return error, returned: %v", err)
	}

	steps[0].Cancel(log)
	steps[0].Finish(log)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDoDryRunOk(t *testing.T) {
	ctx, log, dss, v, mock, err := setupDo("testdata/good/steps/", "migrationok")
	if err != nil {
		t.Errorf("SetupLoad should not returns an error, returned: %v", err)
	}
	rows := sqlmock.NewRows([]string{"NOW()"}).
		AddRow(time.Now())

	mock.ExpectQuery("SELECT NOW()").WillReturnRows(rows)
	mock.ExpectExec("create table if not exists `kamino_user_migrations` \\(`id` varchar\\(255\\) not null primary key, `applied_at` datetime\\) engine=InnoDB charset=UTF8;").WillReturnResult(sqlmock.NewResult(1, 1))
	rows = sqlmock.NewRows([]string{"id", "applied_at"}).
		AddRow("v0001.sql", time.Now())
	mock.ExpectQuery("SELECT \\* FROM `kamino_user_migrations`").WillReturnRows(rows)
	rows = sqlmock.NewRows([]string{"NOW()"}).
		AddRow(time.Now())
	mock.ExpectQuery("SELECT NOW()").WillReturnRows(rows)
	mock.ExpectExec("create table if not exists `kamino_admin_migrations` \\(`id` varchar\\(255\\) not null primary key, `applied_at` datetime\\) engine=InnoDB charset=UTF8;").WillReturnResult(sqlmock.NewResult(1, 1))
	rows = sqlmock.NewRows([]string{"id", "applied_at"}).
		AddRow("v0001.sql", time.Now())
	mock.ExpectQuery("SELECT \\* FROM `kamino_admin_migrations`").WillReturnRows(rows)

	_, steps, err := migration.Load(ctx, log, "testdata/good", "migrationok", 0, v, dss, false, true, nil)
	if err != nil {
		t.Fatalf("Load should not returns an error, returned: %v", err)
	}

	superseed := make(map[string]string)
	superseed["migration.dir"] = "down"
	steps[0].PostLoad(log, superseed)
	err = steps[0].Init(ctx, log)
	if err != nil {
		t.Fatalf("Init should not returns an error, returned: %v", err)
	}

	err = steps[0].Do(context.Background(), log)
	if err != nil {
		t.Errorf("Do should not return error, returned: %v", err)
	}

	steps[0].Cancel(log)
	steps[0].Finish(log)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
