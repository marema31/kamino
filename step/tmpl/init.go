package tmpl

import (
	"context"

	"github.com/Sirupsen/logrus"
)

//Init manage the initialization of the step
func (st *Step) Init(ctx context.Context, log *logrus.Entry) error {
	logStep := log.WithField("name", st.Name).WithField("type", "template")
	//TODO: to be implemented
	logStep.Debug("Initializing step")
	return nil
}
