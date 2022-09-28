package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"franciscoperez.dev/gosqltojson/core"
	"franciscoperez.dev/gosqltojson/database"
	"franciscoperez.dev/gosqltojson/formats"
	"franciscoperez.dev/gosqltojson/io"
	"log"
	"reflect"
)

var (
	configFile     = flag.String("config", "", "Config file like config.json. Default location USER_HOME/.gosql2json/config.json")
	connectionName = flag.String("name", "default", "Database connection name in config file. Default default")
	queryName      = flag.String("query", "default", "Query name in config file. Default default")
	wrapper        = flag.String("wrapper", "", "Use an extra object as a wrapper")
	firstOnly      = flag.Bool("first", false, "Get first row only. Default false")
	keyName        = flag.String("key", "", "Field name used to compute key name")
	valueName      = flag.String("value", "", "Field name used to compute value for key")
	output         = flag.String("output", "", "file name to write the output. example: output.csv, output.json")
	format         = flag.String("format", "json", "Format to write the output. Default json")
	jsonkeys       = flag.String("jsonkeys", "", "Comma separated value to specify column names holding a json column/value")
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	flag.Parse()

	runConfig := core.RunConfig{
		ConfigFile:     *configFile,
		ConnectionName: *connectionName,
		QueryName:      *queryName,
		Wrapper:        *wrapper,
		FirstOnly:      *firstOnly,
		KeyName:        *keyName,
		ValueName:      *valueName,
		Output:         *output,
		Format:         *format,
		JsonKeys:       *jsonkeys,
	}

	listMapOfRows, err := database.RunQuery(runConfig, flag.Args(), database.NewDBConn)
	handleError(err)

	var result interface{}

	if *firstOnly {
		if len(listMapOfRows) > 0 {
			item := listMapOfRows[0]
			singleResult := make(map[string]interface{})

			if *keyName != "" && *valueName != "" {
				newKey := fmt.Sprint(item[*keyName])
				newValue := item[*valueName]

				singleResult[newKey] = newValue
				result = singleResult
			} else if *keyName != "" {
				//Single value by key
				if value, isMapContainsKey := item[*keyName]; isMapContainsKey {
					result = value
				} else {
					//First value, does not matter the key
					for _, v := range item {
						result = v
						break
					}
				}
			} else {
				singleResult = item
				result = singleResult
			}
		} else {
			//No data
			if *keyName != "" && *valueName == "" {
				result = ""
			} else {
				result = make(map[string]interface{})
			}
		}
	} else {
		//No first only
		var newResult []interface{}

		for _, row := range listMapOfRows {
			if *keyName != "" {
				if *valueName != "" {
					keyValueResult := make(map[string]interface{})

					keyValueResult[fmt.Sprint(row[*keyName])] = row[*valueName]
					newResult = append(newResult, keyValueResult)
				} else {
					newResult = append(newResult, row[*keyName])
				}
			} else {
				newResult = append(newResult, row)
			}
		}

		result = newResult
	}

	supportCSV := true

	withWrapper := make(map[string]interface{})

	if *wrapper != "" && runConfig.Format != "csv" {
		withWrapper[*wrapper] = result
		supportCSV = false
	}

	if runConfig.Output != "" {
		if supportCSV && runConfig.Format == "csv" {
			var csvList []map[string]interface{}

			if reflect.ValueOf(result).Kind() == reflect.Map {
				csvList = append(csvList, result.(map[string]interface{}))
			} else if reflect.ValueOf(result).Kind() == reflect.Slice {
				for _, item := range result.([]interface{}) {
					if reflect.ValueOf(item).Kind() == reflect.Map {
						csvList = append(csvList, item.(map[string]interface{}))
					} else {
						var newRow map[string]interface{}

						if *keyName != "" {
							newRow[*keyName] = item
						} else {
							newRow["key"] = item
						}
					}
				}
			} else {
				var firstRow map[string]interface{}

				if *keyName != "" {
					firstRow[*keyName] = result
				} else {
					firstRow["key"] = result
				}
				csvList = append(csvList, firstRow)
			}

			fileName := runConfig.GetOutputFileName()

			err = io.SaveCSV(fileName, csvList)
		} else if runConfig.Format == "json" {
			fileName := runConfig.GetOutputFileName()

			if *wrapper != "" {
				err = io.SaveJSON(fileName, withWrapper)
			} else {
				err = io.SaveJSON(fileName, result)
			}
		}

		handleError(err)
		return
	} else {
		if *wrapper != "" {
			res, err := json.Marshal(withWrapper)
			handleError(err)
			fmt.Println(string(res))
		} else {
			switch result.(type) {
			case string:
				fmt.Println(result)
				break
			default:
				res, err := formats.ToJsonString(result)
				handleError(err)
				fmt.Println(res)
			}
		}
	}
}
