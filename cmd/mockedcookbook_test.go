package cmd_test

import (
	"context"

	"github.com/Sirupsen/logrus"
)

type mockedCookbook struct {
	called        bool
	errorLoad     error
	errorPostLoad error
	doReturnValue bool
}

//Do manage the runnning of the cookbook
func (ck *mockedCookbook) Do(ctx context.Context, log *logrus.Entry) bool {
	ck.called = true
	return ck.doReturnValue
}

// Load the step file and returns the priority and a list of steper for this file
func (ck *mockedCookbook) Load(ctx context.Context, log *logrus.Entry, path string, recipes []string, stepNames []string, stepTypes []string) error {
	ck.called = false
	return ck.errorLoad
}

//PostLoad modify the loaded step values with the values provided in the map in argument
func (ck *mockedCookbook) PostLoad(log *logrus.Entry, superseed map[string]string) error {
	return ck.errorPostLoad
}

// Statistics return statistics on the cookbook
func (ck *mockedCookbook) Statistics() (map[string][]int, int) {
	result := make(map[string][]int)
	var total int
	return result, total

}
