package datasource

import (
	"strings"

	"github.com/Sirupsen/logrus"
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

func (dss *Datasources) lookupOneTag(tag string, dsTypes []Type, engines []Engine) (selected []string) {
	tag = strings.TrimPrefix(tag, "!") // The negation is not useful here, it will be managed by caller

	for _, name := range dss.tagToDatasource[tag] {
		ds := dss.datasources[name]
		if ds.isSelectedEngine(engines) && ds.isSelectedType(dsTypes) {
			selected = append(selected, name)
		}
	}
	return selected
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

func (dss *Datasources) getDsNameFromTagList(log *logrus.Entry, tagList []string, dsTypes []Type, engines []Engine) (selectedNames []string, unselectedNames []string) {
	selected := make(map[string]bool) // Use map to emulate a "set" to avoid duplicates
	unselectedNames = make([]string, 0)
	for _, tagElement := range tagList {
		log.Debugf("Lookup %s", tagElement)
		if !strings.Contains(tagElement, ".") {
			log.Debug("Simple tag")
			if strings.HasPrefix(tagElement, "!") {
				log.Debug("Negative tag")
				unselectedNames = append(unselectedNames, dss.lookupOneTag(tagElement, dsTypes, engines)...)
			} else {
				for _, name := range dss.lookupOneTag(tagElement, dsTypes, engines) {
					selected[name] = true
				}
			}
		} else {
			log.Debug("Composite tag")
			candidates := make(map[string]bool)
			for i, tag := range strings.Split(tagElement, ".") {
				log.Debugf("Lookup for sub-tag %s", tag)
				for _, name := range dss.lookupOneTag(tag, dsTypes, engines) {
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
				log.Debug("Removing datasource not found for this sub-tag")
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

	selectedNames = make([]string, 0, len(selected))
	for dsName := range selected {
		selectedNames = append(selectedNames, dsName)
	}
	return selectedNames, unselectedNames
}

//Lookup return a list of *Datasource that correspond to the
// list of tag expression
func (dss *Datasources) Lookup(log *logrus.Entry, tagList []string, limitedTags []string, dsTypes []Type, engines []Engine) []Datasourcer {
	logLookup := log.WithField("lookup", "tags")
	var selected []string
	unselected := make([]string, 0)
	limited := make([]string, 0)
	if limitedTags != nil {
		limited, _ = dss.getDsNameFromTagList(logLookup, limitedTags, dsTypes, engines)
	}
	if len(tagList) == 0 {
		logLookup.Debug("No tag provided, will only lookup on type and engines")
		// The selection is not based on tag, lookup for all of them
		selected = dss.lookupWithoutTag(dsTypes, engines)
	} else {
		selected, unselected = dss.getDsNameFromTagList(logLookup, tagList, dsTypes, engines)
	}

	logLookup.Debug("Final datasources list:")
	finalDsList := make([]string, 0, len(selected))

	selectedDs := make([]Datasourcer, 0, len(selected))
	for _, dsName := range selected {
		found := false
		for _, name := range unselected {
			if name == dsName {
				found = true
				break
			}
		}
		if !found {
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
			if inLimit {
				finalDsList = append(finalDsList, dsName)
				selectedDs = append(selectedDs, dss.datasources[dsName])
			}
		}
	}
	logLookup.Debug(strings.Join(finalDsList, ","))
	return selectedDs
}
