package core

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"franciscoperez.dev/gosqltojson/io"
)

type Queries map[string]string
type Conections map[string]string

type ConfigFile struct {
	Conections Conections `json:"conections"`
	Queries    Queries    `json:"queries"`
}

type RunConfig struct {
	ConfigFile    string
	ConectionName string
	QueryName     string
	Wrapper       string
	FirstOnly     bool
	KeyName       string
	ValueName     string
	Output        string
	Format        string
}

func (config *ConfigFile) GetConection(conectionName string) string {
	rawQuery := config.Conections[conectionName]

	if rawQuery == "" {
		return conectionName
	}

	return rawQuery
}

func (config *ConfigFile) GetQuery(queryName string) string {
	rawQuery := config.Queries[queryName]

	if rawQuery == "" {
		rawQuery = queryName
	}

	if strings.HasPrefix(rawQuery, "@") {
		var err error
		rawQuery, err = io.LoadFileContentAsString(rawQuery[1:])
		if err != nil {
			panic(err.Error())
		}
	}

	return rawQuery
}

func NewConfigFile(file string) (*ConfigFile, error) {
	config := new(ConfigFile)
	if file == "" {
		return nil, fmt.Errorf("file is required")
	}

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	err = json.NewDecoder(f).Decode(&config)

	if err != nil {
		return nil, err
	}

	return config, nil
}

type GoSqlTo struct {
}
