package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	applyCmd = &cobra.Command{
		Use:                   "apply [flags] <recipe> ... <recipe>",
		Short:                 "Apply will run the recipes provided in arguments",
		Long:                  ``,
		DisableFlagsInUseLine: true,
		RunE: func(_ *cobra.Command, args []string) error {
			return Apply(args)
		},
	}
	names = []string{}
	types = []string{}
)

func init() {
	applyCmd.Flags().StringSliceVarP(&names, "name", "n", []string{}, "comma separated list of recipe step names")
	applyCmd.Flags().StringSliceVarP(&types, "type", "t", []string{}, "comma separated list of recipe step types")
	RootCmd.AddCommand(applyCmd)
}

//Apply will run only the recipes with Apply type
func Apply(args []string) error {
	fmt.Println("TODO: To be implemented")
	return nil
}
