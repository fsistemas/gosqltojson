package formats

import (
	"fmt"
	"reflect"
)

func AnyToListOfMaps(result interface{}, keyName string) []map[string]interface{} {
	var csvList []map[string]interface{}

	if reflect.ValueOf(result).Kind() == reflect.Map {
		csvList = append(csvList, result.(map[string]interface{}))
	} else if reflect.ValueOf(result).Kind() == reflect.Slice {
		for _, item := range result.([]interface{}) {
			if reflect.ValueOf(item).Kind() == reflect.Map {
				csvList = append(csvList, item.(map[string]interface{}))
			} else {
				var newRow map[string]interface{}

				if keyName != "" {
					newRow[keyName] = item
				} else {
					newRow["key"] = item
				}
			}
		}
	} else {
		var firstRow map[string]interface{}

		if keyName != "" {
			firstRow[keyName] = result
		} else {
			firstRow["key"] = result
		}
		csvList = append(csvList, firstRow)
	}

	return csvList
}

func MapToCSVWithHeaders(dataMap []map[string]interface{}) ([][]string, error) {
	var data [][]string

	if len(dataMap) == 0 {
		return [][]string{}, nil
	}

	firstRow := dataMap[0]
	var headers []string

	for key := range firstRow {
		headers = append(headers, key)
	}

	data = append(data, headers)

	for _, row := range dataMap {
		var rowAsList []string
		for _, field := range headers {
			rowAsList = append(rowAsList, fmt.Sprint(row[field]))
		}
		data = append(data, rowAsList)
	}

	return data, nil
}
