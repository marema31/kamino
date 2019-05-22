package sync

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/marema31/kamino/config"
	"github.com/marema31/kamino/provider"
)

func copyData(ctx context.Context, source provider.Loader, destinations []provider.Saver) error {
	for source.Next() {
		record, err := source.Load()
		if err != nil {
			return err
		}

		for _, d := range destinations {
			if err = d.Save(record); err != nil {
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

	source, err := provider.NewLoader(ctx, c.Source)
	if err != nil {
		return err
	}
	defer source.Close()

	var destinations []provider.Saver

	for _, dest := range c.Destinations {
		d, err := provider.NewSaver(ctx, dest)
		if err != nil {
			return err
		}
		defer d.Close()

		destinations = append(destinations, d)
	}

	if c.Cache.File != "" {

		d, err := provider.NewSaver(ctx, map[string]string{"type": c.Cache.Type, "file": c.Cache.File + ".*"})
		if err != nil {
			return err
		}
		defer d.Close()
		defer fmt.Printf("will remove %s\n", d.Name())
		//		defer os.Remove(tempCache.Name())

		destinations = append(destinations, d)
		err = copyData(ctx, source, destinations)
		if err == nil {
			//everything was OK, I just rename the tempfile for cache to its real name
			err = os.Rename(d.Name(), c.Cache.File)
			if err != nil {
				return err
			}
			// For the moment nothing more to do
			return nil
		}
	}
	return copyData(ctx, source, destinations)
}
