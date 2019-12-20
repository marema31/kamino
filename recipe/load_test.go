package recipe_test

import (
	"testing"
)

func TestRecipeLoadOk(t *testing.T) {
	ctx, log, _, ck := setupLoad(false, false, false)

	err := ck.Load(ctx, log, "testdata/good", []string{"recipe1ok", "recipe2ok"}, nil, nil, nil)
	if err != nil {
		t.Errorf("Load should not returns an error, returned: %v", err)
	}

	result, total := ck.Statistics()

	if len(result) != 2 {
		t.Errorf("The cookbook should contains 2 recipes, there was %d", len(result))
	}

	if total != 17 {
		t.Errorf("The cookbook should contains 17 steps, there was %d", total)
	}

	steps, ok := result["recipe2ok"]
	if !ok {
		t.Fatalf("The cookbook should contain the recipe2ok, it was not")
	}
	if len(steps) != 4 {
		t.Fatalf("It should have 4 phases to this recipe, thers was %d", len(steps))
	}

	if steps[1] != 2 {
		t.Errorf("It should have 2 step on this phase of this recipe, thers was %d", steps[1])
	}

}

func TestRecipeLoadKo(t *testing.T) {
	for _, validate := range []bool{true, false} {
		ctx, log, _, ck := setupLoad(false, false, validate)

		err := ck.Load(ctx, log, "testdata/fail", []string{"recipe1", "recipe2"}, nil, nil, nil)
		if err == nil {
			t.Errorf("Load should returns an error")
		}
		err = ck.Load(ctx, log, "testdata/fail", []string{"recipe1", "dserror"}, nil, nil, nil)
		if err == nil {
			t.Errorf("Load should returns an error")
		}
	}
}

func TestRecipeLoadNoFolder(t *testing.T) {
	for _, validate := range []bool{true, false} {
		ctx, log, _, ck := setupLoad(false, false, validate)

		err := ck.Load(ctx, log, "testdata/fail", []string{"recipe1", "unknown"}, nil, nil, nil)
		if err == nil {
			t.Errorf("Load should returns an error")
		}

		err = ck.Load(ctx, log, "testdata/fail", []string{"recipe2", "nosteps"}, nil, nil, nil)
		if err == nil {
			t.Errorf("Load should returns an error")
		}
	}
}
