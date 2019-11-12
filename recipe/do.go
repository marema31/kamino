package recipe

import (
	"context"
	"sort"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/step/common"
)

var mu sync.Mutex
var hadError bool

//watchdog will run has a parallel goroutine until the context is cancelled or the Do sub push to the end channel
func stepWatchdog(ctxRecipe context.Context, log *logrus.Entry, wgWatchdog *sync.WaitGroup, end chan bool, step common.Steper) {
	defer close(end)
	defer wgWatchdog.Done()
	log.Debug("Watchdog entering")

	select {
	case <-ctxRecipe.Done(): // the context has been cancelled
		log.Debug("Watchdog cancel")
		//Step aborted, revert the action of the step if we can
		step.Cancel(log)
		//Even if have cancelled, we should wait for Do.recipe to ask us to quit
		<-end
	case <-end: // All the step of this priority has ended without error, the Do.recipe ask us to quit
		log.Debug("Watchdog finish")
		step.Finish(log)
	}
	log.Debug("Watchdog ending")
}

//stepExecutor will execute the step, cancel the context and raise the global fi
func stepExecutor(ctxRecipe context.Context, log *logrus.Entry, step common.Steper, wgStep *sync.WaitGroup, wgWatchdog *sync.WaitGroup, cancelRecipe func(), end chan bool, hadError chan bool) {
	defer wgStep.Done()

	wgWatchdog.Add(1)
	go stepWatchdog(ctxRecipe, log, wgWatchdog, end, step)

	err := step.Do(ctxRecipe, log)
	if err != nil {
		cancelRecipe()
		hadError <- true
	} else {
		hadError <- false
	}
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
		log.Debugf("Executing step of priority: %d", priority)

		stepsToBeDone := make([]common.Steper, 0, len(ck.Recipes[rname].steps[uint(priority)]))
		for _, step := range ck.Recipes[rname].steps[uint(priority)] {
			yes, err := step.ToSkip(ctx, log)
			if err != nil {
				log.Error("Can not determine if the step a step can be skipped")
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
		nbSteps := len(stepsToBeDone)
		log.Debugf("Will skip %d steps of the %d of this priority", cap(stepsToBeDone)-nbSteps, cap(stepsToBeDone))
		for _, step := range stepsToBeDone {
			err := step.Init(ctxRecipe, log)
			if err != nil {
				//we set the flag for the cookbook, does not execute following priorities for this recipe
				mu.Lock()
				{
					hadError = true
				}
				mu.Unlock()
				log.Errorf("One step of priority %d had error at initialization, skipping the following steps", priority)
				return //We won't execute the following priorities
			}
		}

		//Waitgroup for steps do of this priority level
		var wgStep sync.WaitGroup

		//Waitgroup for watchdog of this priority level
		var wgWatchdog sync.WaitGroup

		//Channel to allow stepExecutor to inform if the step finished in error
		recipeHadError := make(chan bool, nbSteps)
		defer close(recipeHadError)
		//List of end channel for this priority level
		ends := make([]chan bool, 0, nbSteps)
		for _, step := range stepsToBeDone {
			//Prepare waitgroup and end channel for this step
			wgStep.Add(1)
			end := make(chan bool, 1)
			ends = append(ends, end)

			go stepExecutor(ctxRecipe, log, step, &wgStep, &wgWatchdog, cancelRecipe, end, recipeHadError)
		}

		//Wait for all step/stepExecutor to finish
		log.Debugf("Waiting for steps of priority: %d", priority)
		wgStep.Wait()
		log.Debugf("All steps of priority %d ended", priority)
		// All the step of this priority has finished (completed or cancelled), it's time to stop the watchdogs
		for _, end := range ends {
			end <- true
		}

		log.Debug("Waiting for watchdog")
		//Wait for all the watchdog to finish
		wgWatchdog.Wait()
		log.Debug("All watchdog ended")

		//Since there is one end channel by stepExecutor
		//Each stepExecutor will send only one boolean to the recipeHadError channel
		for i := 0; i < len(stepsToBeDone); i++ {
			wasNotOk := <-recipeHadError
			if wasNotOk {
				// One step of this priority finished in error and stepExecutor noticed us
				//we set the flag for the cookbook, does not execute following priorities for this recipe
				mu.Lock()
				{
					hadError = true
				}
				mu.Unlock()
				log.Errorf("One step of priority %d had error, skipping the following steps", priority)
				return //We won't execute the following priorities
			}
		}
	}
	log.Debug("Recipe ended without error")
}

/*Do will start one parallel recipe executor by recipe
Each recipe executor will run all the steps of the recipes by priorities.
All the step of same priority level will be parellelized and the executor
will wait for all them before starting the next batch of step.

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

	default: // Make the poll to ctx.Done() non blocking
		// Do nothing
	}
	return hadError
}
