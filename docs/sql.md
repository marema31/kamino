# SQL

Run templated SQL scripts. Each sql step is defined its in own file (JSON, TOML or YAML) in the steps sub-folder of a recipe. 

Attribute     | Mandatory | Definition | Default
--------------|----------------|------------|-----
admin         | no  | Use the _admin_ user account of datasource to execute the SQL script | false
engines       | no  | Limit the datasource selection to those corresponding to the listed engines | all engines
forceSequential| no | If true all the steps of this priority will be run sequentially
ignoreErrors  | no  | Don't fails on minor errors, warn 
name          | no  | Step name used for step selection by the CLI, more than one step can have the same name
nodb          | no  | Use _admin_ user account of datasource and connect to the database server without a database (useful for database creation) to execute the SQL script | false 
nouser        | no  | If true the step will not apply _user_ migration | false 
priority      | yes | Priority of this step on the recipe execution (ascending order)
queries       | no  | Skip condition queries, see below for more information
tags          | no  | List of tags used for selecting datasource impacted by this step | all
type          | yes | Type of step, in this case _sql_
template      | yes | Path of the SQL script template to be executed
transaction   | no  | If true use transaction if the datasource has transaction defined |false
unique        | no  | If true the script will be run only once by unique datasource URL corresponding to the value of admin and nodb attributes|false

If admin = false and noDb = false, the SQL script will be executed with the _user_ account of the datasource.

## Templates
The SQL script and the attributes `queries` and `template` of the step can contains Golang templates, for list of availables variables refer to this [documentation](/doc/template.md)

In the SQL script each SQL statement can be multiline and must be ended by a semi-column character. If a SQL statement must contain a semi-column (for example a postgres function definition), it must be surrounded by the lines:
```SQL
    --StatementBegin
    ....
    --StatementEnd
```

## Skip queries
In _apply_ mode, Kamino use the `queries` parameter for each selected datasource before executing the migration to determine if the step for this datasource should be skipped. This behavior is disable by using the `--force` CLI flags.

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
