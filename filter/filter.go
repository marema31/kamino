package filter

import (
	"context"
	"fmt"

	"github.com/marema31/kamino/provider/common"
)

//Filter provides way to filter record by record
type Filter interface {
	Filter(common.Record) (common.Record, error)
}

//NewFilter analyze the config map and return object implemnting Filter of the asked type
func NewFilter(ctx context.Context, filterType string, AParam []string, MParam map[string]string) (Filter, error) {
	switch filterType {
	case "replace":
		return newReplaceFilter(ctx, MParam)
	case "only":
		return newOnlyFilter(ctx, AParam)
	default:
		return nil, fmt.Errorf("don't know how to filter %s", filterType)
	}
}
