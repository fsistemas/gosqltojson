package utils

import "strconv"

func IsInteger(value string) bool {
	_, err := strconv.ParseInt(value, 10, 64)
	return err == nil
}

func IsNumber(value string) bool {
	_, err := strconv.ParseFloat(value, 64)
	return err == nil
}
