package parameter

import (
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
		if i >= size {
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
			parameters[name] = parsedValue
		}
	}

	return parameters, nil
}

func parseParameter(v string, format_separator string) (string, error) {
	if format_separator == "" {
		format_separator = "|"
	}

	if v == "" {
		return v, nil
	} else if isNumber(v) {
		return v, nil
	} else if strings.Contains(v, format_separator) {
		parts := strings.Split(v, format_separator)
		field := parts[0]
		date_format := parts[1]

		return parseFormula(field, date_format)
	}

	return parseFormula(v, DEFAULT_DATE_FORMAT)
}

func parseFormula(formula string, date_format string) (string, error) {
	current_date := time.Now()

	for _, date_field := range DATE_FIELDS {
		if strings.Contains(formula, date_field) {
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

				to_add, err := strconv.Atoi(parts[1])

				if err != nil {
					return "", err
				}

				return parseField(parts[0], sign*to_add, current_date, date_format), nil
			}
		}
	}

	return parseField(formula, 0, current_date, date_format), nil
}

func parseField(field string, to_add int, current_date time.Time, date_format string) string {
	if field == "" {
		return field
	}

	if date_format == "" {
		date_format = DEFAULT_DATE_FORMAT
	}

	if field == "CURRENT_DATE" {
		return current_date.AddDate(0, 0, to_add).Format(date_format)
	} else if field == "START_CURRENT_MONTH" {
		currentYear, currentMonth, _ := current_date.Date()
		firstDayOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, current_date.Location())
		return firstDayOfMonth.AddDate(0, to_add, 0).Format(date_format)
	} else if field == "END_CURRENT_MONTH" {
		currentYear, currentMonth, _ := current_date.Date()
		firstDayOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, current_date.Location())
		return firstDayOfMonth.AddDate(0, 1+to_add, -1).Format(date_format)
	} else if field == "START_CURRENT_YEAR" {
		currentYear := current_date.Year()
		firstDayOfMonth := time.Date(currentYear, 1, 1, 0, 0, 0, 0, current_date.Location())
		return firstDayOfMonth.AddDate(to_add, 0, 0).Format(date_format)
	} else if field == "END_CURRENT_YEAR" {
		currentYear := current_date.Year()
		firstDayOfMonth := time.Date(currentYear, 1, 1, 0, 0, 0, 0, current_date.Location())
		return firstDayOfMonth.AddDate(1+to_add, 0, -1).Format(date_format)
	}

	return field
}

func isNumber(value string) bool {
	_, err := strconv.Atoi(value)
	return err == nil
}
