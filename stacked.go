package errors

import (
	"errors"
	"fmt"
	"runtime"
)

const _CAUSED_BY = "\n\tcaused by:\n"

type stackedError struct {
	msg string
	at  string
	err error
}

// Unwrap returns the underlying error of the stackedError.
//
// It returns an error.
func (e *stackedError) Unwrap() error {
	return e.err
}

func Unwrap(err error) error {
	u, ok := err.(*stackedError)
	if !ok {
		return errors.Unwrap(err)
	}
	return u.Unwrap()
}

// Error() returns a string representation of the stackedError.
//
// It returns the last error message and the stack trace location.
// The return type is a string.
func (e *stackedError) Last() string {
	str := e.msg
	if e.at != "" {
		str = fmt.Sprintf("%s at %s", str, e.at)
	}
	return str
}

// String returns a string representation of the given error
// including all stacked underlying errors.
//
// It takes an error as a parameter and checks if it is nil. If the error is nil,
// it returns the string "error is nil - " concatenated with the result of the
// caller function. If the error is of type *stackedError and not nil, it calls
// the String method of the error and returns its result. Otherwise, it calls the
// Error method of the error and returns its result.
func Last(err error) string {
	if err == nil {
		return "error is nil - " + caller(2)
	}
	if e, ok := err.(*stackedError); ok {
		return e.Last()
	}
	return err.Error()
}

// String returns the string representation of the stackedError in full.
func (e *stackedError) Error() string {
	str := e.msg
	if e.at != "" {
		str = fmt.Sprintf("%s at %s", str, e.at)
	}

	if err, ok := e.err.(*stackedError); ok {
		str = fmt.Sprintf("%s%s%s", str, _CAUSED_BY, err.Error())
	} else if e.err != nil {
		str = fmt.Sprintf("%s%s%s", str, _CAUSED_BY, e.err.Error())
	}

	return str
}

// Stack returns the stack trace of an error.
//
// It takes an error as a parameter and returns an interface{}.
func StackTrace(err error) interface{} {
	if err == nil {
		return nil
	}

	out := []map[string]string{}

	for e, ok := err.(*stackedError); ok; {
		m := map[string]string{
			"error": e.msg,
		}
		if e.at != "" {
			m["at"] = e.at
		}
		out = append(out, m)
		tmpErr := e.err
		if tmpErr == nil {
			break
		}
		e, ok = tmpErr.(*stackedError)
		if !ok {
			out = append(out, map[string]string{"error": tmpErr.Error()})
		}
	}

	if len(out) > 0 {
		return out
	}

	return nil
}

// New creates a new stackedError with the given message.
//
// Parameters:
// - msg: the error message.
//
// Returns:
// - error: the newly created stackedError.
func New(msg string) error {
	return &stackedError{
		msg,
		"",
		nil,
	}
}

// TNew creates a new error with location from a given string.
//
// str: the string to create the location from. If empty, a default string will be used.
// error: the created error
func TNew(str string) error {
	return &stackedError{
		str,
		caller(2),
		nil,
	}
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

// TWrap wraps an error with a description and cause.
// it includes the caller location
//
// Parameters:
//   - description: a string representing the description of the error.
//   - cause: an error representing the cause of the error.
//
// Returns:
//   - error: the wrapped error.
func TWrap(description string, cause error) error {
	return &stackedError{
		description,
		caller(2),
		cause,
	}
}

// caller is used to locate the caller that generates this error.
// Args:
//
//	str (string): The string of the error message.
//	skip (int): Defaults to 2 at the moment, but provided for future extensions.
//
// Returns:
//
//	string: representing the location of the caller.
func caller(skip int) string {
	if _, file, line, ok := runtime.Caller(skip); ok {
		return fmt.Sprintf("%s:%d", file, line)
	}

	return fmt.Sprintf("could not locate the caller - skipped %d", skip)
}
