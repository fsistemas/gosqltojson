package config

import (
	"encoding/json"
	"franciscoperez.dev/gosqltojson/internal/io"
	"os"
	"strings"
)

type Queries map[string]string
type Connections map[string]string

type ConfigFile struct {
	Connections Connections `json:"connections"`
	Queries     Queries     `json:"queries"`
}

func (config *ConfigFile) GetConnection(connectionName string) string {
	rawConnection := config.Connections[connectionName]

	if rawConnection == "" {
		return connectionName
	}

	return rawConnection
}

func (config *ConfigFile) GetQuery(queryName string) (string, error) {
	rawQuery := config.Queries[queryName]

	if rawQuery == "" {
		rawQuery = queryName
	}

	if strings.HasPrefix(rawQuery, "@") {
		var err error
		rawQuery, err = io.LoadFileContentAsString(rawQuery[1:])
		if err != nil {
			return "", err
		}
	}

	return rawQuery, nil
}

func NewConfigFile(file string) (*ConfigFile, error) {
	configFileName := file
	config := new(ConfigFile)

	if configFileName == "" {
		configFileName = io.GetEnvOrDefault("HOME", ".") + "/.gosql2json/config.json"

		if !io.FileExists(configFileName) {
			//Fallback to config in current directory
			configFileName = "./config.json"

			if !io.FileExists(configFileName) {
				//Fallback to allow user to test the tool without any configuration
				config.Connections = Connections{
					"default": "sqlite+test.db",
				}

				config.Queries = Queries{
					"default": "SELECT 1 AS a, 2 AS b",
				}

				return config, nil
			}
		}
	}

	f, err := os.Open(configFileName)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(f).Decode(&config)

	if err != nil {
		return nil, err
	}

	err = f.Close()

	if err != nil {
		return nil, err
	}

	return config, nil
}
