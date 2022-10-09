package core

import (
	"fmt"
	"franciscoperez.dev/gosqltojson/internal/config"
	"franciscoperez.dev/gosqltojson/internal/database"
)

func RunQuery(runConfig config.RunConfig, queryParamsFlags []string) (interface{}, error) {
	return RunQueryWithFn(runConfig, queryParamsFlags, database.RunQuery)
}

func RunQueryWithFn(runConfig config.RunConfig, queryParamsFlags []string, queryRowsFn func(runConfig config.RunConfig, queryParamsFlags []string) ([]map[string]interface{}, error)) (interface{}, error) {
	listMapOfRows, err := queryRowsFn(runConfig, queryParamsFlags)

	if err != nil {
		return nil, err
	}

	var result interface{}

	if runConfig.FirstOnly {
		if len(listMapOfRows) > 0 {
			row := listMapOfRows[0]
			singleResult := make(map[string]interface{})

			if runConfig.KeyName != "" && runConfig.ValueName != "" {
				newKey := fmt.Sprint(getKeyValueOrDefault(row, runConfig.KeyName, runConfig.KeyName))
				newValue := getKeyValueOrDefault(row, runConfig.ValueName, nil)

				singleResult[newKey] = newValue
				result = singleResult
			} else if runConfig.KeyName != "" {
				//Single value by key
				result = getKeyValueOrFirstValue(row, runConfig.KeyName)
			} else {
				singleResult = row
				result = singleResult
			}
		} else {
			//No data
			if runConfig.KeyName != "" && runConfig.ValueName == "" {
				result = ""
			} else {
				result = make(map[string]interface{})
			}
		}
	} else {
		//No first only
		var newResult []interface{}

		for _, row := range listMapOfRows {
			if runConfig.KeyName != "" {
				if runConfig.ValueName != "" {
					keyValueResult := make(map[string]interface{})

					newKey := fmt.Sprint(getKeyValueOrDefault(row, runConfig.KeyName, runConfig.KeyName))
					newValue := getKeyValueOrDefault(row, runConfig.ValueName, nil)

					keyValueResult[newKey] = newValue
					newResult = append(newResult, keyValueResult)
				} else {
					newResult = append(newResult, row[runConfig.KeyName])
				}
			} else {
				newResult = append(newResult, row)
			}
		}

		result = newResult
	}

	return result, nil
}

func getKeyValueOrFirstValue(item map[string]interface{}, keyName string) interface{} {
	//Single value by key
	if value, mapContainsKey := item[keyName]; mapContainsKey {
		return value
	} else {
		//First value, does not matter the key
		for _, v := range item {
			return v
		}
	}

	return nil
}

func getKeyValueOrDefault(item map[string]interface{}, keyName string, defaultValue interface{}) interface{} {
	if value, mapContainsKey := item[keyName]; mapContainsKey {
		return value
	}

	return defaultValue
}
