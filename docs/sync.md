# Sync


Synchronize data from database/file to multiple database/file. Each synchronization step is defined its in own file (JSON, TOML or YAML) in the steps sub-folder of a recipe. 

Kamino will determine for each destination, if the synchronization must be applied. This choice depends of the `mode`attribute of the destination.


Attribute     | Mandatory | Definition | Default
--------------|----------------|------------|-----
cache         | no  | Caching attribute for synchronization (see below)
destinations  | yes | List of destinations entry (see below)
filters       | no  | List of filters to be applied on data synchronized (see below)
name          | no  | Step name used for step selection by the CLI, more than one step can have the same name
priority      | yes | Priority of this step on the recipe execution (ascending order)
source        | yes | Source of the synchronization (see below)
type          | yes | Type of step, in this case _sync_


## Source

A synchronization will copy data from the source. This source can be either a file or a database table. Tags must be restrictive enough to select only one datasource or the step will fail.

Attribute     | Mandatory | Definition | Default
--------------|----------------|------------|-----
engines       | no  | Limit the datasource selection to those corresponding to the listed engines (Mysql, Postgres, CSV, JSON, YAML) | all datasource engines
table         | no  | Table to be synchronized. Ignored for files. If missing for database the step will fail.
tags          | no  | List of tags used for selecting datasource impacted by this step | all
types         | no  | Limit the datasource selection to those corresponding to the listed types (Database or File) | all datasource types
where         | no  | SQL WHERE expression to limit the data synchronized (only for databases)


## Destination

A synchronization will copy data to all destination. Theses destinations can be either files or a database tables. A same destination entry can select files and database without problems.

Attribute     | Mandatory | Definition | Default
--------------|----------------|------------|-----
engines       | no  | Limit the datasource selection to those corresponding to the listed engines (Mysql, Postgres, CSV, JSON, YAML) | all datasource engines
key           | no  | Key column, used by some modes to defined if a line already exist.
mode          | yes | Synchronization mode (only for database) (see below)
queries       | no  | Skip condition queries, see below for more information, superseed the mode for skipping the destination
table         | no  | Table to be synchronized. Ignored for files. If missing for database the step will fail.
tags          | no  | List of tags used for selecting datasource impacted by this step | all
types         | no  | Limit the datasource selection to those corresponding to the listed types (Database or File) | all datasource types
where         | no  | SQL WHERE expression to limit the data synchronized (only for databases)


The mode is only valid for database datasource and can be: 
*	exactCopy   : As replace but will remove line with primary key not present in source
*	insert      : Will insert all line from source (may break if primary key already present)
* 	onlyIfEmpty : Will insert only if database was empty
*	replace     : Will update if line with same primary exist or insert the line
*	truncate    : As insert but will truncate the table before
*	update      : Will update if line with same primary exist or skip the line


## Cache

A synchronization will copy data from the source. This source can be either a file or a database table. Tags must be restrictive enough to select only one datasource or the step will fail.

Attribute     | Mandatory | Definition | Default
--------------|----------------|------------|-----
allowonly     | no  | If true, the cache can be used as source even if expired if the real source is not available or if `--cache-only` flags is provided at command line|False
tags          | no  | List of tags used for selecting datasource impacted by this step | all
ttl           | no  | Validity duration of the cache

## Filter

During synchronization, the data can be filter just after being read and before being sent to the write to destinations. **Note** The filter alter the data for all the destinations, you cannot filter data for one destination and not the others. Filter can be combined.

Attribute     | Mandatory | Definition | Default
--------------|----------------|------------|-----
aparameters   | no  | List of parameter for the filter
type          | yes | Type of filter (only or filter)
mparameters   | no  | Dictionary of parameter for the filter

### only
This filter will transform the data to contains only the column listed in the `aparameters` attribute

## replace
This filter will replace the value of columns, impacted columns and the values to be replaced by are listed as a dictionary in the `mparameters` attribute. The replacement value can be a Golang template able to accesses environment variables througth `{{ index .Environments "VARIABLE_NAME" }}`.

## sed
This filter will modify the value of columns, impacted columns and the modification expression to be applied by are listed as a dictionary in the `mparameters` attribute. The expression value is using [Golang regular expression](https://github.com/google/re2/wiki/Syntax) in form `s/PATTERN_TO_FOUND/REPLACE_VALUE/`.

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

## Templates
The attributes `queries` of the step can contains Golang templates, for list of availables variables refer to this [documentation](/doc/template.md)

In the SQL script each SQL statement can be multiline and must be ended by a semi-column character. If a SQL statement must contain a semi-column (for example a postgres function definition), it must be surrounded by the lines:
```SQL
    --StatementBegin
    ....
    --StatementEnd
```
