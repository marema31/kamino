package cmd_test

import (
	"testing"

	"github.com/marema31/kamino/cmd"
)

func TestGetLoggerOk(t *testing.T) {
	cmd.InitConfig()
	cmd.GetLogger()
}

func TestVerboseOk(t *testing.T) {
	cmd.Verbose = true
	cmd.InitConfig()
	cmd.GetLogger()
}

func TestQuietOk(t *testing.T) {
	cmd.Quiet = true
	cmd.InitConfig()
	cmd.GetLogger()
}
