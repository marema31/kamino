package migrate

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newMigrateStatusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "status <recipe> ... <recipe>",
		Short:                 "Show migration status",
		Long:                  ``,
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Status(args)
		},
	}

	return cmd
}

//Status will show the migration status
func Status(args []string) error {
	fmt.Println("TODO: To be implemented")
	return nil
}
