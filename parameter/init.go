package parameter

import (
	"fmt"
	"franciscoperez.dev/gosqltojson/utils"
	"strconv"
	"strings"
	"time"
)

const (
	DEFAULT_DATE_FORMAT = "2006-01-02"
)

var DATE_FIELDS = []string{
	"CURRENT_DATE",
	"START_CURRENT_MONTH",
	"END_CURRENT_MONTH",
	"START_CURRENT_YEAR",
	"END_CURRENT_YEAR",
}

func ParseQueryParams(queryRaw string, args []string) (map[string]interface{}, error) {
	parameters := make(map[string]interface{})

	i := 0
	size := len(args)

	for {
		if i+1 >= size {
			break
		}

		name := strings.ReplaceAll(args[i], "-", "")
		value := args[i+1]

		if strings.HasPrefix(value, "-") {
			value = "true"
			i = i + 1
		} else {
			i = i + 2
		}

		parsedValue, err := parseParameter(value, "|")

		if err != nil {
			return nil, err
		}

		//Only variables in our query
		if strings.Contains(queryRaw, "@"+name) {
			if utils.IsInteger(parsedValue) {
				intValue, _ := strconv.ParseInt(parsedValue, 10, 64)

				parameters[name] = intValue
			} else if utils.IsNumber(parsedValue) {
				floatValue, _ := strconv.ParseFloat(parsedValue, 64)

				parameters[name] = floatValue
			} else {
				parameters[name] = parsedValue
			}
		}
	}

	queryWithAllParams := queryRaw

	for parameter, parameterValue := range parameters {
		queryWithAllParams = strings.ReplaceAll(queryWithAllParams, "@"+parameter, parameter+fmt.Sprint(parameterValue))
	}

	if strings.Contains(queryWithAllParams, "@") {
		return parameters, fmt.Errorf("please provide all required parameters in the query %s", queryWithAllParams)
	}

	return parameters, nil
}

func parseParameter(v string, formatSeparator string) (string, error) {
	if v == "" {
		return v, nil
	} else if utils.IsNumber(v) {
		return v, nil
	} else if strings.Contains(v, formatSeparator) {
		parts := strings.Split(v, formatSeparator)
		field := parts[0]
		dateFormat := parts[1]

		return parseFormula(field, dateFormat)
	}

	return parseFormula(v, DEFAULT_DATE_FORMAT)
}

func parseFormula(formula string, dateFormat string) (string, error) {
	currentDate := time.Now()

	for _, dateField := range DATE_FIELDS {
		if strings.Contains(formula, dateField) {
			var separator string
			sign := 1

			if strings.Contains(formula, "+") {
				separator = "+"
			} else if strings.Contains(formula, "-") {
				separator = "-"
				sign = -1
			}

			if separator != "" {
				parts := strings.Split(formula, separator)

				toAdd, err := strconv.Atoi(parts[1])

				if err != nil {
					return "", err
				}

				return parseField(parts[0], sign*toAdd, currentDate, dateFormat), nil
			}
		}
	}

	return parseField(formula, 0, currentDate, dateFormat), nil
}

func parseField(field string, toAdd int, currentDate time.Time, dateFormat string) string {
	if field == "" {
		return field
	}

	if field == "CURRENT_DATE" {
		return currentDate.AddDate(0, 0, toAdd).Format(dateFormat)
	} else if field == "START_CURRENT_MONTH" {
		currentYear, currentMonth, _ := currentDate.Date()
		firstDayOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentDate.Location())
		return firstDayOfMonth.AddDate(0, toAdd, 0).Format(dateFormat)
	} else if field == "END_CURRENT_MONTH" {
		currentYear, currentMonth, _ := currentDate.Date()
		firstDayOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentDate.Location())
		return firstDayOfMonth.AddDate(0, 1+toAdd, -1).Format(dateFormat)
	} else if field == "START_CURRENT_YEAR" {
		currentYear := currentDate.Year()
		firstDayOfMonth := time.Date(currentYear, 1, 1, 0, 0, 0, 0, currentDate.Location())
		return firstDayOfMonth.AddDate(toAdd, 0, 0).Format(dateFormat)
	} else if field == "END_CURRENT_YEAR" {
		currentYear := currentDate.Year()
		firstDayOfMonth := time.Date(currentYear, 1, 1, 0, 0, 0, 0, currentDate.Location())
		return firstDayOfMonth.AddDate(1+toAdd, 0, -1).Format(dateFormat)
	}

	return field
}
