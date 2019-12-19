package sync_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/marema31/kamino/step/sync"
)

func TestDoOk(t *testing.T) {
	ctx, log, dss, v, prov := setupDo("testdata/good/steps/", "syncok")

	_, steps, err := sync.Load(ctx, log, "testdata/good", "syncok", 0, v, dss, prov, false, false, nil)
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

	sync.MockSourceContent(steps[0], []map[string]string{
		{"id": "1", "name": "Alice"},
		{"id": "2", "name": "Bob"},
	})

	err = steps[0].Do(context.Background(), log)
	if err != nil {
		t.Errorf("Do should not return error, returned: %v", err)
	}

	steps[0].Finish(log)
}

func TestDoLoadError(t *testing.T) {
	ctx, log, dss, v, prov := setupDo("testdata/good/steps/", "cachenotexist")

	_, steps, err := sync.Load(ctx, log, "testdata/good", "cachenotexist", 0, v, dss, prov, false, false, nil)
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

	sync.MockSourceContent(steps[0], []map[string]string{
		{"id": "1", "name": "Alice"},
		{"id": "2", "name": "Bob"},
	})

	sync.MockSourceError(steps[0])
	err = steps[0].Do(context.Background(), log)
	if err == nil {
		t.Errorf("Do should return error")
	}

	steps[0].Finish(log)
}

func TestDoNoCacheOk(t *testing.T) {
	ctx, log, dss, v, prov := setupDo("testdata/good/steps/", "nocache")

	_, steps, err := sync.Load(ctx, log, "testdata/good", "nocache", 0, v, dss, prov, false, false, nil)
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
}

func TestDoCacheNotExistOk(t *testing.T) {
	ctx, log, dss, v, prov := setupDo("testdata/good/steps/", "cachenotexist")

	_, steps, err := sync.Load(ctx, log, "testdata/good", "cachenotexist", 0, v, dss, prov, false, false, nil)
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
}

func TestDoCacheNotExistError(t *testing.T) {
	ctx, log, dss, v, prov := setupDo("testdata/good/steps/", "cachenotexist")

	prov.ErrorSaver = fmt.Errorf("fake error")
	prov.SaverToFail = 3
	_, steps, err := sync.Load(ctx, log, "testdata/good", "cachenotexist", 0, v, dss, prov, false, false, nil)
	if err != nil {
		t.Fatalf("Load should not returns an error, returned: %v", err)
	}

	err = steps[0].Init(ctx, log)
	if err != nil {
		t.Fatalf("Init should not returns an error, returned: %v", err)
	}

	err = steps[0].Do(context.Background(), log)
	if err == nil {
		t.Errorf("Do should return error")
	}

	steps[0].Finish(log)
}

func TestDoCacheError(t *testing.T) {
	ctx, log, dss, v, prov := setupDo("testdata/good/steps/", "syncok")

	_, steps, err := sync.Load(ctx, log, "testdata/good", "syncok", 0, v, dss, prov, false, false, nil)
	if err != nil {
		t.Fatalf("Load should not returns an error, returned: %v", err)
	}

	err = steps[0].Init(ctx, log)
	if err != nil {
		t.Fatalf("Init should not returns an error, returned: %v", err)
	}

	sync.MockSourceContent(steps[0], []map[string]string{
		{"id": "1", "name": "Alice"},
		{"id": "2", "name": "Bob"},
	})

	prov.LoaderToFail = 1
	prov.ErrorLoader = fmt.Errorf("fake error")
	err = steps[0].Do(context.Background(), log)
	if err != nil {
		t.Fatalf("Do should not returns an error, returned: %v", err)
	}

	steps[0].Finish(log)
}

func TestDoForceCacheOk(t *testing.T) {
	ctx, log, dss, v, prov := setupDo("testdata/good/steps/", "syncok")

	_, steps, err := sync.Load(ctx, log, "testdata/good", "syncok", 0, v, dss, prov, false, false, nil)
	if err != nil {
		t.Fatalf("Load should not returns an error, returned: %v", err)
	}

	superseed := make(map[string]string)
	superseed["sync.forceCacheOnly"] = "true"
	steps[0].PostLoad(log, superseed)

	err = steps[0].Init(ctx, log)
	if err != nil {
		t.Fatalf("Init should not returns an error, returned: %v", err)
	}

	err = steps[0].Do(context.Background(), log)
	if err != nil {
		t.Errorf("Do should not return error, returned: %v", err)
	}

	steps[0].Finish(log)
}

func TestDoForceCacheError(t *testing.T) {
	ctx, log, dss, v, prov := setupDo("testdata/good/steps/", "syncok")

	prov.LoaderToFail = 0
	prov.ErrorLoader = fmt.Errorf("fake error")
	_, steps, err := sync.Load(ctx, log, "testdata/good", "syncok", 0, v, dss, prov, false, false, nil)
	if err != nil {
		t.Fatalf("Load should not returns an error, returned: %v", err)
	}

	superseed := make(map[string]string)
	superseed["sync.forceCacheOnly"] = "true"
	steps[0].PostLoad(log, superseed)

	err = steps[0].Init(ctx, log)
	if err == nil {
		t.Fatalf("Init should returns an error")
	}
}

func TestDoAllowCacheOk(t *testing.T) {
	ctx, log, dss, v, prov := setupDo("testdata/good/steps/", "errorallow")

	prov.ErrorLoader = fmt.Errorf("fake error")
	_, steps, err := sync.Load(ctx, log, "testdata/good", "errorallow", 0, v, dss, prov, false, false, nil)
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
}

func TestDoInitLoaderError(t *testing.T) {
	ctx, log, dss, v, prov := setupDo("testdata/good/steps/", "syncok")

	prov.ErrorLoader = fmt.Errorf("fake error")
	_, steps, err := sync.Load(ctx, log, "testdata/good", "syncok", 0, v, dss, prov, false, false, nil)
	if err != nil {
		t.Fatalf("Load should not returns an error, returned: %v", err)
	}

	err = steps[0].Init(ctx, log)
	if err == nil {
		t.Fatalf("Init should returns an error")
	}
}

func TestDoInitSaverError(t *testing.T) {
	ctx, log, dss, v, prov := setupDo("testdata/good/steps/", "syncok")

	prov.ErrorSaver = fmt.Errorf("fake error")
	_, steps, err := sync.Load(ctx, log, "testdata/good", "syncok", 0, v, dss, prov, false, false, nil)
	if err != nil {
		t.Fatalf("Load should not returns an error, returned: %v", err)
	}

	err = steps[0].Init(ctx, log)
	if err == nil {
		t.Fatalf("Init should returns an error")
	}
}

func TestDoCancelCacheLoaderOk(t *testing.T) {
	ctx, log, dss, v, prov := setupDo("testdata/good/steps/", "syncok")

	_, steps, err := sync.Load(ctx, log, "testdata/good", "syncok", 0, v, dss, prov, false, false, nil)
	if err != nil {
		t.Fatalf("Load should not returns an error, returned: %v", err)
	}

	err = steps[0].Init(ctx, log)
	if err != nil {
		t.Fatalf("Init should not returns an error, returned: %v", err)
	}

	sync.MockSourceContent(steps[0], []map[string]string{
		{"id": "1", "name": "Alice"},
		{"id": "2", "name": "Bob"},
	})
	sync.MockSourceError(steps[0])
	err = steps[0].Do(context.Background(), log)
	if err != nil {
		t.Errorf("Do should not return error, returned: %v", err)
	}

	steps[0].Cancel(log)
}

func TestDoCancelCacheSaverOk(t *testing.T) {
	ctx, log, dss, v, prov := setupDo("testdata/good/steps/", "cachenotexist")

	_, steps, err := sync.Load(ctx, log, "testdata/good", "cachenotexist", 0, v, dss, prov, false, false, nil)
	if err != nil {
		t.Fatalf("Load should not returns an error, returned: %v", err)
	}

	err = steps[0].Init(ctx, log)
	if err != nil {
		t.Fatalf("Init should not returns an error, returned: %v", err)
	}

	sync.MockSourceContent(steps[0], []map[string]string{
		{"id": "1", "name": "Alice"},
		{"id": "2", "name": "Bob"},
	})
	sync.MockDestinationError(steps[0])
	err = steps[0].Do(context.Background(), log)
	if err == nil {
		t.Errorf("Do should return error")
	}

	steps[0].Cancel(log)
}

func TestDoUsingCacheOk(t *testing.T) {
	ctx, log, dss, v, prov := setupDo("testdata/good/steps/", "nocache")

	_, steps, err := sync.Load(ctx, log, "testdata/good", "nocache", 0, v, dss, prov, false, false, nil)
	if err != nil {
		t.Fatalf("Load should not returns an error, returned: %v", err)
	}

	err = steps[0].Init(ctx, log)
	if err != nil {
		t.Fatalf("Init should not returns an error, returned: %v", err)
	}

	sync.MockSourceContent(steps[0], []map[string]string{
		{"id": "1", "name": "Alice"},
		{"id": "2", "name": "Bob"},
	})
	sync.MockDestinationError(steps[0])
	err = steps[0].Do(context.Background(), log)
	if err == nil {
		t.Errorf("Do should return error")
	}
}

func TestDoDryRunOk(t *testing.T) {
	ctx, log, dss, v, prov := setupDo("testdata/good/steps/", "nocache")

	_, steps, err := sync.Load(ctx, log, "testdata/good", "nocache", 0, v, dss, prov, false, true, nil)
	if err != nil {
		t.Fatalf("Load should not returns an error, returned: %v", err)
	}

	err = steps[0].Init(ctx, log)
	if err != nil {
		t.Fatalf("Init should not returns an error, returned: %v", err)
	}

	sync.MockSourceContent(steps[0], []map[string]string{
		{"id": "1", "name": "Alice"},
		{"id": "2", "name": "Bob"},
	})
	err = steps[0].Do(context.Background(), log)
	if err != nil {
		t.Errorf("Do should not return error, returned: %v", err)
	}
}
