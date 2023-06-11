package errors

import (
	"errors"
	"fmt"
	"strings"
)

const _NO_DESCRIPTION = "error with no description"

var errDefault = errors.New(_NO_DESCRIPTION)

type joinError struct {
	errs []error
}

func Join(errs ...error) error {
	if j, ok := errors.Join(errs...).(joinError); ok {
		return j
	}
	return errors.Join(errs...)
}
func (e joinError) Error() string {
	var b []byte
	for i, err := range e.errs {
		if i > 0 {
			b = append(b, "\n\t"...)
		}
		b = append(b, err.Error()...)
	}
	return string(b)
}
func (e joinError) Unwrap() []error {
	return e.errs
}

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
		e := joinError{u.Unwrap()}
		start, end := errorJoinedCauseFormat(err.Error())
		if description == "" {
			return fmt.Errorf("%s, %s%w%s", locateAt(errDefault, 2), start, e, end)
		}
		return fmt.Errorf("%s, %s%w%s", description, start, e, end)
	}
	if description == "" {
		return fmt.Errorf("%s, %s%w", locateAt(errDefault, 2), errorCauseFormat(), err)
	}
	return fmt.Errorf("%s, %s%w", description, errorCauseFormat(), err)
}

func errorCauseFormat() string {
	return "caused by:\n\t"
}

func errorJoinedCauseFormat(str string) (string, string) {
	ordinal := strings.Count(str, "joined:[")
	return fmt.Sprintf("caused by:\n%d-joined:[\n\t", ordinal), fmt.Sprintf("\n]-%d", ordinal)
}
