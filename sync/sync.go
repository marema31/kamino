package sync

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/marema31/kamino/config"
	"github.com/marema31/kamino/filter"
	"github.com/marema31/kamino/provider"
)

type kaminoSync struct {
	syncName     string
	source       provider.Loader
	filters      []filter.Filter
	destinations []provider.Saver
}

func copyData(ctx context.Context, ks *kaminoSync) error {

	for ks.source.Next() {
		record, err := ks.source.Load()
		if err != nil {
			return err
		}

		for _, f := range ks.filters {
			if record, err = f.Filter(record); err != nil {
				return err
			}
		}
		for _, d := range ks.destinations {
			if err = d.Save(record); err != nil {
				return err
			}
		}

	}
	return nil
}

// Do will manage a single sync configuration and do tha actual copy
func Do(ctx context.Context, config *config.Config, syncName string) error {

	ks := &kaminoSync{
		syncName:     syncName,
		source:       nil,
		filters:      make([]filter.Filter, 0),
		destinations: make([]provider.Saver, 0),
	}

	end := make(chan bool, 1)

	go func() {
		select {
		case <-ctx.Done(): // the context has been cancelled
			resetAll(ks)
			log.Printf("Synchro %s aborted", ks.syncName)
		case <-end: // the channel has a information
			log.Printf("Synchro %s finished", ks.syncName)
		}
	}()

	c, err := config.Get(syncName)
	if err != nil {
		return fmt.Errorf("sync %s does not defined", syncName)
	}

	for _, fil := range c.Filters {
		f, err := filter.NewFilter(ctx, fil)
		if err != nil {
			return err
		}
		ks.filters = append(ks.filters, f)
	}

	for _, dest := range c.Destinations {
		d, err := provider.NewSaver(ctx, config, dest)
		if err != nil {
			return err
		}

		ks.destinations = append(ks.destinations, d)
	}

	if c.Cache.File == "" {
		ks.source, err = provider.NewLoader(ctx, config, c.Source)
		if err != nil {
			return err
		}
	} else {
		ttlExpired := false
		cacheStat, errFile := os.Stat(c.Cache.File)

		if errFile == nil {
			ttlExpired = time.Since(cacheStat.ModTime()) > c.Cache.TTL
		}
		if os.IsNotExist(errFile) || ttlExpired {
			// The cache file does not exists or older than precised TTL we will (re)create it
			cache, err := provider.NewSaver(ctx, config, map[string]string{"type": c.Cache.Type, "file": c.Cache.File})
			if err != nil {
				return err
			}

			// We will use the source provided
			ks.source, err = provider.NewLoader(ctx, config, c.Source)
			if err == nil {
				// No error on opening the correct source, we continue

				ks.destinations = append(ks.destinations, cache)
				err = copyData(ctx, ks)
				if err == nil {
					closeAll(ks)
					return nil
				}
			}
			//Something goes wrong, we remove the cache from destination since we will want to use it
			ks.destinations = ks.destinations[:len(ks.destinations)-1]
			cache.Reset()
			ks.source.Close()
			ks.source = nil
		}
		// Generation of cache file failed we will try to use the old one if it exists
		if _, err := os.Stat(c.Cache.File); os.IsNotExist(err) {
			return err
		}

		// It exists, we will use the cache file as source
		ks.source, err = provider.NewLoader(ctx, config, map[string]string{"type": c.Cache.Type, "file": c.Cache.File})
		if err != nil {
			return err
		}

	}
	err = copyData(ctx, ks)
	if err != nil {
		return err
	}
	closeAll(ks)
	end <- true
	return nil
}

func closeAll(ks *kaminoSync) {
	if ks.source != nil {
		ks.source.Close()
	}
	for _, d := range ks.destinations {
		d.Close()

	}
}

func resetAll(ks *kaminoSync) {
	if ks.source != nil {
		ks.source.Close()
	}
	for _, d := range ks.destinations {
		d.Reset()

	}
}
