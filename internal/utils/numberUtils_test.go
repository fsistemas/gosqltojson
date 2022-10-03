package utils

import "testing"

func TestIsInteger(t *testing.T) {
	if IsInteger("abc") {
		t.Log("error IsInteger('abc') should be false", nil)
		t.Fail()
	}

	if !IsInteger("123") {
		t.Log("error IsInteger('123') should be true", nil)
		t.Fail()
	}

	if IsInteger("123.45") {
		t.Log("error IsInteger('123.45') should be false", nil)
		t.Fail()
	}
}

func TestIsNumber(t *testing.T) {
	if IsNumber("abc") {
		t.Log("error IsNumber('abc') should be false", nil)
		t.Fail()
	}

	if !IsNumber("123") {
		t.Log("error IsNumber('123') should be true", nil)
		t.Fail()
	}

	if !IsNumber("123.45") {
		t.Log("error IsNumber('123.45') should be true", nil)
		t.Fail()
	}
}
