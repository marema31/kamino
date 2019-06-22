package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func (c *Config) getSource(v *viper.Viper, sync string) (map[string]string, error) {
	source := v.GetStringMapString("source")
	_, ok := source["type"]
	if !ok {
		return nil, fmt.Errorf("source defined for %s sync is invalid", sync)
	}
	return source, nil
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
