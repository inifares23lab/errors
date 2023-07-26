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

func Join(errs ...error) error {
	return errors.Join(errs...)
}

func New(msg string) error {
	return &stackedError{
		msg,
		"",
		nil,
	}
}
