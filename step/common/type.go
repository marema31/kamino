//Package common provides the utility functions and type needed by all specialized step packages
package common

import (
	"context"

	"github.com/Sirupsen/logrus"
)

// Steper Interface that will be used to run the steps
type Steper interface {
	Do(context.Context, *logrus.Entry) error
	Cancel(*logrus.Entry)
	Init(context.Context, *logrus.Entry) error
}
