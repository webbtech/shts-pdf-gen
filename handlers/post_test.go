package handlers

import (
	"encoding/json"
	"strings"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestValidatePost(t *testing.T) {

	t.Run("Missing struct fields", func(t *testing.T) {
		p := &Post{}
		p.request.Body = `{}`

		bson.UnmarshalExtJSON([]byte(p.request.Body), true, &p.input)

		err := validateInput(p.input)
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
		p := &Post{}
		p.request.Body = `{"number": 900, "requestType": "estimat"}`
		json.Unmarshal([]byte(p.request.Body), &p.input)

		err := validateInput(p.input)
		if !strings.HasPrefix(err.Msg, ERR_INVALID_TYPE) {
			t.Fatalf("Expected error to start with: %s", ERR_INVALID_TYPE)
		}
	})
}
