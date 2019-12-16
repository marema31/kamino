package shell_test

import (
	"context"
	"testing"

	"github.com/marema31/kamino/step/shell"
)

func TestToSkipOk(t *testing.T) {
	ctx, log, dss, v := setupDo("testdata/good/steps/", "shellok")

	_, steps, err := shell.Load(ctx, log, "testdata/good", "shellok", 0, v, dss, false, false)
	if err != nil {
		t.Fatalf("Load should not returns an error, returned: %v", err)
	}

	ok, err := steps[0].ToSkip(context.Background(), log)
	if err != nil {
		t.Errorf("ToSkip should not return error, returned: %v", err)
	}

	if ok {
		t.Error("ToSkip should return false")
	}
}
