package filter

import (
	"context"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/provider/types"
)

//Filter provides way to filter record by record
type Filter interface {
	Filter(types.Record) (types.Record, error)
}

//NewFilter analyze the config map and return object implemnting Filter of the asked type
func NewFilter(ctx context.Context, log *logrus.Entry, filterType string, AParam []string, MParam map[string]string) (Filter, error) {
	switch filterType {
	case "replace":
		return newReplaceFilter(ctx, log, MParam)
	case "only":
		return newOnlyFilter(ctx, log, AParam)
	default:
		log.Errorf("Don't know how to filter %s", filterType)
		return nil, fmt.Errorf("don't know how to filter %s", filterType)
	}
}
