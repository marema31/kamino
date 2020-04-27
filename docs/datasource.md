# Datasource

Database connection or file informations needed for a recipe. Each datasource is defined its in own file (JSON, TOML or YAML) in the datasources sub-folder of a recipe. 

Datasource definition contains a list of _tags_ that can be shared among them. These tags will be used by recipes to select on which datasource the operation will occur.

File datasource can only be used for _synchronization_ steps.

Attribute     | Kind/Mandatory | Definition | Default
--------------|----------------|------------|-----
admin         | Database       | Database user with rights needed for admin section of steps | root (mysql) / postgres(postgres)
adminpassword | Database       | Password for the admin user
database      | Database *     | Database name
engine        | All *          | Provider use for the datasource ( mysql, postgres, csv, json or yaml)
file          | File *         | File path for the datasource. Path are relative to recipe folder.
gzip          | File           | If true the source is gziped | false
host          | Database       | Database server (default: localhost)
options       | Database       | Options to the connection string (e.g. sslmode=disable for postgres, tls=skip-verify for mysql)
password      | Database       | Password of database user
port          | Database       | Database server TCP port | 3306 (mysql) / 5432 (postgres)
shema         | Database       | Name of the database schema
tags          | All *          | List of tags that can be used to select this datasource
transaction   | Database       | If true, some step types will use transaction | false
user          | Database       | Database user with rights needed for non-admin section of steps | root (mysql) / postgres(postgres)
zip           | File           | If true the source is ziped | false

Most of the Attribute can take Golang template with the possibility to use environment variables values like so `{{ index .Environments "key"}}`