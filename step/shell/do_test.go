package shell_test

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/marema31/kamino/step/shell"
)

func TestInitOk(t *testing.T) {
	ctx, log, dss, v := setupDo("testdata/good/steps/", "shellok")

	_, steps, err := shell.Load(ctx, log, "testdata/good", "shellok", 0, v, dss, false, false, nil)
	if err != nil {
		t.Fatalf("Load should not returns an error, returned: %v", err)
	}

	dummy, _ := os.Create("testdata/tmp/notify.sh")
	dummy.Close()

	err = steps[0].Init(ctx, log)
	if err != nil {
		t.Fatalf("Init should not returns an error, returned: %v", err)
	}
	os.Remove("testdata/tmp/notify.sh")

	err = steps[0].Init(ctx, log)
	if err == nil {
		t.Fatalf("Init should returns an error")
	}
}

func TestDoFinishOk(t *testing.T) {
	ctx, log, dss, v := setupDo("testdata/good/steps/", "shellok")

	_, steps, err := shell.Load(ctx, log, "testdata/good", "shellok", 0, v, dss, false, false, nil)
	if err != nil {
		t.Fatalf("Load should not returns an error, returned: %v", err)
	}

	st, ok := steps[0].(*shell.Step)
	if !ok {
		t.Fatalf("The step should be a shell step")
	}

	st.FakeCmd(0, []string{})
	err = steps[0].Do(context.Background(), log)
	if err != nil {
		t.Errorf("Do should not return error, returned: %v", err)
	}

	steps[0].Finish(log)
}

func TestDoCancelOk(t *testing.T) {
	ctx, log, dss, v := setupDo("testdata/good/steps/", "shellok")

	_, steps, err := shell.Load(ctx, log, "testdata/good", "shellok", 0, v, dss, false, false, nil)
	if err != nil {
		t.Fatalf("Load should not returns an error, returned: %v", err)
	}

	st, ok := steps[0].(*shell.Step)
	if !ok {
		t.Fatalf("The step should be a shell step")
	}

	st.FakeCmd(1, []string{})
	err = steps[0].Do(context.Background(), log)
	if err == nil {
		t.Errorf("Do should return error")
	}

	steps[0].Cancel(log)
}

func TestDoDryRun(t *testing.T) {
	ctx, log, dss, v := setupDo("testdata/good/steps/", "shellok")

	_, steps, err := shell.Load(ctx, log, "testdata/good", "shellok", 0, v, dss, false, true, nil)
	if err != nil {
		t.Fatalf("Load should not returns an error, returned: %v", err)
	}

	st, ok := steps[0].(*shell.Step)
	if !ok {
		t.Fatalf("The step should be a shell step")
	}

	st.FakeCmd(1, []string{})
	err = steps[0].Do(context.Background(), log)
	if err != nil {
		t.Errorf("Do should not returns an error, returned: %v", err)
	}

	steps[0].Cancel(log)
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	if rc := os.Getenv("TESTHELPEREXIT"); rc != "" {
		ec, _ := strconv.Atoi(rc)
		os.Exit(ec)
	}
	os.Exit(0)
}
