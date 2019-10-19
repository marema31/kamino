package datasource

import (
	"testing"
)

// We are using private function, we must be in same package
func setupDatabaseTest() *Datasources {
	return &Datasources{}
}

func TestLoadMysqlCompleteEngine(t *testing.T) {
	dss := setupDatabaseTest()
	ds, err := dss.load("testdata/good/datasources", "mysqlcomplete")
	if err != nil {
		t.Errorf("Load returns an error %v", err)
	}

	if ds == nil {
		t.Fatalf("Should have been inserted in datasources with the name mysqlcomplete")
	}

	if ds.Type != Database {
		t.Errorf("Should be recognized as database datasource")
	}

	if ds.Engine != Mysql {
		t.Errorf("Should be recognized as Mysql datasource but was recognized as '%s'", ds.GetEngine())
	}

	if ds.Database != "dbmc" {
		t.Errorf("The database is '%s'", ds.Database)
	}

	if ds.Schema != "" {
		t.Errorf("The schema is '%s'", ds.Schema)
	}

	if ds.URL != "bob:123soleil@tcp(hmc:1234)/dbmc" {
		t.Errorf("The user URL is '%s'", ds.URL)
	}

	if ds.URLAdmin != "super:adminpw@tcp(hmc:1234)/dbmc" {
		t.Errorf("The admin URL is '%s'", ds.URLAdmin)
	}

	if ds.URLNoDb != "super:adminpw@tcp(hmc:1234)" {
		t.Errorf("The nodb URL is '%s'", ds.URLNoDb)
	}

	if !ds.Transaction {
		t.Errorf("Should have transaction")
	}

	if len(ds.Tags) != 0 && ds.Tags[0] != "tagmc" {
		t.Errorf("The tag should be found")
	}
}

func TestLoadMysqlMinimalEngine(t *testing.T) {
	dss := setupDatabaseTest()
	ds, err := dss.load("testdata/good/datasources", "mysqlminimal")
	if err != nil {
		t.Errorf("Load returns an error %v", err)
	}

	if ds == nil {
		t.Fatalf("Should have been inserted in datasources with the name mysqlminimal")
	}

	if ds.Type != Database {
		t.Errorf("Should be recognized as database datasource")
	}

	if ds.Engine != Mysql {
		t.Errorf("Should be recognized as Mysql datasource but was recognized as '%s'", ds.GetEngine())
	}

	if ds.Database != "dbmm" {
		t.Errorf("The database is '%s'", ds.Database)
	}

	if ds.Schema != "" {
		t.Errorf("The schema is '%s'", ds.Schema)
	}

	if ds.URL != "root:123soleil@tcp(127.0.0.1:3306)/dbmm" {
		t.Errorf("The user URL is '%s'", ds.URL)
	}

	if ds.URLAdmin != "root:123soleil@tcp(127.0.0.1:3306)/dbmm" {
		t.Errorf("The admin URL is '%s'", ds.URLAdmin)
	}

	if ds.URLNoDb != "root:123soleil@tcp(127.0.0.1:3306)" {
		t.Errorf("The nodb URL is '%s'", ds.URLNoDb)
	}

	if ds.Transaction {
		t.Errorf("Should not have transaction")
	}

	if len(ds.Tags) != 0 && ds.Tags[0] != "tagmm" {
		t.Errorf("The tag should be found")
	}
}

func TestLoadPostgresCompleteEngine(t *testing.T) {
	dss := setupDatabaseTest()
	ds, err := dss.load("testdata/good/datasources", "postgrescomplete")
	if err != nil {
		t.Errorf("Load returns an error %v", err)
	}

	if ds == nil {
		t.Fatalf("Should have been inserted in datasources with the name postgrescomplete")
	}

	if ds.Type != Database {
		t.Errorf("Should be recognized as database datasource")
	}

	if ds.Engine != Postgres {
		t.Errorf("Should be recognized as Postgres datasource but was recognized as '%s'", ds.GetEngine())
	}

	if ds.Database != "dbpc" {
		t.Errorf("The database is '%s'", ds.Database)
	}

	if ds.Schema != "shpc" {
		t.Errorf("The schema is '%s'", ds.Schema)
	}

	if ds.URL != "host=hpc port=1234 user=bob password=123soleil dbname=dbpc sslmode=disable" {
		t.Errorf("The user URL is '%s'", ds.URL)
	}

	if ds.URLAdmin != "host=hpc port=1234 user=super password=adminpw dbname=dbpc sslmode=disable" {
		t.Errorf("The admin URL is '%s'", ds.URLAdmin)
	}

	if ds.URLNoDb != "host=hpc port=1234 user=super password=adminpw sslmode=disable" {
		t.Errorf("The nodb URL is '%s'", ds.URLNoDb)
	}

	if !ds.Transaction {
		t.Errorf("Should have transaction")
	}

	if len(ds.Tags) != 0 && ds.Tags[0] != "tagpc" {
		t.Errorf("The tag should be found")
	}
}

func TestLoadPostgresMinimalEngine(t *testing.T) {
	dss := setupDatabaseTest()
	ds, err := dss.load("testdata/good/datasources", "postgresminimal")
	if err != nil {
		t.Errorf("Load returns an error %v", err)
	}

	if ds == nil {
		t.Fatalf("Should have been inserted in datasources with the name postgresminimal")
	}

	if ds.Type != Database {
		t.Errorf("Should be recognized as database datasource")
	}

	if ds.Engine != Postgres {
		t.Errorf("Should be recognized as Postgres datasource but was recognized as '%s'", ds.GetEngine())
	}

	if ds.Database != "dbpm" {
		t.Errorf("The database is '%s'", ds.Database)
	}

	if ds.Schema != "" {
		t.Errorf("The schema is '%s'", ds.Schema)
	}

	if ds.URL != "host=127.0.0.1 port=5432 user=postgres password=adminpw dbname=dbpm sslmode=disable" {
		t.Errorf("The user URL is '%s'", ds.URL)
	}

	if ds.URLAdmin != "host=127.0.0.1 port=5432 user=postgres password=adminpw dbname=dbpm sslmode=disable" {
		t.Errorf("The admin URL is '%s'", ds.URLAdmin)
	}

	if ds.URLNoDb != "host=127.0.0.1 port=5432 user=postgres password=adminpw sslmode=disable" {
		t.Errorf("The nodb URL is '%s'", ds.URLNoDb)
	}

	if ds.Transaction {
		t.Errorf("Should not have transaction")
	}

	if len(ds.Tags) != 0 && ds.Tags[0] != "tagpm" {
		t.Errorf("The tag should be found")
	}
}

func TestLoadNoDatabase(t *testing.T) {
	dss := setupDatabaseTest()
	_, err := dss.load("testdata/fail/datasources", "nodatabase")
	if err == nil {
		t.Errorf("Load should returns an error")
	}
}

func TestDatabaseOpenWrongType(t *testing.T) {

	ds := Datasource{Engine: JSON, Type: File}
	if _, err := ds.OpenDatabase(false, false); err == nil {
		t.Errorf("OpenDatabase should returns an error")
	}
}

func TestDatabaseOpenMysql(t *testing.T) {
	mockingSQL = true
	ds := Datasource{Engine: Mysql, Type: Database, URL: "bob:123soleil@tcp(localhost:1234)/dbmc", URLAdmin: "URLAdmin", URLNoDb: "URLNoDb"}
	if _, err := ds.OpenDatabase(false, false); err != nil {
		t.Errorf("OpenDatabase should not returns an error, was: %v", err)
	}
}

func TestDatabaseOpenPostgres(t *testing.T) {
	mockingSQL = true
	ds := Datasource{Engine: Postgres, Type: Database, URL: "URL", URLAdmin: "URLAdmin", URLNoDb: "URLNoDb"}
	if _, err := ds.OpenDatabase(false, false); err != nil {
		t.Errorf("OpenDatabase should not returns an error, was: %v", err)
	}
}

func TestDatabaseOpenURLAdmin(t *testing.T) {
	mockingSQL = true
	ds := Datasource{Engine: Mysql, Type: Database, URL: "URL", URLAdmin: "bob:123soleil@tcp(localhost:1234)/dbmc", URLNoDb: "URLNoDb"}
	if _, err := ds.OpenDatabase(true, false); err != nil {
		t.Errorf("OpenDatabase should not returns an error, was: %v", err)
	}
}
func TestDatabaseOpenURLNoDb(t *testing.T) {
	mockingSQL = true
	ds := Datasource{Engine: Mysql, Type: Database, URL: "URL", URLAdmin: "URLAdmin", URLNoDb: "bob:123soleil@tcp(localhost:1234)/dbmc"}
	if _, err := ds.OpenDatabase(false, true); err != nil {
		t.Errorf("OpenDatabase should not returns an error, was: %v", err)
	}
}
