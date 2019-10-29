package recipe

import (
	"context"
	"sort"
)

func (ck *Cookbook) doOneRecipe(ctx context.Context, rname string) error {
	priorities := make([]int, 0, len(ck.Recipes[rname].steps))
	for priority := range ck.Recipes[rname].steps {
		priorities = append(priorities, int(priority))
	}
	sort.Ints(priorities)
	for _, priority := range priorities {
		//TODO: create waitgroup
		for _, step := range ck.Recipes[rname].steps[uint(priority)] {
			//TODO: call step.Do(ctx) in a go routine
			err := step.Do(ctx)
			if err != nil {
				return err
			}

		}
		//Wait on the waitgroup
	}
	return nil
}

// Do run all the steps of the recipes by priorities
func (ck *Cookbook) Do(ctx context.Context) error {
	//TODO: create waitgroup
	for rname := range ck.Recipes {
		//TODO: call in goroutine
		err := ck.doOneRecipe(ctx, rname)
		if err != nil {
			return err
		}
	}
	//TODO: wait on waitgroup
	return nil
}
