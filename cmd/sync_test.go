package cmd_test

import (
	"fmt"
	"testing"

	"github.com/marema31/kamino/cmd"
)

func TestSyncOk(t *testing.T) {
	ck := &mockedCookbook{}
	err := cmd.Sync(ck, nil, []string{"recipe1ok", "recipe2ok"})
	if err != nil {
		t.Errorf("Load should not returns an error, returned: %v", err)
	}
}

func TestSyncLoadError(t *testing.T) {
	ck := &mockedCookbook{}
	ck.errorLoad = fmt.Errorf("fake error")
	err := cmd.Sync(ck, nil, []string{"recipe1ok", "recipe2ok"})
	if err == nil {
		t.Errorf("Sync should returns an error")
	}

	if ck.called {
		t.Errorf("Do should not be called")
	}
}

func TestSyncDoError(t *testing.T) {
	ck := &mockedCookbook{}
	ck.doReturnValue = true
	err := cmd.Sync(ck, nil, []string{"recipe1ok", "recipe2ok"})
	if err == nil {
		t.Errorf("Sync should returns an error")
	}

	if !ck.called {
		t.Errorf("Do should be called")
	}
}

func TestSyncFindRecipesError(t *testing.T) {
	ck := &mockedCookbook{}
	err := cmd.Sync(ck, nil, []string{})
	if err == nil {
		t.Errorf("Sync should returns an error")
	}

	if ck.called {
		t.Errorf("Do should not be called")
	}
}
