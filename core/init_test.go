package core

import "testing"

func TestRunConfig_GetOutputFileName(t *testing.T) {
	runConfigWithoutExt := createRunConfig("my_output")

	if "my_output.json" != runConfigWithoutExt.GetOutputFileName() {
		t.Log("Expected my_output.json")
		t.Fail()
	}

	runConfigWithExt := createRunConfig("my_output.json")

	if "my_output.json" != runConfigWithExt.GetOutputFileName() {
		t.Log("Expected my_output.json")
		t.Fail()
	}
}

func TestConfigFile_GetConnection(t *testing.T) {
	configFile := createConfigFile()

	connString := configFile.GetConnection("CONN_STRING")
	if "CONN_STRING" != connString {
		t.Log("Expected Connection: CONN_STRING")
		t.Fail()
	}

	connString = configFile.GetConnection("CONN_STRING_EXIST")
	if "ANY_STRING" != connString {
		t.Log("Expected Connection: ANY_STRING")
		t.Fail()
	}
}

func TestConfigFile_GetQuery(t *testing.T) {
	configFile := createConfigFile()

	queryString, err := configFile.GetQuery("MY_QUERY")
	if err != nil || "MY_QUERY" != queryString {
		t.Log("Expected Query: MY_QUERY", err)
		t.Fail()
	}

	queryString, err = configFile.GetQuery("QUERY_OK")
	if err != nil || "SELECT 1 AS ok" != queryString {
		t.Log("Expected Query: SELECT 1 AS ok", err)
		t.Fail()
	}

	queryString, err = configFile.GetQuery("@../testdata/test_query.sql")
	if err != nil || "SELECT 1 AS a, 2 AS b" != queryString {
		t.Log("Expected Query: SELECT 1 AS a, 2 as b. But was: ", queryString, err)
		t.Fail()
	}

	queryString, err = configFile.GetQuery("@../testdata/query_does_not_exist.sql")
	if err == nil {
		t.Log("Expected Query: nil with error because file query does not exist", queryString)
		t.Fail()
	}
}

func TestNewConfigFile(t *testing.T) {
	configFile, err := NewConfigFile("../testdata/test_config.json")

	if err != nil {
		t.Log("Expected configFile without error. But was: ", configFile, err)
		t.Fail()
	}

	configFile, err = NewConfigFile("")
	if err != nil {
		t.Log("Expected fallback configFile without error. But was: ", configFile, err)
		t.Fail()
	}

	configFile, err = NewConfigFile("../testdata/config_does_not_exist.json")

	if err == nil {
		t.Log("Expected error, because ../testdata/config_does_not_exist.json does not exist. But got: ", configFile)
		t.Fail()
	}

	configFile, err = NewConfigFile("../testdata/file_exists.txt")

	if err == nil {
		t.Log("Expected error, because ../testdata/file_exists.txt is not a valid config file. But got: ", configFile)
		t.Fail()
	}
}

func createRunConfig(output string) RunConfig {
	return RunConfig{
		ConfigFile:     "test_config.json",
		Format:         "json",
		ConnectionName: "default",
		QueryName:      "default",
		KeyName:        "",
		ValueName:      "",
		FirstOnly:      false,
		Output:         output,
		Wrapper:        "",
	}
}

func createConfigFile() ConfigFile {
	return ConfigFile{
		Connections: Connections{
			"CONN_STRING_EXIST": "ANY_STRING",
		},
		Queries: Queries{
			"QUERY_OK": "SELECT 1 AS ok",
		},
	}
}
