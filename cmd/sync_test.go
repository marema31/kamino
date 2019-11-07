package cmd_test

import (
	"testing"

	"github.com/marema31/kamino/cmd"
)

func TestSyncOk(t *testing.T) {
	err := cmd.Sync([]string{"recipe1ok", "recipe2ok"})
	if err != nil {
		t.Errorf("Load should not returns an error, returned: %v", err)
	}
}
