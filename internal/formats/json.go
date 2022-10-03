package formats

import "encoding/json"

func ToJson(data interface{}) ([]byte, error) {
	return json.MarshalIndent(data, "", " ")
}

func ToJsonString(data interface{}) (string, error) {
	fileBytes, err := ToJson(data)

	if err != nil {
		return "", err
	}

	return string(fileBytes), nil
}

func ToJsonMap(data []byte) (interface{}, error) {
	var newMap interface{}

	err := json.Unmarshal(data, &newMap)

	if err != nil {
		return nil, err
	}

	return newMap, nil
}
