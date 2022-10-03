package formats

import "testing"

func TestToJson(t *testing.T) {
	jsonBytes, err := ToJson([]string{"hello", "world"})

	if err != nil {
		t.Log("Error calling ToJson in an array.", err)
		t.Fail()
	}

	expectedJson := "[\n \"hello\",\n \"world\"\n]"

	if string(jsonBytes) != expectedJson {
		t.Log("Error calling ToJson in an array. expected: " + expectedJson + ", got: " + string(jsonBytes))
		t.Fail()
	}
}

func TestToJsonString(t *testing.T) {
	jsonString, err := ToJsonString([]string{"hello", "world"})

	if err != nil {
		t.Log("Error calling ToJsonString in an array.", err)
		t.Fail()
	}

	expectedJson := "[\n \"hello\",\n \"world\"\n]"

	if jsonString != expectedJson {
		t.Log("Error calling ToJsonString in an array. expected: " + expectedJson + ", got: " + jsonString)
		t.Fail()
	}
}

func TestToJsonMap(t *testing.T) {
	jsonMap, err := ToJsonMap([]byte("{\"a\": \"1\", \"b\": \"x\"}"))

	if err != nil {
		t.Log("Error calling ToJson in an array.", err)
		t.Fail()
	}

	expectedJson := map[string]interface{}{
		"a": "1",
		"b": "x",
	}

	jsonMapWithTypes := jsonMap.(map[string]interface{})

	if jsonMapWithTypes["a"] != expectedJson["a"] || jsonMapWithTypes["b"] != expectedJson["b"] {
		t.Log("Error calling ToJsonMap in an json []byte. expected JSON does not match the result")
		t.Fail()
	}
}
