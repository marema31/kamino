package config

import (
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

// SourceConfig type for source contain all possible fields without verification
type SourceConfig struct {
	Type     string
	File     string
	URL      string
	Inline   string
	Gzip     bool
	Zip      bool
	Database string
	Instance string
	Table    string
	Where    string
}

// DestinationConfig type for destination contain all possible fields without verification
type DestinationConfig struct {
	Type      string
	File      string
	Gzip      bool
	Zip       bool
	Database  string
	Instances []string
	Table     string
	Key       string
	Mode      string
}

// Sync type for characteristics of a sync
type Sync struct {
	Source       SourceConfig
	Filters      []FilterConfig
	Destinations []DestinationConfig
	Cache        CacheConfig
}

// Config implements the config store of kamino
type Config struct {
	v         map[string]*viper.Viper
	databases map[string]map[string]map[string]*kaminodb.KaminoDb
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
