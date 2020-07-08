package sqlscript

import (
	"context"

	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/step/common"
)

//Finish manage the finish of the step (called after all other step of the same priority has ended their Do).
func (st *Step) Finish(log *logrus.Entry) {
	logStep := log.WithField("name", st.Name).WithField("type", "sql")
	logStep.Info("Finish step")

	defer st.datasource.CloseDatabase(log, st.admin, st.noDb) //nolint: errcheck

	if st.tx != nil {
		if err := st.tx.Commit(); err != nil {
			logStep.Error(err)
		}
	}
}

//Cancel manage the cancellation of the step.
func (st *Step) Cancel(log *logrus.Entry) {
	logStep := log.WithField("name", st.Name).WithField("type", "sql")
	logStep.Info("Cancelling step")

	defer st.datasource.CloseDatabase(log, st.admin, st.noDb) //nolint: errcheck

	if st.tx != nil {
		if err := st.tx.Rollback(); err != nil {
			logStep.Error(err)
		}
	}
}

//Do manage the runnning of the step.
func (st *Step) Do(ctx context.Context, log *logrus.Entry) error {
	logStep := log.WithField("name", st.Name).WithField("datasource", st.datasource.GetName()).WithField("type", "sql")
	logStep.Info("Beginning step")

	if st.dryRun {
		for _, stmt := range st.sqlCmds {
			log.Infof("Will run %s", stmt)
		}

		return nil
	}

	var err error

	for _, stmt := range st.sqlCmds {
		logStep.Debug("Executing: ")
		logStep.Debug(stmt)

		if st.tx != nil {
			_, err = st.tx.ExecContext(ctx, stmt)
		} else {
			_, err = st.db.ExecContext(ctx, stmt)
		}

		if err != nil && st.ignoreErrors {
			logStep.Warnf("Ignoring error: %v on statement: %s", err, stmt)
			return nil
		}

		if err != nil {
			logStep.Error("Execution of one statement failed:")
			logStep.Error(stmt)
			logStep.Error(err)

			return err
		}
	}

	logStep.Info("Ending step")

	return nil
}

// ToSkip return true if the step must be skipped (based on the query parameter.
func (st *Step) ToSkip(ctx context.Context, log *logrus.Entry) (bool, error) {
	logStep := log.WithField("name", st.Name).WithField("type", "sql")
	return common.ToSkipDatabase(ctx, logStep, st.datasource, st.admin, st.noDb, st.queries)
}
