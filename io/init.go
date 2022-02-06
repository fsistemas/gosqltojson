package io

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
)

func LoadFileContentAsString(file string) (string, error) {
	if file == "" {
		return "", fmt.Errorf("file is required")
	}

	contentBytes, err := ioutil.ReadFile(file)

	if err != nil {
		return "", err
	}

	return string(contentBytes), nil
}

func GetEnvOrDefault(envName, defaultValue string) string {
	v := os.Getenv(envName)

	if v == "" {
		return defaultValue
	}

	return v
}

func SaveCSV(file string, dataMap []map[string]string) error {
	if file == "" {
		return fmt.Errorf("output file name is required")
	}

	f, err := os.Create(file)
	if err != nil {
		return err
	}

	defer f.Close()

	writer := csv.NewWriter(f)

	var data = [][]string{}

	if len(dataMap) == 0 {
		return nil
	}

	firstRow := dataMap[0]
	headers := []string{}

	for key, _ := range firstRow {
		headers = append(headers, key)
	}

	data = append(data, headers)

	for _, row := range dataMap {
		rowAsList := []string{}
		for _, field := range headers {
			rowAsList = append(rowAsList, row[field])
		}
		data = append(data, rowAsList)
	}

	err = writer.WriteAll(data)
	if err != nil {
		return err
	}

	return nil
}
