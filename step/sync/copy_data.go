package sync

import (
	"context"

	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/provider"
)

func (st *Step) copyData(ctx context.Context, log *logrus.Entry) error {
	source := st.source
	destinations := make([]provider.Saver, len(st.destinations))
	copy(destinations, st.destinations)

	if st.cacheSaver != nil {
		destinations = append(destinations, st.cacheSaver)
	} else if st.cacheLoader != nil {
		source = st.cacheLoader
	}

	log.Infof("Will synchronize %s to", source.Name())

	for _, d := range destinations {
		log.Infof("   - %s", d.Name())
	}

	if st.dryRun {
		return nil
	}

	for source.Next() {
		record, err := source.Load(log)
		if err != nil {
			log.Error("Source reading failed:")
			log.Error(err)

			return err
		}

		for _, f := range st.filters {
			if record, err = f.Filter(record); err != nil {
				log.Error("Filtering failed:")
				log.Error(err)

				return err
			}
		}

		for _, d := range destinations {
			if err = d.Save(log, record); err != nil {
				log.Error("Destination writing failed:")
				log.Error(err)

				return err
			}
		}
		st.count++

		if st.count%1000 == 0 {
			log.Infof("%d rows treated", st.count)
		}

		//Look for cancellation between each data
		select {
		case <-ctx.Done(): //If the context has been cancelled stop the recipe execution here
			log.Debug("Synchronization cancelled")
			return nil

		default: // Make the poll to ctx.Done() non blocking. Do nothing
		}
	}

	return nil
}
