package sync

import (
	"context"
	"fmt"
	"strconv"
	"strings"

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

func getDatasources(log *logrus.Entry, dss datasource.Datasourcers, tags []string, engines []string, dsTypes []string, objectType string, unique bool, dryRun bool, limitedTags []string) ([]datasource.Datasourcer, error) {
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
	datasources := dss.Lookup(log, tags, limitedTags, t, e)
	log.Debugf("Found %d %s", len(datasources), objectType)
	if len(datasources) == 0 {
		log.Errorf("no %s found", objectType)
		return nil, fmt.Errorf("no %s found", objectType)
	}
	if unique && len(datasources) != 1 {
		log.Errorf("too many %ss found", objectType)
		return nil, fmt.Errorf("too many %ss found", objectType)
	}

	return datasources, nil
}

func parseSourceConfig(ctx context.Context, log *logrus.Entry, objectType string, v *viper.Viper, dss datasource.Datasourcers, dryRun bool) (parsedSourceConfig, error) {
	var parsedSource parsedSourceConfig
	var source SourceConfig
	err := v.Unmarshal(&source)
	if err != nil {
		log.Error(err)
		return parsedSource, err
	}

	datasources, err := getDatasources(log, dss, source.Tags, source.Engines, source.Types, objectType, true, dryRun, nil)
	if err != nil {
		return parsedSource, err
	}
	parsedSource.ds = datasources[0]
	parsedSource.table = source.Table
	parsedSource.where = source.Where
	return parsedSource, nil
}

func parseDestConfig(ctx context.Context, log *logrus.Entry, v *viper.Viper, dss datasource.Datasourcers, force bool, dryRun bool, limitedTags []string) ([]parsedDestConfig, error) {
	var dests []DestinationConfig
	parsedDests := make([]parsedDestConfig, 0)

	err := v.UnmarshalKey("destinations", &dests)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	for _, dest := range dests {
		datasources, err := getDatasources(log, dss, dest.Tags, dest.Engines, dest.Types, "destination", false, dryRun, limitedTags)
		if err != nil {
			return nil, err
		}

		for _, datasource := range datasources {
			var p parsedDestConfig
			p.ds = datasource
			p.table = dest.Table
			p.key = dest.Key
			p.mode = strings.ToLower(dest.Mode)
			if p.mode == "onlyifempty" && force {
				p.mode = "truncate"
			}
			parsedDests = append(parsedDests, p)
		}
	}
	return parsedDests, nil
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

//PostLoad modify the loaded step values with the values provided in the map in argument
func (st *Step) PostLoad(log *logrus.Entry, superseed map[string]string) (err error) {
	if value, ok := superseed["sync.forceCacheOnly"]; ok {
		st.forceCacheOnly, err = strconv.ParseBool(value)
	}
	return err
}

//Load data from step file using its viper representation a return priority and list of steps
func Load(ctx context.Context, log *logrus.Entry, recipePath string, name string, nameIndex int, v *viper.Viper, dss datasource.Datasourcers, provider provider.Provider, force bool, dryRun bool, limitedTags []string) (priority uint, steps []common.Steper, err error) {
	var step Step

	step.baseFolder = recipePath
	step.prov = provider
	priority = v.GetUint("priority")

	logStep := log.WithField("name", name).WithField("type", "shell")
	step.Name = fmt.Sprintf("%s:%d", name, nameIndex)
	step.dryRun = dryRun

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

	step.sourceCfg, err = parseSourceConfig(ctx, logStep, "source", sub, dss, dryRun)
	if err != nil {
		return 0, nil, err
	}

	logStep.Debug("Lookup cache")
	if v.IsSet("cache") {
		step.cacheTTL = v.GetDuration("cache.ttl")
		step.allowCacheOnly = v.GetBool("cache.allowonly")
		sub = v.Sub("cache")
		step.cacheCfg, err = parseSourceConfig(ctx, logStep, "cache", sub, dss, dryRun)
		if err != nil {
			return 0, nil, err
		}

	}

	log.Debug("Lookup filters")
	if v.IsSet("filters") {
		step.filters, err = getFilters(ctx, logStep, v)
		if err != nil {
			return 0, nil, err
		}
	}

	log.Debug("Lookup destinations")
	step.destsCfg, err = parseDestConfig(ctx, logStep, v, dss, force, dryRun, limitedTags)
	if err != nil {
		return 0, nil, err
	}

	steps = make([]common.Steper, 0, 1)
	steps = append(steps, &step)
	return priority, steps, nil
}
