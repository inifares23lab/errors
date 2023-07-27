package errors

import (
	"fmt"
	"runtime"
)

const _CAUSED_BY = "\n\tcaused by:\n"

type stackedError struct {
	msg string
	at  string
	err error
}

func (e *stackedError) Unwrap() error {
	return e.err
}

// String returns a string representation of the stackedError.
//
// It returns the error message and the stack trace location.
// The return type is a string.
func (e *stackedError) Error() string {
	str := e.msg
	if e.at != "" {
		str = fmt.Sprintf("%s at %s", str, e.at)
	}
	return str
}

func String(err error) string {
	if err == nil {
		return "error is nil - " + caller(2)
	}
	if e, ok := err.(*stackedError); ok && e != nil {
		return e.String()
	}
	return err.Error()
}

func (e *stackedError) String() string {
	str := e.msg
	if e.at != "" {
		str = fmt.Sprintf("%s at %s", str, e.at)
	}

	if err, ok := e.err.(*stackedError); ok && err != nil {
		str = fmt.Sprintf("%s%s%s", str, _CAUSED_BY, err.String())
	} else if e.err != nil {
		str = fmt.Sprintf("%s%s%s", str, _CAUSED_BY, e.err.Error())
	}

	return str
}

func Stack(err error) interface{} {
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

// TWrap wraps an error with a description and a cause.
//
// The TWrap function takes a description string and an error cause as
// parameters. It checks if the description is empty, and if so, it sets it to
// "_NO_DESCRIPTION". If the cause is nil, it returns a stackedError with a new
// error created from the description, and the location of the error set to the
// calling function. If the cause is not nil, it returns a stackedError with an
// error created by concatenating the description, the "_CAUSED_BY" string, and
// the cause, and the location of the error set to the calling function.
//
// The function returns an error of type stackedError.
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
