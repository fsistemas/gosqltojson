package formats

import "testing"

func TestToCSVWithNullOrEmptyList(t *testing.T) {
	csv, err := ToCSV(nil)

	if err != nil {
		t.Log("Unexpected error")
		t.Fail()
	}

	if csv == nil || len(csv) != 0 {
		t.Log("Expected empty list")
		t.Fail()
	}

	csv, err = ToCSV([]map[string]interface{}{})

	if err != nil {
		t.Log("Unexpected error")
		t.Fail()
	}

	if csv == nil || len(csv) != 0 {
		t.Log("Expected empty list")
		t.Fail()
	}
}

func TestToCSV(t *testing.T) {
	dataMap := []map[string]interface{}{
		{
			"a": "1",
			"b": "2",
		},
		{
			"a": "10",
			"b": "20",
		},
	}

	csv, err := ToCSV(dataMap)

	if err != nil {
		t.Log("Unexpected error")
		t.Fail()
	}

	if csv == nil || len(csv) != 3 {
		t.Log("Expected list with 2 elements => 1 header + 2 data")
		t.Fail()
	}

	if len(csv[0]) != 2 && !(csv[0][0] == "a" && csv[0][1] == "b") {
		t.Log("Expected headers to be a list with 2 elements. [a, b]")
		t.Fail()
	}

	if len(csv[1]) != 2 && !(csv[1][0] == "1" && csv[1][1] == "2") {
		t.Log("Expected first data row to be a list with 2 elements. [1, 2]")
		t.Fail()
	}

	if len(csv[2]) != 2 && !(csv[2][0] == "10" && csv[2][1] == "20") {
		t.Log("Expected seconds data row to be a list with 2 elements. [10, 20]")
		t.Fail()
	}
}
