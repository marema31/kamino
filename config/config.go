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

func (c *Config) getSource(v *viper.Viper, sync string) (map[string]string, error) {
	source := v.GetStringMapString("source")
	_, ok := source["type"]
	if !ok {
		return nil, fmt.Errorf("source defined for %s sync is invalid", sync)
	}
	return source, nil
}

func (c *Config) getFilters(v *viper.Viper, sync string) ([]FilterConfig, error) {
	var filters []FilterConfig
	var ok bool
	fs := v.Get("filters")
	if fs != nil { // There is a filter section

		// filters section is an json array of map
		switch casted := fs.(type) { // Avoid panic if the type is not compatible with the one we want
		case []interface{}:

			// for every element of the filters array that must be a map
			for _, f := range casted {
				var currentfilter FilterConfig
				switch fcasted := f.(type) { // Avoid panic if the type is not compatible with the one we want
				case map[string]interface{}:
					currentfilter.Type, ok = fcasted["type"].(string)
					if !ok {
						return nil, fmt.Errorf("missing type for a filter of %s sync", sync)
					}

					var parameters interface{}
					parameters, ok := fcasted["parameters"]
					if !ok {
						return nil, fmt.Errorf("missing parameters for a filter of %s sync", sync)
					}

					switch pcasted := parameters.(type) { // Avoid panic if the type is not compatible with the one we want
					case []interface{}:
						currentpvalue := make([]string, 0)
						for _, pv := range pcasted {
							currentpvalue = append(currentpvalue, pv.(string))
						}
						currentfilter.AParam = currentpvalue

					case map[string]interface{}:
						ps := make(map[string]string)
						for pk, pv := range pcasted {
							ps[pk] = pv.(string)
						}
						currentfilter.MParam = ps
					}
					filters = append(filters, currentfilter)
				default:
					return nil, fmt.Errorf("one filter defined for %s sync is invalid", sync)
				}
			}
		default:
			return nil, fmt.Errorf("filters defined for %s sync is invalid", sync)
		}
	}
	return filters, nil
}

func (c *Config) getDestinations(v *viper.Viper, sync string) ([]map[string]string, error) {
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
	return dests, nil
}

func (c *Config) getCache(v *viper.Viper, sync string) (CacheConfig, error) {
	var cache CacheConfig
	cache.Type = v.GetString("cache.type")
	cache.TTL = v.GetDuration("cache.ttl")
	cache.File = v.GetString("cache.file")
	return cache, nil
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
