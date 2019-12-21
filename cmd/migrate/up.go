package migrate

import (
	"fmt"

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
			return Up(cookbook, args)
		},
	}

	return cmd
}

//Up will implement the migration process in up direction
func Up(cookbook recipe.Cooker, args []string) error {
	log := common.Logger.WithField("action", "migrate-up")

	superseed, err := createSuperseed()
	if err != nil {
		return err
	}
	superseed["migration.dir"] = "up"

	recipes, err := common.FindRecipes(log, args)
	if err != nil {
		return err
	}
	err = cookbook.Load(common.Ctx, log, common.CfgFolder, recipes, common.Tags, nil, []string{"migration"})
	if err != nil {
		return fmt.Errorf("error while loading the recipes: %v", err)
	}

	err = cookbook.PostLoad(log, superseed)
	if err != nil {
		return fmt.Errorf("error while postloading the recipes: %v", err)
	}

	if cookbook.Do(common.Ctx, log) {
		return fmt.Errorf("a step had an error")
	}
	return nil
}
