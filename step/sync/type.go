//Package sync manages step that synchronize datasets between datasources
package sync

import (
	"time"

	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/step/common"

	"github.com/marema31/kamino/filter"
	"github.com/marema31/kamino/provider"
)

type parsedSourceConfig struct {
	ds    datasource.Datasourcer
	table string
	where string
}

type parsedDestConfig struct {
	ds      datasource.Datasourcer
	table   string
	key     string
	mode    string
	queries []common.SkipQuery
}

// Step informations.
type Step struct {
	Name           string
	baseFolder     string
	source         provider.Loader
	cacheLoader    provider.Loader
	cacheSaver     provider.Saver
	cacheTTL       time.Duration
	allowCacheOnly bool
	forceCacheOnly bool
	destinations   []provider.Saver
	filters        []filter.Filter
	sourceCfg      parsedSourceConfig
	cacheCfg       parsedSourceConfig
	destsCfg       []parsedDestConfig
	prov           provider.Provider
	dryRun         bool
	count          int
	ignoreErrors   bool
}
