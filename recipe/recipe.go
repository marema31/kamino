/*Package recipe provides two objects and manage their workflow:
   - Recipe : set of steps selected by the user (by tags/step type or step name on CLI)
   - Cookbook: set of recipes selected by the user (by name of recipe on the CLI)

The folder of a recipe must contains at least two folders:
	- steps containing the step files of the recipe
	- datasources containing the datasource files to be used by the recipe
but it can also contain whatever files/folders needed for the steps (initial dataset, templates, script, etc...).
All the relative path defined in a step will be relative to the recipe folder.
*/
package recipe

import (
	"context"
	"sort"

	"github.com/marema31/kamino/step"
	"github.com/marema31/kamino/step/common"
)

// Cooker interface for cookbook testing
type Cooker interface {
	Statistics() (map[string][]int, int)
	Load(context.Context, string, []string, []string, []string) error
	Do(context.Context) bool
}

type recipe struct {
	name            string
	steps           map[uint][]common.Steper
	currentPriority uint
}

// Cookbook is a map of recipe indexed by recipe's name
type Cookbook struct {
	Recipes     map[string]recipe
	stepFactory step.Creater
}

// New returns a new
func New(sf step.Creater) *Cookbook {
	rs := make(map[string]recipe)
	return &Cookbook{
		Recipes:     rs,
		stepFactory: sf,
	}
}

// Statistics return number of step by priority by recipes and total number of steps
func (ck *Cookbook) Statistics() (map[string][]int, int) {
	result := make(map[string][]int)
	var total int
	for rname := range ck.Recipes {
		s := make([]int, 0, len(ck.Recipes[rname].steps))
		priorities := make([]int, 0, len(ck.Recipes[rname].steps))
		for priority := range ck.Recipes[rname].steps {
			priorities = append(priorities, int(priority))
		}
		sort.Ints(priorities)
		for _, priority := range priorities {
			nb := len(ck.Recipes[rname].steps[uint(priority)])
			s = append(s, nb)
			total += nb
		}
		result[rname] = s
	}
	return result, total
}
