package errors

// For more info on this, have a look at : https://go.dev/blog/go1.13-errors

import (
	"fmt"
)

// Generic error codes for errors returned by a Lambda function. The code can be
// used by API Gateway to reliably map an error response.
const (
	// CodeApplicationError is a catch-all for internal errors.
	// API Gateway mapping:     500 Internal server error
	CodeApplicationError = "APPLICATION_ERROR"

	// CodeAccessDenied represents an authorization error.
	// API Gateway mapping:     403 forbidden
	CodeAccessDenied = "ACCESS_DENIED"

	// CodeBadInput represents a bad Lambda input error.
	// API Gateway mapping:     400 Bad request
	CodeBadInput = "BAD_INPUT"
)

// StdError struct
type StdError struct {
	Caller     string
	Code       string
	Err        error
	Msg        string
	StatusCode int
}

func (e *StdError) Error() string {
	return fmt.Sprintf("Msg: %s, Caller: %s, Err: %s", e.Msg, e.Caller, e.Err)
}
