//Package common provides the utility functions and type needed by all specialized step packages
package common

import (
	"errors"
)

//ErrEOF error EndOfFile.
var ErrEOF = errors.New("EOF")

//ErrMissingParameter raise when a parameter is missing in provider definition.
var ErrMissingParameter = errors.New("MISSING PARAMETER")

//ErrWrongParameterValue raise when a parameter as a wrong value in provider definition.
var ErrWrongParameterValue = errors.New("WRONG PARAMETER VALUE")
