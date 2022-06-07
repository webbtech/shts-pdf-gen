package handlers

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/webbtech/shts-pdf-gen/config"
	"go.mongodb.org/mongo-driver/bson"
)

func TestValidateDelete(t *testing.T) {

	t.Run("Missing struct fields", func(t *testing.T) {
		d := &Delete{}
		d.request.Body = `{}`

		bson.UnmarshalExtJSON([]byte(d.request.Body), true, &d.input)

		err := validateInput(d.input)
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
		p := &Delete{}
		p.request.Body = `{"number": 900, "requestType": "estimat"}`
		bson.UnmarshalExtJSON([]byte(p.request.Body), true, &p.input)

		err := validateInput(p.input)
		if !strings.HasPrefix(err.Msg, ERR_INVALID_TYPE) {
			t.Fatalf("Expected error to start with: %s", ERR_INVALID_TYPE)
		}
	})
}

func TestProcessDelete(t *testing.T) {

	getConfig(t)
	p := &Delete{Cfg: cfg}
	requestBody := `{"number": 1191, "requestType": "estimate"}`
	json.Unmarshal([]byte(requestBody), &p.input)

	p.process()

	if p.response.StatusCode != 200 {
		t.Fatalf("Expected StatusCode of 200, got %d", p.response.StatusCode)
	}
}

// ================================== Helpers ==========================================

var cfg *config.Config

func getConfig(t *testing.T) {
	t.Helper()

	cfg = &config.Config{}
	cfg.Init()
}
