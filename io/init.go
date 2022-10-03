package io

import (
	"encoding/csv"
	"fmt"
	"franciscoperez.dev/gosqltojson/formats"
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

	writer := csv.NewWriter(f)

	data, err := formats.ToCSV(dataMap)

	if err != nil {
		return err
	}

	err = writer.WriteAll(data)
	if err != nil {
		return err
	}

	err = f.Close()

	return err
}

func SaveJSON(filename string, data interface{}) error {
	if filename == "" {
		return fmt.Errorf("output file name is required")
	}

	fileBytes, err := formats.ToJson(data)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, fileBytes, 0644)
}
