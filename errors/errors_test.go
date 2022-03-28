package errors

import (
	"errors"
	"testing"
)

func stdErrorFunc() error {
	return &StdError{Caller: "stdErrorFunc", Code: "MyStrangeCode", Err: errors.New("My bogus error"), Msg: "Error message", StatusCode: 400}
}

func TestStdErrorFunc(t *testing.T) {

	var err *StdError
	e := stdErrorFunc()

	if ok := errors.As(e, &err); ok {
		if err.Caller != "stdErrorFunc" {
			t.Errorf("got: %s, want: %s", err.Caller, "stdErrorFunc")
		}
	}

	// If we're sure it's the right type
	if errors.As(e, &err) {
		if err.Caller != "stdErrorFunc" {
			t.Errorf("got: %s, want: %s", err.Caller, "stdErrorFunc")
		}
	}

	if errors.Is(e, err) {
		if err.Msg != "Error message" {
			t.Errorf("got: %s, want: %s", err.Msg, "Error message")
		}
	}
}
