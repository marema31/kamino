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
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/step"
	"github.com/marema31/kamino/step/common"
)

// Cooker interface for cookbook testing.
type Cooker interface {
	Statistics() (map[string][]int, int)
	Load(context.Context, *logrus.Entry, string, []string, []string, []string, []string) error
	PostLoad(*logrus.Entry, map[string]string) error
	Do(context.Context, *logrus.Entry) bool
}

type recipe struct {
	name            string
	steps           map[uint][]common.Steper
	currentPriority uint
	dss             datasource.Datasourcers
}

// Cookbook is a map of recipe indexed by recipe's name.
type Cookbook struct {
	Recipes          map[string]recipe
	forcedSequential map[uint]bool
	stepFactory      step.Creater
	conTimeout       time.Duration
	conRetry         int
	force            bool
	sequential       bool
	validate         bool
	dryRun           bool
}

// New returns a new.
func New(sf step.Creater, connectionTimeout time.Duration, connectionRetry int, force bool, sequential bool, validate bool, dryRun bool) *Cookbook {
	return &Cookbook{
		Recipes:          make(map[string]recipe),
		forcedSequential: make(map[uint]bool),
		stepFactory:      sf,
		conTimeout:       connectionTimeout,
		conRetry:         connectionRetry,
		force:            force,
		sequential:       sequential,
		validate:         validate,
		dryRun:           dryRun,
	}
}

// Statistics return number of step by priority by recipes and total number of steps.
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
