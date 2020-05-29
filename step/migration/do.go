package migration

import (
	"context"

	"github.com/Sirupsen/logrus"
	migrate "github.com/rubenv/sql-migrate"

	"github.com/marema31/kamino/step/common"
)

//Finish manage the finish of the step (called after all other step of the same priority has ended their Do).
func (st *Step) Finish(log *logrus.Entry) {
	logStep := log.WithField("name", st.Name).WithField("type", "migration")
	logStep.Debug("Finishing step")
	logStep.Debug("Nothing to do")
}

//Cancel manage the cancellation of the step.
func (st *Step) Cancel(log *logrus.Entry) {
	logStep := log.WithField("name", st.Name).WithField("type", "migration")
	logStep.Info("Cancelling step")
	logStep.Info("Migration is not cancellable")
}

//Do manage the runnning of the step.
func (st *Step) Do(ctx context.Context, log *logrus.Entry) error {
	logStep := log.WithField("name", st.Name).WithField("datasource", st.datasource.GetName()).WithField("type", "migration")
	logStep.Debug("Beginning step")

	limit := st.limit

	if !st.noAdmin && (st.dir == migrate.Up || st.printOnly) {
		logStep.Info("Applying admin migration")

		applied, err := st.applyOrPrint(logStep.WithField("kind", "admin"), true, limit)
		if err != nil {
			return err
		}

		if limit > 0 {
			limit -= applied
			if 0 >= limit && !st.printOnly {
				//Everything was applied
				return nil
			}
		}
	}

	if !st.noUser {
		logStep.Info("Applying user migration")

		applied, err := st.applyOrPrint(logStep.WithField("kind", "user"), false, limit)
		if err != nil {
			return err
		}

		if limit > 0 {
			limit -= applied
			if 0 >= limit {
				//Everything was applied
				return nil
			}
		}
	}

	if st.printOnly {
		//Everything was printed
		return nil
	}

	if !st.noAdmin && st.dir == migrate.Down {
		logStep.Info("Applying admin migration")

		_, err := st.apply(logStep.WithField("kind", "admin"), true, limit)
		if err != nil {
			return err
		}
	}

	return nil
}

func (st *Step) applyOrPrint(log *logrus.Entry, admin bool, limit int) (int, error) {
	if st.printOnly {
		return st.print(log, admin)
	}

	return st.apply(log, admin, limit)
}

// ToSkip return true if the step must be skipped (based on the query parameter.
func (st *Step) ToSkip(ctx context.Context, log *logrus.Entry) (bool, error) {
	logStep := log.WithField("name", st.Name).WithField("type", "migration")

	if !st.noAdmin {
		empty, err := st.datasource.IsTableEmpty(ctx, logStep, st.tableAdmin)
		if empty || err != nil {
			logStep.Debug("Do not skip since tableAdmin is empty")
			return false, err
		}
	}

	if !st.noUser {
		empty, err := st.datasource.IsTableEmpty(ctx, logStep, st.tableUser)
		if empty || err != nil {
			logStep.Debug("Do not skip since tableUser is empty")
			return false, err
		}
	}

	return common.ToSkipDatabase(ctx, logStep, st.datasource, true, false, st.queries)
}
