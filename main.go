package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"franciscoperez.dev/gosqltojson/internal/config"
	"franciscoperez.dev/gosqltojson/internal/core"
	"franciscoperez.dev/gosqltojson/internal/formats"
	"franciscoperez.dev/gosqltojson/internal/io"
	"log"
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

	runConfig := config.RunConfig{
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

	result, err := core.RunQuery(runConfig, flag.Args())
	handleError(err)

	supportCSV := true

	withWrapper := make(map[string]interface{})

	if runConfig.Wrapper != "" && runConfig.Format != "csv" {
		withWrapper[runConfig.Wrapper] = result
		supportCSV = false
	}

	if runConfig.Output != "" {
		if supportCSV && runConfig.Format == "csv" {
			csvList := formats.AnyToListOfMaps(result, runConfig.KeyName)

			fileName := runConfig.GetOutputFileName()

			err = io.SaveCSV(fileName, csvList)
		} else if runConfig.Format == "json" {
			fileName := runConfig.GetOutputFileName()

			if runConfig.Wrapper != "" {
				err = io.SaveJSON(fileName, withWrapper)
			} else {
				err = io.SaveJSON(fileName, result)
			}
		}

		handleError(err)
		return
	} else {
		if runConfig.Wrapper != "" {
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
