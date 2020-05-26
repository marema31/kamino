package migrate

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestAddCommandOk(t *testing.T) {
	command := &cobra.Command{
		Use:   "migrate",
		Short: "Manage schema migration",
		Args:  cobra.NoArgs,
		//		RunE: cmd.ShowHelp(),
	}
	AddCommands(command)
}

func TestNewMigrateCommandOk(t *testing.T) {
	cmd := NewMigrateCommand()
	if cmd == nil {
		t.Errorf("NewMigrateCommand should returns an object")
	}
}

type testEntries struct {
	name       string
	admin      bool
	user       bool
	limit      int
	length     int
	noAdmin    string
	noUser     string
	superLimit string
}

var superseedValues = []testEntries{
	{name: "NoSuperseed", length: 2, superLimit: "0"},
	{name: "NoUser", admin: true, length: 4, noAdmin: "false", noUser: "true", superLimit: "0"},
	{name: "NoAdmin", user: true, length: 4, noAdmin: "true", noUser: "false", superLimit: "0"},
	{name: "ForceLimit", limit: 3, length: 2, superLimit: "3"},
}

func TestCreateSuperseedOk(t *testing.T) {
	for _, tt := range superseedValues {
		t.Run(tt.name, func(t *testing.T) {
			Admin = tt.admin
			User = tt.user
			limit = tt.limit
			superseed, err := createSuperseed("up")
			if err != nil {
				t.Fatalf("createSuperseed should not returns an error, returned: %v", err)
			}
			if len(superseed) != tt.length {
				t.Errorf("createSuperseed should returns an %d elements, returned: %d, %v", tt.length, len(superseed), superseed)
			}
			if superseed["migration.noAdmin"] != tt.noAdmin {
				t.Errorf("createSuperseed should returns noAdmin=%s, returned: %s", tt.noAdmin, superseed["migration.noAdmin"])
			}
			if superseed["migration.noUser"] != tt.noUser {
				t.Errorf("createSuperseed should returns noUser=%s, returned: %s", tt.noAdmin, superseed["migration.noUser"])
			}
			if superseed["migration.limit"] != tt.superLimit {
				t.Errorf("createSuperseed should returns limit=%s, returned: %s", tt.superLimit, superseed["migration.limit"])
			}
		})
	}
}

func TestCreateSuperseedError(t *testing.T) {
	Admin = true
	User = true
	_, err := createSuperseed("up")
	if err == nil {
		t.Fatalf("createSuperseed should returns an error")
	}
}
