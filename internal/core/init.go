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
			item := listMapOfRows[0]
			singleResult := make(map[string]interface{})

			if runConfig.KeyName != "" && runConfig.ValueName != "" {
				newKey := fmt.Sprint(item[runConfig.KeyName])
				newValue := item[runConfig.ValueName]

				singleResult[newKey] = newValue
				result = singleResult
			} else if runConfig.KeyName != "" {
				//Single value by key
				if value, isMapContainsKey := item[runConfig.KeyName]; isMapContainsKey {
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

					keyValueResult[fmt.Sprint(row[runConfig.KeyName])] = row[runConfig.ValueName]
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
