package datasource

import (
	"testing"

	"github.com/Sirupsen/logrus"
)

// We are using private function, we must be in same package
func setupDatabaseTest() *Datasources {
	return &Datasources{datasources: make(map[string]*Datasource)}
}

func TestLoadMysqlCompleteEngine(t *testing.T) {
	dss := setupDatabaseTest()
	ds, err := dss.load("testdata/good", "mysqlcomplete")
	if err != nil {
		t.Errorf("Load returns an error %v", err)
	}

	if ds.dstype != Database {
		t.Errorf("Should be recognized as database datasource")
	}

	if ds.engine != Mysql {
		t.Errorf("Should be recognized as Mysql datasource but was recognized as '%s'", EngineToString(ds.GetEngine()))
	}

	if ds.database != "dbmc" {
		t.Errorf("The database is '%s'", ds.database)
	}

	if ds.schema != "" {
		t.Errorf("The schema is '%s'", ds.schema)
	}

	if ds.url != "bob:123soleil@tcp(hmc:1234)/dbmc?parseTime=true" {
		t.Errorf("The user url is '%s'", ds.url)
	}

	if ds.urlAdmin != "super:adminpw@tcp(hmc:1234)/dbmc?parseTime=true" {
		t.Errorf("The admin url is '%s'", ds.urlAdmin)
	}

	if ds.urlNoDb != "super:adminpw@tcp(hmc:1234)/mysql?parseTime=true" {
		t.Errorf("The nodb url is '%s'", ds.urlNoDb)
	}

	if !ds.transaction {
		t.Errorf("Should have transaction")
	}

	if len(ds.tags) != 0 && ds.tags[0] != "tagmc" {
		t.Errorf("The tag should be found")
	}
}

func TestLoadMysqlMinimalEngine(t *testing.T) {
	dss := setupDatabaseTest()
	ds, err := dss.load("testdata/good", "mysqlminimal")
	if err != nil {
		t.Errorf("Load returns an error %v", err)
	}

	if ds.dstype != Database {
		t.Errorf("Should be recognized as database datasource")
	}

	if ds.engine != Mysql {
		t.Errorf("Should be recognized as Mysql datasource but was recognized as '%s'", EngineToString(ds.GetEngine()))
	}

	if ds.database != "dbmm" {
		t.Errorf("The database is '%s'", ds.database)
	}

	if ds.schema != "" {
		t.Errorf("The schema is '%s'", ds.schema)
	}

	if ds.url != "root:123soleil@tcp(127.0.0.1:3306)/dbmm?parseTime=true" {
		t.Errorf("The user url is '%s'", ds.url)
	}

	if ds.urlAdmin != "root:123soleil@tcp(127.0.0.1:3306)/dbmm?parseTime=true" {
		t.Errorf("The admin url is '%s'", ds.urlAdmin)
	}

	if ds.urlNoDb != "root:123soleil@tcp(127.0.0.1:3306)/mysql?parseTime=true" {
		t.Errorf("The nodb url is '%s'", ds.urlNoDb)
	}

	if ds.transaction {
		t.Errorf("Should not have transaction")
	}

	if len(ds.tags) != 0 && ds.tags[0] != "tagmm" {
		t.Errorf("The tag should be found")
	}
}

func TestLoadPostgresCompleteEngine(t *testing.T) {
	dss := setupDatabaseTest()
	ds, err := dss.load("testdata/good", "postgrescomplete")
	if err != nil {
		t.Errorf("Load returns an error %v", err)
	}

	if ds.dstype != Database {
		t.Errorf("Should be recognized as database datasource")
	}

	if ds.engine != Postgres {
		t.Errorf("Should be recognized as Postgres datasource but was recognized as '%s'", EngineToString(ds.GetEngine()))
	}

	if ds.database != "dbpc" {
		t.Errorf("The database is '%s'", ds.database)
	}

	if ds.schema != "shpc" {
		t.Errorf("The schema is '%s'", ds.schema)
	}

	if ds.url != "host=hpc port=1234 user=bob password=123soleil dbname=dbpc sslmode=disable" {
		t.Errorf("The user url is '%s'", ds.url)
	}

	if ds.urlAdmin != "host=hpc port=1234 user=super password=adminpw dbname=dbpc sslmode=disable" {
		t.Errorf("The admin url is '%s'", ds.urlAdmin)
	}

	if ds.urlNoDb != "host=hpc port=1234 user=super password=adminpw sslmode=disable" {
		t.Errorf("The nodb url is '%s'", ds.urlNoDb)
	}

	if !ds.transaction {
		t.Errorf("Should have transaction")
	}

	if len(ds.tags) != 0 && ds.tags[0] != "tagpc" {
		t.Errorf("The tag should be found")
	}
}

func TestLoadPostgresMinimalEngine(t *testing.T) {
	dss := setupDatabaseTest()
	ds, err := dss.load("testdata/good", "postgresminimal")
	if err != nil {
		t.Errorf("Load returns an error %v", err)
	}

	if ds.dstype != Database {
		t.Errorf("Should be recognized as database datasource")
	}

	if ds.engine != Postgres {
		t.Errorf("Should be recognized as Postgres datasource but was recognized as '%s'", EngineToString(ds.GetEngine()))
	}

	if ds.database != "dbpm" {
		t.Errorf("The database is '%s'", ds.database)
	}

	if ds.schema != "" {
		t.Errorf("The schema is '%s'", ds.schema)
	}

	if ds.url != "host=127.0.0.1 port=5432 user=postgres password=adminpw dbname=dbpm sslmode=disable" {
		t.Errorf("The user url is '%s'", ds.url)
	}

	if ds.urlAdmin != "host=127.0.0.1 port=5432 user=postgres password=adminpw dbname=dbpm sslmode=disable" {
		t.Errorf("The admin url is '%s'", ds.urlAdmin)
	}

	if ds.urlNoDb != "host=127.0.0.1 port=5432 user=postgres password=adminpw sslmode=disable" {
		t.Errorf("The nodb url is '%s'", ds.urlNoDb)
	}

	if ds.transaction {
		t.Errorf("Should not have transaction")
	}

	if len(ds.tags) != 0 && ds.tags[0] != "tagpm" {
		t.Errorf("The tag should be found")
	}
	if ds.IsTransaction() {
		t.Errorf("The datasource should not have transaction")
	}
}

func TestLoadNoDatabase(t *testing.T) {
	dss := setupDatabaseTest()
	_, err := dss.load("testdata/fail", "nodatabase")
	if err == nil {
		t.Errorf("Load should returns an error")
	}
}

func TestDatabaseOpenWrongType(t *testing.T) {

	ds := Datasource{engine: JSON, dstype: File}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")

	if _, err := ds.OpenDatabase(log, false, false); err == nil {
		t.Errorf("OpenDatabase should returns an error")
	}
}

func TestDatabaseOpenMysql(t *testing.T) {
	mockingSQL = true
	ds := Datasource{engine: Mysql, dstype: Database, url: "bob:123soleil@tcp(localhost:1234)/dbmc", urlAdmin: "urlAdmin", urlNoDb: "urlNoDb"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	if _, err := ds.OpenDatabase(log, false, false); err != nil {
		t.Errorf("OpenDatabase should not returns an error, was: %v", err)
	}
}

func TestDatabaseOpenPostgres(t *testing.T) {
	mockingSQL = true
	ds := Datasource{engine: Postgres, dstype: Database, url: "url", urlAdmin: "urlAdmin", urlNoDb: "urlNoDb"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	if _, err := ds.OpenDatabase(log, false, false); err != nil {
		t.Errorf("OpenDatabase should not returns an error, was: %v", err)
	}
}

func TestDatabaseOpenAdmin(t *testing.T) {
	mockingSQL = true
	ds := Datasource{engine: Mysql, dstype: Database, url: "url", urlAdmin: "bob:123soleil@tcp(localhost:1234)/dbmc", urlNoDb: "urlNoDb"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	if _, err := ds.OpenDatabase(log, true, false); err != nil {
		t.Errorf("OpenDatabase should not returns an error, was: %v", err)
	}
}
func TestDatabaseOpenNoDb(t *testing.T) {
	mockingSQL = true
	ds := Datasource{engine: Mysql, dstype: Database, url: "url", urlAdmin: "urlAdmin", urlNoDb: "bob:123soleil@tcp(localhost:1234)/dbmc"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	if _, err := ds.OpenDatabase(log, false, true); err != nil {
		t.Errorf("OpenDatabase should not returns an error, was: %v", err)
	}
}

func TestDatabaseReOpen(t *testing.T) {
	mockingSQL = true
	ds := Datasource{engine: Mysql, dstype: Database, url: "bob:123soleil@tcp(localhost:1234)/dbmc", urlAdmin: "urlAdmin", urlNoDb: "urlNoDb"}
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	db1, err := ds.OpenDatabase(log, false, false)
	if err != nil {
		t.Errorf("OpenDatabase should not returns an error, was: %v", err)
	}
	db2, err := ds.OpenDatabase(log, false, false)
	if err != nil {
		t.Errorf("OpenDatabase should not returns an error, was: %v", err)
	}

	if db1 != db2 {
		t.Errorf("Reopenning database should return the same object")
	}
}
func TestLoadNoTags(t *testing.T) {
	dss := setupDatabaseTest()
	_, err := dss.load("testdata/good", "mysqlnotag")
	if err != nil {
		t.Errorf("Load should not returns an error")
	}
}

func TestDatabaseCloseAll(t *testing.T) {
	mockingSQL = true
	dss := setupDatabaseTest()
	ds, err := dss.load("testdata/good", "postgresminimal")
	if err != nil {
		t.Errorf("Load returns an error %v", err)
	}
	dss.datasources["test"] = &ds
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	_, err = ds.OpenDatabase(log, false, false)
	if err != nil {
		t.Errorf("OpenDatabase should not returns an error, was: %v", err)
	}
	_, err = ds.OpenDatabase(log, true, false)
	if err != nil {
		t.Errorf("OpenDatabase should not returns an error, was: %v", err)
	}
	_, err = ds.OpenDatabase(log, false, true)
	if err != nil {
		t.Errorf("OpenDatabase should not returns an error, was: %v", err)
	}

	dss.CloseAll(log)
}
