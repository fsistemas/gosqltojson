# gosqltojson

## What is gosqltojson
gosqltojson is a command line tool use to execute sql queries and write the result to a JSON file

gosqltojson is a go implementation of [sql2json](https://github.com/fsistemas/sql2json), the first tool originally writen in python3.

## Feature parity with sql2json
So far some features of sql2json are implemented:
- databases: Only: mysql, postgres, sqlite
- formats: JSON only. Like first release of sql2json
- Multiline sql from external file using @full_path_sql_file
- output: screen only
- No custom variables and operations

## Why gosqltojson if sql2json does well the job?
For fun and to practice learning go with a real project

go run . -config config.json -query @query.sql -wrapper data -- -x 1 -y 2 -z 3

go run . -config config.json -name golang -query @mysqlquery.sql -wrapper data  -- -a 11 | jq

