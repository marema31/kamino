package filter

import (
	"fmt"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/provider/types"
)

// ReplaceFilter specific type for replace filter operation.
type ReplaceFilter struct {
	columns map[string]string
}

func newReplaceFilter(log *logrus.Entry, mParam map[string]string) (Filter, error) {
	logFilter := log.WithField("filter", "replace")

	if mParam == nil {
		logFilter.Error("Missing MParameters")
		return nil, fmt.Errorf("no parameter to filter replace: %w", errMissingParameter)
	}

	if len(mParam) == 0 {
		logFilter.Error("Refuse to filter nothing")
		return nil, fmt.Errorf("filter replace refuse to replace nothing: %w", errWrongParameterValue)
	}

	envVar := make(map[string]string)
	columns := make(map[string]string)

	for _, v := range os.Environ() {
		splitV := strings.Split(v, "=")
		envVar[splitV[0]] = splitV[1]
	}

	data := tmplEnv{Environments: envVar}

	logFilter.Info("Will apply replace filter on:")

	for name, value := range mParam {
		parsed, err := parseField(name, value, data)
		if err != nil {
			log.Errorf("unable to parse the template for %s (%s): %v", name, value, err)
			return nil, err
		}

		columns[name] = parsed

		logFilter.Infof("   - %s : %s", name, parsed)
	}

	return &ReplaceFilter{columns: columns}, nil
}

// Filter : replace the content of column by provided values (insert the column if not present).
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
