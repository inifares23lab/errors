package errors

import (
	"errors"
	"fmt"
	"runtime"
	"runtime/debug"
)

const (
	_NO_DESCRIPTION = "error with no description"
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
func Wrap(description string, err error) error {
	errors.Join()
	if err == nil {
		if description == "" {
			return nil
		}
		return locateAt(description, 2)
	}
	if description == "" {
		return fmt.Errorf("%s, caused by:\n\t(%w)", locateAt(_NO_DESCRIPTION, 2), err)
	}
	return fmt.Errorf("%s, caused by:\n\t(%w)", description, err)
}

// Locate the error's position.
// This function should not be called from within this package
// else it will not return the correct result.
// Args:
//
//	err (error): The error to locate.
//
// Returns:
//
//	error: The error with the location of the caller.
func Locate(err error) error {
	if err != nil {
		return locateAt(err.Error(), 2)
	}
	return nil
}

// locateAt is used to locate the caller that generates this error.
// It should always be the second function in the stack inside this package.
// The skip parameter should be set to 2.
// In this way it can always get the location of the first function outside the package
// which is where the business logic is.
// Args:
//
//	str (string): The string of the error message.
//	skip (int): Unused at the moment, but provided for future extensions. Defaults to 2.
//
// Returns:
//
//	error: The error with the location of the caller.
func locateAt(str string, skip int) error {
	if _, file, line, ok := runtime.Caller(skip); ok {
		return fmt.Errorf("%s at \"%s:%d\"", str, file, line)
	}
	// this should never happen but if it does it adds the goroutine stacktrace with a little extra overhead
	return fmt.Errorf("%s at \"could not locate the error, getting stacktrace:\n\t(%s)\"",
		str, debug.Stack())
}

func Is(err error, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}

func Join(errs ...error) error {
	return errors.Join(errs...)
}

func New(s string) error {
	return errors.New(s)
}
