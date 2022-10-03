# gosqltojson

## What is gosqltojson
gosqltojson is a command line tool to automate execution of sql queries and show the result in your screen or save it to a file(json or csv). 
Output formats: json, csv.

## Limitations
- CSV works only with output file flag --output

gosqltojson is a go implementation of [sql2json](https://github.com/fsistemas/sql2json), sql2json was written in python3.

## Why create gosqltojson if sql2json(python version) does well the job?
- Short version: For fun
- Long version: To practice and to improve my skills with go programming language creating a useful and real project

# How to use gosql2json

## Available options and help
gosql2json -help

## Default output format
The default output format is json.

## Configuration file
gosql2json by default use a config file located at USER_HOME/.sql2json/config.json

config.json structure:

```
{
    "connections": {
        "default": "sqlite+test.db",
        "postgress": "postgres+scott:tiger@localhost:5432/mydatabase",
        "mysql": "mysql+scott:tiger@localhost/foo"
    },
    "queries": {
        "default": "SELECT 1 AS a, 2 AS b",
        "sales_month_since": "SELECT inv.month, SUM(inv.amount) AS sales FROM invoices inv WHERE inv.date >= @date_from ",
        "total_sales_since": "SELECT SUM(inv.amount) AS sales FROM invoices inv WHERE inv.date >= @date_from ",
        "long_query": "@FULL_PATH_TO_SQL_FILE",
		"json": "SELECT JSON_OBJECT('id', 87, 'name', 'carrot') AS json",
		"jsonarray": "SELECT JSON_ARRAY(1, 'abc', NULL, TRUE) AS jsonarray, JSON_OBJECT('id', 87, 'name', 'carrot') AS jsonobject",
        "operation_parameters": "@/Users/myusername/myproject/my-super-test_query.sql"
    }
}
```

## Use a config.json in a different path

You can use gosql2json --config PATH_TO_YOUR_CONFIG_FILE

## Available variables to do your life easy:
- START_CURRENT_MONTH: Date the first day of current month
- CURRENT_DATE: Current Date
- END_CURRENT_MONTH: Date the last day of current month
- START_CURRENT_YEAR: First day of current year
- END_CURRENT_YEAR: First day of current year

## Date formats to CURRENT_DATE, START_CURRENT_MONTH, END_CURRENT_MONTH, START_CURRENT_YEAR, END_CURRENT_YEAR
You can use date format supported by go, default is 2006-01-02(YYYY-MM-DD)

For example:

```
gosqltojson -query "SELECT @firstDayLastMonth AS dateFrom, @lastDayThisYear AS dateTo" - --firstDayLastMonth "START_CURRENT_MONTH-1" --lastDayThisYear "END_CURRENT_YEAR"
```

Query result:
```
[
    {
        "dateFrom": "2022-08-01",
        "dateTo": "2022-12-31"
    }
]
```

## How to run queries using gosqltojson:

### Run query sales_month in database connection mysql:

gosqltojson --name mysql --query sales_month_since --date_from "START_CURRENT_MONTH-1"

Output:

```
[
    {
        "month": "January",
        "sales": 5000
    },
    {
        "month": "February",
        "sales": 3000
    }
]
```

### I don't wat an array, I want an object with an attribute with the results, useful to generate in format to post to geckoboard

gosqltojson -name mysql -query sales_month_since -wrapper data - --date_from "START_CURRENT_MONTH-1"

Output:

```
{
    "data": [
        {
            "month": "January",
            "sales": 5000
        },
        {
            "month": "February",
            "sales": 3000
        }
    ]
}
```

### Run query sales_month in database connection mysql, use month as key, sales as value:

gosqltojson -name mysql -query sales_month_since -key month -value sales - --date_from "START_CURRENT_MONTH-1"

Output:

```
[
    {
        "January": 5000
    },
    {
        "sales": 3000
    }
]
```

### Run query sales_month in database connection mysql, get the first row and only sales amount:

gosqltojson -name mysql -query total_sales_since -first -key sales - --date_from "CURRENT_DATE-10"

Output: 500 or the amount of money you have sold last 10 days

### When I use gosqltojson with result of JSON functions I get escaped strings as value

gosqltojson as a flag to allow you to specify your JSON columns

gosqltojson -name mysql -query json -jsonkeys "json, jsonarray"

Result:

```
[
    {
        "json": {
            "id":  87,
            "name", "carrot"
        }
        "jsonarray": [1, "abc", null, true],
    }
]
```

## This is only a single row I want first row only, no array.

gosqltojson -name mysql -query json -jsonkeys "json, jsonarray" -first

Result:

```
    {
        "json": {
            "id":  87,
            "name", "carrot"
        }
        "jsonarray": [1, "abc", null, true],
    }
```

### I have a long query. How to run sql query in external sql file?

query "operation_parameters"
Path "Users/myusername/myproject/my-super-query.sql"

Content of my-super-query.sql:

```
SELECT 
p.name,
p.age
FROM persons p
WHERE p.age > :min_age
AND p.creation_date > :min_date
ORDER BY p.age DESC
LIMIT 10
```

min_age: 18
min_date: Today YYYY-MM-DD

gosqltojson -name mysql -query operation_parameters - --min_age 18 --min_date "CURRENT_DATE"

```
[
    {
        "age": "40",
        "name": "p4"
    },
    {
        "age": "30",
        "name": "p3"
    },
    {
        "age": "20",
        "name": "p2"
    }
]
```

min_age: 18
min_date: First day, current year YYYY-01-01 00:00:00

gosqltojson -name mysql -query operation_parameters - --min_age 10 --min_date START_CURRENT_YEAR

```
[
    {
        "age": "40",
        "name": "p4"
    },
    {
        "age": "30",
        "name": "p3"
    },
    {
        "age": "20",
        "name": "p2"
    },
    {
        "age": "12",
        "name": "p1"
    }
]
```

### How to run external SQL query not defined in config file?

gosqltojson -name mysql -query "@/Users/myusername/myproject/my-super-query.sql" - --min_age 18 --min_date START_CURRENT_YEAR

```
[
    {
        "age": "40",
        "name": "p4"
    },
    {
        "age": "30",
        "name": "p3"
    },
    {
        "age": "20",
        "name": "p2"
    }
]
```

### Run custom query inline

You don't need to have all your queries in config file

gosqltojson -name mysql -query "SELECT CURRENT_DATE() AS date" -first -key date

```
2022-09-25
```


### Write sql query result to a CSV file

gosqltojson -name mysql -query sales_month_since -format=csv -output Sales - --date_from "START_CURRENT_MONTH-1"

```
Output:
Sales.csv
```

### Write sql query result to an Excel file

That feature is not available yet, maybe in the future.
If you really need excel pleas try [sql2json](https://github.com/fsistemas/sql2json)

### Write sql query result to a json file

gosqltojson -name mysql -query sales_month_since -format=json -output Sales - --date_from "START_CURRENT_MONTH-1"

Output:
```
Sales.json
```

# How to run during development

go run . -config config.json -query @test_query.sql -wrapper data -- -x 1 -y 2 -z 3 -a 12

go run . -config config.json -name golang -query @mysqlquery.sql -wrapper data  -- -a 11

go run . -config config.json -name golang -query @mysqlquery.sql -wrapper data  -- -date_from CURRENT_DATE

go run . -config config.json -name golang -query @mysqlquery.sql -wrapper data  -- -date_from START_CURRENT_MONT-1

./gosqltojson -query @test_query.sql -format=json -output output.json - --a 2
./gosqltojson -query @test_query.sql -format=csv -output output.csv - --a 2

## How to build
go build
go build -ldflags="-s -w"

## Run Unit Test
go test ./...
