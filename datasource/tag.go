package datasource

import (
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/mb0/glob"
)

func (ds *Datasource) isSelectedEngine(engines []Engine) bool {
	if len(engines) == 0 {
		return true
	}

	for _, e := range engines {
		if e == ds.engine {
			return true
		}
	}

	return false
}

func (ds *Datasource) isSelectedType(dsTypes []Type) bool {
	if len(dsTypes) == 0 {
		return true
	}

	for _, t := range dsTypes {
		if t == ds.dstype {
			return true
		}
	}

	return false
}

func (dss *Datasources) lookupGlobTag(log *logrus.Entry, tag string, dsTypes []Type, engines []Engine) (selected []string, err error) {
	selected = make([]string, 0)

	g, err := glob.New(glob.Default())
	if err != nil {
		log.Errorf("Unable to initialize the globbing engine: %v", err)
		return selected, err
	}

	for _, names := range dss.tagToDatasource {
		for _, name := range names {
			matched, err := g.Match(tag, name)
			if err != nil {
				log.Errorf("Using the tag %s failed: %v", tag, err)
				return selected, err
			}

			if matched {
				ds := dss.datasources[name]
				if ds.isSelectedEngine(engines) && ds.isSelectedType(dsTypes) {
					selected = append(selected, name)
				}
			}
		}
	}

	return selected, nil
}

func (dss *Datasources) lookupOneTag(log *logrus.Entry, tag string, dsTypes []Type, engines []Engine) (selected []string, err error) {
	tag = strings.TrimPrefix(tag, "!") // The negation is not useful here, it will be managed by caller

	if strings.ContainsAny(tag, "*[") {
		return dss.lookupGlobTag(log, tag, dsTypes, engines)
	}

	for _, name := range dss.tagToDatasource[tag] {
		ds := dss.datasources[name]
		if ds.isSelectedEngine(engines) && ds.isSelectedType(dsTypes) {
			selected = append(selected, name)
		}
	}

	return selected, nil
}

func (dss *Datasources) lookupWithoutTag(dsTypes []Type, engines []Engine) (selectedNames []string) {
	selected := make(map[string]bool) // Use map to emulate a "set" to avoid duplicates

	for _, names := range dss.tagToDatasource {
		for _, name := range names {
			ds := dss.datasources[name]
			if ds.isSelectedEngine(engines) && ds.isSelectedType(dsTypes) {
				selected[name] = true
			}
		}
	}

	selectedNames = make([]string, 0, len(selected))

	for dsName := range selected {
		selectedNames = append(selectedNames, dsName)
	}

	return selectedNames
}

func (dss *Datasources) findCandidatesTags(log *logrus.Entry, tagElement string, dsTypes []Type, engines []Engine) (map[string]bool, error) {
	candidates := make(map[string]bool)

	for i, tag := range strings.Split(tagElement, ".") {
		log.Debugf("Lookup for sub-tag %s", tag)

		names, err := dss.lookupOneTag(log, tag, dsTypes, engines)
		if err != nil {
			return candidates, err
		}

		for _, name := range names {
			if i == 0 {
				log.Debugf("Found: %s", name)

				candidates[name] = true
			} else {
				if _, ok := candidates[name]; ok {
					log.Debugf("Found: %s, present from previous sub-tag", name)

					candidates[name] = true
				} else {
					log.Debugf("Skipped: %s, not present for previous sub-tag", name)
				}
			}
		}

		log.Debugf("Removing datasource not found for this sub-tag: %s", tag)

		for name, viewed := range candidates {
			if !viewed {
				log.Debugf("Removing: %s", name)
				delete(candidates, name)
			} else {
				// Prepare the map for the next tag
				candidates[name] = false
			}
		}
	}

	return candidates, nil
}

//nolint:nestif
func (dss *Datasources) getDsNameFromTagList(log *logrus.Entry, tagList []string, dsTypes []Type, engines []Engine) ([]string, error) {
	selected := make(map[string]bool) // Use map to emulate a "set" to avoid duplicates
	unselectedNames := make([]string, 0)

	for _, tagElement := range tagList {
		log.Debugf("Lookup %s", tagElement)

		if !strings.Contains(tagElement, ".") {
			log.Debug("Simple tag")

			names, err := dss.lookupOneTag(log, tagElement, dsTypes, engines)
			if err != nil {
				return nil, err
			}

			if strings.HasPrefix(tagElement, "!") {
				log.Debug("Negative tag")

				unselectedNames = append(unselectedNames, names...)
			} else {
				for _, name := range names {
					selected[name] = true
				}
			}
		} else {
			log.Debug("Composite tag")
			candidates, err := dss.findCandidatesTags(log, tagElement, dsTypes, engines)

			if err != nil {
				return nil, err
			}
			log.Debugf("Final list for %s", tagElement)
			for name := range candidates {
				if strings.HasPrefix(tagElement, "!") {
					unselectedNames = append(unselectedNames, name)
				} else {
					log.Debugf("  - %s", name)
					selected[name] = true
				}
			}
		}
	}

	selectedNames := make([]string, 0, len(selected))
	for dsName := range selected {
		selectedNames = append(selectedNames, dsName)
	}

	return removeFromList(selectedNames, unselectedNames), nil
}

func removeFromList(selected []string, unselected []string) (filtered []string) {
	for _, dsName := range selected {
		found := false

		for _, name := range unselected {
			if name == dsName {
				found = true
				break
			}
		}

		if !found {
			filtered = append(filtered, dsName)
		}
	}

	return filtered
}

func (dss *Datasources) lookupLimited(log *logrus.Entry, tagList []string, limitedTags []string, dsTypes []Type, engines []Engine) ([]string, error) {
	log.Debug("Lookup limited tag list")

	limited := make([]string, 0)

	if limitedTags == nil {
		return limited, nil
	}

	limited, err := dss.getDsNameFromTagList(log, limitedTags, dsTypes, engines)
	if err != nil {
		return nil, err
	}

	if len(limited) == 0 {
		allNegation := true

		for _, tag := range limitedTags {
			if !strings.HasPrefix(tag, "!") {
				allNegation = false
				break
			}
		}

		if allNegation {
			limited, err = dss.getDsNameFromTagList(log, append(tagList, limitedTags...), dsTypes, engines)
			if err != nil {
				return nil, err
			}
		}
	}

	return limited, nil
}

//Lookup return a list of *Datasource that correspond to the
// list of tag expression.
func (dss *Datasources) Lookup(log *logrus.Entry, tagList []string, limitedTags []string, dsTypes []Type, engines []Engine) (selectedDs []Datasourcer, notLimitedDs []Datasourcer, err error) {
	logLookup := log.WithField("lookup", "tags")

	var selected []string

	limited, err := dss.lookupLimited(logLookup, tagList, limitedTags, dsTypes, engines)
	if err != nil {
		return nil, nil, err
	}

	if len(tagList) == 0 {
		logLookup.Debug("No tag provided, will only lookup on type and engines")
		// The selection is not based on tag, lookup for all of them
		selected = dss.lookupWithoutTag(dsTypes, engines)
	} else {
		selected, err = dss.getDsNameFromTagList(logLookup, tagList, dsTypes, engines)
		if err != nil {
			return nil, nil, err
		}
	}

	logLookup.Debug("Final datasources list:")

	finalDsList := make([]string, 0, len(selected))
	selectedDs = make([]Datasourcer, 0, len(selected))
	notLimitedDs = make([]Datasourcer, 0, len(selected))

	for _, dsName := range selected {
		inLimit := true
		if limitedTags != nil {
			inLimit = false

			for _, name := range limited {
				if name == dsName {
					inLimit = true
					break
				}
			}
		}

		notLimitedDs = append(notLimitedDs, dss.datasources[dsName])

		if inLimit {
			finalDsList = append(finalDsList, dsName)
			selectedDs = append(selectedDs, dss.datasources[dsName])
		}
	}

	logLookup.Debug(strings.Join(finalDsList, ","))

	return selectedDs, notLimitedDs, nil
}
