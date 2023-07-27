package errors

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestStackedError_String(t *testing.T) {
	// Testing the string representation of a stackedError without additional information
	err := &stackedError{
		msg: "Error message",
	}
	expected := "Error message"
	if result := err.String(); !strings.EqualFold(result, expected) {
		t.Errorf("Expected %s, but got %s", expected, result)
	}

	// Testing the string representation of a stackedError with additional information
	err = &stackedError{
		msg: "Error message",
		at:  "some location",
	}
	expected = "Error message at some location"
	if result := err.String(); !strings.EqualFold(result, expected) {
		t.Errorf("Expected %s, but got %s", expected, result)
	}

	// Testing the string representation of a stackedError with a caused by error
	err = &stackedError{
		msg: "Error message",
		err: New("Caused by error"),
	}
	expected = "Error message\n\tcaused by:\nCaused by error"
	if result := err.String(); !strings.EqualFold(result, expected) {
		t.Errorf("Expected %s, but got %s", expected, result)
	}

	// Testing the string representation of a stackedError with a caused by error
	err = &stackedError{
		msg: "Error message",
		err: New("Caused by error"),
		at:  "some location",
	}
	expected = "Error message at some location\n\tcaused by:\nCaused by error"
	if result := err.String(); !strings.EqualFold(result, expected) {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestTWrap(t *testing.T) {
	cause := New("This is the cause error")

	// Test case 1: Wrapping an error with a description
	err := TWrap("This is the description", cause)
	if !strings.Contains(err.Error(), "This is the description") {
		t.Errorf("Expected error description 'This is the description', got '%s'", err.Error())
	}
	if err.(*stackedError).err != cause {
		t.Errorf("Expected cause error '%v', got '%v'", cause, err.(*stackedError).err)
	}

	// Test case 2: Wrapping an error without a description
	err = TWrap("", cause)
	desc := strings.Trim(strings.Split(err.Error(), "at")[0], " ")
	if desc != "" {
		t.Errorf("Expected empty error description, got '%s'", desc)
	}
	if err.(*stackedError).err != cause {
		t.Errorf("Expected cause error '%v', got '%v'", cause, err.(*stackedError).err)
	}
}

func TestStack(t *testing.T) {
	// Testing when error is nil

	out := Stack(nil)
	if out != nil {
		t.Errorf("Expected nil or empty, got %v", out)
	}

	// Testing when error is not a stackedError
	err := New("Test Error")
	expectedOutput := []map[string]string{
		{"error": "Test Error"},
	}
	out = Stack(err)
	if !reflect.DeepEqual(out, expectedOutput) {
		t.Errorf("Expected %v, got %v", expectedOutput, out)
	}

	// Testing when error is a stackedError with nested stackedErrors
	err1 := &stackedError{
		msg: "Error 1",
		err: &stackedError{
			msg: "Error 2",
			err: errors.New("Error 3"),
		},
	}
	expectedOutput = []map[string]string{
		{"error": "Error 1"},
		{"error": "Error 2"},
		{"error": "Error 3"},
	}
	out = Stack(err1)
	if !reflect.DeepEqual(out, expectedOutput) {
		t.Errorf("Expected %v, got %v", expectedOutput, out)
	}

	// Testing when error is a stackedError with nested stackedErrors and at field
	err2 := &stackedError{
		msg: "Error 1",
		err: &stackedError{
			msg: "Error 2",
			at:  "main.go:10",
			err: errors.New("Error 3"),
		},
	}
	expectedOutput = []map[string]string{
		{"error": "Error 1"},
		{"error": "Error 2", "at": "main.go:10"},
		{"error": "Error 3"},
	}
	out = Stack(err2)
	if !reflect.DeepEqual(out, expectedOutput) {
		t.Errorf("Expected %v, got %v", expectedOutput, out)
	}
}
