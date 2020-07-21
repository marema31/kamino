package sync

import (
	"context"
	"errors"
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

// SourceConfig type for source contain all possible fields without verification.
type SourceConfig struct {
	Tags    []string
	Engines []string
	Types   []string
	Table   string
	Where   string
}

// DestinationConfig type for destination contain all possible fields without verification.
type DestinationConfig struct {
	Tags    []string
	Engines []string
	Types   []string
	Table   string
	Key     string
	Mode    string
	Queries []string
}

// FilterConfig type for filter contain all possible fields without verification.
type FilterConfig struct {
	Aparameters []string
	Mparameters map[string]string
	Type        string
}

var errDatasource = errors.New("NOT CORRECT NUMBER OF DATASOURCES")

func getDatasources(log *logrus.Entry, dss datasource.Datasourcers, tags []string, engines []string, dsTypes []string, objectType string, unique bool, limitedTags []string) (limited []datasource.Datasourcer, notLimited []datasource.Datasourcer, err error) {
	if len(tags) == 0 {
		tags = []string{""}
	}

	e, err := datasource.StringsToEngines(engines)
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}

	t, err := datasource.StringsToTypes(dsTypes)
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}

	limited, notLimited, err = dss.Lookup(log, tags, limitedTags, t, e)
	if err != nil {
		return nil, nil, err
	}

	log.Debugf("Found %d %s (from %d available)", len(limited), objectType, len(notLimited))

	if len(limited) == 0 && len(notLimited) == 0 {
		log.Errorf("no %s found", objectType)
		return nil, nil, fmt.Errorf("no %s found: %w", objectType, errDatasource)
	}

	if unique && len(limited) != 1 {
		log.Errorf("too many %ss found", objectType)
		return nil, nil, fmt.Errorf("too many %ss found: %w", objectType, errDatasource)
	}

	return limited, notLimited, nil
}

func parseSourceConfig(log *logrus.Entry, objectType string, v *viper.Viper, dss datasource.Datasourcers) (parsedSourceConfig, parsedSourceConfig, error) {
	var (
		parsedLimitedSource    parsedSourceConfig
		parsedNotLimitedSource parsedSourceConfig
		source                 SourceConfig
	)

	err := v.Unmarshal(&source)
	if err != nil {
		log.Error(err)
		return parsedLimitedSource, parsedNotLimitedSource, err
	}

	limited, notLimited, err := getDatasources(log, dss, source.Tags, source.Engines, source.Types, objectType, true, nil)
	if err != nil {
		return parsedLimitedSource, parsedNotLimitedSource, err
	}

	parsedLimitedSource.ds = limited[0]
	parsedLimitedSource.table = source.Table
	parsedLimitedSource.where = source.Where
	parsedNotLimitedSource.ds = limited[0]
	parsedNotLimitedSource.table = source.Table
	parsedNotLimitedSource.where = source.Where

	if len(notLimited) != 0 {
		parsedNotLimitedSource.ds = notLimited[0]
	}

	return parsedLimitedSource, parsedNotLimitedSource, nil
}

func addParsedDest(log *logrus.Entry, parseDests []parsedDestConfig, datasource datasource.Datasourcer, dest DestinationConfig, tqueries []common.TemplateSkipQuery, force bool) ([]parsedDestConfig, error) {
	var p parsedDestConfig
	p.ds = datasource
	p.table = dest.Table
	p.key = dest.Key

	p.mode = strings.ToLower(dest.Mode)
	if p.mode == "onlyifempty" && force {
		p.mode = "truncate"
	}

	tmplValues := datasource.FillTmplValues()

	queries, err := common.RenderQueries(log, tqueries, tmplValues)
	if err != nil {
		return parseDests, err
	}

	p.queries = queries

	return append(parseDests, p), err
}

func parseDestConfig(log *logrus.Entry, v *viper.Viper, dss datasource.Datasourcers, force bool, limitedTags []string) ([]parsedDestConfig, []parsedDestConfig, error) {
	var dests []DestinationConfig

	parsedLimitedDests := make([]parsedDestConfig, 0)
	parsedNotLimitedDests := make([]parsedDestConfig, 0)

	err := v.UnmarshalKey("destinations", &dests)
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}

	for _, dest := range dests {
		tqueries, err := common.ParseQueries(log, dest.Queries)
		if err != nil {
			return nil, nil, err
		}

		limited, notLimited, err := getDatasources(log, dss, dest.Tags, dest.Engines, dest.Types, "destination", false, limitedTags)
		if err != nil {
			return nil, nil, err
		}

		for _, datasource := range limited {
			parsedLimitedDests, err = addParsedDest(log, parsedLimitedDests, datasource, dest, tqueries, force)
			if err != nil {
				return nil, nil, err
			}
		}

		for _, datasource := range notLimited {
			parsedNotLimitedDests, err = addParsedDest(log, parsedNotLimitedDests, datasource, dest, tqueries, force)
			if err != nil {
				return nil, nil, err
			}
		}
	}

	return parsedLimitedDests, parsedNotLimitedDests, nil
}

func getFilters(log *logrus.Entry, v *viper.Viper) ([]filter.Filter, error) {
	fcs := make([]FilterConfig, 0)
	filters := make([]filter.Filter, 0)

	err := v.UnmarshalKey("filters", &fcs)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	for _, fc := range fcs {
		f, err := filter.NewFilter(log, fc.Type, fc.Aparameters, fc.Mparameters)
		if err != nil {
			return nil, err
		}

		filters = append(filters, f)
	}

	log.Debugf("Found %d filters", len(filters))

	return filters, nil
}

//PostLoad modify the loaded step values with the values provided in the map in argument.
func (st *Step) PostLoad(log *logrus.Entry, superseed map[string]string) (err error) {
	if value, ok := superseed["sync.forceCacheOnly"]; ok {
		st.forceCacheOnly, err = strconv.ParseBool(value)
	}

	return err
}

//Load data from step file using its viper representation a return priority and list of steps.
func Load(ctx context.Context, log *logrus.Entry, recipePath string, name string, nameIndex int, v *viper.Viper, dss datasource.Datasourcers, provider provider.Provider, force bool, dryRun bool, limitedTags []string) (priority uint, steps []common.Steper, err error) {
	var step Step

	step.baseFolder = recipePath
	step.prov = provider
	priority = v.GetUint("priority")
	step.ignoreErrors = v.GetBool("ignoreErrors")

	logStep := log.WithField("name", name).WithField("type", "shell")
	step.Name = fmt.Sprintf("%s:%d", name, nameIndex)
	step.dryRun = dryRun

	if !v.IsSet("source") {
		logStep.Error("No source provided")
		return 0, nil, fmt.Errorf("no source definition: %w", common.ErrMissingParameter)
	}

	if !v.IsSet("destinations") {
		logStep.Error("No destinations provided")
		return 0, nil, fmt.Errorf("no destinations definition: %w", common.ErrMissingParameter)
	}

	logStep.Debug("Lookup source")

	sub := v.Sub("source")

	parsedLimitedSrcCfg, parsedNotLimitedSrcCfg, err := parseSourceConfig(logStep, "source", sub, dss)
	if err != nil {
		return 0, nil, err
	}

	logStep.Debug("Lookup cache")

	if v.IsSet("cache") {
		step.cacheTTL = v.GetDuration("cache.ttl")
		step.allowCacheOnly = v.GetBool("cache.allowonly")
		sub = v.Sub("cache")

		_, step.cacheCfg, err = parseSourceConfig(logStep, "cache", sub, dss)
		if err != nil {
			return 0, nil, err
		}
	}

	log.Debug("Lookup filters")

	if v.IsSet("filters") {
		step.filters, err = getFilters(logStep, v)
		if err != nil {
			return 0, nil, err
		}
	}

	log.Debug("Lookup destinations")

	parsedLimitedDestsCfg, parsedNotLimitedDestsCfg, err := parseDestConfig(logStep, v, dss, force, limitedTags)
	if err != nil {
		return 0, nil, err
	}

	step.sourceCfg = parsedLimitedSrcCfg
	step.destsCfg = parsedLimitedDestsCfg

	steps = make([]common.Steper, 0, 1)

	if len(limitedTags) != 0 {
		switch {
		case parsedLimitedSrcCfg.ds == nil && len(parsedLimitedDestsCfg) == 0:
			log.Debug("Skipping since no tags selected")
			return priority, steps, nil
		case parsedLimitedSrcCfg.ds != nil && len(parsedLimitedDestsCfg) == 0:
			step.sourceCfg = parsedLimitedSrcCfg
			step.destsCfg = parsedNotLimitedDestsCfg
		case parsedLimitedSrcCfg.ds == nil && len(parsedLimitedDestsCfg) != 0:
			step.sourceCfg = parsedNotLimitedSrcCfg
			step.destsCfg = parsedLimitedDestsCfg
		}
	}

	if step.sourceCfg.ds == nil {
		log.Error("No source found")
		return 0, nil, fmt.Errorf("no source found: %w", errDatasource)
	}

	if len(step.destsCfg) == 0 {
		log.Error("No destination found")
		return 0, nil, fmt.Errorf("no destination found: %w", errDatasource)
	}

	steps = append(steps, &step)

	return priority, steps, nil
}
