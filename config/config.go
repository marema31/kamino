package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// CacheConfig type for cache configuration with duration interpretation by viper
type CacheConfig struct {
	Type string
	File string
	TTL  time.Duration
}

// Sync type for characteristics of a sync
type Sync struct {
	Source       map[string]string
	Destinations []map[string]string
	Cache        CacheConfig
}

// Config implements the config store of kamino
type Config struct {
	v map[string]*viper.Viper
}

// New initialize the config store
func New(path string, filename string) (*Config, error) {
	v := viper.New()

	v.SetConfigName(filename) // The file will be named [filename].json, [filename].yaml or [filename.toml]
	v.AddConfigPath(path)
	v.SetEnvPrefix("kamino")
	v.AutomaticEnv()
	err := v.ReadInConfig()

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

	config := &Config{v: subs}
	return config, err
}

// Get return the configuration block for a synchronization
func (c *Config) Get(sync string) (*Sync, error) {

	v, ok := c.v[sync]
	if !ok {
		return nil, fmt.Errorf("the configuration block for %s sync does not exist", sync)
	}

	source := v.GetStringMapString("source")
	_, ok = source["type"]
	if !ok {
		return nil, fmt.Errorf("source defined for %s sync is invalid", sync)
	}

	var cache CacheConfig
	cache.Type = v.GetString("cache.type")
	cache.TTL = v.GetDuration("cache.ttl")
	cache.File = v.GetString("cache.file")

	var dests []map[string]string
	ds := v.Get("destinations")

	switch casted := ds.(type) { // Avoid panic if the type is not compatible with the one we want
	case []interface{}:
		for _, d := range casted {
			switch dcasted := d.(type) { // Avoid panic if the type is not compatible with the one we want
			case map[string]interface{}:
				currentdest := make(map[string]string)
				for dk, dv := range dcasted {
					currentdest[dk] = dv.(string)
				}
				dests = append(dests, currentdest)
			default:
				return nil, fmt.Errorf("one destination defined for %s sync is invalid", sync)
			}
		}
	default:
		return nil, fmt.Errorf("destinations defined for %s sync is invalid", sync)
	}

	s := &Sync{
		Source:       source,
		Destinations: dests,
		Cache:        cache,
	}
	return s, nil
}
