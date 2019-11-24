package cmd_test

import (
	"testing"

	"github.com/marema31/kamino/cmd"
	"github.com/marema31/kamino/cmd/common"
)

func TestGetLoggerOk(t *testing.T) {
	cmd.InitConfig()
	_ = cmd.GetLogger()
}

func TestVerboseOk(t *testing.T) {
	common.Verbose = true
	cmd.InitConfig()
	_ = cmd.GetLogger()
}

func TestQuietOk(t *testing.T) {
	common.Quiet = true
	cmd.InitConfig()
	_ = cmd.GetLogger()
}
