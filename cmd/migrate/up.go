package migrate

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newMigrateUpCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "up [--admin|--user] <recipe> ... <recipe>",
		Short: "Apply up migration",
		Long:  `Apply the up block of schema migration.`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Up(args)
		},
	}

	return cmd
}

//Up will implement the migration process in up direction
func Up(args []string) error {
	fmt.Println("TODO: To be implemented")
	return nil
}
