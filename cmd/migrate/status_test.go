package migrate_test

import (
	"fmt"
	"testing"

	"github.com/marema31/kamino/cmd/migrate"
)

func TestStatusOk(t *testing.T) {
	ck := &mockedCookbook{}
	migrate.Admin = false
	migrate.User = false
	err := migrate.Status(ck, []string{"recipe1ok", "recipe2ok"})
	if err != nil {
		t.Errorf("Status should not returns an error, returned: %v", err)
	}

	if !ck.called {
		t.Errorf("Do should be called")
	}
}

func TestStatusFindRecipesError(t *testing.T) {
	ck := &mockedCookbook{}
	err := migrate.Status(ck, []string{})
	if err == nil {
		t.Errorf("Status should returns an error")
	}

	if ck.called {
		t.Errorf("Do should not be called")
	}
}

func TestStatuscreateSStatuserseedError(t *testing.T) {
	ck := &mockedCookbook{}
	migrate.Admin = true
	migrate.User = true
	err := migrate.Status(ck, []string{"recipe1ok", "recipe2ok"})
	if err == nil {
		t.Errorf("Status should returns an error")
	}

	if ck.called {
		t.Errorf("Do should not be called")
	}
}
func TestStatusLoadError(t *testing.T) {
	ck := &mockedCookbook{}
	migrate.Admin = false
	migrate.User = false
	ck.errorLoad = fmt.Errorf("fake error")
	err := migrate.Status(ck, []string{"recipe1ok", "recipe2ok"})
	if err == nil {
		t.Errorf("Status should returns an error")
	}

	if ck.called {
		t.Errorf("Do should not be called")
	}
}

func TestStatusPostLoadError(t *testing.T) {
	ck := &mockedCookbook{}
	migrate.Admin = false
	migrate.User = false
	ck.errorPostLoad = fmt.Errorf("fake error")
	err := migrate.Status(ck, []string{"recipe1ok", "recipe2ok"})
	if err == nil {
		t.Errorf("Status should returns an error")
	}

	if ck.called {
		t.Errorf("Do should not be called")
	}
}
func TestStatusDoError(t *testing.T) {
	ck := &mockedCookbook{}
	migrate.Admin = false
	migrate.User = false
	ck.doReturnValue = true
	err := migrate.Status(ck, []string{"recipe1ok", "recipe2ok"})
	if err == nil {
		t.Errorf("Status should returns an error")
	}

	if !ck.called {
		t.Errorf("Do should be called")
	}
}
