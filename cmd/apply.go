package cmd

import (
	"fmt"

	"github.com/marema31/kamino/cmd/common"
	"github.com/marema31/kamino/recipe"
	"github.com/marema31/kamino/step"
	"github.com/spf13/cobra"
)

var (
	applyCmd = &cobra.Command{
		Use:                   "apply [flags] <recipe> ... <recipe>",
		Short:                 "Apply will run the recipes provided in arguments",
		Long:                  ``,
		DisableFlagsInUseLine: true,
		RunE: func(_ *cobra.Command, args []string) error {
			cookbook := recipe.New(&step.Factory{}, common.Timeout, common.Retry, common.Force, common.Sequential, false, common.DryRun)
			return Apply(cookbook, names, types, args)
		},
	}
	names = []string{}
	types = []string{}
)

func init() {
	applyCmd.Flags().StringSliceVarP(&names, "name", "n", []string{}, "comma separated list of recipe step names")
	applyCmd.Flags().StringSliceVarP(&types, "type", "t", []string{}, "comma separated list of recipe step types")
	rootCmd.AddCommand(applyCmd)
}

//Apply will run only the recipes with Apply type.
func Apply(cookbook recipe.Cooker, names []string, types []string, args []string) error {
	log := common.Logger.WithField("action", "apply")

	recipes, err := common.FindRecipes(log, args)
	if err != nil {
		return err
	}

	err = cookbook.Load(common.Ctx, log, common.CfgFolder, recipes, common.Tags, names, types)
	if err != nil {
		return fmt.Errorf("error while loading the recipes: %w", err)
	}

	err = cookbook.PostLoad(log, common.CreateSuperseed())
	if err != nil {
		return fmt.Errorf("error while postloading the recipes: %w", err)
	}

	if cookbook.Do(common.Ctx, log) {
		return fmt.Errorf("a step had an error: %w", common.ErrStep)
	}

	return nil
}
