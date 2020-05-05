package migrate

import (
	"fmt"
	"strconv"

	"github.com/marema31/kamino/cmd/common"
	"github.com/marema31/kamino/recipe"
	"github.com/spf13/cobra"
)

// Admin only migration.
var Admin bool

// User only migration.
var User bool
var limit int

//AddCommands adds all subcommands to RootCmd.
func AddCommands(cmd *cobra.Command) {
	cmd.AddCommand(
		NewMigrateCommand(),
	)
}

//NewMigrateCommand declare the migration sub commands.
func NewMigrateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Manage schema migration",
		Args:  cobra.NoArgs,
		//		RunE: cmd.ShowHelp(),
	}

	cmd.AddCommand(
		newMigrateUpCommand(),
		newMigrateDownCommand(),
		newMigrateStatusCommand(),
	)

	cmd.PersistentFlags().BoolVarP(&Admin, "admin", "a", false, "Only admin migration (if relevant)")
	cmd.PersistentFlags().BoolVarP(&User, "user", "u", false, "Only user migration")
	cmd.PersistentFlags().IntVarP(&limit, "limit", "l", 0, "Max number of migration(0 for all)")

	return cmd
}

func createSuperseed() (map[string]string, error) {
	superseed := common.CreateSuperseed()

	if Admin && User {
		return superseed, fmt.Errorf("option --admin and --user are mutually exclusive: %w", common.ErrWrongParameterValue)
	}

	if Admin {
		superseed["migration.noAdmin"] = "false"
		superseed["migration.noUser"] = "true"
	}

	if User {
		superseed["migration.noAdmin"] = "true"
		superseed["migration.noUser"] = "false"
	}

	if limit != 0 {
		superseed["migration.limit"] = strconv.Itoa(limit)
	}

	return superseed, nil
}

// UpDown manage the actual migration.
func UpDown(direction string, cookbook recipe.Cooker, args []string) error {
	log := common.Logger.WithField("action", "migrate-"+direction)

	superseed, err := createSuperseed()
	if err != nil {
		return err
	}

	superseed["migration.dir"] = direction

	recipes, err := common.FindRecipes(log, args)
	if err != nil {
		return err
	}

	err = cookbook.Load(common.Ctx, log, common.CfgFolder, recipes, common.Tags, nil, []string{"migration"})
	if err != nil {
		return fmt.Errorf("error while loading the recipes: %w", err)
	}

	err = cookbook.PostLoad(log, superseed)
	if err != nil {
		return fmt.Errorf("error while postloading the recipes: %w", err)
	}

	if cookbook.Do(common.Ctx, log) {
		return fmt.Errorf("a step had an error: %w", common.ErrStep)
	}

	return nil
}
