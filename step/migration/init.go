package migration

import (
	"context"
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
)

//Init manage the initialization of the step.
func (st *Step) Init(ctx context.Context, log *logrus.Entry) error {
	logStep := log.WithField("name", st.Name).WithField("type", "migration")
	logStep.Debug("Initializing step")

	if _, err := os.Stat(st.folder); err != nil {
		if os.IsNotExist(err) {
			log.Errorf("The migration folder %s does not exists", st.folder)
			return err
		}
	}

	if st.noAdmin {
		return nil
	}

	adminFolder := filepath.Join(st.folder, "admin")
	if _, err := os.Stat(adminFolder); err != nil {
		if os.IsNotExist(err) && st.noUser {
			log.Errorf("The admin migration folder %s does not exists and you asked me to apply it only", adminFolder)
			return err
		}

		log.Infof("The admin migration folder %s does not exists ... skipping", adminFolder)

		st.noAdmin = true
	}

	return nil
}
