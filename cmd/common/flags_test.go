package common_test

import (
	"testing"

	"github.com/marema31/kamino/cmd/common"
)

func TestSuperseed(t *testing.T) {
	common.Force = false
	common.DryRun = false
	superseed := common.CreateSuperseed()
	if v, ok := superseed["kamino.force"]; ok {
		t.Errorf("Should no return force, returned %s", v)
	}
	if v, ok := superseed["kamino.dryrun"]; ok {
		t.Errorf("Should no return dryrun, returned %s", v)
	}

	common.Force = true
	common.DryRun = false
	superseed = common.CreateSuperseed()
	if _, ok := superseed["kamino.force"]; !ok {
		t.Errorf("Should return force")
	}
	if v, ok := superseed["kamino.force"]; ok && v != "true" {
		t.Errorf("Should return force = true, returned %s", v)
	}
	if v, ok := superseed["kamino.dryrun"]; ok {
		t.Errorf("Should no return dryrun, returned %s", v)
	}

	common.Force = false
	common.DryRun = true
	superseed = common.CreateSuperseed()
	if v, ok := superseed["kamino.force"]; ok {
		t.Errorf("Should no return force, returned %s", v)
	}
	if _, ok := superseed["kamino.dryrun"]; !ok {
		t.Errorf("Should return dryrun")
	}
	if v, ok := superseed["kamino.dryrun"]; ok && v != "true" {
		t.Errorf("Should return dryrun = true, returned %s", v)
	}

}
