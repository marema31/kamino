package datasource

import "strings"

func (ds *Datasource) isSelectedEngine(engines []Engine) bool {
	if len(engines) == 0 {
		return true
	}

	for _, e := range engines {
		if e == ds.Engine {
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
		if t == ds.Type {
			return true
		}
	}
	return false
}

func (dss Datasources) lookupOneTag(tag string, dsTypes []Type, engines []Engine) (selected []string) {
	//TODO: Implement "!tag" for all but this tag
	for _, name := range dss.tagToDatasource[tag] {
		ds := dss.datasources[name]
		if ds.isSelectedEngine(engines) && ds.isSelectedType(dsTypes) {
			selected = append(selected, name)
		}
	}
	return selected
}

func (dss Datasources) lookupWithoutTag(dsTypes []Type, engines []Engine) (selected []string) {
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
func (dss Datasources) Lookup(tagList []string, dsTypes []Type, engines []Engine) []*Datasource {
	selected := make(map[string]*Datasource) // Use map to emulate a "set" to avoid duplicates
	if len(tagList) == 0 {
		// The selection is not based on tag, lookup for all of them
		for _, name := range dss.lookupWithoutTag(dsTypes, engines) {
			selected[name] = dss.datasources[name]
		}
	} else {
		for _, tagElement := range tagList {
			if !strings.Contains(tagElement, ".") {
				// Simple tag, all corresponding datasource are selected

				for _, name := range dss.lookupOneTag(tagElement, dsTypes, engines) {
					selected[name] = dss.datasources[name]
				}
			} else {
				// Composite tag, only datasource corresponding to all tags are selected
				candidates := make(map[string]bool)
				for i, tag := range strings.Split(tagElement, ".") {
					for _, name := range dss.lookupOneTag(tag, dsTypes, engines) {
						if i == 0 {
							candidates[name] = true
						} else {
							if _, ok := candidates[name]; ok {
								candidates[name] = true
							}
						}
					}
					for name, viewed := range candidates {
						if !viewed {
							// This name as not be found for this tag part, remove it
							delete(candidates, name)
						} else {
							// Prepare the map for the next tag
							candidates[name] = false
						}
					}
				}
				for name := range candidates {
					selected[name] = dss.datasources[name]
				}
			}
		}
	}

	selectedDs := make([]*Datasource, 0, len(selected))
	for _, ds := range selected {
		selectedDs = append(selectedDs, ds)
	}
	return selectedDs
}
