package cmd

import (
	"fmt"

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
			cookbook := recipe.New(step.Factory{})
			return Apply(cookbook, cfgFolder, names, types, args)
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

//Apply will run only the recipes with Apply type
func Apply(cookbook recipe.Cooker, cfgFolder string, names []string, types []string, args []string) error {
	err := cookbook.Load(ctx, cfgFolder, args, names, types)
	if err != nil {
		return fmt.Errorf("error while loading the recipes: %v", err)
	}
	if cookbook.Do(ctx) {
		return fmt.Errorf("a step had an error")
	}
	return nil
}
