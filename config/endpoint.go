package config

import (
	"log"

	"github.com/spf13/viper"
)

func (c *Config) getSource(v *viper.Viper, sync string) (SourceConfig, error) {
	var source SourceConfig
	sv := v.Sub("source")
	err := sv.Unmarshal(&source)
	if err != nil {
		return source, err
	}
	return source, nil
}

func (c *Config) getDestinations(v *viper.Viper, sync string) ([]DestinationConfig, error) {
	var dests []DestinationConfig
	err := v.UnmarshalKey("destinations", &dests)
	if err != nil {
		log.Fatal(err)
		return []DestinationConfig{}, err
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
