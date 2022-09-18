package database

import (
	"franciscoperez.dev/gosqltojson/core"
	"franciscoperez.dev/gosqltojson/parameter"
)

func RunQuery(runConfig core.RunConfig, queryParamsFlags []string, connectToDB func(*core.ConfigFile, string) (DBConn, error)) ([]map[string]interface{}, error) {
	configFile, err := core.NewConfigFile(runConfig.ConfigFile)

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

	return dbConn.Query(queryRaw, queryParams)
}
