# Kamino usage

`kamino [Flags] <command> [command-options] [recipe1 recipe2 ... recipen]`

The recipes can contains "globbing characters" (like in shell), you may have to protect them from your shell (for example ` '*dev[123]*' ` ). It is not possible to use path as recipe names (for example `../myOtherRecipe` or `/tmp/myOtherRecipe`). The recipes must be in the config path (by default the current directory).

## Flags
Flags (short form)              | meaning
--------------------------------|----------------------------------------------------------
`--config <path> (-c)          `| Lookup for recipe in the provided path (default current directory)
`--connection-retry nb         `| Maximum number of database connection retries (default 1)
`--connection-timeout duration `| Timeout of each database connection retry (default 2ms)
`--dry-run (-d)                `| List action that should be done but do not do them
`--force (-f)                  `| Execute steps without verifying the skip query
`--help (-h)                   `| Help for the selected action
`--quiet (-q)                  `| Do not print to screen
`--sequential                  `| Run the step one by one removing the parallelization
`--tags tag1,tag2 (-T)         `| Run the recipe only for datasources corresponding to the provided tags, tags can contains "globbing characters" (like in shell), you may have to protect them from your shell (for example `--tags "p[a-z]k?m*n"`)
`--verbose (-v)                `| Be more verbose for log message

## Available Commands:
*  `apply`       Apply will run the recipes provided in arguments
*  `help`        Help about any command
*  `migrate`     Manage schema migration
*  `synchronize` Manage data synchronization
*  `validate`    Validate the configuration files in recipes (datasources and steps)
*  `version`     Output the current build information

### apply
Run all the steps of provided recipes (or all the recipes of the config path).
The steps will be run by order of priority (smaller priority first). All steps of same
priority will be run in parallel (if the flags --sequential is not provided).

If more than one recipe is selected to run, the execution of the recipe will be parallelized. For each recipe:
1. Load all datasource definition
2. Load all steps and determine for each of them the datasources that will be impacted by the step by using the tags of step and datasource
3. For each priority 
    1. Make skip query (if relevant and `-f` not provided) for each steps of this priority
    2. Initialize the step not skipped (open file or database connection)
    3. Execute the step
    4. Finish the execution (commit of database transaction or file content)

If a step fails, all the steps of same priority will be cancelled (rollback of database transaction or file content)
and the next priorities will not be executed.

By default all steps of a recipe will be executed (or skipped), the `apply` action had two options to limit the list of steps:

* `--name <stepName1>,<stepName2>... (-n)` Only execute the step of provided names, the name can contains "globbing characters" (like in shell), you may have to protect them from your shell (for example ` --name '*synchronization[123]*' ` )
* `--type <type1>,<type2>... (-t)` Only execute the step of provided types


### migrate
Run only the `migration` steps of provided recipes (or all the recipes of the config path).
The execution workflow will be the same as `apply` but the skip queries will not be taken in account.

The migration that will be applied can be limited by
* `--admin (-a)`  Only admin migration (if relevant)
* `--limit <nb>`  Max number of migration (0 by default for all migration)
* `--user`        Only user migration

The `migrate` action need a mandatory sub-action (up, down or status)

#### up
Execute all the _up_ section of all _admin_ migrations and after all the _up_ section of all _user_ migrations.

Migration are applied in order defined by the file name.
Only migration that has not be already applied will be applied.

If one migration fails, the process exit with an error

#### down
Execute all the _down_ section of all _user_ migrations and after all the _down_ section of all _admin_ migrations.

Migration are applied in reverse order defined by the file name.
Only migration that has been applied will be unapplied.

If one migration fails, the process exit with an error

#### status
Show the list of all available migration and when they were applied (or not).



### synchronize
Run only the `synchronization` steps of provided recipes (or all the recipes of the config path).
The execution workflow will be the same as `apply`.

The flag `--cache-only (-C)` force kamino to use the datasource defined as synchronization cache instead of the defined source datasource.
