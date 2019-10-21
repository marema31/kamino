package sync

import (
	"context"
	"fmt"

	"github.com/marema31/kamino/filter"
	"github.com/spf13/viper"
)

func getFilters(ctx context.Context, v *viper.Viper, sync string) ([]filter.Filter, error) {
	filters := make([]filter.Filter, 0)
	var ok bool
	fs := v.Get("filters")
	if fs != nil { // There is a filter section

		// filters section is an json array of map
		switch casted := fs.(type) { // Avoid panic if the type is not compatible with the one we want
		case []interface{}:

			// for every element of the filters array that must be a map
			for _, f := range casted {
				var filterType string
				var AParam []string
				var MParam map[string]string

				switch fcasted := f.(type) { // Avoid panic if the type is not compatible with the one we want
				case map[string]interface{}:
					filterType, ok = fcasted["type"].(string)
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
						AParam = make([]string, 0)
						for _, pv := range pcasted {
							AParam = append(AParam, pv.(string))
						}

					case map[string]interface{}:
						MParam = make(map[string]string)
						for pk, pv := range pcasted {
							MParam[pk] = pv.(string)
						}
					}

					f, err := filter.NewFilter(ctx, filterType, AParam, MParam)
					if err != nil {
						return nil, err
					}
					filters = append(filters, f)

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
