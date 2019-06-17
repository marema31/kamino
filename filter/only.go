package filter

import (
	"context"
	"fmt"

	"github.com/marema31/kamino/config"
	"github.com/marema31/kamino/provider/common"
)

// OnlyFilter specific type for Only filter operation
type OnlyFilter struct {
	columns []string
}

func newOnlyFilter(ctx context.Context, config config.FilterConfig) (Filter, error) {
	if config.AParam == nil {
		return nil, fmt.Errorf("no parameter to filter only")
	}
	return &OnlyFilter{columns: config.AParam}, nil
}

// Filter : Only the content of column by provided values (insert the column if not present)
func (of *OnlyFilter) Filter(in common.Record) (common.Record, error) {
	out := make(common.Record, len(in))

	for _, col := range of.columns {
		value, ok := in[col]
		if ok {
			out[col] = value
		}
	}
	return out, nil
}
