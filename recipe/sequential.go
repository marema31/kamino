package recipe

import (
	"context"
	"sort"

	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/step/common"
)

func (ck *Cookbook) doSequentialOneRecipe(ctx context.Context, log *logrus.Entry, rname string) {
	// Create an ordered list of priorities of the recipe
	priorities := make([]int, 0, len(ck.Recipes[rname].steps))
	for priority := range ck.Recipes[rname].steps {
		priorities = append(priorities, int(priority))
	}
	sort.Ints(priorities)

	for _, priority := range priorities {
		log.Debugf("Executing step of priority: %d", priority)
		stepsToBeDone := make([]common.Steper, 0, len(ck.Recipes[rname].steps[uint(priority)]))
		if ck.force {
			stepsToBeDone = append(stepsToBeDone, ck.Recipes[rname].steps[uint(priority)]...)
			log.Debugf("Force mode, will do all the %d steps of this priority", cap(stepsToBeDone))
		} else {
			for _, step := range ck.Recipes[rname].steps[uint(priority)] {
				yes, err := step.ToSkip(ctx, log)
				if err != nil {
					log.Error("Can not determine if the step a step can be skipped")
					hadError = true
					return
				}
				if !yes {
					stepsToBeDone = append(stepsToBeDone, step)
				}
			}
			nbSteps := len(stepsToBeDone)
			log.Debugf("Will skip %d steps of the %d of this priority", cap(stepsToBeDone)-nbSteps, cap(stepsToBeDone))
		}
		for _, step := range stepsToBeDone {
			err := step.Init(ctx, log)
			if err != nil {
				//we set the flag for the cookbook, does not execute following priorities for this recipe
				hadError = true
				log.Errorf("One step of priority %d had error at initialization, skipping the following steps", priority)
				return //We won't execute the following priorities
			}
		}

		for i, step := range stepsToBeDone {
			err := step.Do(ctx, log)
			//Look for cancellation between each steps
			select {
			case <-ctx.Done(): // the context has been cancelled before, since all the goroutine has also be notified
				//  via context inheritance, we can afford to take this event in account after their termination (via the wgRecipe.Wait)
				hadError = true
				for _, step := range stepsToBeDone[:i+1] {
					step.Cancel(log)
				}
				log.Debug("Recipe execution cancelled")
				return

			default: // Make the poll to ctx.Done() non blocking
				// Do nothing
			}
			if err != nil {
				hadError = true
				for _, step := range stepsToBeDone[:i+1] {
					step.Cancel(log)
				}
				log.Debug("Recipe execution failed")
				return
			}
		}

		for _, step := range stepsToBeDone {
			step.Finish(log)
		}
	}
	log.Debug("Recipe ended without error")
}

/* sequentialDo will start the recipe sequentially, each step will be run sequentially too.
If an error occurs in one of the steps or user CTRL+C , all the same priority level steps will
receive an cancelation that they could use to rollback by example and all the step
with a priority level not already launched will not be runned.
*/
func (ck *Cookbook) sequentialDo(ctx context.Context, log *logrus.Entry) bool {
	hadError = false
	for rname := range ck.Recipes {
		logRecipe := log.WithField("recipe", rname)
		logRecipe.Debug("Executing recipe")
		ck.doSequentialOneRecipe(ctx, log, rname)
	}
	return hadError
}
