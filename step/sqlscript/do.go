package sqlscript

import (
	"context"

	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/step/common"
)

//Cancel manage the cancellation of the step
func (st *Step) Cancel(log *logrus.Entry) {
	logStep := log.WithField("name", st.Name).WithField("type", "sql")
	//TODO: to be implemented
	logStep.Info("Cancelling step")
}

//Do manage the runnning of the step
func (st *Step) Do(ctx context.Context, log *logrus.Entry) error {
	logStep := log.WithField("name", st.Name).WithField("type", "sql")
	//TODO: to be implemented
	logStep.Debug("Beginning step")
	return nil
}

// ToSkip return true if the step must be skipped (based on the query parameter
func (st *Step) ToSkip(ctx context.Context, log *logrus.Entry) (bool, error) {
	logStep := log.WithField("name", st.Name).WithField("type", "sql")
	return common.ToSkipDatabase(ctx, logStep, st.datasource, st.admin, st.noDb, st.query)
}
