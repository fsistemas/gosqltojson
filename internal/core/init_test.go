package core

import (
	"fmt"
	"franciscoperez.dev/gosqltojson/internal/config"
	"reflect"
	"testing"
)

func TestRunQueryWithFnNoFlags(t *testing.T) {
	runConfig := config.RunConfig{}
	var queryParamsFlags []string

	result, err := RunQueryWithFn(runConfig, queryParamsFlags, queryRowsFn)

	if err != nil {
		t.Error("Expected result. No error", err)
	}

	if result == nil || reflect.TypeOf(result).Kind() != reflect.Slice {
		t.Error("Expected not null array as result. type: ", reflect.TypeOf(result).Kind())
	}

	rows := result.([]interface{})

	if len(rows) != 3 {
		t.Error("Expected 3 rows.")
	}
}

func TestRunQueryWithFnConfigKeyName(t *testing.T) {
	runConfig := config.RunConfig{
		KeyName: "a",
	}
	var queryParamsFlags []string

	result, err := RunQueryWithFn(runConfig, queryParamsFlags, queryRowsFn)

	if err != nil {
		t.Error("Expected result. No error", err)
	}

	if result == nil || reflect.TypeOf(result).Kind() != reflect.Slice {
		t.Error("Expected not null array/slice as result. type: ", reflect.TypeOf(result).Kind())
	}

	rows := result.([]interface{})

	if len(rows) != 3 {
		t.Error("Expected 3 rows.")
	}

	if reflect.TypeOf(rows[0]).Kind() != reflect.Int && rows[0].(int) != 1 {
		t.Error("Expected first row to be int value 1. But type is: ", reflect.TypeOf(rows[0]).Kind())
	}

	if reflect.TypeOf(rows[1]).Kind() != reflect.Int && rows[1].(int) != 10 {
		t.Error("Expected second row to be int value 10. But type is: ", reflect.TypeOf(rows[1]).Kind())
	}

	if reflect.TypeOf(rows[2]).Kind() != reflect.Int && rows[1].(int) != 100 {
		t.Error("Expected third row to be int value 100. But type is: ", reflect.TypeOf(rows[2]).Kind())
	}
}

func TestRunQueryWithFnConfigKeyNameValueName(t *testing.T) {
	runConfig := config.RunConfig{
		KeyName:   "b",
		ValueName: "a",
	}

	var queryParamsFlags []string

	result, err := RunQueryWithFn(runConfig, queryParamsFlags, queryRowsFn)

	if err != nil {
		t.Error("Expected result. No error", err)
	}

	if result == nil || reflect.TypeOf(result).Kind() != reflect.Slice {
		t.Error("Expected not null array/slice as result. type: ", reflect.TypeOf(result).Kind())
	}

	rows := result.([]interface{})

	if len(rows) != 3 {
		t.Error("Expected 3 rows.")
	}

	if reflect.TypeOf(rows[0]).Kind() != reflect.Map {
		t.Error("Expected first row to be map. But type is: ", reflect.TypeOf(rows[0]).Kind())
	}

	firstRow := rows[0].(map[string]interface{})

	if len(firstRow) != 1 || firstRow["x"] != 1 {
		t.Errorf("Expected first row to be map with a single key x. Size: %d %d", len(firstRow), firstRow["x"])
	}

	if reflect.TypeOf(rows[1]).Kind() != reflect.Map {
		t.Error("Expected second row to be map. But type is: ", reflect.TypeOf(rows[1]).Kind())
	}

	secondRow := rows[1].(map[string]interface{})

	if len(secondRow) != 1 || secondRow["y"] != 10 {
		t.Errorf("Expected second row to be map with a single key y. Size: %d %d", len(secondRow), secondRow["y"])
	}

	if reflect.TypeOf(rows[2]).Kind() != reflect.Map {
		t.Error("Expected third row to be map. But type is: ", reflect.TypeOf(rows[2]).Kind())
	}

	thirdRow := rows[2].(map[string]interface{})

	if len(thirdRow) != 1 || thirdRow["z"] != 100 {
		t.Errorf("Expected third row to be map with a single key z. Size: %d %s", len(thirdRow), thirdRow)
	}
}

func TestRunQueryWithFnConfigKeyNameDoesNotExistValueNameDoesNotExist(t *testing.T) {
	runConfig := config.RunConfig{
		KeyName:   "bx",
		ValueName: "ax",
	}

	var queryParamsFlags []string

	result, err := RunQueryWithFn(runConfig, queryParamsFlags, queryRowsFn)

	if err != nil {
		t.Error("Expected result. No error", err)
	}

	if result == nil || reflect.TypeOf(result).Kind() != reflect.Slice {
		t.Error("Expected not null array/slice as result. type: ", reflect.TypeOf(result).Kind())
	}

	rows := result.([]interface{})

	if len(rows) != 3 {
		t.Error("Expected 3 rows.")
	}

	if reflect.TypeOf(rows[0]).Kind() != reflect.Map {
		t.Error("Expected first row to be map. But type is: ", reflect.TypeOf(rows[0]).Kind())
	}

	firstRow := rows[0].(map[string]interface{})

	if len(firstRow) != 1 || firstRow["bx"] != nil {
		t.Errorf("Expected first row to be map with a single key bx. Size: %d %s", len(firstRow), firstRow)
	}

	if reflect.TypeOf(rows[1]).Kind() != reflect.Map {
		t.Error("Expected second row to be map. But type is: ", reflect.TypeOf(rows[1]).Kind())
	}

	secondRow := rows[1].(map[string]interface{})

	if len(secondRow) != 1 || secondRow["x"] != nil {
		t.Errorf("Expected second row to be map with a single key bx. Size: %d %s", len(secondRow), secondRow)
	}

	if reflect.TypeOf(rows[2]).Kind() != reflect.Map {
		t.Error("Expected third row to be map. But type is: ", reflect.TypeOf(rows[2]).Kind())
	}

	thirdRow := rows[2].(map[string]interface{})

	if len(thirdRow) != 1 || thirdRow["bx"] != nil {
		t.Errorf("Expected third row to be map with a single key bx. Size: %d %s", len(thirdRow), thirdRow)
	}
}

func TestRunQueryWithFnFirstOnly(t *testing.T) {
	runConfig := config.RunConfig{
		FirstOnly: true,
	}

	var queryParamsFlags []string

	result, err := RunQueryWithFn(runConfig, queryParamsFlags, queryRowsFn)

	if err != nil {
		t.Error("Expected result. No error", err)
	}

	if result == nil || reflect.TypeOf(result).Kind() != reflect.Map {
		t.Error("Expected not null map as result. type: ", reflect.TypeOf(result).Kind())
	}

	row := result.(map[string]interface{})

	if len(row) != 2 {
		t.Error("Expected map with 2 keys.")
	}

	if row["a"] != 1 || row["b"] != "x" {
		t.Errorf("Expected map with 2 keys. a = 1, b = 'x'. Got: %s", row)
	}
}

func TestRunQueryWithFnFirstOnlyKeyName(t *testing.T) {
	runConfig := config.RunConfig{
		FirstOnly: true,
		KeyName:   "a",
	}

	var queryParamsFlags []string

	result, err := RunQueryWithFn(runConfig, queryParamsFlags, queryRowsFn)

	if err != nil {
		t.Error("Expected result. No error", err)
	}

	if result == nil || reflect.TypeOf(result).Kind() != reflect.Int {
		t.Error("Expected not null int as result. type: ", reflect.TypeOf(result).Kind())
	}

	if result.(int) != 1 {
		t.Errorf("Expected int 1. Got: %s", result)
	}
}

func TestRunQueryWithFnFirstOnlyKeyNameDoesNotExist(t *testing.T) {
	runConfig := config.RunConfig{
		FirstOnly: true,
		KeyName:   "ax",
	}

	var queryParamsFlags []string

	result, err := RunQueryWithFn(runConfig, queryParamsFlags, queryRowsFn)

	if err != nil {
		t.Error("Expected result. No error", err)
	}

	if result == nil || reflect.TypeOf(result).Kind() != reflect.Int {
		t.Error("Expected not null int as result. type: ", reflect.TypeOf(result).Kind())
	}

	if result.(int) != 1 {
		t.Errorf("Expected int 1(first value in first row). Got: %s", result)
	}
}

func TestRunQueryWithFnFirstOnlyKeyNameDoesNotExistValueNameDoesNotExist(t *testing.T) {
	runConfig := config.RunConfig{
		FirstOnly: true,
		KeyName:   "ax",
		ValueName: "bx",
	}

	var queryParamsFlags []string

	result, err := RunQueryWithFn(runConfig, queryParamsFlags, queryRowsFn)

	if err != nil {
		t.Error("Expected result. No error", err)
	}

	if result == nil || reflect.TypeOf(result).Kind() != reflect.Map {
		t.Error("Expected not null map as result. type: ", reflect.TypeOf(result).Kind())
	}

	rowMap := result.(map[string]interface{})

	if len(rowMap) != 1 || rowMap["bx"] != nil {
		t.Errorf("Expected map with single key. key: bx, value: nil. Got: %s", rowMap)
	}
}

func TestRunQueryWithFnFirstOnlyKeyNameValueName(t *testing.T) {
	runConfig := config.RunConfig{
		FirstOnly: true,
		KeyName:   "b",
		ValueName: "a",
	}

	var queryParamsFlags []string

	result, err := RunQueryWithFn(runConfig, queryParamsFlags, queryRowsFn)

	if err != nil {
		t.Error("Expected result. No error", err)
	}

	if result == nil || reflect.TypeOf(result).Kind() != reflect.Map {
		t.Error("Expected not null Map as result. type: ", reflect.TypeOf(result).Kind())
	}

	rowMap := result.(map[string]interface{})

	if len(rowMap) != 1 {
		t.Errorf("Expected map with a key. Got map of: %d keys", len(rowMap))
	}

	if rowMap["x"] != 1 {
		t.Errorf("Expected key x to be 1. Got: %s", rowMap["x"])
	}
}

func TestRunQueryWithFnFirstOnlyKeyNameNoRows(t *testing.T) {
	runConfig := config.RunConfig{
		FirstOnly: true,
		KeyName:   "b",
	}

	var queryParamsFlags []string

	result, err := RunQueryWithFn(runConfig, queryParamsFlags, queryRowsFnNoRows)

	if err != nil {
		t.Error("Expected result. No error", err)
	}

	if result == nil || reflect.TypeOf(result).Kind() != reflect.String || result.(string) != "" {
		t.Errorf("Expected empty string as result. got type: %s value: %s", reflect.TypeOf(result).Kind(), result)
	}
}

func TestRunQueryWithFnFirstOnlyKeyNameValueNameNoRows(t *testing.T) {
	runConfig := config.RunConfig{
		FirstOnly: true,
		KeyName:   "b",
		ValueName: "a",
	}

	var queryParamsFlags []string

	result, err := RunQueryWithFn(runConfig, queryParamsFlags, queryRowsFnNoRows)

	if err != nil {
		t.Error("Expected result. No error", err)
	}

	if result == nil || reflect.TypeOf(result).Kind() != reflect.Map || len(result.(map[string]interface{})) != 0 {
		t.Error("Expected empty Map as result. type: ", reflect.TypeOf(result).Kind(), result)
	}
}

func TestRunQueryWithFnError(t *testing.T) {
	runConfig := config.RunConfig{}

	var queryParamsFlags []string

	_, err := RunQueryWithFn(runConfig, queryParamsFlags, queryRowsFnError)

	if err == nil {
		t.Error("Expected a fake error. But got nil")
	}
}

func queryRowsFn(runConfig config.RunConfig, queryParamsFlags []string) ([]map[string]interface{}, error) {
	rows := []map[string]interface{}{
		{
			"a": 1,
			"b": "x",
		},
		{
			"a": 10,
			"b": "y",
		},
		{
			"a": 100,
			"b": "z",
		},
	}

	return rows, nil
}

func queryRowsFnNoRows(runConfig config.RunConfig, queryParamsFlags []string) ([]map[string]interface{}, error) {
	var rows []map[string]interface{}

	return rows, nil
}

func queryRowsFnError(runConfig config.RunConfig, queryParamsFlags []string) ([]map[string]interface{}, error) {
	return nil, fmt.Errorf("fake error")
}
