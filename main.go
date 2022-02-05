package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	configFile    = flag.String("config", GetenvOrDefault("HOME", ".")+"/.sql2json/config.json", "Config file")
	conectionName = flag.String("name", "default", "Database connection name in config file")
	queryName     = flag.String("query", "default", "Query name in config file")
	useWrapper    = flag.String("wrapper", "", "Use an extra object as a wrapper")
	firstOnly     = flag.Bool("first", false, "Get first row only")
	keyName       = flag.String("key", "", "Field name used to compute key name")
	valueName     = flag.String("value", "", "Field name used to compute value for key")
)

type Queries map[string]string
type Conections map[string]string

type Config struct {
	Conections Conections `json:"conections"`
	Queries    Queries    `json:"queries"`
}

type ConectionConfig struct {
	ConectionType   string
	ConectionString string
}

func NewConectionConfig(rawConectionString string) (*ConectionConfig, error) {

	parts := strings.Split(rawConectionString, "+")

	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid conection string")
	}

	return &ConectionConfig{ConectionType: parts[0], ConectionString: parts[1]}, nil
}

func (config *Config) getConection(conectionName string) string {
	rawQuery := config.Conections[conectionName]

	if rawQuery == "" {
		return conectionName
	}

	return rawQuery
}

func (config *Config) getQuery(queryName string) string {
	rawQuery := config.Queries[queryName]

	if rawQuery == "" {
		return queryName
	}

	return rawQuery
}

func GetenvOrDefault(envName, defaultValue string) string {
	v := os.Getenv(envName)

	if v == "" {
		return defaultValue
	}

	return v
}

func loadQueryFromFile(file string) (string, error) {
	if file == "" {
		return "", fmt.Errorf("file is required")
	}

	query, err := ioutil.ReadFile(file)

	if err != nil {
		return "", err
	}

	return string(query), nil
}

func loadConfigFile(file string) (*Config, error) {
	config := new(Config)
	if file == "" {
		return nil, fmt.Errorf("file is required")
	}

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	err = json.NewDecoder(f).Decode(&config)

	if err != nil {
		return nil, err
	}

	return config, nil
}

type QueryResult map[string]interface{}

func parseRowsToMaps(rows *sql.Rows) ([]QueryResult, error) {
	var results []QueryResult

	columns, err := rows.Columns()

	if err != nil {
		log.Fatal(err.Error())
	}

	values := make([]interface{}, len(columns))
	scan_args := make([]interface{}, len(columns))

	for i := range values {
		scan_args[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scan_args...)

		if err != nil {
			panic(err.Error())
		}

		rowMap := make(map[string]interface{})

		for i, colum := range columns {
			rowMap[colum] = fmt.Sprintf("%s", values[i])
		}

		results = append(results, rowMap)
	}

	return results, nil
}

func parseQueryParams(args []string) (map[string]interface{}, error) {
	parameters := make(map[string]interface{})

	i := 0
	size := len(args)

	for {
		if i >= size {
			break
		}

		name := strings.ReplaceAll(args[i], "-", "")
		value := args[i+1]

		if strings.HasPrefix(value, "-") {
			value = "true"
			i = i + 1
		} else {
			i = i + 2
		}

		parameters[name] = value
	}

	return parameters, nil
}

func parseUsedQueryParams(queryRaw string, queryParamsMap map[string]interface{}) (map[string]interface{}, error) {
	usedQueryParams := make(map[string]interface{})

	for key, value := range queryParamsMap {
		if strings.Contains(queryRaw, "@"+key) {
			usedQueryParams[key] = value
		}
	}

	return usedQueryParams, nil
}

func main() {
	//myFlagSet := flag.NewFlagSet("", flag.ContinueOnError)
	//err := myFlagSet.Parse(os.Args[1:])

	/*
		//flag.CommandLine.Init("", flag.ContinueOnError)
		//err := flag.CommandLine.Parse(os.Args[1:])

		if err != nil {
			log.Fatal(err.Error())
		}
	*/

	flag.Parse()

	//fmt.Println(flag.Args())

	config, err := loadConfigFile(*configFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	conectionString := config.getConection(*conectionName)
	queryRaw := config.getQuery(*queryName)

	if strings.HasPrefix(queryRaw, "@") {
		queryRaw, err = loadQueryFromFile((queryRaw[1:]))
		if err != nil {
			panic(err.Error())
		}
	}

	conectionConfig, err := NewConectionConfig(conectionString)

	if err != nil {
		panic("failed to parse database conection string. valid format: type+conString. example: sqlite+file.db, cause: " + err.Error())
	}

	var sqlDB *sql.DB
	var db *gorm.DB

	if conectionConfig.ConectionType == "mysql" {
		sqlDB, err = sql.Open(conectionConfig.ConectionType, conectionConfig.ConectionString)

		if err != nil {
			log.Fatal(err.Error())
		}

		db, err = gorm.Open(mysql.New(mysql.Config{
			Conn: sqlDB,
		}), &gorm.Config{})

		if err != nil {
			log.Fatal(err.Error())
		}
	} else if conectionConfig.ConectionType == "postgres" {
		sqlDB, err = sql.Open(conectionConfig.ConectionType, conectionConfig.ConectionString)

		if err != nil {
			log.Fatal(err.Error())
		}

		db, err = gorm.Open(postgres.New(postgres.Config{
			Conn: sqlDB,
		}), &gorm.Config{})

		if err != nil {
			log.Fatal(err.Error())
		}
	} else if strings.HasPrefix(conectionConfig.ConectionType, "sqlite") {
		db, err = gorm.Open(sqlite.Open(conectionConfig.ConectionString), &gorm.Config{})

		if err != nil {
			log.Fatal(err.Error())
		}
	} else {
		log.Fatal(fmt.Errorf("invalid conection type: %s. valid types: mysql, postgres, sqlite", conectionConfig.ConectionType))
	}

	if sqlDB != nil {
		defer sqlDB.Close()
	}

	//map[string]interface{}{"name": "jinzhu", "name2": "jinzhu2"})

	queryParams, err := parseQueryParams(flag.Args())

	if err != nil {
		log.Fatal(err.Error())
	}

	usedQueryParams, err := parseUsedQueryParams(queryRaw, queryParams)

	//fmt.Println(usedQueryParams)

	if err != nil {
		log.Fatal(err.Error())
	}

	var dbQuery *gorm.DB

	if len(usedQueryParams) > 0 {
		dbQuery = db.Raw(queryRaw, usedQueryParams)
	} else {
		dbQuery = db.Raw(queryRaw)
	}

	rows, err := dbQuery.Rows()

	if err != nil {
		log.Fatal(err.Error())
	}

	defer rows.Close()

	mapOfRows, err := parseRowsToMaps(rows)

	if err != nil {
		log.Fatal(err.Error())
	}

	var result interface{}

	if *firstOnly {
		item := mapOfRows[0]

		if len(mapOfRows) > 0 {
			if *keyName != "" && *valueName != "" {
				itemResult := make(map[string]interface{})

				key := fmt.Sprintf("%s", item[*keyName])
				value := item[*valueName]

				itemResult[key] = value

				result = itemResult
			} else if *keyName != "" || *valueName != "" {
				var key string

				if *keyName != "" {
					key = *keyName
				} else {
					key = *valueName
				}

				result = item[key]
			} else {
				if *keyName != "" && *valueName == "" {
					result = ""
				} else {
					result = item
				}
			}
		} else {
			result = make(map[string]interface{})
		}
	} else {
		if *keyName != "" && *valueName != "" {
			var newRows []QueryResult

			for _, row := range mapOfRows {
				key := fmt.Sprintf("%v", row[*keyName])
				value := row[*valueName]
				item := make(map[string]interface{})
				item[key] = value
				newRows = append(newRows, item)
			}

			result = newRows
		} else if *keyName != "" || *valueName != "" {
			var key string

			if *keyName != "" {
				key = *keyName
			} else {
				key = *valueName
			}

			var newRows []string

			for _, row := range mapOfRows {
				value := fmt.Sprintf("%v", row[key])
				newRows = append(newRows, value)
			}

			result = newRows
		} else {
			result = mapOfRows
		}
	}

	if *useWrapper != "" {
		withWrapper := make(map[string]interface{})
		withWrapper[*useWrapper] = result

		res, err := json.Marshal(withWrapper)
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Println(string(res))
	} else {
		res, err := json.Marshal(result)
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Println(string(res))
	}
}
