package filter

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/provider/types"
)

//Filter provides way to filter record by record.
type Filter interface {
	Filter(types.Record) (types.Record, error)
}

//NewFilter analyze the config map and return object implemnting Filter of the asked type.
func NewFilter(log *logrus.Entry, filterType string, aParam []string, mParam map[string]string) (Filter, error) {
	switch filterType {
	case "replace":
		return newReplaceFilter(log, mParam)
	case "only":
		return newOnlyFilter(log, aParam)
	default:
		log.Errorf("Don't know how to filter %s", filterType)
		return nil, fmt.Errorf("don't know how to filter %s: %w", filterType, errWrongParameterValue)
	}
}
