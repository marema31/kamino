package sync

import (
	"fmt"

	"github.com/marema31/kamino/mockprovider"
	"github.com/marema31/kamino/step/common"
)

func MockSourceContent(step common.Steper, content []map[string]string) error {
	//For test purpose we must see what is inside the step and for this convert the interface to the presumed type
	st, ok := step.(*Step)
	if !ok {
		return fmt.Errorf("The step should be a sync step")
	}

	s, ok := st.source.(*mockprovider.MockLoader)

	if !ok {
		return fmt.Errorf("The source should be a mockloader")
	}

	s.Content = content
	return nil
}

func MockSourceError(step common.Steper) error {
	//For test purpose we must see what is inside the step and for this convert the interface to the presumed type
	st, ok := step.(*Step)
	if !ok {
		return fmt.Errorf("The step should be a sync step")
	}

	s, ok := st.source.(*mockprovider.MockLoader)

	if !ok {
		return fmt.Errorf("The source should be a mockloader")
	}

	s.ErrorLoad = fmt.Errorf("fake error")
	return nil
}

func MockDestinationError(step common.Steper) error {
	//For test purpose we must see what is inside the step and for this convert the interface to the presumed type
	st, ok := step.(*Step)
	if !ok {
		return fmt.Errorf("The step should be a sync step")
	}

	s, ok := st.destinations[0].(*mockprovider.MockSaver)

	if !ok {
		return fmt.Errorf("The source should be a mockSaver")
	}

	s.ErrorSave = fmt.Errorf("fake error")
	return nil
}
