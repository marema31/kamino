package migration

import (
	"context"
	"os"

	"github.com/Sirupsen/logrus"
)

//Init manage the initialization of the step
func (st *Step) Init(ctx context.Context, log *logrus.Entry) error {
	logStep := log.WithField("name", st.Name).WithField("type", "migration")
	//TODO: to be implemented
	logStep.Debug("Initializing step")
	if _, err := os.Stat(st.folder); err != nil {
		if os.IsNotExist(err) {
			log.Errorf("The migration folder %s does not exists", st.folder)
			return err
		}
	}
	return nil
}
