package cmd_test

import (
	"fmt"
	"testing"

	"github.com/marema31/kamino/cmd"
)

func TestValidateOk(t *testing.T) {
	ck := &mockedCookbook{}
	err := cmd.Validate(ck, nil, []string{"recipe1ok", "recipe2ok"})
	if err != nil {
		t.Errorf("validate should not returns an error, returned: %v", err)
	}
}

func TestValidateLoadError(t *testing.T) {
	ck := &mockedCookbook{}
	ck.errorLoad = fmt.Errorf("fake error")
	err := cmd.Validate(ck, nil, []string{"recipe1ok", "recipe2ok"})
	if err == nil {
		t.Errorf("validate should returns an error")
	}

	if ck.called {
		t.Errorf("Do should not be called")
	}
}

func TestValidateFindRecipesError(t *testing.T) {
	ck := &mockedCookbook{}
	err := cmd.Validate(ck, nil, []string{})
	if err == nil {
		t.Errorf("validate should returns an error")
	}

	if ck.called {
		t.Errorf("Do should not be called")
	}
}
