package sync

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/viper"

	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/filter"
	"github.com/marema31/kamino/provider"
	"github.com/marema31/kamino/step/common"
)

// SourceConfig type for source contain all possible fields without verification
type SourceConfig struct {
	Tags    []string
	Engines []string
	Types   []string
	Table   string
	Where   string
}

// DestinationConfig type for destination contain all possible fields without verification
type DestinationConfig struct {
	Tags    []string
	Engines []string
	Types   []string
	Table   string
	Key     string
	Mode    string
}

// FilterConfig type for filter contain all possible fields without verification
type FilterConfig struct {
	Aparameters []string
	Mparameters map[string]string
	Type        string
}

func getDatasource(dss datasource.Datasourcers, tags []string, engines []string, dsTypes []string) ([]datasource.Datasourcer, error) {
	if len(tags) == 0 {
		tags = []string{""}
	}

	e, err := datasource.StringsToEngines(engines)
	if err != nil {
		return nil, err
	}

	t, err := datasource.StringsToTypes(dsTypes)
	if err != nil {
		return nil, err
	}
	return dss.Lookup(tags, t, e), nil
}

func getLoader(ctx context.Context, name string, objectType string, v *viper.Viper, dss datasource.Datasourcers, prov provider.Provider) (provider.Loader, error) {
	var source SourceConfig
	err := v.Unmarshal(&source)
	if err != nil {
		return nil, err
	}

	datasources, err := getDatasource(dss, source.Tags, source.Engines, source.Types)
	if err != nil {
		return nil, err
	}
	if len(datasources) == 0 {
		return nil, fmt.Errorf("no %s found for synchronization %s", objectType, name)
	}
	if len(datasources) != 1 {
		return nil, fmt.Errorf("too many %ss found for synchronization %s", objectType, name)
	}
	return prov.NewLoader(ctx, datasources[0], source.Table, source.Where)
}

func getSavers(ctx context.Context, name string, objectType string, v *viper.Viper, dss datasource.Datasourcers, prov provider.Provider) ([]provider.Saver, error) {
	var dests []DestinationConfig
	savers := make([]provider.Saver, 0)

	err := v.UnmarshalKey(objectType, &dests)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	for _, dest := range dests {
		datasources, err := getDatasource(dss, dest.Tags, dest.Engines, dest.Types)
		if err != nil {
			return nil, err
		}

		for _, datasource := range datasources {
			saver, err := prov.NewSaver(ctx, datasource, dest.Table, dest.Key, dest.Mode)
			if err != nil {
				return nil, err
			}
			savers = append(savers, saver)
		}
	}
	if len(savers) == 0 {
		return nil, fmt.Errorf("no %s found for synchronization %s", objectType, name)
	}
	return savers, nil
}

func getFilters(ctx context.Context, v *viper.Viper, sync string) ([]filter.Filter, error) {
	fcs := make([]FilterConfig, 0)
	filters := make([]filter.Filter, 0)
	err := v.UnmarshalKey("filters", &fcs)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	for _, fc := range fcs {
		f, err := filter.NewFilter(ctx, fc.Type, fc.Aparameters, fc.Mparameters)
		if err != nil {
			return nil, err
		}
		filters = append(filters, f)

	}
	return filters, nil
}

//Load data from step file using its viper representation a return priority and list of steps
func Load(ctx context.Context, filename string, v *viper.Viper, dss datasource.Datasourcers, provider provider.Provider) (priority uint, steps []common.Steper, err error) {
	var step Step

	priority = v.GetUint("priority")

	name := v.GetString("name")
	step.Name = name

	if !v.IsSet("source") {
		return 0, nil, fmt.Errorf("synchronization %s does not have a source definition", name)
	}
	if !v.IsSet("destinations") {
		return 0, nil, fmt.Errorf("synchronization %s does not have a destinations definition", name)
	}

	//Lookup source
	sub := v.Sub("source")

	step.source, err = getLoader(ctx, name, "source", sub, dss, provider)
	if err != nil {
		return 0, nil, err
	}

	//Lookup cache
	if v.IsSet("cache") {
		step.cacheTTL = v.GetDuration("cache.ttl")
		sub = v.Sub("cache")
		step.cacheLoader, err = getLoader(ctx, name, "cache", sub, dss, provider)
		if err != nil {
			return 0, nil, err
		}
		cs, err := getSavers(ctx, name, "cache", v, dss, provider)
		if err != nil {
			return 0, nil, err
		}
		step.cacheSaver = cs[0]

	}

	//Lookup filters
	if v.IsSet("filters") {
		step.filters, err = getFilters(ctx, v, name)
		if err != nil {
			return 0, nil, err
		}
	}
	//Lookup destinations
	step.destinations, err = getSavers(ctx, name, "destinations", v, dss, provider)
	if err != nil {
		return 0, nil, err
	}

	steps = make([]common.Steper, 0, 1)
	steps = append(steps, step)
	return priority, steps, nil
}
