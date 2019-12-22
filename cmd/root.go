//Package cmd manage the first level of command of the CLI
package cmd

import (
	"context"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/marema31/kamino/cmd/common"
	"github.com/marema31/kamino/cmd/migrate"
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

// Execute the corresponding cobra sub-command
func Execute(c context.Context) error {
	common.Ctx = c //Store the context for all sub-command definition
	return rootCmd.Execute()
}

// Called at package initialization, before main execution
func init() {
	cobra.OnInitialize(InitConfig)

	rootCmd.PersistentFlags().StringVarP(&common.CfgFolder, "config", "c", "", "config folder")
	rootCmd.PersistentFlags().BoolVarP(&common.DryRun, "dry-run", "d", false, "list action only do not do them")
	rootCmd.PersistentFlags().BoolVarP(&common.Force, "force", "f", false, "execute steps without verifying the skip query")
	rootCmd.PersistentFlags().BoolVarP(&common.Quiet, "quiet", "q", false, "do not print to screen")
	rootCmd.PersistentFlags().IntVar(&common.Retry, "connection-retry", 1, "number maximum of database connection retries")
	rootCmd.PersistentFlags().BoolVar(&common.Sequential, "sequential", false, "run the step one by one")
	rootCmd.PersistentFlags().StringSliceVarP(&common.Tags, "tags", "T", []string{}, "comma separated list of tags to filter the calculated impacted datasources")
	rootCmd.PersistentFlags().DurationVar(&common.Timeout, "connection-timeout", time.Millisecond*2, "timeout of each database connection retry")
	rootCmd.PersistentFlags().BoolVarP(&common.Verbose, "verbose", "v", false, "logs more verbose")
	migrate.AddCommands(rootCmd)
}

// GetLogger returns the logger instancied at initialization phase
func GetLogger() *logrus.Logger {
	return common.Logger
}

// InitConfig reads in config file and ENV variables if set.
func InitConfig() {
	// Log handling
	Formatter := new(logrus.TextFormatter)
	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	Formatter.FullTimestamp = true
	common.Logger.SetFormatter(Formatter)

	/* TODO: Configuration file ?
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

		viper.SetDefault("KAMINOPATH", "~/.kaminorc")
		// Search config in home directory with name ".kaminorc" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".kaminorc")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	*/
	if common.Verbose {
		common.Logger.SetLevel(logrus.DebugLevel)
	}

	if common.Quiet {
		common.Logger.SetLevel(logrus.PanicLevel)
	}
}
