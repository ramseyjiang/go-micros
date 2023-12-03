package apierror

import (
	"errors"
	"runtime"
)

// convertWrappedToAPIError converts a chain of wrapped errors
// to an APIError. Supports StatusHTTP and Source methods.
// This allows other error handling packages to be compatible
// with APIError without requiring a dependency on APIError.
// They do require the common method signatures noted below.
func convertWrappedToAPIError(err error) *APIError {
	// supported interfaces
	type httpError interface{ StatusHTTP() (int, string) }
	type traceableError interface {
		error

		Source() (pc uintptr, file string, line int, ok bool)
	}

	if err == nil {
		return nil
	}

	// shortcut
	if already, ok := err.(*APIError); ok {
		return already
	}

	var converted *APIError

	switch v := err.(type) {
	case *APIError:
		converted = v
	case httpError:
		code, msg := v.StatusHTTP()
		converted = &APIError{
			ErrorCode:    code,
			ErrorField:   msg,
			originErr:    err,
			ErrorMessage: err.Error(),
			App:          ApplicationName,
		}

	}

	var v traceableError
	ok := errors.As(err, &v)
	if ok {
		if pc, file, line, ok := v.Source(); ok {
			runtimeFuncPtr := runtime.FuncForPC(pc)
			converted.traceFunc = runtimeFuncPtr.Name()
			converted.traceFile = file
			converted.traceLine = line
		}
	}

	for err := errors.Unwrap(err); err != nil; err = errors.Unwrap(err) {
		converted.Stack = convertWrappedToAPIError(err)
	}

	return converted
}
