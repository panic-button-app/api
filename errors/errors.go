package errors

import "net/http"

// Code represents a canonical error code.
type Code = int

const (
	// CodeInternal ...
	CodeInternal = Code(iota)
	// CodeUnauthorized ...
	CodeUnauthorized
	// CodeNotFound ...
	CodeNotFound
	// CodeDependentServiceFailure ...
	CodeDependentServiceFailure
)

// HTTPMapping is a mapping between error codes and http status codes.
var HTTPMapping = map[Code]int{
	CodeInternal:                http.StatusInternalServerError,
	CodeUnauthorized:            http.StatusUnauthorized,
	CodeNotFound:                http.StatusNotFound,
	CodeDependentServiceFailure: http.StatusInternalServerError,
}

// Error is an implementation of the error interface.
type Error struct {
	original error
	Code     Code
}

func (err Error) Error() string {
	return err.original.Error()
}

// Annotate adds an error code to the given error.
func Annotate(err error, code Code) error {
	return Error{
		original: err,
		Code:     code,
	}
}
