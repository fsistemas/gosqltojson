package formats

import "fmt"

func ToCSV(dataMap []map[string]interface{}) ([][]string, error) {
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
