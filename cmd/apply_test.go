package cmd_test

import (
	"fmt"
	"testing"

	"github.com/marema31/kamino/cmd"
)

func TestApplyOk(t *testing.T) {
	ck := &mockedCookbook{}
	err := cmd.Apply(ck, nil, nil, []string{"recipe1ok", "recipe2ok"})
	if err != nil {
		t.Errorf("Apply should not returns an error, returned: %v", err)
	}

	if !ck.called {
		t.Errorf("Do should be called")
	}
}

func TestApplyLoadError(t *testing.T) {
	ck := &mockedCookbook{}
	ck.errorLoad = fmt.Errorf("fake error")
	err := cmd.Apply(ck, nil, nil, []string{"recipe1ok", "recipe2ok"})
	if err == nil {
		t.Errorf("Apply should returns an error")
	}

	if ck.called {
		t.Errorf("Do should not be called")
	}
}

func TestApplyDoError(t *testing.T) {
	ck := &mockedCookbook{}
	ck.doReturnValue = true
	err := cmd.Apply(ck, nil, nil, []string{"recipe1ok", "recipe2ok"})
	if err == nil {
		t.Errorf("Apply should returns an error")
	}

	if !ck.called {
		t.Errorf("Do should be called")
	}
}

func TestApplyFindRecipesError(t *testing.T) {
	ck := &mockedCookbook{}
	err := cmd.Apply(ck, nil, nil, []string{})
	if err == nil {
		t.Errorf("Apply should returns an error")
	}

	if ck.called {
		t.Errorf("Do should not be called")
	}
}
