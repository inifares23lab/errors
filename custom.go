package errors

import (
	"errors"
	"fmt"
)

const _NO_DESCRIPTION = "error with no description"

var errDefault = errors.New(_NO_DESCRIPTION)

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
func Wrap(description string, err error) error {
	if err == nil {
		if description == "" {
			return nil
		}
		return locateAt(errors.New(description), 2)
	}
	if u, ok := err.(interface {
		Unwrap() []error
	}); ok && len(u.Unwrap()) > 1 {
		start, end := errorJoinedCauseFormat()
		if description == "" {
			return fmt.Errorf("%s%s%w%s", locateAt(errDefault, 2), start, err, end)
		}
		return fmt.Errorf("%s%s%w%s", description, start, err, end)
	}
	if description == "" {
		return fmt.Errorf("%s%s%w", locateAt(errDefault, 2), errorCauseFormat(), err)
	}
	return fmt.Errorf("%s%s%w", description, errorCauseFormat(), err)
}

func errorCauseFormat() string {
	return "\n\tcaused by:\n"
}

func errorJoinedCauseFormat() (string, string) {
	return "\n\tcaused by:\n[\n", "\n]"
}
