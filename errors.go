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

func Unwrap(err error) error {
	return errors.Unwrap(err)
}

func New(s string) error {
	return errors.New(s)
}

func Join(errs ...error) error {
	return errors.Join(errs...)
}
