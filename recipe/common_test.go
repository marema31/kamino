package recipe_test

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/provider"
	"github.com/marema31/kamino/recipe"
	"github.com/marema31/kamino/step/common"
	"github.com/spf13/viper"
)

func setupLoad() (context.Context, *logrus.Entry, *MockedStepFactory, *recipe.Cookbook) {
	ctx := context.Background()
	sf := &MockedStepFactory{}
	ck := recipe.New(sf)
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	return ctx, log, sf, ck
}

// ***** MOCK STEP and STEP CREATER
type mockedStep struct {
	name        string
	Initialized bool
	Called      bool
	Canceled    bool
	HasError    bool
	InitError   error
	StepError   error
	Priority    uint
}

//Init manage the initialization of the step
func (st *mockedStep) Init(ctx context.Context, log *logrus.Entry) error {
	log.WithField("name", st.name).WithField("error", st.InitError).Info("Initializing")
	st.Initialized = true
	return st.InitError
}

//Cancel manage the cancellation of the step
func (st *mockedStep) Cancel(log *logrus.Entry) {
	log.WithField("name", st.name).Info("Cancelling")
	st.Canceled = true
}

//Do manage the runnning of the step
func (st *mockedStep) Do(ctx context.Context, log *logrus.Entry) error {
	log.WithField("name", st.name).WithField("error", st.StepError).Info("Doing")
	st.Called = true

	time.Sleep(1 * time.Second) // It is moking we are doing nothing
	if st.StepError != nil {
		st.HasError = true
		log.WithField("name", st.name).Info("Error")
	}
	log.WithField("name", st.name).Info("End doing")
	return st.StepError
}

// MockedStepFactory implement the step.Creater interface for testing purpose returning step.Steper that are doing no I/O
type MockedStepFactory struct {
	Steps map[string][]*mockedStep
}

// Load the step file and returns the priority and a list of steper for this file
func (sf *MockedStepFactory) Load(ctx context.Context, log *logrus.Entry, recipePath string, filename string, dss datasource.Datasourcers, prov provider.Provider, stepNames []string, stepTypes []string) (uint, []common.Steper, error) {
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
	stepInitError := v.GetBool("stepiniterror")

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
		if stepInitError {
			m.InitError = fmt.Errorf("fake error")
		}

		steps = append(steps, m)
		sf.Steps[rname] = append(sf.Steps[rname], m)
	}
	return priority, steps, nil
}
