package recipe_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/provider"
	"github.com/marema31/kamino/recipe"
	"github.com/marema31/kamino/step"
	"github.com/marema31/kamino/step/common"
	"github.com/spf13/viper"
)

type mockedStep struct {
	name string
}

//Do manage the runnning of the step
func (st *mockedStep) Do(ctx context.Context) error {
	//TODO: to be implemented
	fmt.Printf("Will do shell for %s\n", st.name)
	return nil
}

// MockedStepFactory implement the step.Creater interface for testing purpose returning step.Steper that are doing no I/O
type MockedStepFactory struct {
}

// Load the step file and returns the priority and a list of steper for this file
func (sf MockedStepFactory) Load(ctx context.Context, path string, filename string, dss datasource.Datasourcers, prov provider.Provider) (uint, []common.Steper, error) {
	v := viper.New()

	v.SetConfigName(filename) // The file will be named [filename].json, [filename].yaml or [filename.toml]
	v.AddConfigPath(path)
	err := v.ReadInConfig()
	if err != nil {
		return 0, nil, err
	}
	priority := v.GetUint("priority")
	nbSteps := v.GetInt("steps")
	generateError := v.GetBool("error")

	if generateError {
		return 0, nil, fmt.Errorf("fake error")
	}

	steps := make([]common.Steper, 0, nbSteps)
	for i := 0; nbSteps > i; i++ {
		m := &mockedStep{name: fmt.Sprintf("%s_%d", filename, i)}
		steps = append(steps, m)
	}
	//TODO: read the step file, this one must contain : name, priority and
	return priority, steps, nil
}

func setupLoad() (context.Context, step.Creater, *recipe.Cookbook) {
	ctx := context.Background()
	sf := MockedStepFactory{}
	ck := recipe.New(sf)
	return ctx, sf, ck
}

func TestRecipeLoadOk(t *testing.T) {
	ctx, _, ck := setupLoad()

	err := ck.Load(ctx, "testdata/good", []string{"recipe1ok", "recipe2ok"}, nil, nil)
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
	ctx, _, ck := setupLoad()

	err := ck.Load(ctx, "testdata/fail", []string{"recipe1", "recipe2"}, nil, nil)
	if err == nil {
		t.Errorf("Load should returns an error")
	}
}

func TestRecipeLoadNoFolder(t *testing.T) {
	ctx, _, ck := setupLoad()

	err := ck.Load(ctx, "testdata/fail", []string{"recipe1", "unknown"}, nil, nil)
	if err == nil {
		t.Errorf("Load should returns an error")
	}

	err = ck.Load(ctx, "testdata/fail", []string{"recipe1", "nosteps"}, nil, nil)
	if err == nil {
		t.Errorf("Load should returns an error")
	}
}
