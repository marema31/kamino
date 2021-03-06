package sync

import (
	"context"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/provider"
	"github.com/marema31/kamino/step/common"
)

//Init manage the initialization of the step.
func (st *Step) Init(ctx context.Context, log *logrus.Entry) error {
	logStep := log.WithField("name", st.Name).WithField("type", "sync")

	var err error

	logStep.Debug("Initializing step")
	logStep.Debug("Creating loader instance for source")

	if st.forceCacheOnly && st.cacheCfg.ds != nil {
		logStep.Info("Cache usage forced")

		st.source, err = st.prov.NewLoader(ctx, log, st.cacheCfg.ds, st.cacheCfg.table, "")
		st.cacheCfg.ds = nil
	} else {
		st.source, err = st.prov.NewLoader(ctx, log, st.sourceCfg.ds, st.sourceCfg.table, st.sourceCfg.where)
		if err != nil && st.allowCacheOnly && st.cacheCfg.ds != nil {
			logStep.Info("Source not available, I will use the cache")

			st.source, err = st.prov.NewLoader(ctx, log, st.cacheCfg.ds, st.cacheCfg.table, "")
			st.cacheCfg.ds = nil
		}
	}

	if err != nil {
		return err
	}

	log.Debug("Creating saver instances for destinations")

	savers := make([]provider.Saver, 0, len(st.destsCfg))

	for _, dest := range st.destsCfg {
		skip, err := common.ToSkipDatabase(ctx, logStep, dest.ds, false, false, dest.queries)

		if err != nil {
			return fmt.Errorf("unable to determine is this destination must be skipped, %w", err)
		}

		if !skip {
			saver, err := st.prov.NewSaver(ctx, log, dest.ds, dest.table, dest.key, dest.mode)
			if err != nil {
				return err
			}

			savers = append(savers, saver)
		}
	}

	st.destinations = savers

	return nil
}
