# kamino
[![Golang Documentation](https://godoc.org/github.com/marema31/kamino?status.svg)](https://godoc.org/github.com/marema31/kamino) 
[![Build Status](https://travis-ci.org/marema31/kamino.svg)](https://travis-ci.org/marema31/kamino) 
[![GitHub release](http://img.shields.io/github/release/marema31/kamino.svg?style=flat-square)](https://github.com/marema31/kamino/releases/latest)
[![Golang CI](https://golangci.com/badges/github.com/marema31/kamino.svg)](https://golangci.com/r/github.com/marema31/kamino)
[![Go Report Card](https://goreportcard.com/badge/github.com/marema31/kamino)](https://goreportcard.com/report/github.com/marema31/kamino)
[![CodeCoverage](https://codecov.io/gh/marema31/kamino/branch/master/graph/badge.svg)](https://codecov.io/gh/marema31/kamino/branch/master)
[![BSD-2 License](http://img.shields.io/badge/license-BSD--2-blue.svg?style=flat)](LICENSE)

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
1. Run a sql script as admin to create roles and/or users
1. Run a sql script as admin to create a schema
1. Run sql-migrate to populate the database schema
1. Synchronize tables from other database or files
1. Run a sql script as admin to give rights on tables/schemas
1. Generate a configuration file from a golang template with information of the database
1. Call a shell script to post a completion message to slack


## Concepts

Kamino combine several type of 'objects' to manage the database lifecycle

### Datasource
Each environment can be composed of several 'datasources'. This datasource is connection informations needed for a recipe (server, users/password, engine type, database name and/or schema if relevant). Each datasource is described in [files](/docs/datasource.md) in JSON, TOML or YAML format

### Step
A step is an action to be done on the selected datasources. The action could be:
   * [_sql_](/docs/sql.md)     : run a templated sql script
   * [_migrate_](/docs/migration.md) : apply sql-migrate migration, the actions will be different if Kamino is in apply or migrate mode.
   * [_sync_](/docs/sync.md)    : Synchronize data from an other database or files
   * [_shell_](/docs/shell.md)   : run a shell script
   * [_template_](/docs/template.md): Create a file from a template

More than one step can share the same 'name'. When user runs Kamino in _apply_ mode, he can select a liste of step name or a step type to be run, only the steps corresponding to this criterias will be executed.

Each steps are described in files in JSON,TOML or YAML format. The needed attributes depend on the type of action.

**Note:** All relative path provided in attribute of a step are relative to the recipe folder that contains the recipe.


#### Idempotency
In _apply_ and _sync_ mode, before running a step of type _migration_, _sql_ or _sync_, Kamino will determine if the action is necessary by running a provided SQL query, if this query returns a count different of 0, the step will be skipped. ~~Once the action done, the query is re-run and the results must be different or Kamino will stop with an error.~~ 
This condition can be removed with the --force option.

In _migrate_ mode, the condition will never be enforced.

For step of type _shell_, it is the script responsability to be idempotent.

The idempotency is not garanteed for step of type _template_

### Tags
Each datasource have _tags_ (list of tag) that can be shared among them. These tags are used by steps to select on which datasource the action will occur.

Steps can have a list of _tag_ selectors, if a datasource correspond to one selector of the list, it will be choosed by the step, the negation will remove datasource from the list of choosen ones.

Tags containing one ":" character will be considered has "named tags", their usage for selection is not different, but for template they can be referenced by Tags array or by the NamedTags map where the tag is splitted in two (name of tag before the ":", value of tag after).

A tag selector can be:
* A single tag (myTag): All datasources with this tag will be selected
* Intersection of tag (myTag1.myTag2): Only datasource with all this tag will be selected
* A negation of a tag selector (!ta1.ta2): All datasources complying to this tag selector will be excluded of the final list of datasource impacted by the step. 

Tag selector can contains "globbing characters" (like in shell)

### Recipe
A recipe is a collection of steps (a folders that contains the steps), the orders of step will be determine by the priority of the step (lower will be the first one), if more than one step have the same priority, they will be run in parallel. 


## Usage
Create a folder named after the name of your recipe, in this folder populate at least two folders (_datasources_ and _steps_) with corresponding files and run `kamino apply`. 

Refer to the [documentation](/docs/cli.md) for a comprehensive description of the tool options

### Example

The [_testdata_](/testdata) folder of this repository contains an exemple of a recipe (pokemon) with all possible actions/options and a _docker-compose.yml_ file.  To see kamino in action run 

    docker-compose -f testdata/pokemon/docker-compose.yml up -d
    kamino --connection-timeout 10s --connection-retry 5   -c testdata  apply pokemon
  
  That will spin-up 
  * three mysql containers (root/adminpw)
  * one postgres (postgres/adminpw) 
  * phpmyadmin (http://localhost:12348) 
  * phppgadmin (http://localhost:12349). 
  
  That will also start kamino using the testdata folder as base folder for recipe and applying the pokemon recipe. 
  
  The `--connection--*` parameters force to kamino to retry the first connection to datasources 5 times at 10 seconds interval, which is needed on this case because docker-compose will return when all containers are up, but database engines accept network connections only several seconds after the start of the containers.

  This example is meant to show all the possibilities of kamino, your recipe does not have to be so complicated.

## Contribution
I've made this project as a real use case to learn Golang.
I've tried to adopt the Go mindset but I'm sure that other gophers could do better. 

If this project can be useful for you, feel free to open issues, propose documentation updates or even pull request.

Contributions are of course always welcome!

1. Fork marema31/kamino (https://github.com/marema31/kamino/fork)
2. Create a feature branch
3. Commit your changes
4. Run test using `go test ./...`
5. Create a Pull Request

See [`CONTRIBUTING.md`](https://github.com/marema31/kamino/blob/master/CONTRIBUTING.md) for details.

## License

Copyright (c) 2019-present [Marc Carmier](https://github.com/marema31)

Licensed under [BSD-2 Clause License](./LICENSE)