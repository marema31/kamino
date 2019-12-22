package migrate

import (
	"github.com/marema31/kamino/cmd/common"
	"github.com/marema31/kamino/recipe"
	"github.com/marema31/kamino/step"
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
			cookbook := recipe.New(&step.Factory{}, common.Timeout, common.Retry, true, common.Sequential, false, common.DryRun)
			return UpDown("down", cookbook, args)
		},
	}

	return cmd
}
