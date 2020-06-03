package sqlscript

import (
	"context"

	"github.com/Sirupsen/logrus"
)

//Init manage the initialization of the step.
func (st *Step) Init(ctx context.Context, log *logrus.Entry) error {
	logStep := log.WithField("name", st.Name).WithField("type", "sql")
	logStep.Debug("Initializing step")

	db, err := st.datasource.OpenDatabase(logStep, st.admin, st.noDb)
	if err != nil {
		log.Errorf("Failed to open database, %v", err)
		return err
	}

	st.db = db

	if st.transaction {
		tx, err := st.db.Begin()
		if err != nil {
			log.Errorf("Failed to open new transaction, %v", err)
			return err
		}

		st.tx = tx
	}

	return nil
}
