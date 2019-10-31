package sqlscript

import (
	"context"
	"fmt"
)

//Cancel manage the cancellation of the step
func (st *Step) Cancel() {
	fmt.Printf("Will cancel sqlscript for %s\n", st.Name)
}

//Do manage the runnning of the step
func (st *Step) Do(ctx context.Context) error {
	//TODO: to be implemented
	fmt.Printf("Will do sqlscript for %s\n", st.Name)
	return nil
}
