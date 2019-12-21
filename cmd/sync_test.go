package cmd_test

import (
	"fmt"
	"testing"

	"github.com/marema31/kamino/cmd"
	"github.com/marema31/kamino/cmd/common"
)

func TestSyncOk(t *testing.T) {
	common.CfgFolder = "testdata/good"
	ck := &mockedCookbook{}
	err := cmd.Sync(ck, nil, []string{"recipe1ok", "recipe2ok"})
	if err != nil {
		t.Errorf("Load should not returns an error, returned: %v", err)
	}
}

func TestSyncLoadError(t *testing.T) {
	common.CfgFolder = "testdata/good"
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

func TestSyncPostLoadError(t *testing.T) {
	common.CfgFolder = "testdata/good"
	ck := &mockedCookbook{}
	ck.errorPostLoad = fmt.Errorf("fake error")
	err := cmd.Sync(ck, nil, []string{"recipe1ok", "recipe2ok"})
	if err == nil {
		t.Errorf("Sync should returns an error")
	}

	if ck.called {
		t.Errorf("Do should not be called")
	}
}

func TestSyncSuperseed(t *testing.T) {
	common.CfgFolder = "testdata/good"
	ck := &mockedCookbook{}
	cmd.CacheOnly = false
	err := cmd.Sync(ck, nil, []string{"recipe1ok", "recipe2ok"})
	if err != nil {
		t.Errorf("Sync should not returns an error, returned: %v", err)
	}

	if v, ok := ck.superseed["sync.forceCacheOnly"]; ok {
		t.Errorf("Should no return forceCacheOnly, returned %s", v)
	}

	cmd.CacheOnly = true
	err = cmd.Sync(ck, nil, []string{"recipe1ok", "recipe2ok"})
	if err != nil {
		t.Errorf("Sync should not returns an error, returned: %v", err)
	}

	if _, ok := ck.superseed["sync.forceCacheOnly"]; !ok {
		t.Errorf("Should return forceCacheOnly")
	}
	if v, ok := ck.superseed["sync.forceCacheOnlye"]; ok && v != "true" {
		t.Errorf("Should return forceCacheOnly= true, returned %s", v)
	}

}

func TestSyncDoError(t *testing.T) {
	common.CfgFolder = "testdata/good"
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
	common.CfgFolder = "testdata"
	ck := &mockedCookbook{}
	err := cmd.Sync(ck, nil, []string{})
	if err == nil {
		t.Errorf("Sync should returns an error")
	}

	if ck.called {
		t.Errorf("Do should not be called")
	}
}
