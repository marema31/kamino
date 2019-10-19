package migrate

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var admin, user bool
var ctx context.Context

//AddCommands adds all subcommands to RootCmd
func AddCommands(c context.Context, cmd *cobra.Command) {
	ctx = c //Save context only to make it accessible to sub-command
	cmd.AddCommand(
		NewMigrateCommand(),
	)
	//TODO: remove this line when ctx will be used elsewhere
	fmt.Print(ctx)
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
