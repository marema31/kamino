# Shell

Execute a shell script. Each shell step is defined its in own file (JSON, TOML or YAML) in the steps sub-folder of a recipe. 

Kamino will always run the script. If the script must does some action only once by datasource, it must be able to determine by itself.

Attribute     | Mandatory | Definition | Default
--------------|----------------|------------|-----
arguments     | no  | List of arguments to provides to the executed script
engines       | no  | Limit the datasource selection to those corresponding to the listed engines | all engines
environment   | no  | List of environment variables to be setted
forceSequential| no | If true all the steps of this priority will be run sequentially
ignoreErrors  | no  | Don't fails on minor errors, warn 
name          | no  | Step name used for step selection by the CLI, more than one step can have the same name
path          | no  | Current folder when executing the script
priority      | yes | Priority of this step on the recipe execution (ascending order)
scripts       | yes | Path of the script to be executed
tags          | no  | List of tags used for selecting datasource impacted by this step | all
type          | yes | Type of step, in this case _shell_

## Templates
The attributes `arguments`, `environment` and `path` can contains Golang templates, for list of availables variables refer to this [documentation](/doc/template.md)
