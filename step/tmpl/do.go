package tmpl

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/Sirupsen/logrus"
)

//Finish manage the finish of the step (called after all other step of the same priority has ended their Do).
func (st *Step) Finish(log *logrus.Entry) {
	logStep := log.WithField("name", st.Name).WithField("type", "sql")
	logStep.Info("Finishing step")

	if err := st.output.CloseFile(logStep); err != nil {
		logStep.Error(err)
	}
}

//Cancel manage the cancellation of the step.
func (st *Step) Cancel(log *logrus.Entry) {
	logStep := log.WithField("name", st.Name).WithField("type", "template")
	logStep.Info("Cancelling step")

	if err := st.output.ResetFile(logStep); err != nil {
		logStep.Error(err)
	}
}

//Do manage the runnning of the step.
func (st *Step) Do(ctx context.Context, log *logrus.Entry) error {
	logStep := log.WithField("name", st.Name).WithField("type", "template")
	logStep.Debug("Beginning step")

	if st.mode == Unique && st.input.FilePath != "" {
		return st.doUnique(logStep)
	}

	return st.doNonUnique(logStep)
}

func (st *Step) doNonUnique(log *logrus.Entry) error {
	for _, ds := range st.datasources {
		logDs := log.WithField("datasource", ds.GetName())
		logDs.Debug("Rendering template")

		if st.dryRun {
			log.Infof("Rendering template to %s in append/replace mode", st.output.FilePath)
			return nil
		}

		if err := st.template.Execute(st.outputHandle, ds.FillTmplValues()); err != nil {
			logDs.Error("Template rendering failed")
			logDs.Error(err)

			return err
		}
	}

	return nil
}

func (st *Step) doUnique(log *logrus.Entry) error {
	rendered := bytes.NewBuffer(make([]byte, 0, 1024))

	fileContent, err := ioutil.ReadFile(st.input.FilePath)
	if err != nil {
		log.Error("Reading original content failed")
		log.Error(err)

		return err
	}

	for _, ds := range st.datasources {
		logDs := log.WithField("datasource", ds.GetName())
		logDs.Debug("Rendering template")

		if st.dryRun {
			log.Infof("Rendering template to %s in unique mode", st.output.FilePath)
			continue
		}

		if err := st.template.Execute(rendered, ds.FillTmplValues()); err != nil {
			logDs.Error("Template rendering failed")
			logDs.Error(err)

			return err
		}

		renderedString := rendered.String()

		isExist, err := regexp.Match(renderedString, fileContent)
		if err != nil {
			logDs.Error("Looking for previous occurrence failed")
			logDs.Error(err)

			return err
		}

		if !isExist {
			if _, err := rendered.WriteTo(st.outputHandle); err != nil {
				logDs.Error("Writing failed")
				logDs.Error(err)

				return err
			}

			fileContent = append(fileContent, rendered.Bytes()...)
		}

		rendered.Reset()
	}

	return nil
}

// ToSkip return true if the step must be skipped.
func (st *Step) ToSkip(ctx context.Context, log *logrus.Entry) (bool, error) {
	logStep := log.WithField("name", st.Name).WithField("type", "template")
	logStep.Debug("Do we need to skip the step ?")

	if st.onlyIfNotExists {
		if _, err := os.Stat(st.destination); !os.IsNotExist(err) {
			return true, nil
		}
	}

	return false, nil
}
