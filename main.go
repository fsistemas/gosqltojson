package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"franciscoperez.dev/gosqltojson/database"
	"franciscoperez.dev/gosqltojson/io"
)

var (
	configFile    = flag.String("config", io.GetEnvOrDefault("HOME", ".")+"/.sql2json/config.json", "Config file")
	conectionName = flag.String("name", "default", "Database connection name in config file")
	queryName     = flag.String("query", "default", "Query name in config file")
	wrapper       = flag.String("wrapper", "", "Use an extra object as a wrapper")
	firstOnly     = flag.Bool("first", false, "Get first row only")
	keyName       = flag.String("key", "", "Field name used to compute key name")
	valueName     = flag.String("value", "", "Field name used to compute value for key")
	output        = flag.String("output", "", "file name to write the output")
	format        = flag.String("format", "json", "Format to write the output")
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	flag.Parse()

	runConfig := database.RunConfig{
		ConfigFile:    *configFile,
		ConectionName: *conectionName,
		QueryName:     *queryName,
		Wrapper:       *wrapper,
		FirstOnly:     *firstOnly,
		KeyName:       *keyName,
		ValueName:     *valueName,
		Output:        *output,
		Format:        *format,
	}

	mapOfRows, err := database.RunQuery(runConfig)
	handleError(err)

	if runConfig.Output != "" && runConfig.Format == "csv" {
		err = io.SaveCSV(runConfig.Output, mapOfRows)
		handleError(err)
		return
	}

	var result interface{}

	if *firstOnly {
		var item map[string]string

		if len(mapOfRows) > 0 {
			item = mapOfRows[0]
		}

		if len(mapOfRows) > 0 {
			if *keyName != "" && *valueName != "" {
				itemResult := make(map[string]interface{})

				key := item[*keyName]
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
			result = make(map[string]string)
		}
	} else {
		if *keyName != "" && *valueName != "" {
			var newRows []map[string]string

			for _, row := range mapOfRows {
				key := fmt.Sprintf("%v", row[*keyName])
				value := row[*valueName]
				item := make(map[string]string)
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

	if *wrapper != "" {
		withWrapper := make(map[string]interface{})
		withWrapper[*wrapper] = result

		res, err := json.Marshal(withWrapper)
		handleError(err)
		fmt.Println(string(res))
	} else {
		res, err := json.Marshal(result)
		handleError(err)
		fmt.Println(string(res))
	}
}
