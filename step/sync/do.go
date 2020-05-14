package sync

import (
	"context"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
)

//Finish manage the finish of the step (called after all other step of the same priority has ended their Do).
func (st *Step) Finish(log *logrus.Entry) {
	logStep := log.WithField("name", st.Name).WithField("type", "sql")
	logStep.Info("Finishing step")

	if st.source != nil {
		st.source.Close(logStep)
	}

	if st.cacheLoader != nil {
		st.cacheLoader.Close(logStep)
	}

	if st.cacheSaver != nil {
		st.cacheSaver.Close(logStep)
	}

	for _, d := range st.destinations {
		if err := d.Close(logStep); err != nil {
			logStep.Error(err)
		}
	}
}

//Cancel manage the cancellation of the step.
func (st *Step) Cancel(log *logrus.Entry) {
	logStep := log.WithField("name", st.Name).WithField("type", "sync")
	logStep.Info("Cancelling step")

	if st.source != nil {
		st.source.Close(logStep)
	}

	if st.cacheLoader != nil {
		st.cacheLoader.Close(logStep)
	}

	if st.cacheSaver != nil {
		if err := st.cacheSaver.Reset(logStep); err != nil {
			logStep.Error(err)
		}
	}

	for _, d := range st.destinations {
		if err := d.Reset(logStep); err != nil {
			logStep.Error(err)
		}
	}
}

func (st *Step) useCache(ctx context.Context, logStep *logrus.Entry) error {
	ttlExpired := false

	cacheStat, errFile := st.cacheCfg.ds.Stat()
	if errFile == nil {
		ttlExpired = time.Since(cacheStat.ModTime()) > st.cacheTTL
	}

	if os.IsNotExist(errFile) || ttlExpired {
		logStep.Info("Cache file does not exist or too old, recreating it")

		var err error

		st.cacheSaver, err = st.prov.NewSaver(ctx, logStep, st.cacheCfg.ds, st.cacheCfg.table, "", "")
		if err != nil {
			return err
		}

		err = st.copyData(ctx, logStep)
		if err == nil {
			logStep.Info("Synchronization ok, cache file created")
			return nil
		}

		// Something goes wrong, we will try from the cache even it is expired
		err = st.cacheSaver.Reset(logStep) // reset it will only remove the temporary file
		if err != nil {
			return err
		}

		st.cacheSaver = nil //Since we closed it, ensure that it will not be closed afterwards

		logStep.Info("Cache refresh failed")
	}

	// Generation of cache file failed we will try to use the old one if it exists
	if _, err := st.cacheCfg.ds.Stat(); os.IsNotExist(err) {
		logStep.Error("Cache does not exist, and synchronization from source failed")
		return err
	}

	logStep.Info("Using cache as source")

	cacheLoader, err := st.prov.NewLoader(ctx, logStep, st.cacheCfg.ds, st.cacheCfg.table, "")
	if err != nil {
		logStep.Error("Opening cache file failed .. skipping it")
		logStep.Error(err)
	} else {
		st.cacheLoader = cacheLoader
	}

	return nil
}

//Do manage the runnning of the step.
func (st *Step) Do(ctx context.Context, log *logrus.Entry) error {
	logStep := log.WithField("name", st.Name).WithField("datasource", st.sourceCfg.ds.GetName()).WithField("type", "sync")
	logStep.Debug("Beginning step")

	if st.cacheCfg.ds != nil {
		err := st.useCache(ctx, log)
		if err != nil {
			return err
		}
	}

	err := st.copyData(ctx, logStep)
	if err != nil {
		logStep.Error("Synchronization failed")
		return err
	}

	if st.cacheCfg.ds == nil {
		logStep.Info("Synchronization ok, no cache file created")
	} else {
		logStep.Info("Synchronization ok")
	}

	return nil
}

// ToSkip return true if the step must be skipped.
func (st *Step) ToSkip(ctx context.Context, log *logrus.Entry) (bool, error) {
	logStep := log.WithField("name", st.Name).WithField("type", "sync")
	logStep.Debug("Step always executed, the synchronization mode will determine if something will be done")

	return false, nil
}
