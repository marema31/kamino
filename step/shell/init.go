package shell

import (
	"context"
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
)

//Init manage the initialization of the step
func (st *Step) Init(ctx context.Context, log *logrus.Entry) error {
	logStep := log.WithField("name", st.Name).WithField("type", "shell")
	logStep.Debug("Initializing step")

	script := filepath.Join(st.cmd.Dir, st.cmd.Path)
	if _, err := os.Stat(script); err != nil {
		if os.IsNotExist(err) {
			log.Errorf("The script %s does not exists at execution phase", st.cmd.Path)
			return err
		}
	}

	return nil
}
