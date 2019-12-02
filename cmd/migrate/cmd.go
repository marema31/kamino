package migrate

import (
	"fmt"
	"strconv"

	"github.com/marema31/kamino/cmd/common"
	"github.com/spf13/cobra"
)

// Admin only migration
var Admin bool

// User only migration
var User bool
var limit int

//AddCommands adds all subcommands to RootCmd
func AddCommands(cmd *cobra.Command) {
	cmd.AddCommand(
		NewMigrateCommand(),
	)
}

//NewMigrateCommand declare the migration sub commands
func NewMigrateCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Manage schema migration",
		Args:  cobra.NoArgs,
		//		RunE: cmd.ShowHelp(),
	}

	cmd.AddCommand(
		newMigrateUpCommand(),
		newMigrateDownCommand(),
		newMigrateStatusCommand(),
	)

	cmd.PersistentFlags().BoolVarP(&Admin, "admin", "a", false, "Only admin migration (if relevant)")
	cmd.PersistentFlags().BoolVarP(&User, "user", "u", false, "Only user migration")
	cmd.PersistentFlags().IntVarP(&limit, "limit", "l", 0, "Max number of migration(0 for all)")
	return cmd
}

func createSuperseed() (map[string]string, error) {
	superseed := common.CreateSuperseed()
	if Admin && User {
		return superseed, fmt.Errorf("option --admin and --user are mutually exclusive")
	}
	if Admin {
		superseed["migration.noAdmin"] = "false"
		superseed["migration.noUser"] = "true"
	}
	if User {
		superseed["migration.noAdmin"] = "true"
		superseed["migration.noUser"] = "false"
	}
	if limit != 0 {
		superseed["migration.limit"] = strconv.Itoa(limit)
	}
	return superseed, nil
}
