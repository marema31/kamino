package datasource

import "strings"

// Datasource tag dictionnary for lookup
var tagToDatasource = make(map[string][]string)

// Insert the datasource name in all entry of the dictionnary
// that correspond to one tag of the tag list
func insertTag(tagList []string, name string) {
	for _, tag := range tagList {
		if _, ok := tagToDatasource[tag]; ok {
			tagToDatasource[tag] = append(tagToDatasource[tag], name)
		} else {
			dl := make([]string, 0, 1)
			dl = append(dl, name)
			tagToDatasource[tag] = dl
		}
	}
}

func isSelectedEngine(ds *Datasource, engines []Engine) bool {
	if engines == nil || len(engines) == 0 {
		return true
	}

	for _, e := range engines {
		if e == ds.Engine {
			return true
		}
	}
	return false
}

func isSelectedType(ds *Datasource, dsTypes []Type) bool {
	if dsTypes == nil || len(dsTypes) == 0 {
		return true
	}

	for _, t := range dsTypes {
		if t == ds.Type {
			return true
		}
	}
	return false
}

func lookupOneTag(tag string, dsTypes []Type, engines []Engine) (selected []string) {
	for _, name := range tagToDatasource[tag] {
		ds := datasources[name]
		if isSelectedEngine(ds, engines) && isSelectedType(ds, dsTypes) {
			selected = append(selected, name)
		}
	}
	return selected
}

func lookupWithoutTag(dsTypes []Type, engines []Engine) (selected []string) {
	for _, names := range tagToDatasource {
		for _, name := range names {
			ds := datasources[name]
			if isSelectedEngine(ds, engines) && isSelectedType(ds, dsTypes) {
				selected = append(selected, name)
			}
		}
	}
	return selected
}

//Lookup return a list of *Datasource that correspond to the
// list of tag expression
func Lookup(tagList []string, dsTypes []Type, engines []Engine) []*Datasource {
	//TODO: implement tagList empty
	selected := make(map[string]*Datasource) // Use map to emulate a "set" to avoid duplicates
	if len(tagList) == 0 {
		// The selection is not based on tag, lookup for all of them
		for _, name := range lookupWithoutTag(dsTypes, engines) {
			selected[name] = datasources[name]
		}
	} else {
		for _, tagElement := range tagList {
			if strings.Index(tagElement, ".") == -1 {
				// Simple tag, all corresponding datasource are selected

				for _, name := range lookupOneTag(tagElement, dsTypes, engines) {
					selected[name] = datasources[name]
				}
			} else {
				// Composite tag, only datasource corresponding to all tags are selected
				candidates := make(map[string]bool)
				for i, tag := range strings.Split(tagElement, ".") {
					for _, name := range lookupOneTag(tag, dsTypes, engines) {
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
					selected[name] = datasources[name]
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
