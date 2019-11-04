package shell

import (
	"context"

	"github.com/Sirupsen/logrus"
)

//Cancel manage the cancellation of the step
func (st *Step) Cancel(log *logrus.Entry) {
	logStep := log.WithField("name", st.Name).WithField("type", "shell")
	//TODO: to be implemented
	logStep.Info("Cancelling step")
}

//Do manage the runnning of the step
func (st *Step) Do(ctx context.Context, log *logrus.Entry) error {
	logStep := log.WithField("name", st.Name).WithField("type", "shell")
	//TODO: to be implemented
	logStep.Debug("Beginning step")
	return nil
}
