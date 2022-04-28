package handlers

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestValidateRequest(t *testing.T) {

	t.Run("Missing struct fields", func(t *testing.T) {
		p := &Pdf{}
		p.request.Body = `{}`

		err := p.validateInput()
		nLines := strings.Split(err.Msg, "\n")
		expectedNumErrs := 2
		haveLines := len(nLines)
		if expectedNumErrs != haveLines {
			t.Fatalf("Number of Msg errors should be: %d, have: %d", expectedNumErrs, haveLines)
		}

		expectedError1 := "Field validation for 'EstimateNumber' failed on the 'required' tag"
		haveError1 := nLines[0]
		if expectedError1 != haveError1 {
			t.Fatalf("First error message should be: %s, have: %s", expectedError1, haveError1)
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
