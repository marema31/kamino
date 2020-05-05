package datasource

import (
	"errors"
)

var errMissingParameter = errors.New("MISSING PARAMETER")
var errWrongParameterValue = errors.New("WRONG PARAMETER VALUE")
var errWrongType = errors.New("WRONG TYPE")
var errNoConfiguration = errors.New("NO CONFIGURATION FOUND")
