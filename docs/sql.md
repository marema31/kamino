# SQL

Run templated SQL scripts. Each sql step is defined its in own file (JSON, TOML or YAML) in the steps sub-folder of a recipe. 

In _apply_ mode, Kamino use the `query` parameter for each selected datasource before executing the migration, if for a datasource this query returns a count different of 0, the step will be skipped for this datasource. This behavior is disable by using the `--force` CLI flags.


Attribute     | Mandatory | Definition | Default
--------------|----------------|------------|-----
admin         | no  | Use the _admin_ user account of datasource to execute the SQL script | false
engines       | no  | Limit the datasource selection to those corresponding to the listed engines | all engines
name          | no  | Step name used for step selection by the CLI, more than one step can have the same name
nodb          | no  | Use _admin_ user account of datasource and connect to the database server without a database (useful for database creation) to execute the SQL script | false 
nouser        | no  | If true the step will not apply _user_ migration | false 
priority      | yes | Priority of this step on the recipe execution (ascending order)
queries       | yes | Skip condition queries, each query should returns only one column/line, if this result is different of 0, the step will be skipped, if there is more than one query they will be executed in order until one returns a 0 or all have a different result
tags          | no  | List of tags used for selecting datasource impacted by this step | all
type          | yes | Type of step, in this case _sql_
template      | yes | Path of the SQL script template to be executed
transaction   | no  | If true use transaction if the datasource has transaction defined |false
unique        | no  | If true the script will be run only once by unique datasource URL corresponding to the value of admin and nodb attributes|false

If admin = false and noDb = false, the SQL script will be executed with the _user_ account of the datasource.

The SQL script and the attributes `query` and `template` of the step can contains Golang templates, for list of availables variables refer to this [documentation](/doc/template.md)

In the SQL script each SQL statement can be multiline and must be ended by a semi-column character. If a SQL statement must contain a semi-column (for example a postgres function definition), it must be surrounded by the lines:
```SQL
    --StatementBegin
    ....
    --StatementEnd
```