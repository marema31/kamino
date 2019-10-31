package recipe

import (
	"context"
	"sort"
	"sync"

	"github.com/marema31/kamino/step/common"
)

var mu sync.Mutex
var hadError bool

//watchdog will run has a parallel goroutine until the context is cancelled or the Do sub push to the end channel
func stepWatchdog(ctxRecipe context.Context, wgWatchdog *sync.WaitGroup, end chan bool, step common.Steper) {
	defer close(end)
	defer wgWatchdog.Done()

	select {
	case <-ctxRecipe.Done(): // the context has been cancelled
		//Step aborted, revert the action of the step if we can
		step.Cancel()
		//Even if have cancelled, we should wait for Do.recipe to ask us to quit
		<-end
	case <-end: // the channel has a information we stop this goroutine
	}
}

//stepExecutor will execute the step, cancel the context and raise the global fi
func stepExecutor(ctxRecipe context.Context, step common.Steper, wgStep *sync.WaitGroup, wgWatchdog *sync.WaitGroup, cancelRecipe func(), end chan bool, hadError chan bool) {
	defer wgStep.Done()

	wgWatchdog.Add(1)
	go stepWatchdog(ctxRecipe, wgWatchdog, end, step)

	err := step.Do(ctxRecipe)
	if err != nil {
		hadError <- true
		cancelRecipe()
	} else {
		hadError <- false
	}

}
func (ck *Cookbook) doOneRecipe(ctx context.Context, wgRecipe *sync.WaitGroup, rname string) {
	defer wgRecipe.Done()

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
		nbSteps := len(ck.Recipes[rname].steps[uint(priority)])

		//Waitgroup for steps of this priority level
		var wgStep sync.WaitGroup

		//Waitgroup for watchdog of this priority level
		var wgWatchdog sync.WaitGroup

		//Channel to allow stepExecutor to inform if the step finished in error
		recipeHadError := make(chan bool, nbSteps)
		defer close(recipeHadError)
		//List of end channel for this priority level
		ends := make([]chan bool, 0, nbSteps)
		for _, step := range ck.Recipes[rname].steps[uint(priority)] {
			//Prepare waitgroup and end channel for this step
			wgStep.Add(1)
			end := make(chan bool, 1)
			ends = append(ends, end)

			go stepExecutor(ctxRecipe, step, &wgStep, &wgWatchdog, cancelRecipe, end, recipeHadError)
		}

		//Wait for all step/stepExecutor to finish
		wgStep.Wait()
		// All the step of this priority has finished (completed or cancelled), it's time to stop the watchdogs
		for _, end := range ends {
			end <- true
		}

		//Wait for all the watchdog to finish
		wgWatchdog.Wait()

		//Since there is one end channel by stepExecutor
		//Each stepExecutor will send only one boolean to the recipeHadError channel
		for i := 0; i < nbSteps; i++ {
			wasNotOk := <-recipeHadError
			if wasNotOk {
				// One step of this priority finished in error and stepExecutor noticed us
				//we set the flag for the cookbook, does not execute following priorities for this recipe
				mu.Lock()
				{
					hadError = true
				}
				mu.Unlock()
				return //We won't execute the following priorities
			}
		}
	}
}

/*Do will start one parallel recipe executor by recipe
Each recipe executor will run all the steps of the recipes by priorities.
All the step of same priority level will be parellelized and the executor
will wait for all them before starting the next batch of step.

If an error occurs in one of the steps or user CTRL+C , all the same priority level steps will
receive an cancelation that they could use to rollback by example and all the step
with a priority level not already launched will not be runned.
*/
func (ck *Cookbook) Do(ctx context.Context) bool {
	// Waitgroup for the recipes
	var wgRecipe sync.WaitGroup

	for rname := range ck.Recipes {
		wgRecipe.Add(1)
		go ck.doOneRecipe(ctx, &wgRecipe, rname)
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
