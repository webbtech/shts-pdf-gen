package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/webbtech/shts-pdf-gen/handlers"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var h handlers.Handler

	switch request.Path {
	default:
		h = &handlers.Ping{}
	}

	return h.Response(request)
}

func main() {
	lambda.Start(handler)
}
