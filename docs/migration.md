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
name          | no  | Step name used for step selection by the CLI, more than one step can have the same name
noadmin       | no  | If true the step will not apply _admin_ migration | false 
nouser        | no  | If true the step will not apply _user_ migration | false 
priority      | yes | Priority of this step on the recipe execution (ascending order)
queries       | yes | Skip condition queries, each query should returns only one column/line, if this result is different of 0, the step will be skipped, if there is more than one query they will be executed in order until one returns a 0 or all have a different resulttags          | no  | List of tags used for selecting datasource impacted by this step | all
type          | yes | Type of step, in this case _migration_
usertable     | no  | Table used for store applied _user_ migration informations | kamino_user_migrations


The attributes `folder` and `query` can contains Golang templates, for list of availables variables refer to this [documentation](/doc/template.md)
