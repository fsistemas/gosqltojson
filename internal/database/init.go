package database

import (
	"fmt"
	"franciscoperez.dev/gosqltojson/internal/config"
	"franciscoperez.dev/gosqltojson/internal/formats"
	"franciscoperez.dev/gosqltojson/internal/parameter"
	"strings"
)

func RunQuery(runConfig config.RunConfig, queryParamsFlags []string) ([]map[string]interface{}, error) {
	return runQuery(runConfig, queryParamsFlags, NewDBConn)
}

func runQuery(runConfig config.RunConfig, queryParamsFlags []string, connectToDB func(*config.ConfigFile, string) (DBConn, error)) ([]map[string]interface{}, error) {
	configFile, err := config.NewConfigFile(runConfig.ConfigFile)

	if err != nil {
		return nil, err
	}

	queryRaw, err := configFile.GetQuery(runConfig.QueryName)

	if err != nil {
		return nil, err
	}

	dbConn, err := connectToDB(configFile, runConfig.ConnectionName)

	if err != nil {
		return nil, err
	}

	defer dbConn.Close()

	queryParams, err := parameter.ParseQueryParams(queryRaw, queryParamsFlags)
	if err != nil {
		return nil, err
	}

	rows, err := dbConn.Query(queryRaw, queryParams)

	if err != nil {
		return nil, err
	}

	return parseRowsJsonColumns(runConfig, rows)
}

func parseRowsJsonColumns(runConfig config.RunConfig, rows []map[string]interface{}) ([]map[string]interface{}, error) {
	if runConfig.JsonKeys == "" {
		return rows, nil
	}

	var newRows []map[string]interface{}

	jsonKeyListMap := jsonKeysToMap(runConfig.JsonKeys)

	for _, row := range rows {
		newRow, err := parseRowJsonColumns(row, jsonKeyListMap)

		if err != nil {
			return nil, err
		}

		newRows = append(newRows, newRow)
	}

	return newRows, nil
}

func parseRowJsonColumns(row map[string]interface{}, jsonKeyListMap map[string]bool) (map[string]interface{}, error) {
	newRow := map[string]interface{}{}

	for column, value := range row {
		if _, isMapContainsKey := jsonKeyListMap[column]; isMapContainsKey {

			valueByes := []byte(fmt.Sprint(value))

			newValue, err := formats.ToJsonMap(valueByes)

			if err != nil {
				return nil, err
			}

			newRow[column] = newValue
		} else {
			newRow[column] = value
		}
	}

	return newRow, nil
}

func jsonKeysToMap(csvNames string) map[string]bool {
	jsonKeyListMap := map[string]bool{}

	jsonKeyList := strings.Split(csvNames, ",")

	for _, jsonKeyName := range jsonKeyList {
		jsonKeyListMap[strings.Trim(jsonKeyName, " ")] = true
	}

	return jsonKeyListMap
}
