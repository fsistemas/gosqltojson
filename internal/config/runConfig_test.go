package config

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
