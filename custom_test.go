package errors

import (
	"errors"
	"strings"
	"testing"
)

// TestWrap tests the Wrap function.
//
// It provides tests for various scenarios where Wrap is provided with different
// input parameters. The function tests if Wrap behaves as expected and if the
// error message contains the expected output. The function doesn't return
// anything.
func TestWrap(t *testing.T) {
	t.Run("Test when description and err are both nil/empty", func(t *testing.T) {
		err := Wrap("", nil)
		if err == nil || (err != nil && !strings.Contains(err.Error(), _NO_DESCRIPTION)) {
			t.Errorf("Expected default error but got %v", err)
		}
	})

	t.Run("Test when only description is provided", func(t *testing.T) {
		description := "test error"
		err := Wrap(description, nil)
		if err == nil {
			t.Errorf("Expected non-nil error but got nil")
		}
		if !strings.Contains(err.Error(), description) {
			t.Errorf("Expected error to contain %s but got %s", description, err.Error())
		}
	})

	t.Run("Test when only err is provided", func(t *testing.T) {
		innerErr := errors.New("inner error")
		err := Wrap("", innerErr)
		if err == nil {
			t.Errorf("Expected non-nil error but got nil")
		}
		if !errors.Is(err, innerErr) {
			t.Errorf("Expected error to contain wrapped %v but got %v", innerErr, err)
		}
		if !strings.Contains(err.Error(), _NO_DESCRIPTION) {
			t.Errorf("Expected error to contain %s but got %s", _NO_DESCRIPTION, err.Error())
		}
		causedBy := "caused by:"
		if !strings.Contains(err.Error(), causedBy) {
			t.Errorf("Expected error to contain %s but got %s", causedBy, err.Error())
		}
	})

	t.Run("Test when description and err are both provided", func(t *testing.T) {
		description := "test error"
		innerErr := errors.New("inner error")
		err := Wrap(description, innerErr)
		if err == nil {
			t.Errorf("Expected non-nil error but got nil")
		}
		if !errors.Is(err, innerErr) {
			t.Errorf("Expected error to contain wrapped %v but got %v", innerErr, err)
		}
		if !strings.Contains(err.Error(), description) {
			t.Errorf("Expected error to contain %s but got %s", description, err.Error())
		}
		causedBy := "caused by:"
		if !strings.Contains(err.Error(), causedBy) {
			t.Errorf("Expected error to contain %s but got %s", causedBy, err.Error())
		}
	})
}

func TestWrapLocate(t *testing.T) {
	t.Run("Test when description and err are both nil/empty", func(t *testing.T) {
		err := WrapLocate("", nil)
		if err == nil || (err != nil && !strings.Contains(err.Error(), _NO_DESCRIPTION)) {
			t.Errorf("Expected default error but got %v", err)
		}
		at := "at "
		if !strings.Contains(err.Error(), at) {
			t.Errorf("Expected error to contain %s but got %s", at, err.Error())
		}
	})

	t.Run("Test when only description is provided", func(t *testing.T) {
		description := "test error"
		err := WrapLocate(description, nil)
		if err == nil {
			t.Errorf("Expected non-nil error but got nil")
		}
		if !strings.Contains(err.Error(), description) {
			t.Errorf("Expected error to contain %s but got %s", description, err.Error())
		}
		at := "at "
		if !strings.Contains(err.Error(), at) {
			t.Errorf("Expected error to contain %s but got %s", at, err.Error())
		}

	})

	t.Run("Test when only err is provided", func(t *testing.T) {
		innerErr := errors.New("inner error")
		err := WrapLocate("", innerErr)
		if err == nil {
			t.Errorf("Expected non-nil error but got nil")
		}
		if !errors.Is(err, innerErr) {
			t.Errorf("Expected error to contain wrapped %v but got %v", innerErr, err)
		}
		if !strings.Contains(err.Error(), _NO_DESCRIPTION) {
			t.Errorf("Expected error to contain %s but got %s", _NO_DESCRIPTION, err.Error())
		}
		causedBy := "caused by:"
		if !strings.Contains(err.Error(), causedBy) {
			t.Errorf("Expected error to contain %s but got %s", causedBy, err.Error())
		}
		at := "at "
		if !strings.Contains(err.Error(), at) {
			t.Errorf("Expected error to contain %s but got %s", at, err.Error())
		}
	})
}
