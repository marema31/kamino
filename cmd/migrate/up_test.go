package migrate_test

import (
	"fmt"
	"testing"

	"github.com/marema31/kamino/cmd/common"
	"github.com/marema31/kamino/cmd/migrate"
)

func TestUpOk(t *testing.T) {
	common.CfgFolder = "testdata/good"
	ck := &mockedCookbook{}
	migrate.Admin = false
	migrate.User = false
	err := migrate.Up(ck, []string{"recipe1ok", "recipe2ok"})
	if err != nil {
		t.Errorf("Up should not returns an error, returned: %v", err)
	}

	if !ck.called {
		t.Errorf("Do should be called")
	}
}

func TestUpFindRecipesError(t *testing.T) {
	common.CfgFolder = "testdata"
	ck := &mockedCookbook{}
	err := migrate.Up(ck, []string{})
	if err == nil {
		t.Errorf("Up should returns an error")
	}

	if ck.called {
		t.Errorf("Do should not be called")
	}
}

func TestUpcreateSuperseedError(t *testing.T) {
	common.CfgFolder = "testdata/good"
	ck := &mockedCookbook{}
	migrate.Admin = true
	migrate.User = true
	err := migrate.Up(ck, []string{"recipe1ok", "recipe2ok"})
	if err == nil {
		t.Errorf("Up should returns an error")
	}

	if ck.called {
		t.Errorf("Do should not be called")
	}
}
func TestUpLoadError(t *testing.T) {
	common.CfgFolder = "testdata/good"
	ck := &mockedCookbook{}
	migrate.Admin = false
	migrate.User = false
	ck.errorLoad = fmt.Errorf("fake error")
	err := migrate.Up(ck, []string{"recipe1ok", "recipe2ok"})
	if err == nil {
		t.Errorf("Up should returns an error")
	}

	if ck.called {
		t.Errorf("Do should not be called")
	}
}

func TestUpPostLoadError(t *testing.T) {
	common.CfgFolder = "testdata/good"
	ck := &mockedCookbook{}
	migrate.Admin = false
	migrate.User = false
	ck.errorPostLoad = fmt.Errorf("fake error")
	err := migrate.Up(ck, []string{"recipe1ok", "recipe2ok"})
	if err == nil {
		t.Errorf("Up should returns an error")
	}

	if ck.called {
		t.Errorf("Do should not be called")
	}
}
func TestUpDoError(t *testing.T) {
	common.CfgFolder = "testdata/good"
	ck := &mockedCookbook{}
	migrate.Admin = false
	migrate.User = false
	ck.doReturnValue = true
	err := migrate.Up(ck, []string{"recipe1ok", "recipe2ok"})
	if err == nil {
		t.Errorf("Up should returns an error")
	}

	if !ck.called {
		t.Errorf("Do should be called")
	}
}
