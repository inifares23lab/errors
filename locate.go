package errors

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

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
		return fmt.Errorf("%s\n\tat \"%s:%d\"", str, file, line)
	}
	// this should never happen but if it does it adds the goroutine stacktrace with a little extra overhead
	return fmt.Errorf("%s\n\tat \"could not locate the error, getting stacktrace:\n(%s)\"",
		str, debug.Stack())
}
