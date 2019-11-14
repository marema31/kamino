//Package cmd manage the first level of command of the CLI
package cmd

import (
	"context"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/marema31/kamino/cmd/migrate"
)

var (
	ctx    context.Context
	logger = logrus.New()
	// CfgFolder configuration folder
	CfgFolder string
	// DryRun would not really do the action but logs
	DryRun bool
	// Quiet no logs
	Quiet bool
	// Verbose add debug logs
	Verbose bool
	//TODO: use this for the subcommands
	tags []string
	// Timeout of each database ping try
	conTimeout time.Duration
	// Number of retry of database ping
	conRetry int
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
func Execute(c context.Context) error {
	ctx = c //Store the context for all sub-command definition
	return rootCmd.Execute()
}

// Called at package initialization, before main execution
func init() {
	cobra.OnInitialize(InitConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVarP(&CfgFolder, "config", "c", "", "config folder")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "logs more verbose")
	rootCmd.PersistentFlags().BoolVarP(&DryRun, "dry-run", "d", false, "list action only do not do them")
	rootCmd.PersistentFlags().BoolVarP(&Quiet, "quiet", "q", false, "do not print to screen")
	rootCmd.PersistentFlags().StringSliceVarP(&tags, "tags", "T", []string{}, "comma separated list of tags to filter the calculated impacted datasources")
	rootCmd.PersistentFlags().DurationVar(&conTimeout, "connection-timeout", time.Millisecond*2, "timeout of each database connection retry")
	rootCmd.PersistentFlags().IntVar(&conRetry, "connection-retry", 1, "number maximum of database connection retries")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	migrate.AddCommands(ctx, rootCmd)
}

// GetLogger returns the logger instancied at initialization phase
func GetLogger() *logrus.Logger {
	return logger
}

// InitConfig reads in config file and ENV variables if set.
func InitConfig() {
	// Log handling
	Formatter := new(logrus.TextFormatter)
	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	Formatter.FullTimestamp = true
	logger.SetFormatter(Formatter)

	if Verbose {
		logger.SetLevel(logrus.DebugLevel)
	}
	if Quiet {
		logger.SetLevel(logrus.PanicLevel)
	}

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
