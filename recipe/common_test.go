package recipe_test

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/provider"
	"github.com/marema31/kamino/recipe"
	"github.com/marema31/kamino/step/common"
	"github.com/spf13/viper"
)

func setupLoad() (context.Context, *MockedStepFactory, *recipe.Cookbook) {
	ctx := context.Background()
	sf := &MockedStepFactory{}
	ck := recipe.New(sf)
	return ctx, sf, ck
}

// ***** MOCK STEP and STEP CREATER
type mockedStep struct {
	name      string
	Called    bool
	Canceled  bool
	HasError  bool
	StepError error
	Priority  uint
}

//Cancel manage the cancellation of the step
func (st *mockedStep) Cancel() {
	st.Canceled = true
}

//Do manage the runnning of the step
func (st *mockedStep) Do(ctx context.Context) error {
	st.Called = true

	time.Sleep(1 * time.Second) // It is moking we are doing nothing
	if st.StepError != nil {
		st.HasError = true
	}
	return st.StepError
}

// MockedStepFactory implement the step.Creater interface for testing purpose returning step.Steper that are doing no I/O
type MockedStepFactory struct {
	Steps map[string][]*mockedStep
}

// Load the step file and returns the priority and a list of steper for this file
func (sf *MockedStepFactory) Load(ctx context.Context, recipePath string, filename string, dss datasource.Datasourcers, prov provider.Provider) (uint, []common.Steper, error) {
	v := viper.New()

	v.SetConfigName(filename) // The file will be named [filename].json, [filename].yaml or [filename.toml]
	stepsFolder := filepath.Join(recipePath, "steps")
	v.AddConfigPath(stepsFolder)
	err := v.ReadInConfig()
	if err != nil {
		return 0, nil, err
	}
	priority := v.GetUint("priority")
	nbSteps := v.GetInt("steps")
	generateError := v.GetBool("generateerror")
	stepError := v.GetBool("steperror")

	if generateError {
		return 0, nil, fmt.Errorf("fake error")
	}

	rname := filepath.Base(recipePath)
	steps := make([]common.Steper, 0, nbSteps)
	if sf.Steps == nil {
		sf.Steps = make(map[string][]*mockedStep)
	}

	if len(sf.Steps[rname]) == 0 {
		sf.Steps[rname] = make([]*mockedStep, 0, nbSteps)
	}
	for i := 0; nbSteps > i; i++ {
		m := &mockedStep{name: fmt.Sprintf("%s_%d", filename, i), Priority: priority}
		if stepError {
			m.StepError = fmt.Errorf("fake error")
		}
		steps = append(steps, m)
		sf.Steps[rname] = append(sf.Steps[rname], m)
	}
	return priority, steps, nil
}