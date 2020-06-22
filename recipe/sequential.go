package recipe

import (
	"context"

	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/step/common"
)

func (ck *Cookbook) doSequentialOneRecipe(ctx context.Context, log *logrus.Entry, stepsToBeDone []common.Steper) {
	for i, step := range stepsToBeDone {
		err := step.Do(ctx, log)
		if err != nil {
			hadError = true

			for _, step := range stepsToBeDone[:i+1] {
				step.Cancel(log)
			}

			log.Debug("Recipe execution failed")

			return
		}
		//Look for cancellation between each steps
		select {
		case <-ctx.Done(): // If the context has been cancelled stop the recipe execution here
			hadError = true

			for _, step := range stepsToBeDone[:i+1] {
				step.Cancel(log)
			}

			log.Debug("Recipe execution cancelled")

			return

		default: // Make the poll to ctx.Done() non blocking. Do nothing
			continue
		}
	}

	for _, step := range stepsToBeDone {
		step.Finish(log)
	}
}
