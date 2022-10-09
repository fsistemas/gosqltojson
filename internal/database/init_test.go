package database

import (
	"fmt"
	"franciscoperez.dev/gosqltojson/internal/config"
	"testing"
)

func TestRunQuery(t *testing.T) {
	runConfig := config.RunConfig{}

	var queryParamFlags []string

	rows, err := runQuery(runConfig, queryParamFlags, createFakeDBConn)

	if err != nil || len(rows) == 0 {
		t.Log("runQuery should not return error. But got: ", rows, err)
		t.Fail()
	}

	runConfigConfigFileDoesNotExists := config.RunConfig{
		ConfigFile: "config_file_does_not_exists.json",
	}

	rows, err = runQuery(runConfigConfigFileDoesNotExists, queryParamFlags, createFakeDBConn)

	if err == nil {
		t.Log("runQuery with config: config_file_does_not_exists.json that does not exist should return error. But got: ", rows)
		t.Fail()
	}

	runConfigQueryFileDoesNotExists := config.RunConfig{
		QueryName: "@query_file_does_not_exists.json",
	}

	rows, err = runQuery(runConfigQueryFileDoesNotExists, queryParamFlags, createFakeDBConn)

	if err == nil {
		t.Log("runQuery with query: query_file_does_not_exists.json that does not exist should return error. But got: ", rows)
		t.Fail()
	}

	runConfigQueryErrorToConnectToDB := config.RunConfig{}

	rows, err = runQuery(runConfigQueryErrorToConnectToDB, queryParamFlags, createDBConnError)

	if err == nil {
		t.Log("runQuery with error to connect to the database should return error. But got: ", rows)
		t.Fail()
	}

	runConfigQueryMissingParameters := config.RunConfig{
		QueryName: "SELECT 1 AS a WHERE @a > 10",
	}

	rows, err = runQuery(runConfigQueryMissingParameters, queryParamFlags, createFakeDBConn)

	if err == nil {
		t.Log("runQuery for a query with missing parameters should return error. But got: ", rows)
		t.Fail()
	}

	runConfigQueryWithJsonColumns := config.RunConfig{
		JsonKeys:  "jsoncolumn",
		QueryName: "UN_USED_FAKE_QUERY",
	}

	rows, err = runQuery(runConfigQueryWithJsonColumns, queryParamFlags, createFakeDBConn)
}

type FakeDBConn struct {
}

func (coon FakeDBConn) Close() error {
	return nil
}

func (coon FakeDBConn) Query(queryRaw string, parameters map[string]interface{}) ([]map[string]interface{}, error) {
	rows := []map[string]interface{}{
		{
			"a":          1,
			"b":          "2",
			"jsoncolumn": "{ \"x\": 1, \"y\": \"z\" }",
		},
		{
			"a":          10,
			"b":          "20",
			"jsoncolumn": "{ \"x\": 10, \"y\": \"zz\" }",
		},
	}
	return rows, nil
}

func createFakeDBConn(_ *config.ConfigFile, _ string) (DBConn, error) {
	return FakeDBConn{}, nil
}

func createDBConnError(_ *config.ConfigFile, _ string) (DBConn, error) {
	return nil, fmt.Errorf("fake error for testing")
}
