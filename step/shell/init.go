package shell

import (
	"context"
	"os"

	"github.com/Sirupsen/logrus"
)

//Init manage the initialization of the step
func (st *Step) Init(ctx context.Context, log *logrus.Entry) error {
	logStep := log.WithField("name", st.Name).WithField("type", "shell")
	//TODO: to be implemented
	logStep.Debug("Initializing step")
	if _, err := os.Stat(st.script); err != nil {
		if os.IsNotExist(err) {
			log.Errorf("The script %s does not exists at execution phase", st.script)
			return err
		}
	}

	return nil
}
