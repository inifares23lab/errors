package errors

import (
	"errors"
)

func Is(err error, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}

func Unwrap(err error) *stackedError {
	if err == nil {
		return nil
	}

	stacked := errors.Unwrap(err)

	if stacked == nil {
		return nil
	}

	if e, ok := stacked.(*stackedError); ok {
		return e
	}

	return &stackedError{
		stacked,
		"",
	}
}

func New(s string) error {
	if s == "" {
		s = _NO_DESCRIPTION
	}
	return &stackedError{
		errors.New(s),
		"",
	}
}

func Join(errs ...error) error {
	return errors.Join(errs...)
}
