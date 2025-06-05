#! /bin/bash

# Perform migrations using goose (https://github.com/pressly/goose). The script assumes that:
# 1. you are in the root directory of the project
# 2. there is a .env file in the currect directory
# 3. migrations are stored in sql/schema

prog=$(basename $0)
action=up

if [[ $# == 1 ]]; then
    if [[ $1 == "up" || $1 == "down" ]]; then
        action=$1
    else
        echo "usage: $prog [up/down]"
        exit 1
    fi
elif [[ $# > 1 ]]; then
    echo "usage: $prog [up/down]"
    exit 1
fi

if [[ ! -f .env ]] then echo "$prog: No .env found";  exit 1; fi
url=$(grep ^DB_URL .env | cut -d= -f2)
if [[ -z $url ]] then echo "$prog: no database URL in environment"; exit 1; fi
if [[ ! -d sql/schema ]]; then echo "$prog: schema directory not found"; exit 1; fi

cd sql/schema
goose postgres $url $action
