#! /bin/sh

shift
DATABASE=$1
shift
shift
USER=$1
shift


echo "Was notified for ${DATABASE} database with ${USER} user due to the following tags: $*" |sed -e 's/-t //g'