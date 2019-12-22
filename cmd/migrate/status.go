package migrate

import (
	"fmt"

	"github.com/marema31/kamino/cmd/common"
	"github.com/marema31/kamino/recipe"
	"github.com/marema31/kamino/step"
	"github.com/spf13/cobra"
)

func newMigrateStatusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "status <recipe> ... <recipe>",
		Short:                 "Show migration status",
		Long:                  ``,
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cookbook := recipe.New(&step.Factory{}, common.Timeout, common.Retry, true, true, false, common.DryRun)
			return Status(cookbook, args)
		},
	}

	return cmd
}

//Status will show the migration status
func Status(cookbook recipe.Cooker, args []string) error {
	log := common.Logger.WithField("action", "migrate-status")

	superseed, err := createSuperseed()
	if err != nil {
		return err
	}

	superseed["migration.dir"] = "status"

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
