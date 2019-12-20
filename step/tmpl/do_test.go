package tmpl_test

import (
	"bufio"
	"context"
	"os"
	"testing"

	"github.com/marema31/kamino/step/tmpl"
)

func TestDoFinishOk(t *testing.T) {
	ctx, log, dss, v := setupDo("testdata/good/steps/", "tmplok")

	_, steps, err := tmpl.Load(ctx, log, "testdata/good", "tmplok", 0, v, dss, false, false, nil)
	if err != nil {
		t.Fatalf("Load should not returns an error, returned: %v", err)
	}

	err = steps[0].Init(ctx, log)
	if err != nil {
		t.Fatalf("Init should not returns an error, returned: %v", err)
	}

	toskip, err := steps[0].ToSkip(ctx, log)
	if err != nil {
		t.Fatalf("ToSkip should not returns an error, returned: %v", err)
	}
	if toskip == true {
		t.Fatalf("ToSkip should always returns false, returned: %v", toskip)
	}

	err = steps[0].Do(context.Background(), log)
	if err != nil {
		t.Errorf("Do should not return error, returned: %v", err)
	}

	steps[0].Finish(log)
	os.Remove("testdata/tmp/db3.cfg")
	os.Remove("testdata/tmp/db4.cfg")
}

func TestDoCancelOk(t *testing.T) {
	ctx, log, dss, v := setupDo("testdata/good/steps/", "tmplok")

	_, steps, err := tmpl.Load(ctx, log, "testdata/good", "tmplok", 0, v, dss, false, false, nil)
	if err != nil {
		t.Fatalf("Load should not returns an error, returned: %v", err)
	}

	err = steps[0].Init(ctx, log)
	if err != nil {
		t.Fatalf("Init should not returns an error, returned: %v", err)
	}

	toskip, err := steps[0].ToSkip(ctx, log)
	if err != nil {
		t.Fatalf("ToSkip should not returns an error, returned: %v", err)
	}
	if toskip == true {
		t.Fatalf("ToSkip should always returns false, returned: %v", toskip)
	}

	err = steps[0].Do(context.Background(), log)
	if err != nil {
		t.Errorf("Do should not return error, returned: %v", err)
	}

	steps[0].Cancel(log)
}

func nbLineByFile(fileName string) int {
	file, _ := os.Open(fileName)
	fileScanner := bufio.NewScanner(file)
	lineCount := 0
	for fileScanner.Scan() {
		lineCount++
	}
	return lineCount
}

func helperDoTest(t *testing.T, configFile string, createdFile string, nblines1 int, nblines2 int, dryRun bool) {
	ctx, log, dss, v := setupDo("testdata/good/steps/", configFile)

	_, steps, err := tmpl.Load(ctx, log, "testdata/good", configFile, 0, v, dss, false, dryRun, nil)
	if err != nil {
		t.Fatalf("Load should not returns an error, returned: %v", err)
	}

	err = steps[0].Init(ctx, log)
	if err != nil {
		t.Fatalf("Init should not returns an error, returned: %v", err)
	}
	err = steps[0].Do(context.Background(), log)
	if err != nil {
		t.Errorf("Do should not return error, returned: %v", err)
	}
	steps[0].Finish(log)

	if nb := nbLineByFile(createdFile); nb != nblines1 {
		t.Errorf("%s should have %d lines after first run but it have %d", createdFile, nblines1, nb)
	}

	_, steps, err = tmpl.Load(ctx, log, "testdata/good", configFile, 0, v, dss, false, dryRun, nil)
	if err != nil {
		t.Fatalf("Load should not returns an error, returned: %v", err)
	}

	err = steps[0].Init(ctx, log)
	if err != nil {
		t.Fatalf("Init should not returns an error, returned: %v", err)
	}
	err = steps[0].Do(context.Background(), log)
	if err != nil {
		t.Errorf("Do should not return error, returned: %v", err)
	}
	steps[0].Finish(log)

	if nb := nbLineByFile(createdFile); nb != nblines2 {
		t.Errorf("%s should have %d lines after second run but it have %d", createdFile, nblines2, nb)
	}

	os.Remove(createdFile)
}

func TestDoReplace(t *testing.T) {
	helperDoTest(t, "replace", "testdata/tmp/replace.cfg", 1, 1, false)
}

func TestDoAppend(t *testing.T) {
	helperDoTest(t, "notags", "testdata/tmp/db2.cfg", 1, 2, false)
}

func TestDoUnique(t *testing.T) {
	helperDoTest(t, "fixeddest", "testdata/tmp/fixed.cfg", 2, 2, false)
}

func TestDoDryRun(t *testing.T) {
	helperDoTest(t, "notags", "testdata/tmp/db2.cfg", 0, 0, true)
}
