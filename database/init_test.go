package database

import (
	"fmt"
	"franciscoperez.dev/gosqltojson/core"
	"testing"
)

type FakeDBConn struct {
}

func (coon FakeDBConn) Close() error {
	return nil
}

func (coon FakeDBConn) Query(queryRaw string, parameters map[string]interface{}) ([]map[string]interface{}, error) {
	rows := []map[string]interface{}{
		{
			"a": 1,
			"b": "2",
		},
		{
			"a": 10,
			"b": "20",
		},
	}
	return rows, nil
}

func TestRunQuery(t *testing.T) {
	runConfig := core.RunConfig{}

	var queryParamFlags []string

	rows, err := RunQuery(runConfig, queryParamFlags, createFakeDBConn)

	if err != nil || len(rows) == 0 {
		t.Log("RunQuery should not return error. But got: ", rows, err)
		t.Fail()
	}

	runConfigConfigFileDoesNotExists := core.RunConfig{
		ConfigFile: "config_file_does_not_exists.json",
	}

	rows, err = RunQuery(runConfigConfigFileDoesNotExists, queryParamFlags, createFakeDBConn)

	if err == nil {
		t.Log("RunQuery with config: config_file_does_not_exists.json that does not exist should return error. But got: ", rows)
		t.Fail()
	}

	runConfigQueryFileDoesNotExists := core.RunConfig{
		QueryName: "@query_file_does_not_exists.json",
	}

	rows, err = RunQuery(runConfigQueryFileDoesNotExists, queryParamFlags, createFakeDBConn)

	if err == nil {
		t.Log("RunQuery with query: query_file_does_not_exists.json that does not exist should return error. But got: ", rows)
		t.Fail()
	}

	runConfigQueryErrorToConnectToDB := core.RunConfig{}

	rows, err = RunQuery(runConfigQueryErrorToConnectToDB, queryParamFlags, createDBConnError)

	if err == nil {
		t.Log("RunQuery with error to connect to the database should return error. But got: ", rows)
		t.Fail()
	}

	runConfigQueryMissingParameters := core.RunConfig{
		QueryName: "SELECT 1 AS a WHERE @a > 10",
	}

	rows, err = RunQuery(runConfigQueryMissingParameters, queryParamFlags, createFakeDBConn)

	if err == nil {
		t.Log("RunQuery for a query with missing parameters should return error. But got: ", rows)
		t.Fail()
	}
}

func createFakeDBConn(_ *core.ConfigFile, _ string) (DBConn, error) {
	return FakeDBConn{}, nil
}

func createDBConnError(_ *core.ConfigFile, _ string) (DBConn, error) {
	return nil, fmt.Errorf("fake error for testing")
}
