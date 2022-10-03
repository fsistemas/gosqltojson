package parameter

import (
	"fmt"
	"strings"
	"testing"
)

func TestParseQueryParams(t *testing.T) {
	args := []string{
		"-x",
		"-a",
		"2",
		"-b",
		"-dateFrom",
		"CURRENT_DATE-5|2006-01-02",
		"-dateTo",
		"CURRENT_DATE+5",
		"-c",
		"1.5",
		"-startCurrentMonth",
		"START_CURRENT_MONTH+1",
		"-endCurrentMonth",
		"END_CURRENT_MONTH-1",
		"-startCurrentYear",
		"START_CURRENT_YEAR-1",
		"-endCurrentYear",
		"END_CURRENT_YEAR+1",
	}

	parameters, err := ParseQueryParams("SELECT * FROM (SELECT 1 AS a, 2 AS b, @c AS c , DATE() AS currentDate) data WHERE data.a > @a AND data.currentDate > @dateFrom AND data.currentDate < @dateTo", args)
	if err != nil {
		t.Log("ParseQueryParams with args will not produce an error", err)
		t.Fail()
	}

	if len(parameters) != 4 {
		t.Log("ParseQueryParams with ", len(args), " args and an only 4 used parameter will produce 4 parameters but was ", len(parameters), parameters)
		t.Fail()
	}

	valueA, containsA := parameters["a"]
	if !containsA || fmt.Sprint(valueA) != "2" {
		t.Log("Result parameters will contain parameter 'a' with value 2, but value was ", valueA)
		t.Fail()
	}

	valueDateFrom, containsDateFrom := parameters["dateFrom"]
	if !containsDateFrom {
		t.Log("Result parameters will contain parameter 'dateFrom', but value was ", valueDateFrom)
		t.Fail()
	}

	valueDateTo, containsDateTo := parameters["dateTo"]
	if !containsDateTo {
		t.Log("Result parameters will contain parameter 'valueDateTo', but value was ", valueDateTo)
		t.Fail()
	}

	_, err = ParseQueryParams("SELECT @a AS a, @b AS b, @missing_parameter AS missing_parameter", args)
	if err == nil || strings.Contains(err.Error(), "please provide all required parameters in the query %s") {
		t.Log("ParseQueryParams with missing parameters will generate an error", err)
		t.Fail()
	}

	invalidArgs := []string{
		"-a",
		"1",
		"-dateFrom",
		"CURRENT_DATE-asdf|2006-01-02",
	}

	_, err = ParseQueryParams("SELECT * FROM (SELECT 1 AS a, 2 AS b, DATE() AS currentDate) data WHERE data.a > @a AND data.currentDate > @dateFrom", invalidArgs)
	if err == nil {
		t.Log("ParseQueryParams with invalid args/parameters will not produce an error", err)
		t.Fail()
	}
}
