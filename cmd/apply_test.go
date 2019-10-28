package cmd_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/marema31/kamino/cmd"
)

type mockedCookbook struct {
	called    bool
	errorLoad error
	errorDo   error
}

//Do manage the runnning of the cookbook
func (ck *mockedCookbook) Do(ctx context.Context) error {
	ck.called = true
	return ck.errorDo
}

// Load the step file and returns the priority and a list of steper for this file
func (ck *mockedCookbook) Load(ctx context.Context, path string, recipes []string, stepNames []string, stepTypes []string) error {
	ck.called = false
	return ck.errorLoad
}

// Statistics return statistics on the cookbook
func (ck *mockedCookbook) Statistics() (map[string][]int, int) {
	result := make(map[string][]int)
	var total int
	return result, total

}

func TestApplyOk(t *testing.T) {
	ck := &mockedCookbook{}
	err := cmd.Apply(ck, "testdata/good", []string{"recipe1ok", "recipe2ok"}, nil, nil)
	if err != nil {
		t.Errorf("Load should not returns an error, returned: %v", err)
	}

	if !ck.called {
		t.Errorf("Do should be called")
	}
}

func TestApplyLoadError(t *testing.T) {
	ck := &mockedCookbook{}
	ck.errorLoad = fmt.Errorf("fake error")
	err := cmd.Apply(ck, "testdata/good", []string{"recipe1ok", "recipe2ok"}, nil, nil)
	if err == nil {
		t.Errorf("Load should returns an error")
	}

	if ck.called {
		t.Errorf("Do should not be called")
	}
}

func TestApplyDoError(t *testing.T) {
	ck := &mockedCookbook{}
	ck.errorDo = fmt.Errorf("fake error")
	err := cmd.Apply(ck, "testdata/good", []string{"recipe1ok", "recipe2ok"}, nil, nil)
	if err == nil {
		t.Errorf("Load should returns an error")
	}

	if !ck.called {
		t.Errorf("Do should be called")
	}
}
