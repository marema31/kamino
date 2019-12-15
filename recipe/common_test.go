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

func setupLoad(force bool, sequential bool, validate bool) (context.Context, *logrus.Entry, *MockedStepFactory, *recipe.Cookbook) {
	ctx := context.Background()
	sf := &MockedStepFactory{}
	ck := recipe.New(sf, time.Millisecond*2, 2, force, sequential, validate)
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	log := logger.WithField("appname", "kamino")
	return ctx, log, sf, ck
}

// ***** MOCK STEP and STEP CREATER
type mockedStep struct {
	name          string
	Initialized   bool
	Called        bool
	Canceled      bool
	Finished      bool
	HasError      bool
	PostLoaded    bool
	InitError     error
	PostLoadError error
	StepError     error
	Priority      uint
	ToBeSkipped   bool
	ToSkipError   error
}

//Init manage the initialization of the step
func (st *mockedStep) Init(ctx context.Context, log *logrus.Entry) error {
	log.WithField("name", st.name).WithField("error", st.InitError).Info("Initializing")
	st.Initialized = true
	return st.InitError
}

//Finish manage the finish of the step (called after all other step of the same priority has ended their Do)
func (st *mockedStep) Finish(log *logrus.Entry) {
	log.WithField("name", st.name).WithField("error", st.InitError).Info("Finishing")
	st.Finished = true
}

//Cancel manage the cancellation of the step
func (st *mockedStep) Cancel(log *logrus.Entry) {
	log.WithField("name", st.name).Info("Cancelling")
	st.Canceled = true
}

// ToSkip return true if the step must be skipped
func (st *mockedStep) ToSkip(ctx context.Context, log *logrus.Entry) (bool, error) {
	log.WithField("name", st.name).Info("Do we need to skip the step ?")
	return st.ToBeSkipped, st.ToSkipError
}

//Do manage the runnning of the step
func (st *mockedStep) Do(ctx context.Context, log *logrus.Entry) error {
	log.WithField("name", st.name).WithField("error", st.StepError).Info("Doing")
	st.Called = true

	time.Sleep(5 * time.Millisecond) // It is moking we are doing nothing
	if st.StepError != nil {
		st.HasError = true
		log.WithField("name", st.name).Info("Error")
	}
	log.WithField("name", st.name).Info("End doing")
	return st.StepError
}

//PostLoad modify the loaded step values with the values provided in the map in argument
func (st *mockedStep) PostLoad(log *logrus.Entry, superseed map[string]string) error {
	// Nothing to do
	st.PostLoaded = true
	return st.PostLoadError
}

// MockedStepFactory implement the step.Creater interface for testing purpose returning step.Steper that are doing no I/O
type MockedStepFactory struct {
	Steps map[string][]*mockedStep
}

// Load the step file and returns the priority and a list of steper for this file
func (sf *MockedStepFactory) Load(ctx context.Context, log *logrus.Entry, recipePath string, filename string, dss datasource.Datasourcers, prov provider.Provider, stepNames []string, stepTypes []string, force bool) (uint, []common.Steper, error) {
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
	postLoadError := v.GetBool("postloaderror")

	ToBeSkipped := v.GetBool("tobeskipped")
	ToSkipError := v.GetBool("toskiperror")

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

		if ToSkipError {
			m.ToSkipError = fmt.Errorf("fake error")
		}

		if postLoadError {
			m.PostLoadError = fmt.Errorf("fake error")
		}

		m.ToBeSkipped = ToBeSkipped

		steps = append(steps, m)
		sf.Steps[rname] = append(sf.Steps[rname], m)
	}
	return priority, steps, nil
}
