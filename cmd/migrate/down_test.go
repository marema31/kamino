package migrate_test

import (
	"fmt"
	"testing"

	"github.com/marema31/kamino/cmd/common"
	"github.com/marema31/kamino/cmd/migrate"
)

func TestDownOk(t *testing.T) {
	common.CfgFolder = "testdata/good"
	ck := &mockedCookbook{}
	migrate.Admin = false
	migrate.User = false
	err := migrate.Down(ck, []string{"recipe1ok", "recipe2ok"})
	if err != nil {
		t.Errorf("Down should not returns an error, returned: %v", err)
	}

	if !ck.called {
		t.Errorf("Do should be called")
	}
}

func TestDownFindRecipesError(t *testing.T) {
	common.CfgFolder = "testdata"
	ck := &mockedCookbook{}
	err := migrate.Down(ck, []string{})
	if err == nil {
		t.Errorf("Down should returns an error")
	}

	if ck.called {
		t.Errorf("Do should not be called")
	}
}

func TestDowncreateSDownerseedError(t *testing.T) {
	common.CfgFolder = "testdata/good"
	ck := &mockedCookbook{}
	migrate.Admin = true
	migrate.User = true
	err := migrate.Down(ck, []string{"recipe1ok", "recipe2ok"})
	if err == nil {
		t.Errorf("Down should returns an error")
	}

	if ck.called {
		t.Errorf("Do should not be called")
	}
}
func TestDownLoadError(t *testing.T) {
	common.CfgFolder = "testdata/good"
	ck := &mockedCookbook{}
	migrate.Admin = false
	migrate.User = false
	ck.errorLoad = fmt.Errorf("fake error")
	err := migrate.Down(ck, []string{"recipe1ok", "recipe2ok"})
	if err == nil {
		t.Errorf("Down should returns an error")
	}

	if ck.called {
		t.Errorf("Do should not be called")
	}
}

func TestDownPostLoadError(t *testing.T) {
	common.CfgFolder = "testdata/good"
	ck := &mockedCookbook{}
	migrate.Admin = false
	migrate.User = false
	ck.errorPostLoad = fmt.Errorf("fake error")
	err := migrate.Down(ck, []string{"recipe1ok", "recipe2ok"})
	if err == nil {
		t.Errorf("Down should returns an error")
	}

	if ck.called {
		t.Errorf("Do should not be called")
	}
}
func TestDownDoError(t *testing.T) {
	common.CfgFolder = "testdata/good"
	ck := &mockedCookbook{}
	migrate.Admin = false
	migrate.User = false
	ck.doReturnValue = true
	err := migrate.Down(ck, []string{"recipe1ok", "recipe2ok"})
	if err == nil {
		t.Errorf("Down should returns an error")
	}

	if !ck.called {
		t.Errorf("Do should be called")
	}
}
