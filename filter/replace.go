package filter

import (
	"context"
	"fmt"

	"github.com/marema31/kamino/provider/common"
)

// ReplaceFilter specific type for replace filter operation
type ReplaceFilter struct {
	columns map[string]string
}

func newReplaceFilter(ctx context.Context, MParam map[string]string) (Filter, error) {
	if MParam == nil {
		return nil, fmt.Errorf("no parameter to filter replace")
	}
	return &ReplaceFilter{columns: MParam}, nil
}

// Filter : replace the content of column by provided values (insert the column if not present)
func (rf *ReplaceFilter) Filter(in common.Record) (common.Record, error) {
	out := make(common.Record, len(in))

	for col, value := range in {
		out[col] = value
	}
	for col, value := range rf.columns {
		out[col] = value
	}
	return out, nil
}
