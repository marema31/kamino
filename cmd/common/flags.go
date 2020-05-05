package common

import (
	"context"
	"time"

	"github.com/Sirupsen/logrus"
)

var (
	// CfgFolder configuration folder
	CfgFolder string
	// Ctx root context of the application
	Ctx context.Context
	// DryRun would not really do the action but logs
	DryRun bool
	// Force the step execution without executing skip queries
	Force bool
	// Logger log engine of the application
	Logger = logrus.New()
	// Quiet no logs
	Quiet bool
	// Retry number of database ping
	Retry int
	// Sequential do not run step in parallel
	Sequential bool
	// Tags user provided tag list to limit the concerned datasource
	Tags []string
	// Timeout of each database ping try
	Timeout time.Duration
	// Verbose add debug logs
	Verbose bool
)

// CreateSuperseed creates the postload configuration map.
func CreateSuperseed() map[string]string {
	superseed := make(map[string]string)

	if Force {
		superseed["kamino.force"] = "true"
	}

	if DryRun {
		superseed["kamino.dryrun"] = "true"
	}

	return superseed
}
