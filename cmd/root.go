//Package cmd manage the first level of command of the CLI
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/marema31/kamino/cmd/migrate"

	"github.com/spf13/cobra"
)

var (
	ctx         context.Context
	cfgFolder   string
	dryRun      bool
	quiet       bool
	environment string
	instances   []string
)

// RootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kamino [OPTIONS] COMMAND <recipe> ... <recipe>",
	Short: "Development database manager",
	Long: `
			Manage development databases lifecycle described in 'recipes'

			It can be used to automatically:
			  * create database instances
			  * create the database schema (via sql-migrate migration)
			  * import initial dataset (from other database or files)
			  * generate configuration file for tools using these databases
			  * call shell script with information on these databases

			It can be used to manage databases for development or testing environment
			using the 'as-code' devops motto in simple and idempotent way`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(c context.Context) {
	ctx = c //Store the context for all sub-command definition
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVarP(&cfgFolder, "config", "c", "", "config folder (default is $HOME/.kamino.d)")
	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false, "list action only do not do them")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "do not print to screen")
	rootCmd.PersistentFlags().StringVarP(&environment, "environment", "e", "", "database environment (by default the only existing environment)")
	rootCmd.PersistentFlags().StringSliceVarP(&instances, "instances", "i", []string{}, "comma separated list of instance (default is all the instances)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	migrate.AddCommands(ctx, rootCmd)
}

//TODO: called by init function
// initConfig reads in config file and ENV variables if set.
func initConfig() {
	/*
		if cfgFolder != "" {
			// Use config file from the flag.
			viper.SetConfigFile(cfgFile)
		} else {
			// Find home directory.
			home, err := homedir.Dir()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			viper.SetDefault("SAMPATH", "~/.sam")
			// Search config in home directory with name ".sam" (without extension).
			viper.AddConfigPath(home)
			viper.SetConfigName(".sam")
		}

		viper.AutomaticEnv() // read in environment variables that match

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	*/
}
