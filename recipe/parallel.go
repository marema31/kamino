package recipe

import (
	"context"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/step/common"
)

//watchdog will run has a parallel goroutine until the context is cancelled or the Do sub push to the end channel.
func stepWatchdog(ctxRecipe context.Context, log *logrus.Entry, wgWatchdog *sync.WaitGroup, end chan bool, step common.Steper) {
	defer close(end)
	defer wgWatchdog.Done()

	select {
	case <-ctxRecipe.Done(): // the context has been cancelled
		log.Debug("Watchdog cancel")
		//Step aborted, revert the action of the step if we can
		step.Cancel(log)
		//Even if have cancelled, we should wait for Do.recipe to ask us to quit
		<-end
	case <-end: // All the step of this priority has ended without error, the Do.recipe ask us to quit
		//		log.Debug("Watchdog finish")
		step.Finish(log)
	}
}

//stepExecutor will execute the step, cancel the context and raise the global flag.
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

func (ck *Cookbook) doParallelOneRecipe(ctx context.Context, cancelRecipe context.CancelFunc, log *logrus.Entry, stepsToBeDone []common.Steper) {
	nbSteps := len(stepsToBeDone)

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

		go stepExecutor(ctx, log, step, &wgStep, &wgWatchdog, cancelRecipe, end, recipeHadError)
	}

	//Wait for all step/stepExecutor to finish
	log.Debug("Waiting for steps of this priority")
	wgStep.Wait()
	log.Debug("All steps of this priority ended")
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
			log.Error("One step of this priority had error, skipping the following steps")

			return //We won't execute the following priorities
		}
	}
}
