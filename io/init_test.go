package io

import (
	"os"
	"testing"
)

func TestFileExists(t *testing.T) {
	if FileExists(".") {
		t.Log("Current directory . exists but it is not a file")
		t.Fail()
	}

	if FileExists("file_does_not_exist.ext") {
		t.Log("file_does_not_exist.ext should not exist")
		t.Fail()
	}

	if !FileExists("../testdata/file_exists.txt") {
		t.Log("testdata/file_exists.txt should exist")
		t.Fail()
	}
}

func TestLoadFileContentAsString(t *testing.T) {
	_, err := LoadFileContentAsString("../testdata/file_exists.txt")

	if err != nil {
		t.Log("LoadFileContentAsString should be able to read ../testdata/file_exists.txt", err)
		t.Fail()
	}

	_, err = LoadFileContentAsString("")

	if err == nil {
		t.Log("LoadFileContentAsString with empty filename should be fail")
		t.Fail()
	}

	_, err = LoadFileContentAsString("../testdata/file_does_not_exists.txt")

	if err == nil {
		t.Log("LoadFileContentAsString with ../testdata/file_does_not_exists.txt should be fail")
		t.Fail()
	}
}

func TestGetEnvOrDefault(t *testing.T) {
	if os.Getenv("HOME") != GetEnvOrDefault("HOME", "default") {
		t.Log("In theory HOME variable always exist")
		t.Fail()
	}

	if "default" != GetEnvOrDefault("ENV_VAR_DOES_NOT_EXIST", "default") {
		t.Log("The variable ENV_VAR_DOES_NOT_EXIST does not exist. Excepted to fallback to 'default'")
		t.Fail()
	}
}
