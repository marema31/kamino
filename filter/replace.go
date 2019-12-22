package filter

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/provider/types"
)

// ReplaceFilter specific type for replace filter operation
type ReplaceFilter struct {
	columns map[string]string
}

func newReplaceFilter(log *logrus.Entry, mParam map[string]string) (Filter, error) {
	logFilter := log.WithField("filter", "replace")

	if mParam == nil {
		logFilter.Error("Missing MParameters")
		return nil, fmt.Errorf("no parameter to filter replace")
	}

	if len(mParam) == 0 {
		logFilter.Error("Refuse to filter nothing")
		return nil, fmt.Errorf("filter replace refuse to replace nothing")
	}

	return &ReplaceFilter{columns: mParam}, nil
}

// Filter : replace the content of column by provided values (insert the column if not present)
func (rf *ReplaceFilter) Filter(in types.Record) (types.Record, error) {
	out := make(types.Record, len(in))

	for col, value := range in {
		out[col] = value
	}

	for col, value := range rf.columns {
		out[col] = value
	}

	return out, nil
}
