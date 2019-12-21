package cmd

import (
	"fmt"

	"github.com/marema31/kamino/cmd/common"
	"github.com/marema31/kamino/recipe"
	"github.com/marema31/kamino/step"
	"github.com/spf13/cobra"
)

var (
	validateCmd = &cobra.Command{
		Use:     "validate [flags] <recipe> ... <recipe>",
		Short:   "validate the datasources and recipe files",
		Long:    ``,
		Aliases: []string{"validate"},
		RunE: func(_ *cobra.Command, args []string) error {
			cookbook := recipe.New(&step.Factory{}, common.Timeout, common.Retry, common.Force, common.Sequential, false, common.DryRun)
			return Validate(cookbook, names, args)
		},
	}
)

func init() {
	rootCmd.AddCommand(validateCmd)
}

//Validate will validate the datasources and recipes files
func Validate(cookbook recipe.Cooker, names []string, args []string) error {
	log := common.Logger.WithField("action", "sync")

	recipes, err := common.FindRecipes(log, args)
	if err != nil {
		return err
	}

	err = cookbook.Load(common.Ctx, log, common.CfgFolder, recipes, common.Tags, names, []string{})
	if err != nil {
		return fmt.Errorf("error while loading the recipes: %v", err)
	}

	return nil
}
