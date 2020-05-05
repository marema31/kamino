package recipe

import (
	"context"

	"github.com/Sirupsen/logrus"
)

var hadError bool

// Do will execute the step either in parallel or sequentially (defined by the Cookbook.sequential flags).
func (ck *Cookbook) Do(ctx context.Context, log *logrus.Entry) bool {
	if ck.sequential {
		return ck.sequentialDo(ctx, log)
	}

	return ck.parallelDo(ctx, log)
}
