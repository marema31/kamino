package recipe

import (
	"context"
	"sort"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/step/common"
)

var (
	mu       sync.Mutex
	hadError bool
)

/*Do will start one parallel recipe executor by recipe
Each recipe executor will run all the steps of the recipes by priorities.
All the step of same priority level will be parellelized or sequentialized (defined by the Cookbook.sequential flags),
and the executor will wait for all them before starting the next batch of step.
If an error occurs in one of the steps or user CTRL+C , all the same priority level steps will
receive an cancelation that they could use to rollback by example and all the step
with a priority level not already launched will not be runned.
*/
func (ck *Cookbook) Do(ctx context.Context, log *logrus.Entry) bool {
	// Waitgroup for the recipes
	var wgRecipe sync.WaitGroup

	hadError = false

	for rname := range ck.Recipes {
		wgRecipe.Add(1)

		logRecipe := log.WithField("recipe", rname)

		defer ck.Recipes[rname].dss.CloseAll(logRecipe)

		go ck.doOneRecipe(ctx, logRecipe, &wgRecipe, rname)
	}

	wgRecipe.Wait()
	//We wont treat the event before all recipe goroutine are finished
	select {
	case <-ctx.Done(): // the context has been cancelled before, since all the goroutine has also be notified
		//  via context inheritance, we can afford to take this event in account after their termination (via the wgRecipe.Wait)
		return true // we want the shell return code to be not ok

	default: // Make the poll to ctx.Done() non blocking. Do nothing
	}

	return hadError
}

func (ck *Cookbook) doOneRecipe(ctx context.Context, log *logrus.Entry, wgRecipe *sync.WaitGroup, rname string) {
	defer wgRecipe.Done()

	log.Debug("Executing recipe")
	// create a new context specific to this recipe.
	// ctx cancellation will propagate to it,
	// but its own cancel will stop to this context
	ctxRecipe, cancelRecipe := context.WithCancel(ctx)
	defer cancelRecipe()

	// Create an ordered list of priorities of the recipe
	priorities := make([]int, 0, len(ck.Recipes[rname].steps))
	for priority := range ck.Recipes[rname].steps {
		priorities = append(priorities, int(priority))
	}

	sort.Ints(priorities)

	for _, priority := range priorities {
		logPriority := log.WithField("priority", priority)
		logPriority.Debug("Determining step of this priority")

		stepsToBeDone := make([]common.Steper, 0, len(ck.Recipes[rname].steps[uint(priority)]))
		if ck.force {
			stepsToBeDone = append(stepsToBeDone, ck.Recipes[rname].steps[uint(priority)]...)
			logPriority.Debugf("Force mode, will do all the %d steps of this priority", cap(stepsToBeDone))
		} else {
			for _, step := range ck.Recipes[rname].steps[uint(priority)] {
				yes, err := step.ToSkip(ctx, logPriority)
				if err != nil {
					logPriority.Error("Can not determine if the step a step can be skipped")
					mu.Lock()
					{
						hadError = true
					}
					mu.Unlock()
					return
				}
				if !yes {
					stepsToBeDone = append(stepsToBeDone, step)
				}
			}
		}

		nbSteps := len(stepsToBeDone)
		logPriority.Infof("Will skip %d steps of the %d of this priority", cap(stepsToBeDone)-nbSteps, cap(stepsToBeDone))

		for _, step := range stepsToBeDone {
			err := step.Init(ctxRecipe, logPriority)
			if err != nil {
				//we set the flag for the cookbook, does not execute following priorities for this recipe
				mu.Lock()
				{
					hadError = true
				}

				mu.Unlock()
				logPriority.Error("One step of this priority had error at initialization, skipping the following steps")

				return //We won't execute the following priorities
			}
		}

		logPriority.Debug("Executing step of this priority")

		if ck.sequential || ck.forcedSequential[uint(priority)] {
			ck.doSequentialOneRecipe(ctxRecipe, logPriority, stepsToBeDone)
		} else {
			ck.doParallelOneRecipe(ctxRecipe, cancelRecipe, logPriority, stepsToBeDone)
		}

		if hadError {
			return
		}
	}

	log.Debug("Recipe ended without error")
}
