package handlers

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	lerrors "github.com/webbtech/shts-pdf-gen/errors"
)

type Any struct {
	response events.APIGatewayProxyResponse
}

func (c *Any) Response(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	c.process()
	return c.response, nil
}

func (c *Any) process() {
	rb := responseBody{Message: "Invalid Verb", Code: lerrors.CodeAccessDenied}
	body, _ := json.Marshal(&rb)

	c.response = events.APIGatewayProxyResponse{
		Body:       string(body),
		Headers:    headers,
		StatusCode: 403,
	}
}
