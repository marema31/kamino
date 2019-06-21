package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func (c *Config) getFilters(v *viper.Viper, sync string) ([]FilterConfig, error) {
	var filters []FilterConfig
	var ok bool
	fs := v.Get("filters")
	if fs != nil { // There is a filter section

		// filters section is an json array of map
		switch casted := fs.(type) { // Avoid panic if the type is not compatible with the one we want
		case []interface{}:

			// for every element of the filters array that must be a map
			for _, f := range casted {
				var currentfilter FilterConfig
				switch fcasted := f.(type) { // Avoid panic if the type is not compatible with the one we want
				case map[string]interface{}:
					currentfilter.Type, ok = fcasted["type"].(string)
					if !ok {
						return nil, fmt.Errorf("missing type for a filter of %s sync", sync)
					}

					var parameters interface{}
					parameters, ok := fcasted["parameters"]
					if !ok {
						return nil, fmt.Errorf("missing parameters for a filter of %s sync", sync)
					}

					switch pcasted := parameters.(type) { // Avoid panic if the type is not compatible with the one we want
					case []interface{}:
						currentpvalue := make([]string, 0)
						for _, pv := range pcasted {
							currentpvalue = append(currentpvalue, pv.(string))
						}
						currentfilter.AParam = currentpvalue

					case map[string]interface{}:
						ps := make(map[string]string)
						for pk, pv := range pcasted {
							ps[pk] = pv.(string)
						}
						currentfilter.MParam = ps
					}
					filters = append(filters, currentfilter)
				default:
					return nil, fmt.Errorf("one filter defined for %s sync is invalid", sync)
				}
			}
		default:
			return nil, fmt.Errorf("filters defined for %s sync is invalid", sync)
		}
	}
	return filters, nil
}
