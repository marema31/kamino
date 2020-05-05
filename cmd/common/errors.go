//Package common provides the utility functions and type needed by all specialized step packages
package common

import (
	"errors"
)

//ErrStep raise when a step fininished in error.
var ErrStep = errors.New("STEP")

//ErrMissingParameter raise when a parameter is missing in command line.
var ErrMissingParameter = errors.New("MISSING PARAMETER")

//ErrWrongParameterValue raise when a parameter as a wrong value in command line.
var ErrWrongParameterValue = errors.New("WRONG PARAMETER VALUE")

//ErrNoConfiguration raise when no configuration is found.
var ErrNoConfiguration = errors.New("NO CONFIGURATION FOUND")
