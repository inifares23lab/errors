package errors

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

type stackedError struct {
	msg string
	at  string
	err error
}

func (err *stackedError) Unwrap() error {
	return errors.Unwrap(err)
}

// String returns a string representation of the stackedError.
//
// It returns the error message and the stack trace location.
// The return type is a string.
func (e *stackedError) Error() string {
	if e == nil {
		return "error is nil - " + locateAt(1)
	}
	return fmt.Sprintf(
		"%s at: %s",
		e.msg,
		e.at,
	)
}

func head(err error) string {
	if e, ok := err.(*stackedError); ok {
		return e.msg
	}
	if err == nil {
		return ""
	}
	stacked := errors.Unwrap(err)
	if stacked == nil {
		return err.Error()
	}
	return strings.Replace(err.Error(), stacked.Error(), "", 1)
}

func Stack(err error) interface{} {
	if err == nil {
		return nil
	}

	type stackString struct {
		msg string
		at  string
	}

	out := []stackString{}
	e, ok := err.(*stackedError)

LOOP:
	for {
		switch {
		case ok:
			out = append(out, stackString{e.msg, e.at})
			fallthrough
		case ok && e.err != nil:
			e, ok = e.err.(*stackedError)
		default:
			break LOOP
		}
	}

	if len(out) > 0 {
		return out
	}

	return nil
}

// Wrap wraps an error with a description and an optional cause.
//
// Parameters:
// - description: the description of the error. If empty, a default description will be used.
// - cause: the cause of the error. If nil, the error will not have a cause.
//
// Returns:
// - error: the wrapped error.
func Wrap(description string, cause error) error {
	return &stackedError{
		description,
		"",
		cause,
	}
}

// WrapLocate wraps an error with a description and a cause.
//
// The WrapLocate function takes a description string and an error cause as
// parameters. It checks if the description is empty, and if so, it sets it to
// "_NO_DESCRIPTION". If the cause is nil, it returns a stackedError with a new
// error created from the description, and the location of the error set to the
// calling function. If the cause is not nil, it returns a stackedError with an
// error created by concatenating the description, the "_CAUSED_BY" string, and
// the cause, and the location of the error set to the calling function.
//
// The function returns an error of type stackedError.
func WrapLocate(description string, cause error) error {
	return &stackedError{
		description,
		locateAt(2),
		cause,
	}
}

// NewLocate creates a new error with location from a given string.
//
// str: the string to create the location from. If empty, a default string will be used.
// error: the created error
func NewLocate(str string) error {
	return &stackedError{
		str,
		locateAt(2),
		nil,
	}
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
func locateAt(skip int) string {
	if _, file, line, ok := runtime.Caller(skip); ok {
		return fmt.Sprintf("\"%s:%d\"", file, line)
	}

	return fmt.Sprintf("could not locate the caller - skipped %d", skip)
}
