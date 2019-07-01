package config

import (
	"fmt"

	"github.com/marema31/kamino/kaminodb"
)

// Get return the configuration block for a synchronization
func (c *Config) Get(sync string) (*Sync, error) {

	v, ok := c.v[sync]
	if !ok {
		return nil, fmt.Errorf("the configuration block for %s sync does not exist", sync)
	}

	source, err := c.getSource(v, sync)
	if err != nil {
		return nil, err
	}

	filters, err := c.getFilters(v, sync)
	if err != nil {
		return nil, err
	}

	cache, err := c.getCache(v, sync)
	if err != nil {
		return nil, err
	}

	dests, err := c.getDestinations(v, sync)
	if err != nil {
		return nil, err
	}

	s := &Sync{
		Source:       source,
		Filters:      filters,
		Destinations: dests,
		Cache:        cache,
	}
	return s, nil
}

//GetDbs return array of the kaminodb object from a database name, environment and instances
func (c *Config) GetDbs(name string, environment string, instances []string) ([]*kaminodb.KaminoDb, error) {
	var kdbs []*kaminodb.KaminoDb
	_, ok := c.databases[name]
	if !ok {
		return nil, fmt.Errorf("database %s does not have configuration file", name)
	}

	if environment == "" {
		if len(c.databases[name]) > 1 {
			return nil, fmt.Errorf("database %s have more than one environment, you must provided the one you want to use", name)
		}
		for e := range c.databases[name] {
			environment = e
		}
	}

	if instances == nil {
		for _, kdb := range c.databases[name][environment] {
			kdbs = append(kdbs, kdb)
		}

		return kdbs, nil
	}

	for _, instance := range instances {
		kdb, ok := c.databases[name][environment][instance]
		if !ok {
			return nil, fmt.Errorf("database %s does not have %s instance", name, instance)
		}
		kdbs = append(kdbs, kdb)
	}
	return kdbs, nil
}
