package tmpl_test

import (
	"context"
	"os"
	"testing"

	"github.com/marema31/kamino/step/tmpl"
)

func TestToSkipOnlyIfNotExists(t *testing.T) {
	ctx, log, dss, v := setupDo("testdata/good/steps/", "ifnotexist")

	_, steps, err := tmpl.Load(ctx, log, "testdata/good", "ifnotexist", 0, v, dss, false, false)
	if err != nil {
		t.Fatalf("Load should not returns an error, returned: %v", err)
	}

	dummy, _ := os.Create("testdata/tmp/alreadyexists.yaml")
	dummy.Close()

	ok, err := steps[0].ToSkip(context.Background(), log)
	if err != nil {
		t.Errorf("ToSkip should not return error, returned: %v", err)
	}

	if !ok {
		t.Error("ToSkip should return true")
	}

	os.Remove("testdata/tmp/alreadyexists.yaml")
	ok, err = steps[0].ToSkip(context.Background(), log)
	if err != nil {
		t.Errorf("ToSkip should not return error, returned: %v", err)
	}

	if ok {
		t.Error("ToSkip should return false")
	}

}

func TestToSkipOk(t *testing.T) {
	ctx, log, dss, v := setupDo("testdata/good/steps/", "ifnotexist")

	_, steps, err := tmpl.Load(ctx, log, "testdata/good", "tmplok", 0, v, dss, false, false)
	if err != nil {
		t.Fatalf("Load should not returns an error, returned: %v", err)
	}

	dummy, _ := os.Create("testdata/tmp/replace.cfg")
	dummy.Close()

	ok, err := steps[0].ToSkip(context.Background(), log)
	if err != nil {
		t.Errorf("ToSkip should not return error, returned: %v", err)
	}

	if ok {
		t.Error("ToSkip should return false")
	}

	os.Remove("testdata/tmp/replace.cfg")
}
