package tmpl

import (
	"context"

	"github.com/Sirupsen/logrus"
)

//Finish manage the finish of the step (called after all other step of the same priority has ended their Do)
func (st *Step) Finish(log *logrus.Entry) {
	logStep := log.WithField("name", st.Name).WithField("type", "sql")
	logStep.Info("Finishing step")
	//TODO: to be implemented
}

//Cancel manage the cancellation of the step
func (st *Step) Cancel(log *logrus.Entry) {
	logStep := log.WithField("name", st.Name).WithField("type", "template")
	//TODO: to be implemented
	logStep.Info("Cancelling step")
}

//Do manage the runnning of the step
func (st *Step) Do(ctx context.Context, log *logrus.Entry) error {
	logStep := log.WithField("name", st.Name).WithField("type", "template")
	//TODO: to be implemented
	logStep.Debug("Beginning step")
	return nil
}

// ToSkip return true if the step must be skipped
func (st *Step) ToSkip(ctx context.Context, log *logrus.Entry) (bool, error) {
	logStep := log.WithField("name", st.Name).WithField("type", "template")
	//TODO: to be implemented
	logStep.Debug("Do we need to skip the step ?")
	return true, nil
}
