package errors

import (
	"errors"
	"fmt"
	"runtime"
)

const (
	_NO_DESCRIPTION = "error with no description"
	_CAUSED_BY      = "\n\tcaused by:\n"
)

type stackedError struct {
	errorString string
	at          string
}

func (e *stackedError) At() string {
	return e.at
}

func (e *stackedError) Error() string {
	if e == nil {
		return "error is nil - " + locateAt(1)
	}
	last := Head(e)
	return fmt.Sprintf(
		"%s at: %s",
		last.errorString,
		last.at,
	)
}

func (e *stackedError) String() string {
	if e == nil {
		return "error is nil - " + locateAt(1)
	}
	var str string
	for _, err := range Stack(e).([]string) {
		str += err + "\n"
	}
	return str
}

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

	err := &stackedError{
		errorString: description,
	}

	if cause == nil {
		return err
	}

	return fmt.Errorf("%s%s%w", err, _CAUSED_BY, cause)
}

// WrapLocate formats an error message with an optional underlying error wrapped as its cause.
// It also locates the error in order to improve the debugging experience.
func WrapLocate(description string, cause error) error {
	if description == "" {
		description = _NO_DESCRIPTION
	}
	err := &stackedError{
		errorString: description,
		at:          locateAt(2),
	}
	if cause == nil {
		return err
	}
	return fmt.Errorf("%s%s%w", err, _CAUSED_BY, err)
}

// NewLocate creates a new error with location from a given string.
//
// str: the string to create the location from. If empty, a default string will be used.
// error: the created error
func NewLocate(str string) error {
	if str == "" {
		str = _NO_DESCRIPTION
	}

	return &stackedError{
		errorString: str,
		at:          locateAt(2),
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

func Head(err error) *stackedError {
	if err == nil {
		return nil
	}

	stacked := errors.Unwrap(err)

	if stacked == nil {
		if e, ok := err.(*stackedError); ok {
			return e
		}
		return &stackedError{errorString: err.Error()}
	}

	if e, ok := stacked.(*stackedError); ok {
		return e
	}

	return &stackedError{
		errorString: stacked.Error(),
	}
}

func Stack(err error) interface{} {
	if err == nil {
		return nil
	}

	out := []string{Head(err).Error()}

	for e := errors.Unwrap(err); e != nil; e = errors.Unwrap(e) {
		out = append(out, Head(e).Error())
	}

	if len(out) > 0 {
		return out
	}

	return nil
}
