# Migration

Apply sql-migrate migration, the actions will be different if Kamino is in apply or migrate mode. Each migration step is defined its in own file (JSON, TOML or YAML) in the steps sub-folder of a recipe. 

Kamino use the [sql-migrate library](https://github.com/rubenv/sql-migrate) to manage the migrations.

Migration is a list of ordered (by filename) sql files containing _up_ and _down_ snippets to manage the schema of a database. Each of this files correspond to a version of the schema of the database. Kamino can manage two levels of migrations that use different user with different rights. The migration files available on the migration folder are applied with the _user_ of the datasource, the migration files on the _admin_ subfolder of migration are applied with the _admin_ user of the datasource. _Admin_ migration are not mandatory.

In _apply_ mode, Kamino use the `query` parameter for each selected datasource before executing the migration, if for a datasource this query returns a count different of 0, the step will be skipped for this datasource. This behavior is disable in _migrate_ mode or by using the `--force` CLI flags.

Attribute     | Mandatory | Definition | Default
--------------|----------------|------------|-----
admintable    | no  | Table used for store applied _admin_ migration informations | kamino_admin_migrations
engines       | no  | Limit the datasource selection to those corresponding to the listed engines | all engines
folder        | yes | Folder containing the migration files (template)
forceSequential| no | If true all the steps of this priority will be run sequentially
ignoreErrors  | no  | Don't fails on minor errors, warn 
limit         | no  | Number of migration applied | 0 (all)
name          | no  | Step name used for step selection by the CLI, more than one step can have the same name
noadmin       | no  | If true the step will not apply _admin_ migration | false 
noforce       | no  | If true the step will be skipped for `kamino migrate` sub-command or for `kamino apply --force` 
nouser        | no  | If true the step will not apply _user_ migration | false 
priority      | yes | Priority of this step on the recipe execution (ascending order)
queries       | no  | Skip condition queries, see below for more information
tags          | no  | List of tags used for selecting datasource impacted by this step | all
type          | yes | Type of step, in this case _migration_
usertable     | no  | Table used for store applied _user_ migration informations | kamino_user_migrations

## Templates 
The attributes `folder` and `query` can contains Golang templates, for list of availables variables refer to this [documentation](/doc/template.md)

## Skip queries
In _apply_ mode, Kamino use the `queries` parameter for each selected datasource before executing the migration to determine if the step for this datasource should be skipped. This behavior is disable by using the `--force` CLI flags.

***Note:*** the step will never be skipped if the migration table(s) does not exists or are empty, tested tables (`admintable`/`usertable`) will depend of `noadmin`/`nouser` flags of the step.

`queries` parameter is a list of templated SQL query. Theses queries will be rendered for each datasource with values specific to this datasource.

Each queries: 
  * can be prefixed by condition with a form of `!`, `=<number>:` or `!=<number>:`, if no condition is provided it will be `=0:`, `!` is an abbreviation for `!=0:`
  * must returns only one column/line, the return value is compare to the condition,
  * are runned one by one until one validates the corresponding condition.
  
If all queries does not validate their conditions, the step will be skipped for this datasource.

#### example 
```yaml
- queries:
  - "SELECT COUNT(id) from table1"
  - "!SELECT COUNT(id) from table2"
  - "=10:SELECT COUNT(id) from table3"
  - "!=10:SELECT COUNT(id) from table4"
```

The step will be skipped only if: 
  * table1 has at least one row,
  * table2 is empty
  * table3 does not have 10 rows,
  * table4 has 10 rows