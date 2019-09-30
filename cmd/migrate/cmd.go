package migrate

import (
	"github.com/spf13/cobra"
)

var admin, user bool

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

	cmd.PersistentFlags().BoolVarP(&admin, "admin", "a", false, "Only admin migration (if relevant)")
	cmd.PersistentFlags().BoolVarP(&user, "user", "u", false, "Only user migration")
	return cmd
}
