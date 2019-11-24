package cmd

import (
	"fmt"

	"github.com/marema31/kamino/cmd/common"
	"github.com/marema31/kamino/recipe"
	"github.com/marema31/kamino/step"
	"github.com/spf13/cobra"
)

var (
	syncCmd = &cobra.Command{
		Use:     "synchronize [flags] <recipe> ... <recipe>:<table>,...,<table>",
		Short:   "Synchronize will run only the recipes (or the table of recipe) with sync type",
		Long:    ``,
		Aliases: []string{"sync"},
		RunE: func(_ *cobra.Command, args []string) error {
			cookbook := recipe.New(&step.Factory{}, common.Timeout, common.Retry, common.Force, common.Sequential)
			return Sync(cookbook, names, args)
		},
	}
	cacheOnly bool
)

func init() {
	syncCmd.Flags().BoolVarP(&cacheOnly, "cache-only", "C", false, "Use only cache as source")
	rootCmd.AddCommand(syncCmd)
}

//Sync will run only the recipes with sync type
func Sync(cookbook recipe.Cooker, names []string, args []string) error {
	fmt.Println("TODO: To be implemented")
	log := common.Logger.WithField("action", "sync")

	recipes, err := common.FindRecipes(args)
	if err != nil {
		return err
	}

	err = cookbook.Load(common.Ctx, log, common.CfgFolder, recipes, names, []string{"sync"})
	if err != nil {
		return fmt.Errorf("error while loading the recipes: %v", err)
	}

	if cookbook.Do(common.Ctx, log) {
		return fmt.Errorf("a step had an error")
	}
	return nil
}
