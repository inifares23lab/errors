package errors

import (
	"errors"
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
)

const (
	_NO_DESCRIPTION = "error with no description"
	_CAUSED_BY      = "\n\tcaused by:\n"
)

// Wrap formats an error message with an optional underlying error wrapped as its cause.
// Args:
//
//	description (string): The description of the error.
//	err (error): The underlying error to wrap.
//
// Returns:
//
//	error: The wrapped errors.
func Wrap(description string, cause error) error {
	if description == "" {
		description = _NO_DESCRIPTION
	}
	if cause == nil {
		return fmt.Errorf("%s", description)
	}
	return fmt.Errorf("%s%s%w", description, _CAUSED_BY, cause)
}

// WrapLocate formats an error message with an optional underlying error wrapped as its cause.
// It also locates the error in order to improve the debugging experience.
func WrapLocate(description string, cause error) error {
	if description == "" {
		description = _NO_DESCRIPTION
	}
	if cause == nil {
		return locateAt(description, 2)
	}
	return fmt.Errorf("%s%s%w", locateAt(description, 2), _CAUSED_BY, cause)
}

// NewLocate creates a new error with location from a given string.
//
// str: the string to create the location from. If empty, a default string will be used.
// error: the created error
func NewLocate(str string) error {
	if str == "" {
		str = _NO_DESCRIPTION
	}
	return locateAt(str, 2)
}

// locateAt is used to locate the caller that generates this error.
// Args:
//
//	str (string): The string of the error message.
//	skip (int): Defaults to 2 at the moment, but provided for future extensions.
//
// Returns:
//
//	error: The error with the location of the caller.
func locateAt(str string, skip int) error {
	if _, file, line, ok := runtime.Caller(skip); ok {
		return fmt.Errorf("%s\n\tat \"%s:%d\"", str, file, line)
	}
	// this should never happen but if it does it adds the goroutine stacktrace with a little extra overhead
	return fmt.Errorf("%s\n\tat \"could not locate the error, getting stacktrace:\n(%s)\"",
		str, debug.Stack())
}

func Last(err error) error {
	if err == nil {
		return nil
	}
	if e := errors.Unwrap(err); e != nil {
		err = fmt.Errorf(strings.ReplaceAll(err.Error(), e.Error(), ""))
	}
	return err
}

func Stack(err error) interface{} {
	if err == nil {
		return nil
	}
	out := []string{Last(err).Error()}
	for e := errors.Unwrap(err); e != nil; e = errors.Unwrap(e) {
		err = e
		out = append(out, Last(err).Error())
	}
	if len(out) > 0 {
		return out
	}
	return nil
}
