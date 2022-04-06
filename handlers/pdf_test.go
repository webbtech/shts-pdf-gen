package handlers

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestValidateRequest(t *testing.T) {
	t.Run("Missing request body", func(t *testing.T) {
		p := &Pdf{}
		requestBody := ""
		json.Unmarshal([]byte(requestBody), &p.input)

		expectedErr := ERR_MISSING_REQUEST_BODY
		err := p.validateInput()

		if err.Msg != expectedErr {
			t.Fatalf("Expected error should be: %s, have: %s", expectedErr, err.Msg)
		}
	})

	t.Run("Missing type input", func(t *testing.T) {
		p := &Pdf{}
		requestBody := `{"number": 900}`
		json.Unmarshal([]byte(requestBody), &p.input)

		expectedErr := ERR_MISSING_TYPE
		err := p.validateInput()
		if err.Msg != expectedErr {
			t.Fatalf("Expected error should be: %s, have: %s", expectedErr, err.Msg)
		}
	})

	t.Run("Missing number input", func(t *testing.T) {
		p := &Pdf{}
		requestBody := `{"requestType": "estimate"}`
		json.Unmarshal([]byte(requestBody), &p.input)

		expectedErr := ERR_MISSING_NUMBER
		err := p.validateInput()
		if expectedErr != err.Msg {
			t.Fatalf("Expected error should be: %s, have: %s", expectedErr, err.Msg)
		}
	})

	t.Run("Invalid type input", func(t *testing.T) {
		p := &Pdf{}
		requestBody := `{"number": 900, "requestType": "estimat"}`
		json.Unmarshal([]byte(requestBody), &p.input)

		err := p.validateInput()
		if !strings.HasPrefix(err.Msg, ERR_INVALID_TYPE) {
			t.Fatalf("Expected error to start with: %s", ERR_INVALID_TYPE)
		}
	})
}
