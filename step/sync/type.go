//Package sync manages step that synchronize datasets between datasources
package sync

import (
	"time"

	"github.com/marema31/kamino/filter"
	"github.com/marema31/kamino/provider"
)

// Step informations
type Step struct {
	Name         string
	source       provider.Loader
	cacheLoader  provider.Loader
	cacheSaver   provider.Saver
	cacheTTL     time.Duration
	destinations []provider.Saver
	filters      []filter.Filter
}
