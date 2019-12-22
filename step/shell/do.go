package shell

import (
	"context"
	"io"
	"os"
	"sync"

	"github.com/Sirupsen/logrus"
)

//Finish manage the finish of the step (called after all other step of the same priority has ended their Do)
func (st *Step) Finish(log *logrus.Entry) {
	logStep := log.WithField("name", st.Name).WithField("type", "shell")
	logStep.Info("Finishing step")
}

//Cancel manage the cancellation of the step
func (st *Step) Cancel(log *logrus.Entry) {
	logStep := log.WithField("name", st.Name).WithField("type", "shell")
	logStep.Info("Cancelling step")
}

//Do manage the runnning of the step
func (st *Step) Do(ctx context.Context, log *logrus.Entry) error {
	logStep := log.WithField("name", st.Name).WithField("type", "shell")
	logStep.Debug("Beginning step")

	var wg sync.WaitGroup

	if st.dryRun {
		log.Infof("Will run %s", st.cmd.String())
		return nil
	}

	stdout, err := st.cmd.StdoutPipe()
	if err != nil {
		logStep.Error("Script output gathering failed")
		logStep.Error(err)

		return err
	}

	stderr, err := st.cmd.StderrPipe()
	if err != nil {
		logStep.Error("Script output gathering failed")
		logStep.Error(err)

		return err
	}

	err = st.cmd.Start()
	if err != nil {
		logStep.Error("Script execution failed")
		logStep.Error(err)

		return err
	}

	logStep.Debug("Waiting for command to finish...")

	wg.Add(2)

	go func() {
		defer wg.Done()
		if _, err := io.Copy(os.Stdout, stdout); err != nil {
			logrus.Error(err)
		}
	}()

	go func() {
		defer wg.Done()
		if _, err := io.Copy(os.Stderr, stderr); err != nil {
			logrus.Error(err)
		}
	}()

	wg.Wait()

	err = st.cmd.Wait()
	if err != nil {
		logStep.Error("Script finished with error")
		logStep.Error(err)
	}

	return err
}

// ToSkip return true if the step must be skipped
func (st *Step) ToSkip(ctx context.Context, log *logrus.Entry) (bool, error) {
	logStep := log.WithField("name", st.Name).WithField("type", "shell")
	logStep.Debug("Do we need to skip the step ?")

	return false, nil
}
