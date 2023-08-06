package errors

import (
	"errors"
	"fmt"
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
	if result := err.Error(); !strings.EqualFold(result, expected) {
		t.Errorf("Expected %s, but got %s", expected, result)
	}

	// Testing the string representation of a stackedError with additional information
	err = &stackedError{
		msg: "Error message",
		at:  "some location",
	}
	expected = "Error message at some location"
	if result := err.Error(); !strings.EqualFold(result, expected) {
		t.Errorf("Expected %s, but got %s", expected, result)
	}

	// Testing the string representation of a stackedError with a caused by error
	err = &stackedError{
		msg: "Error message",
		err: New("Caused by error"),
	}
	expected = "Error message\n\tcaused by:\nCaused by error"
	if result := err.Error(); !strings.EqualFold(result, expected) {
		t.Errorf("Expected %s, but got %s", expected, result)
	}

	// Testing the string representation of a stackedError with a caused by error
	err = &stackedError{
		msg: "Error message",
		err: New("Caused by error"),
		at:  "some location",
	}
	expected = "Error message at some location\n\tcaused by:\nCaused by error"
	if result := err.Error(); !strings.EqualFold(result, expected) {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestTWrap(t *testing.T) {
	cause := New("This is the cause error")

	// Test case 1: Wrapping an error with a description
	err := TWrap("This is the description", cause)
	if !strings.Contains(Last(err), "This is the description") {
		t.Errorf("Expected error description 'This is the description', got '%s'", err.Error())
	}
	if err.(*stackedError).err != cause {
		t.Errorf("Expected cause error '%v', got '%v'", cause, err.(*stackedError).err)
	}

	// Test case 2: Wrapping an error without a description
	err = TWrap("", cause)
	desc := strings.Trim(strings.Split(Last(err), "at")[0], " ")
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

// tests that the local Unwrap function is backward compatble with standard library
// and correctly unwrapps erros generated with fmt.Errorf("%w", err)
func TestUnwrap(t *testing.T) {
	// Test case 1: Unwrapping an error generated with fmt.Errorf("%w", err)
	err := New("This is the wrapped error")
	wrappedErr := fmt.Errorf("hey %w", err)
	unwrappedErr := Unwrap(wrappedErr)
	if unwrappedErr != err {
		t.Errorf("Unwrap() returned incorrect error, expected: %v, got: %v", err, unwrappedErr)
	}

	// Test case 2: Unwrapping an error that is not generated with fmt.Errorf("%w", err)
	err = New("This is the wrapped stacked error")
	unwrappedErr = Unwrap(err)
	if unwrappedErr != nil {
		t.Errorf("Unwrap() returned incorrect error, expected: %v, got: %v", err, unwrappedErr)
	}

	err = TNew("This is the wrapped stackedError - level 0")
	err1 := fmt.Errorf("this is standard wrapped - level 1, %w", err)
	err2 := TWrap("this is not standard wrapped - level 2", err)
	unwrapped1 := Unwrap(err1)
	unwrapped2 := Unwrap(err2)
	if unwrapped1 != err {
		t.Errorf("Unwrap() returned incorrect error, expected: %v, got: %v", err, unwrapped1)
	}
	if unwrapped2 != err {
		t.Errorf("Unwrap() returned incorrect error, expected: %v, got: %v", err1, unwrapped2)
	}
}

// tests that the local Is function is backward compatble with standard library
// and correctly handles erros generated with fmt.Errorf("%w", err)
func TestIs(t *testing.T) {
	// Test case 1: comparing an error generated with fmt.Errorf("%w", err)
	err := New("This is the wrapped error")
	wrappedErr := fmt.Errorf("hey %w", err)
	if !Is(wrappedErr, err) {
		t.Errorf("Is() failed with %v and %v\n", err, wrappedErr)
	}

	// Test case 2: comparing a nil error
	err = New("This is the wrapped stacked error")
	if Is(nil, err) {
		t.Errorf("Is() did not fail with %v and %v\n", nil, err)
	}

	// Test case 3: comparing errors generated with this and the standard library
	err = New("This is the wrapped stackedError - level 0")
	err1 := fmt.Errorf("this is standard wrapped - level 1, %w", err)
	err2 := Wrap("this is not standard wrapped - level 2", err1)
	if !Is(err1, err) {
		t.Errorf("Is() failed with %v and %v\n", err, err1)
	}
	if !Is(err2, err1) {
		t.Errorf("Is() failed with %v and %v\n", err1, err2)
	}
	if !Is(err2, err) {
		t.Errorf("Is() failed with %v and %v\n", err, err2)
	}
}

// tests that the local As function is backward compatble with standard library
// and correctly handles erros generated with fmt.Errorf("%w", err)
func TestAs(t *testing.T) {
	// Test case 1: As() on error generated with fmt.Errorf("%w", err)
	err := New("This is the wrapped error")
	wrappedErr := fmt.Errorf("hey %w", err)
	if !As(wrappedErr, &err) {
		t.Errorf("As() failed with %v and %v\n", err, wrappedErr)
	}

	// Test case 2: As() on error generated with this and the standard library
	err = New("This is the wrapped stackedError - level 0")
	err1 := fmt.Errorf("this is standard wrapped - level 1, %w", err)
	err2 := Wrap("this is not standard wrapped - level 2", err1)
	if !As(err1, &err) {
		t.Errorf("As() failed with %v and %v\n", err, err1)
	}
	if !As(err2, &err1) {
		t.Errorf("As() failed with %v and %v\n", err1, err2)
	}
	if !As(err2, &err) {
		t.Errorf("As() failed with %v and %v\n", err, err2)
	}
}
