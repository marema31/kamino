package sync

import (
	"context"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
)

//Finish manage the finish of the step (called after all other step of the same priority has ended their Do)
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
		d.Close(logStep)

	}
}

//Cancel manage the cancellation of the step
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
		st.cacheSaver.Reset(logStep)
	}
	for _, d := range st.destinations {
		d.Reset(logStep)

	}
}

//Do manage the runnning of the step
func (st *Step) Do(ctx context.Context, log *logrus.Entry) error {
	logStep := log.WithField("name", st.Name).WithField("type", "sync")
	logStep.Debug("Beginning step")

	if st.cacheCfg.ds != nil {
		ttlExpired := false
		cacheStat, errFile := st.cacheCfg.ds.Stat()
		if errFile == nil {
			ttlExpired = time.Since(cacheStat.ModTime()) > st.cacheTTL
		}

		if os.IsNotExist(errFile) || ttlExpired {
			logStep.Info("Cache file does not exist or too old, recreating it")
			err := st.copyData(ctx, logStep)
			if err == nil {
				logStep.Info("Synchronization ok, cache file created")
				return nil
			}

			// Something goes wrong, we will try from the cache even it is expired
			st.cacheSaver.Reset(logStep) // reset it will only remove the temporary file
			st.cacheSaver = nil          //Since we closed it, ensure that it will not be closed afterwards
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
			logStep.Error("Opening cache file failed")
			logStep.Error(err)
			return err
		}
		st.cacheLoader = cacheLoader
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

// ToSkip return true if the step must be skipped
func (st *Step) ToSkip(ctx context.Context, log *logrus.Entry) (bool, error) {
	logStep := log.WithField("name", st.Name).WithField("type", "sync")
	//TODO: to be implemented
	logStep.Debug("Do we need to skip the step ?") //TODO : may be at saver level ?
	return false, nil
}
