# kamino

Manage development databases lifecycle described in 'recipes'

Kamino can be used to automatically:
  * create database instances
  * create the database schema (via sql-migrate migration)
  * import initial dataset (from other database or files)
  * generate configuration file for tools using these databases
  * call shell script with information on these databases

Kamino is intend to manage databases for development, testing or CI/CD environments
using the 'as-code' devops motto in simple and idempotent way

A common 'recipe' would implies the following steps:
	1. Run a sql script as admin to create the database
	2. Run a sql script as admin to create roles and/or users
  3. Run a sql script as admin to create a schema
	4. Run sql-migrate to populate the database schema
	5. Synchronize tables from other database or files
	6. Run a sql script as admin to give rights on tables/schemas
	7. Generate a configuration file from a golang template with information of the database
	8. Call a shell script to post a completion message to slack


## Concepts

Kamino combine several type of "objects" to manage the database lifecycle

### Environment
All execution of Kamino will occur on an unique "environment" ( development, testing, etc... this names are not enforced). If environment is not provided by the user Kamino will try to determine a "default" environment.

### Datasource
Each environment can be composed of several "datasources". This datasource is connection informations needed for a recipe (server, users/password, engine type, database name and/or schema if relevant). Each of this datasource has also have "tags" that can be shared among them. This tag will be used by recipes to select on which datasource the operation will occur.

### Instance
User can provide a list of "instances" when running Kamino. A instance correspond to a datasource tag that will be combined to the tag selection of recipes to narrow the list of datasources on which the operation will occur. 

By example, if a recipe select all the datasource with a tag "pokemon" and user provides "fr,us" as instances list, only the datasource that have the tuple of tag ("pokemon","fr") or ("pokemon","us") will be selected for this execution, the databsource ("pokemon","uk") will be skipped, if the user does not provide a instances list, the three datasources will be selected.

### Step
A step is an action to be done on the selected datasources. The action could be:
   * sql     : run a templated sql script
   * migrate : apply sql-migrate migration, the actions will be different if Kamino is in apply or migrate mode.
   * sync    : Synchronize data from an other database or files
   * shell   : run a shell script
   * template: Create a file from a template

More than one step can share the same "name". When user runs Kamino in "apply mode", he can select a liste of step name or a step type to be run, only the steps corresponding to this criterias will be executed.

#### Idempotency
In "apply mode", before running a step of type "sql" or "sync", Kamino will determine if the action is necessary by running a provided SQL query, if this query returns a count different of O, the step will be skipped. Once the action done, the query is re-run and the results must be different or Kamino will stop with an error. This condition will not be enforced in "sync" mode.

In "apply mode", step of type "migration" will be run only if the database contains no tables. This condition will not be enforced in "migrate mode".

For step of type "shell", it is the script responsability to be idempotent.

The idempotency is not garanteed for step of type "template"

### Recipe
A recipe is a collection of steps (a folders that contains the steps), the orders of step will be determine by the priority of the step (lower will be the first one), if more than one step have the same priority, they will be run in parallel. 

### Migration
List of ordered sql files containing "up" and "down" snippets to manage the schema of a database. Each of this files correspond to a version of the schema of the database. Kamino use the [sql-migrate library](https://github.com/rubenv/sql-migrate) to manage the migrations.



