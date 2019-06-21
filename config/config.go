package config

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/marema31/kamino/kaminodb"
	"github.com/spf13/viper"
)

// CacheConfig type for cache configuration with duration interpretation by viper
type CacheConfig struct {
	Type string
	File string
	TTL  time.Duration
}

// FilterConfig type for filter configuration to allow map and array parameters
type FilterConfig struct {
	Type   string
	AParam []string
	MParam map[string]string
}

// Sync type for characteristics of a sync
type Sync struct {
	Source       map[string]string
	Filters      []FilterConfig
	Destinations []map[string]string
	Cache        CacheConfig
}

// Config implements the config store of kamino
type Config struct {
	v         map[string]*viper.Viper
	databases map[string]*kaminodb.KaminoDb
}

// New initialize the config store
func New(path string, filename string) (*Config, error) {
	v := viper.New()

	v.SetConfigName(filename) // The file will be named [filename].json, [filename].yaml or [filename.toml]
	v.AddConfigPath(path)
	v.SetEnvPrefix("kamino")
	v.AutomaticEnv()
	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	subs := make(map[string]*viper.Viper)

	for _, k := range v.AllKeys() {
		//AllKeys return all the keys and subkeys as "sync1.source.engine" if we do not have subkey the sync entry is not valid
		if idx := strings.IndexByte(k, '.'); idx >= 0 {
			sync := k[:idx]

			_, ok := subs[sync]
			if !ok {
				subs[sync] = v.Sub(sync)
			}
		}
	}

	databases, err := getDatabases(path)
	if err != nil {
		log.Fatal(err)

		return nil, err
	}

	config := &Config{v: subs, databases: databases}
	return config, err
}

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

//GetDb return the kaminodb object from a database name
func (c *Config) GetDb(name string) (*kaminodb.KaminoDb, error) {
	kdb, ok := c.databases[name]
	if !ok {
		return nil, fmt.Errorf("databse %s does not have configuration file", name)
	}

	return kdb, nil
}
