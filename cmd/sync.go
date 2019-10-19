package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	syncCmd = &cobra.Command{
		Use:     "synchronize [flags] <recipe> ... <recipe>:<table>,...,<table>",
		Short:   "Synchronize will run only the recipes (or the table of recipe) with sync type",
		Long:    ``,
		Aliases: []string{"sync"},
		RunE: func(_ *cobra.Command, args []string) error {
			return Sync(args)
		},
	}
	cacheOnly bool
)

func init() {
	syncCmd.Flags().BoolVarP(&cacheOnly, "cache-only", "C", false, "Use only cache as source")
	rootCmd.AddCommand(syncCmd)
}

//Sync will run only the recipes with sync type
func Sync(args []string) error {
	fmt.Println("TODO: To be implemented")
	return nil
}
