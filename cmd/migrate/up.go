package migrate

import (
	"github.com/marema31/kamino/cmd/common"
	"github.com/marema31/kamino/recipe"
	"github.com/marema31/kamino/step"
	"github.com/spf13/cobra"
)

func newMigrateUpCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "up <recipe> ... <recipe>",
		Short: "Apply up migration",
		Long: `Apply the up block of schema migration.
		
		All non-applied Admin's migrations will be done first, then user's ones`,
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cookbook := recipe.New(&step.Factory{}, common.Timeout, common.Retry, true, common.Sequential, false, common.DryRun)
			return UpDown("up", cookbook, args)
		},
	}

	return cmd
}
