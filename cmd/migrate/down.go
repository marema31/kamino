package migrate

import (
	"fmt"

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
			cookbook := recipe.New(&step.Factory{}, common.Timeout, common.Retry, true, common.Sequential)
			return Down(cookbook, args)
		},
	}

	return cmd
}

//Down will implement the migration process in down direction
func Down(cookbook recipe.Cooker, args []string) error {
	log := common.Logger.WithField("action", "migrate-down")

	superseed, err := createSuperseed()
	if err != nil {
		return err
	}
	superseed["migration.dir"] = "down"

	recipes, err := common.FindRecipes(args)
	if err != nil {
		return err
	}
	err = cookbook.Load(common.Ctx, log, common.CfgFolder, recipes, nil, []string{"migration"})
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
