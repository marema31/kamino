package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/marema31/kamino/cmd/migrate"

	"github.com/spf13/cobra"
)

var (
	cfgFolder   string
	dryRun      bool
	quiet       bool
	environment string
	instances   []string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
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
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVarP(&cfgFolder, "config", "c", "", "config folder (default is $HOME/.kamino.d)")
	RootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false, "list action only do not do them")
	RootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "do not print to screen")
	RootCmd.PersistentFlags().StringVarP(&environment, "environment", "e", "", "database environment (by default the only existing environment)")
	var sinstance string
	RootCmd.PersistentFlags().StringVarP(&sinstance, "instances", "i", "", "comma separated list of instance (default is all the instances)")
	instances = strings.Split(sinstance, ",")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	migrate.AddCommands(RootCmd)
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
