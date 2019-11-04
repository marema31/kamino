package sync

import (
	"context"
	"fmt"

	"github.com/Sirupsen/logrus"
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

func getDatasource(log *logrus.Entry, dss datasource.Datasourcers, tags []string, engines []string, dsTypes []string) ([]datasource.Datasourcer, error) {
	if len(tags) == 0 {
		tags = []string{""}
	}

	e, err := datasource.StringsToEngines(engines)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	t, err := datasource.StringsToTypes(dsTypes)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return dss.Lookup(log, tags, t, e), nil
}

func getLoader(ctx context.Context, log *logrus.Entry, objectType string, v *viper.Viper, dss datasource.Datasourcers, prov provider.Provider) (provider.Loader, error) {
	var source SourceConfig
	err := v.Unmarshal(&source)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	datasources, err := getDatasource(log, dss, source.Tags, source.Engines, source.Types)
	if err != nil {
		return nil, err
	}
	if len(datasources) == 0 {
		return nil, fmt.Errorf("no %s found", objectType)
	}
	if len(datasources) != 1 {
		return nil, fmt.Errorf("too many %ss found", objectType)
	}
	log.Debugf("Found 1 datasource for %s", objectType)

	log.Debugf("Creating loader instance for %s", objectType)
	return prov.NewLoader(ctx, log, datasources[0], source.Table, source.Where)
}

func getSavers(ctx context.Context, log *logrus.Entry, objectType string, v *viper.Viper, dss datasource.Datasourcers, prov provider.Provider) ([]provider.Saver, error) {
	var dests []DestinationConfig
	savers := make([]provider.Saver, 0)

	err := v.UnmarshalKey(objectType, &dests)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	for _, dest := range dests {
		datasources, err := getDatasource(log, dss, dest.Tags, dest.Engines, dest.Types)
		if err != nil {
			return nil, err
		}

		for _, datasource := range datasources {
			log.Debugf("Creating saver instances for %s", objectType)
			saver, err := prov.NewSaver(ctx, log, datasource, dest.Table, dest.Key, dest.Mode)
			if err != nil {
				return nil, err
			}
			savers = append(savers, saver)
		}
	}
	if len(savers) == 0 {
		log.Errorf("No %s found", objectType)
		return nil, fmt.Errorf("no %s found", objectType)
	}
	log.Debugf("Found %d %s", len(savers), objectType)
	return savers, nil
}

func getFilters(ctx context.Context, log *logrus.Entry, v *viper.Viper) ([]filter.Filter, error) {
	fcs := make([]FilterConfig, 0)
	filters := make([]filter.Filter, 0)
	err := v.UnmarshalKey("filters", &fcs)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	for _, fc := range fcs {
		f, err := filter.NewFilter(ctx, log, fc.Type, fc.Aparameters, fc.Mparameters)
		if err != nil {
			return nil, err
		}
		filters = append(filters, f)

	}
	log.Debugf("Found %d filters", len(filters))
	return filters, nil
}

//Load data from step file using its viper representation a return priority and list of steps
func Load(ctx context.Context, log *logrus.Entry, filename string, v *viper.Viper, dss datasource.Datasourcers, provider provider.Provider) (priority uint, steps []common.Steper, err error) {
	var step Step

	priority = v.GetUint("priority")

	name := v.GetString("name")
	logStep := log.WithField("name", name).WithField("type", "shell")
	step.Name = name

	if !v.IsSet("source") {
		logStep.Error("No source provided")
		return 0, nil, fmt.Errorf("no source definition")
	}
	if !v.IsSet("destinations") {
		logStep.Error("No destinations provided")
		return 0, nil, fmt.Errorf("no destinations definition")
	}

	logStep.Debug("Lookup source")
	sub := v.Sub("source")

	step.source, err = getLoader(ctx, logStep, "source", sub, dss, provider)
	if err != nil {
		return 0, nil, err
	}

	logStep.Debug("Lookup cache")
	if v.IsSet("cache") {
		step.cacheTTL = v.GetDuration("cache.ttl")
		sub = v.Sub("cache")
		step.cacheLoader, err = getLoader(ctx, logStep, "cache", sub, dss, provider)
		if err != nil {
			return 0, nil, err
		}
		cs, err := getSavers(ctx, logStep, "cache", v, dss, provider)
		if err != nil {
			return 0, nil, err
		}
		step.cacheSaver = cs[0]

	}

	log.Debug("Lookup filters")
	if v.IsSet("filters") {
		step.filters, err = getFilters(ctx, logStep, v)
		if err != nil {
			return 0, nil, err
		}
	}
	log.Debug("Lookup destinations")
	step.destinations, err = getSavers(ctx, logStep, "destinations", v, dss, provider)
	if err != nil {
		return 0, nil, err
	}

	steps = make([]common.Steper, 0, 1)
	steps = append(steps, &step)
	return priority, steps, nil
}
