package io

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func FileExists(fileName string) bool {
	info, err := os.Stat(fileName)

	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

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

func SaveCSV(file string, dataMap []map[string]interface{}) error {
	if file == "" {
		return fmt.Errorf("output file name is required")
	}

	f, err := os.Create(file)
	if err != nil {
		return err
	}

	defer f.Close()

	writer := csv.NewWriter(f)

	var data [][]string

	if len(dataMap) == 0 {
		return nil
	}

	firstRow := dataMap[0]
	var headers []string

	for key, _ := range firstRow {
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

	err = writer.WriteAll(data)
	if err != nil {
		return err
	}

	return nil
}

func SaveJSON(filename string, data interface{}) error {
	if filename == "" {
		return fmt.Errorf("output file name is required")
	}

	fileBytes, err := json.MarshalIndent(data, "", " ")

	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, fileBytes, 0644)
}
