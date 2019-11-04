package filter

import (
	"context"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/provider/types"
)

// OnlyFilter specific type for Only filter operation
type OnlyFilter struct {
	columns []string
}

func newOnlyFilter(ctx context.Context, log *logrus.Entry, AParam []string) (Filter, error) {
	logFilter := log.WithField("filter", "only")
	if AParam == nil {
		logFilter.Error("Missing AParameters")
		return nil, fmt.Errorf("no parameter to filter only")
	}
	if len(AParam) == 0 {
		logFilter.Error("Refuse to filter everything")
		return nil, fmt.Errorf("filter only refuse to filter everything")
	}

	return &OnlyFilter{columns: AParam}, nil
}

// Filter : Only the content of column by provided values (insert the column if not present)
func (of *OnlyFilter) Filter(in types.Record) (types.Record, error) {
	out := make(types.Record, len(in))

	for _, col := range of.columns {
		value, ok := in[col]
		if ok {
			out[col] = value
		}
	}
	return out, nil
}
