package filter

import (
	"context"
	"fmt"

	"github.com/marema31/kamino/config"
	"github.com/marema31/kamino/provider/common"
)

//Filter provides way to filter record by record
type Filter interface {
	Filter(common.Record) (common.Record, error)
}

//NewFilter analyze the config map and return object implemnting Filter of the asked type
func NewFilter(ctx context.Context, config config.FilterConfig) (Filter, error) {
	switch config.Type {
	case "replace":
		return newReplaceFilter(ctx, config)
	case "only":
		return newOnlyFilter(ctx, config)
	default:
		return nil, fmt.Errorf("don't know how to filter %s", config.Type)
	}
}
