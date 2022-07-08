package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

// NOTE: this test requires that config defaults are uploaded to s3, which isn't ideal for unit testing unfortunately
// This also then means that the db connection is to the production db

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

func TestAnyHandler(t *testing.T) {
	os.Setenv("Stage", "test")
	var msg string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
	}))
	defer ts.Close()

	requestBody := `{"number": 1191, "requestType": "estimate"}`
	r, err := handler(events.APIGatewayProxyRequest{HTTPMethod: "PUT", Body: requestBody})

	expectedMsg := "Invalid Verb"
	msg = extractMessage(r.Body)
	if msg != expectedMsg {
		t.Fatalf("Expected error message: %s received: %s", expectedMsg, msg)
	}
	if err != nil {
		t.Fatal("Everything should be ok")
	}
}

func TestGetHandler(t *testing.T) {

	os.Setenv("Stage", "test")
	var msg string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer ts.Close()

	r, err := handler(events.APIGatewayProxyRequest{HTTPMethod: "GET"})

	expectedMsg := "Healthy"
	msg = extractMessage(r.Body)
	if msg != expectedMsg {
		t.Fatalf("Expected error message: %s received: %s", expectedMsg, msg)
	}
	if err != nil {
		t.Fatal("Everything should be ok")
	}
}

// NOTE: here we should use mocks to avoid having to use the mongodb and Pdf packages
func TestPostHandler(t *testing.T) {

	os.Setenv("Stage", "test")
	var msg string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
	}))
	defer ts.Close()

	requestBody := `{"number": 1191, "requestType": "estimate"}`
	r, err := handler(events.APIGatewayProxyRequest{HTTPMethod: "POST", Body: requestBody})

	expectedMsg := "Success"
	msg = extractMessage(r.Body)
	if msg != expectedMsg {
		t.Fatalf("Expected error message: %s received: %s", expectedMsg, msg)
	}
	if err != nil {
		t.Fatal("Everything should be ok")
	}
}

func TestDeleteHandler(t *testing.T) {
	os.Setenv("Stage", "test")
	var msg string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer ts.Close()

	requestBody := `{"number": 1191, "requestType": "estimate"}`
	r, err := handler(events.APIGatewayProxyRequest{HTTPMethod: "DELETE", Body: requestBody})

	expectedMsg := "File deleted"
	msg = extractMessage(r.Body)
	if msg != expectedMsg {
		t.Fatalf("Expected error message: %s received: %s", expectedMsg, msg)
	}
	if err != nil {
		t.Fatal("Everything should be ok")
	}
}

// ============================== Helpers =====================================

func extractMessage(b string) (msg string) {
	var dat map[string]string
	_ = json.Unmarshal([]byte(b), &dat)
	return dat["message"]
}
