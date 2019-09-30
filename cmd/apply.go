package cmd

import (
	"fmt"
	"strings"

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
	names []string
	types []string
)

func init() {
	var (
		sname string
		stype string
	)
	applyCmd.Flags().StringVarP(&sname, "name", "n", "", "comma separated list of recipe step names")
	names = strings.Split(sname, ",")
	applyCmd.Flags().StringVarP(&stype, "type", "t", "", "comma separated list of recipe step types")
	types = strings.Split(stype, ",")
	RootCmd.AddCommand(applyCmd)
}

//Apply will run only the recipes with Apply type
func Apply(args []string) error {
	fmt.Println("TODO: To be implemented")
	return nil
}
