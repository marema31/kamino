# Template

Create templated files. Each template step is defined its in own file (JSON, TOML or YAML) in the steps sub-folder of a recipe. 

Kamino will always do the step, since the order is not fixed, if a destination is used for more than one datasource, the order of the content of this destination may change at each Kamino run.


Attribute      | Mandatory | Definition | Default
---------------|----------------|------------|-----
destination    | yes | Path of the file to be rendered
engines        | no  | Limit the datasource selection to those corresponding to the listed engines | all engines
forceSequential| no | If true all the steps of this priority will be run sequentially
gzip           | no  | If true the destination will be gziped | false
name           | no  | Step name used for step selection by the CLI, more than one step can have the same name
replacemode    | no  | How the step will manage the fact that the destination file already exists | replace 
priority       | yes | Priority of this step on the recipe execution (ascending order)
tags           | no  | List of tags used for selecting datasource impacted by this step | all
type           | yes | Type of step, in this case _template_
template       | yes | Path of the template to be parsed
zip            | no  | If true the destination will be ziped | false


The attribute `replacemode` impact the result produced by the templating system if the destination already exists before the step
  - `append` Append the result of the step at the end of the destination file
  - `replace` Replace the content of the destination file by the result of the step
  - `skip` Skip the step
  - `unique` Append the result of the step at the end of the destination file but avoid duplicate the template rendering for a datasource if the same block is already present in the destination file.


The template file and the attributes `destination`, `query` and `template` of the step can contains Golang templates, for list of availables variables refer below

## Template format

For template, Kamino provides variables corresponding to the current datasource. To use a variable, the template should only contains `{{.MyVariable}}`, for the NamedTags you should use {{ index .NamedTags "myKey"}}. Available variables are:


Name         | Definition
-------------|------------
Database     | Database name
Engine       | Datasource engine ( CSV, JSON, MySQL, Postgres, YAML)
Environments | Environment variables usable by `{{ index .Environments "key"}}`
FilePath     | Path of the datasource (if it is a file)
Host         | Database server hostname
Name         | Datasource name
NamedTags    | Map of named tags (tags with key:value definition), usable by `{{ index .NamedTags "key"}}`
Password     | Database user password
Port         | Database server TCP port
Schema       | Schema name if relevant
Tags         | Tags of the datasource
Transaction  | True if datasource database use transaction
Type         | Datasource type (database or file)
User         | Database user

Golang templating system provides [more features](https://golang.org/pkg/html/template/). Kamino also use [Sprig](http://masterminds.github.io/sprig/) library that offer more features to the templating system