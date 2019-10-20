package migration

import (
	"context"
	"fmt"
)

//Do manage the runnning of the step
func (st Step) Do(ctx context.Context) error {
	//TODO: to be implemented
	fmt.Printf("Will do shell for %s\n", st.Name)
	return nil
}
