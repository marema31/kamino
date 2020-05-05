package tmpl

import (
	"context"
	"io"
	"os"

	"github.com/Sirupsen/logrus"
)

//Init manage the initialization of the step.
func (st *Step) Init(ctx context.Context, log *logrus.Entry) error {
	logStep := log.WithField("name", st.Name).WithField("type", "template")
	logStep.Debug("Initializing step")

	var err error

	if st.outputHandle, err = st.output.OpenWriteFile(logStep); err != nil {
		logStep.Error("Opening destination file failed")
		logStep.Error(err)

		return err
	}

	if st.input.FilePath != "" {
		if _, err := st.input.Stat(); os.IsNotExist(err) {
			st.input.FilePath = "" // The file does not exists, it is not a input
			return nil
		}

		logStep.Debug("Copying original content to destination file")

		if st.inputHandle, err = st.input.OpenReadFile(logStep); err != nil {
			logStep.Error("Opening original destination file failed")
			logStep.Error(err)

			return err
		}

		// we are in append or uniqe mode therefore we must ensure the content of the file is in the temporary file
		if _, err = io.Copy(st.outputHandle, st.inputHandle); err != nil {
			logStep.Error("Copying original content to destination file failed")
			logStep.Error(err)

			return err
		}

		st.inputHandle.Close()
		st.inputHandle = nil
	}

	return nil
}
