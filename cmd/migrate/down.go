package migrate

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newMigrateDownCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "down <recipe> ... <recipe>",
		Short: "Apply down migration",
		Long: `Apply the down block of schema migration.
		
		The last applied user's migrations will be done first, if none, will applied last applied admin's ones`,
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Down(args)
		},
	}

	return cmd
}

//Down will implement the migration process in down direction
func Down(args []string) error {
	fmt.Println("TODO: To be implemented")
	return nil
}
