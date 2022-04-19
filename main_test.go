package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestPingHandler(t *testing.T) {

	os.Setenv("Stage", "test")

	var msg string
	t.Run("Successful ping", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
		}))
		defer ts.Close()

		r, err := handler(events.APIGatewayProxyRequest{Path: "/"})

		expectedMsg := "Healthy"
		msg = extractMessage(r.Body)
		if msg != expectedMsg {
			t.Fatalf("Expected error message: %s received: %s", expectedMsg, msg)
		}
		if err != nil {
			t.Fatal("Everything should be ok")
		}
	})
}

func TestEnvVars(t *testing.T) {
	os.Setenv("PARAM1", "VALUE12")
	t.Run("Successful ping", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
		}))
		defer ts.Close()

		_, _ = handler(events.APIGatewayProxyRequest{Path: "/"})
		p, exists := os.LookupEnv("PARAM1")
		if !exists {
			t.Fatalf("Expected value for PARAM1 to be: %s", p)
		}
	})
}

// TODO: here we should use mocks to avoid having to use the mongodb and Pdf packages
func TestPdfHandler(t *testing.T) {
	os.Setenv("Stage", "test")

	var msg string
	t.Run("Successful pdf", func(t *testing.T) {

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(201)
		}))
		defer ts.Close()

		requestBody := `{"number": 1011, "requestType": "estimate"}`
		r, err := handler(events.APIGatewayProxyRequest{Path: "/pdf", Body: requestBody})

		expectedMsg := "Success"
		msg = extractMessage(r.Body)
		if msg != expectedMsg {
			t.Fatalf("Expected error message: %s received: %s", expectedMsg, msg)
		}
		if err != nil {
			t.Fatal("Everything should be ok")
		}
	})
}

func extractMessage(b string) (msg string) {
	var dat map[string]string
	_ = json.Unmarshal([]byte(b), &dat)
	return dat["message"]
}
