package sync

import (
	"context"

	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/provider"
)

//Init manage the initialization of the step
func (st *Step) Init(ctx context.Context, log *logrus.Entry) error {
	logStep := log.WithField("name", st.Name).WithField("type", "sync")
	var err error

	logStep.Debug("Initializing step")
	logStep.Debug("Creating loader instance for source")

	//TODO: add a parameter allowFromCacheOnly, if its true and the NewLoader get back in error, open the cache as source and st.cacheCfg = nil
	st.source, err = st.prov.NewLoader(ctx, log, st.sourceCfg.ds, st.sourceCfg.table, st.sourceCfg.where)
	if err != nil {
		return err
	}

	if st.cacheCfg.ds != nil {
		// Not creating loader instance for cache since the file may not exists if the cache has never be generated
		logStep.Debug("Creating saver instance for cache")
		st.cacheSaver, err = st.prov.NewSaver(ctx, log, st.cacheCfg.ds, st.cacheCfg.table, "", "")
		if err != nil {
			return err
		}
	}

	log.Debug("Creating saver instances for destinations")
	savers := make([]provider.Saver, 0, len(st.destsCfg))
	for _, dest := range st.destsCfg {
		saver, err := st.prov.NewSaver(ctx, log, dest.ds, dest.table, dest.key, dest.mode)
		if err != nil {
			return err
		}
		savers = append(savers, saver)
	}
	st.destinations = savers
	return nil
}
