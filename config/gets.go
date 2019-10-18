package config

import (
	"fmt"
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
