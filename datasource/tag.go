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

func (dss *Datasources) lookupWithoutTag(dsTypes []Type, engines []Engine) (selected []string) {
	for _, names := range dss.tagToDatasource {
		for _, name := range names {
			ds := dss.datasources[name]
			if ds.isSelectedEngine(engines) && ds.isSelectedType(dsTypes) {
				selected = append(selected, name)
			}
		}
	}
	return selected
}

//Lookup return a list of *Datasource that correspond to the
// list of tag expression
func (dss *Datasources) Lookup(log *logrus.Entry, tagList []string, dsTypes []Type, engines []Engine) []Datasourcer {
	logLookup := log.WithField("lookup", "tags")
	selected := make(map[string]*Datasource) // Use map to emulate a "set" to avoid duplicates
	unselected := make([]string, 0)
	if len(tagList) == 0 {
		logLookup.Debug("No tag provided, will only lookup on type and engines")
		// The selection is not based on tag, lookup for all of them
		for _, name := range dss.lookupWithoutTag(dsTypes, engines) {
			logLookup.Debugf("Found : %s", name)
			selected[name] = dss.datasources[name]
		}
	} else {
		for _, tagElement := range tagList {
			logLookup.Debugf("Lookup %s", tagElement)
			if !strings.Contains(tagElement, ".") {
				logLookup.Debug("Simple tag")
				if strings.HasPrefix(tagElement, "!") {
					logLookup.Debug("Negative tag")
					unselected = append(unselected, dss.lookupOneTag(tagElement, dsTypes, engines)...)
				} else {
					for _, name := range dss.lookupOneTag(tagElement, dsTypes, engines) {
						selected[name] = dss.datasources[name]
					}
				}
			} else {
				logLookup.Debug("Composite tag")
				candidates := make(map[string]bool)
				for i, tag := range strings.Split(tagElement, ".") {
					logLookup.Debugf("Lookup for sub-tag %s", tag)
					for _, name := range dss.lookupOneTag(tag, dsTypes, engines) {
						if i == 0 {
							logLookup.Debugf("Found: %s", name)
							candidates[name] = true
						} else {
							if _, ok := candidates[name]; ok {
								logLookup.Debugf("Found: %s, present from previous sub-tag", name)
								candidates[name] = true
							} else {
								logLookup.Debugf("Skipped: %s, not present for previous sub-tag", name)
							}
						}
					}
					logLookup.Debug("Removing datasource not found for this sub-tag")
					for name, viewed := range candidates {
						if !viewed {
							logLookup.Debugf("Removing: %s", name)
							delete(candidates, name)
						} else {
							// Prepare the map for the next tag
							candidates[name] = false
						}
					}
				}
				logLookup.Debugf("Final list for %s", tagElement)
				for name := range candidates {
					if strings.HasPrefix(tagElement, "!") {
						unselected = append(unselected, name)
					} else {
						logLookup.Debugf("  - %s", name)
						selected[name] = dss.datasources[name]
					}
				}
			}
		}
	}

	logLookup.Debug("Final datasources list:")
	finalDsList := make([]string, 0, len(selected))

	selectedDs := make([]Datasourcer, 0, len(selected))
	for _, ds := range selected {
		found := false
		for _, name := range unselected {
			if name == ds.name {
				found = true
				break
			}
		}
		if !found {
			finalDsList = append(finalDsList, ds.name)
			selectedDs = append(selectedDs, ds)
		}
	}
	logLookup.Debug(strings.Join(finalDsList, ","))
	return selectedDs
}
