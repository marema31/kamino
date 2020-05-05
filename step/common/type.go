//Package common provides the utility functions and type needed by all specialized step packages
package common

import (
	"context"
	"errors"

	"github.com/Sirupsen/logrus"
)

// Steper Interface that will be used to run the steps.
type Steper interface {
	Do(context.Context, *logrus.Entry) error
	Cancel(*logrus.Entry)
	Finish(log *logrus.Entry)
	Init(context.Context, *logrus.Entry) error
	ToSkip(context.Context, *logrus.Entry) (bool, error)
	PostLoad(*logrus.Entry, map[string]string) error
}

//ErrMissingParameter raise when a parameter is missing in step definition.
var ErrMissingParameter = errors.New("MISSING PARAMETER")

//ErrWrongParameterValue raise when a parameter as a wrong value in step definition.
var ErrWrongParameterValue = errors.New("WRONG PARAMETER VALUE")
