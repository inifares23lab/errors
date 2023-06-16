package errors

import (
	"errors"
	"fmt"
)

const (
	_NO_DESCRIPTION = "error with no description"
	_CAUSED_BY      = "\n\tcaused by:\n"
)

// Wrap formats an error message with an optional underlying error wrapped as its cause.
// It locates the error only if it is the first in the chain or if a descripion is missing
// to improve the debugging experience without adding too much overhead.
// Args:
//
//	description (string): The description of the error.
//	err (error): The underlying error to wrap.
//
// Returns:
//
//	error: The wrapped error with or without the location of the caller.
func Wrap(description string, cause error) error {
	if description == "" {
		description = _NO_DESCRIPTION
	}
	if cause == nil {
		return errors.New(description)
	}
	return fmt.Errorf("%s%s%w", description, _CAUSED_BY, cause)
}

func WrapLocate(description string, cause error) error {
	if description == "" {
		description = _NO_DESCRIPTION
	}
	if cause == nil {
		return locateAt(description, 2)
	}
	return fmt.Errorf("%s%s%w", locateAt(description, 2), _CAUSED_BY, cause)
}

func NewLocate(str string) error {
	return locateAt(str, 2)
}
