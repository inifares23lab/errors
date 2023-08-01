package errors

import (
	"errors"
)

// Wraps standard errors.Is
func Is(err error, target error) bool {
	return errors.Is(err, target)
}

// Wraps standard errors.As
func As(err error, target any) bool {
	return errors.As(err, target)
}

// // go >= 1.20
// // Wraps standard errors.Join
// func Join(errs ...error) error {
// 	return errors.Join(errs...)
// }
