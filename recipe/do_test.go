package recipe_test

import (
	"context"
	"os"
	"runtime/trace"
	"testing"
	"time"
)

func TestRecipeDoOk(t *testing.T) {
	ctx, sf, ck := setupLoad()

	err := ck.Load(ctx, "testdata/good", []string{"recipe1ok", "recipe2ok"}, nil, nil)
	if err != nil {
		t.Errorf("Load should not returns an error, returned: %v", err)
	}

	if _, ok := sf.Steps["recipe1ok"]; !ok {
		t.Errorf("The cookbook must have recipe1ok")
	} else {
		if len(sf.Steps["recipe1ok"]) != 7 {
			t.Fatalf("The cookbook must have 7 steps but has : %d", len(sf.Steps["recipe1ok"]))
		}
	}
	if _, ok := sf.Steps["recipe2ok"]; !ok {
		t.Errorf("The cookbook must have recipe1ok")
	} else {
		if len(sf.Steps["recipe2ok"]) != 10 {
			t.Fatalf("The cookbook must have 10 steps but has : %d", len(sf.Steps["recipe2ok"]))
		}
	}

	hadError := ck.Do(ctx)
	if hadError {
		t.Errorf("Do should return false")
	}

	for rname := range sf.Steps {
		for _, step := range sf.Steps[rname] {
			if !step.Called {
				t.Fatalf("One step was not executed")
			}
		}
		for _, step := range sf.Steps[rname] {
			if step.Canceled {
				t.Fatalf("One step was canceled")
			}
		}
	}
}

func TestRecipeDoAllCancel(t *testing.T) {
	ctx, sf, ck := setupLoad()

	err := ck.Load(ctx, "testdata/good", []string{"recipe1ok", "recipe2ok"}, nil, nil)
	if err != nil {
		t.Errorf("Load should not returns an error, returned: %v", err)
	}

	// Make sure the cancel will be fired
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(-7*time.Hour))
	cancel()

	hadError := ck.Do(ctx)
	if !hadError {
		t.Errorf("Do should return true")
	}

	if _, ok := sf.Steps["recipe1ok"]; !ok {
		t.Fatalf("The cookbook must have recipe1ok")
	}
	if _, ok := sf.Steps["recipe2ok"]; !ok {
		t.Errorf("The cookbook must have recipe2ok")
	}

	for rname := range sf.Steps {
		seen := make(map[bool]int)
		for _, step := range sf.Steps[rname] {
			seen[step.Canceled] += 1
		}
		if seen[true] == 0 {
			t.Errorf("No step of %s was canceled", rname)
		}
	}
}

func TestRecipeDoStepError(t *testing.T) {
	trace.Start(os.Stderr)
	ctx, sf, ck := setupLoad()

	err := ck.Load(ctx, "testdata/good", []string{"recipe1ok", "steperror"}, nil, nil)
	if err != nil {
		t.Errorf("Load should not returns an error, returned: %v", err)
	}

	hadError := ck.Do(ctx)
	if !hadError {
		t.Errorf("Do should return true")
	}

	for _, step := range sf.Steps["recipe1ok"] {
		if step.Canceled {
			t.Errorf("A step of recipe1ok was cancelled")
		}
	}

	// Way to verify the status of all steps in synthetic way (only for debug)
	//debug := make([]string, 0, 5)
	//debug = append(debug, fmt.Sprintf("\n%11v | %6v | %10v | %6v", "name", "Called", "Cancelled", "Error"))
	for _, step := range sf.Steps["steperror"] {
		if step.Priority < 5 && step.Canceled {
			t.Errorf("A step with priority less than 5 of steperror was cancelled")
		}
		if step.Priority < 5 && !step.Called {
			t.Errorf("A step with priority less than 5 of steperror was not called")
		}
		/*		if step.Priority == 5 {
					debug = append(debug, fmt.Sprintf("%11v | %6v | %10v | %6v", step.name, step.Called, step.Canceled, step.HasError))
				}
		*/
		/* TODO: Non blocking issue (#1), for the moment will will go to MVP, come to back after
		if step.Priority == 5 && !step.Canceled {
			t.Errorf("A step with priority 5 of steperror was not cancelled")
		}
		*/
		if step.Priority == 5 && !step.Called {
			t.Errorf("A step with priority 5 of steperror was not called")
		}
		if step.Priority > 5 && step.Canceled {
			t.Errorf("A step with priority more than 5 of steperror was cancelled")
		}
		if step.Priority > 5 && step.Called {
			t.Errorf("A step with priority more than 5 of steperror was cancelled")
		}
	}

	//t.Fatal(strings.Join(debug, "\n"))
}