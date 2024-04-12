package forms

import "testing"

func TestCustomError_Add(t *testing.T) {
	c := customError{}

	c.Add("hello", "hello")

	if c.Get("hello") == "" {
		t.Error("should have the field value after add field")
	}
}

func TestCustomError_Get(t *testing.T) {
	c := customError{}

	if c.Get("notExist") != "" {
		t.Error("should not have the field value after add field")
	}
}
