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
	admin      bool
	user       bool
	limit      int
	length     int
	noAdmin    string
	noUser     string
	superLimit string
}

var superseedValues = []testEntries{
	{length: 0},
	{admin: true, length: 2, noAdmin: "false", noUser: "true"},
	{user: true, length: 2, noAdmin: "true", noUser: "false"},
	{limit: 3, length: 1, superLimit: "3"},
}

func TestCreateSuperseedOk(t *testing.T) {
	for _, entry := range superseedValues {
		Admin = entry.admin
		User = entry.user
		limit = entry.limit
		superseed, err := createSuperseed()
		if err != nil {
			t.Fatalf("createSuperseed should not returns an error, returned: %v", err)
		}
		if len(superseed) != entry.length {
			t.Errorf("createSuperseed should returns an %d elements, returned: %d", entry.length, len(superseed))
		}
		if superseed["migration.noAdmin"] != entry.noAdmin {
			t.Errorf("createSuperseed should returns noAdmin=%s, returned: %s", entry.noAdmin, superseed["migration.noAdmin"])
		}
		if superseed["migration.noUser"] != entry.noUser {
			t.Errorf("createSuperseed should returns noUser=%s, returned: %s", entry.noAdmin, superseed["migration.noUser"])
		}
		if superseed["migration.limit"] != entry.superLimit {
			t.Errorf("createSuperseed should returns limit=%s, returned: %s", entry.superLimit, superseed["migration.limit"])
		}
	}
}
func TestCreateSuperseedError(t *testing.T) {
	Admin = true
	User = true
	_, err := createSuperseed()
	if err == nil {
		t.Fatalf("createSuperseed should returns an error")
	}
}
