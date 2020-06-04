package filter

import (
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/provider/types"
	"github.com/rwtodd/Go.Sed/sed"
)

// SedFilter specific type for sed filter operation.
type SedFilter struct {
	engines map[string]*sed.Engine
	regexps map[string]string
	log     *logrus.Entry
}

func newSedFilter(log *logrus.Entry, mParam map[string]string) (Filter, error) {
	logFilter := log.WithField("filter", "sed")

	if mParam == nil {
		logFilter.Error("Missing MParameters")
		return nil, fmt.Errorf("no parameter to filter sed: %w", errMissingParameter)
	}

	if len(mParam) == 0 {
		logFilter.Error("Refuse to filter nothing")
		return nil, fmt.Errorf("filter Sed refuse to sed nothing: %w", errWrongParameterValue)
	}

	engines := make(map[string]*sed.Engine)
	regexps := make(map[string]string)

	logFilter.Info("Will apply sed filter on:")

	for name, value := range mParam {
		engine, err := sed.New(strings.NewReader(value))
		if err != nil {
			log.Errorf("unable to parse the sed expression for %s (%s): %v", name, value, err)
			return nil, err
		}

		engines[name] = engine
		regexps[name] = value

		logFilter.Infof("   - %s : %s", name, value)
	}

	return &SedFilter{engines: engines, regexps: regexps, log: logFilter}, nil
}

// Filter : Sed the content of column by provided values (insert the column if not present).
func (sf *SedFilter) Filter(in types.Record) (types.Record, error) {
	out := make(types.Record, len(in))

	for col, value := range in {
		out[col] = value
	}

	for col, engine := range sf.engines {
		value, err := engine.RunString(in[col])
		if err != nil {
			sf.log.Errorf("unable to execute sed expression (%s) on %s: %v", sf.regexps[col], in[col], err)
		}

		out[col] = value
	}

	return out, nil
}
