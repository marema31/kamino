package sync

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/marema31/kamino/config"
	"github.com/marema31/kamino/filter"
	"github.com/marema31/kamino/provider"
)

func copyData(ctx context.Context, source provider.Loader, filters []filter.Filter, destinations []provider.Saver) error {
	for source.Next() {
		record, err := source.Load()
		if err != nil {
			return err
		}

		for _, f := range filters {
			if record, err = f.Filter(record); err != nil {
				return err
			}
		}
		for _, d := range destinations {
			if err = d.Save(record); err != nil {
				d.Reset()
				return err
			}
		}

	}
	return nil
}

// Do will manage a single sync configuration and do tha actual copy
func Do(ctx context.Context, config *config.Config, syncName string) error {
	c, err := config.Get(syncName)
	if err != nil {
		log.Fatal(err)
	}

	var source provider.Loader
	var filters []filter.Filter
	var destinations []provider.Saver

	for _, fil := range c.Filters {
		f, err := filter.NewFilter(ctx, fil)
		if err != nil {
			return err
		}
		filters = append(filters, f)
	}

	for _, dest := range c.Destinations {
		d, err := provider.NewSaver(ctx, dest)
		if err != nil {
			return err
		}
		defer d.Close()

		destinations = append(destinations, d)
	}

	if c.Cache.File == "" {
		source, err = provider.NewLoader(ctx, c.Source)
		if err != nil {
			return err
		}
		defer source.Close()
	} else {
		ttlExpired := false
		cacheStat, errFile := os.Stat(c.Cache.File)

		if errFile == nil {
			ttlExpired = time.Since(cacheStat.ModTime()) > c.Cache.TTL
		}
		if os.IsNotExist(errFile) || ttlExpired {
			// The cache file does not exists or older than precised TTL we will (re)create it
			d, err := provider.NewSaver(ctx, map[string]string{"type": c.Cache.Type, "file": c.Cache.File})
			if err != nil {
				return err
			}
			defer d.Close()

			// We will use the source provided
			source, err := provider.NewLoader(ctx, c.Source)
			if err == nil {
				// No error on opening the correct source, we continue
				defer source.Close()

				destinations = append(destinations, d)
				err = copyData(ctx, source, filters, destinations)
				if err == nil {
					// For the moment nothing more to do
					return nil
				}
			}
			//Something goes wrong, we remove the cache from destination since we will want to use it
			destinations = destinations[:len(destinations)-1]
		}
		// Generation of cache file failed we will try to use the old one if it exists
		if _, err := os.Stat(c.Cache.File); os.IsNotExist(err) {
			return err
		}

		// It exists, we will use the cache file as source
		source, err = provider.NewLoader(ctx, map[string]string{"type": c.Cache.Type, "file": c.Cache.File})
		if err != nil {
			return err
		}
		defer source.Close()
	}
	return copyData(ctx, source, filters, destinations)
}
