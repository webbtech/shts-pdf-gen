package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"path"

	"github.com/aws/aws-lambda-go/events"
	"github.com/webbtech/shts-pdf-gen/config"
	lerrors "github.com/webbtech/shts-pdf-gen/errors"
	"github.com/webbtech/shts-pdf-gen/model"
	"github.com/webbtech/shts-pdf-gen/services"
	"go.mongodb.org/mongo-driver/bson"
)

// docs for response to delete: https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods/DELETE

type Delete struct {
	Cfg      *config.Config
	input    *model.DocRequest
	request  events.APIGatewayProxyRequest
	response events.APIGatewayProxyResponse
}

func (d *Delete) Response(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	d.request = request
	d.process()
	return d.response, nil
}

func (d *Delete) process() {

	rb := responseBody{}
	var body []byte
	var err error
	var statusCode int = 200
	var stdError *lerrors.StdError

	// we're getting ExtJSON from the Realm function: createPDF,
	// so unmarshaling must be done on an EJSON formatted doc
	// see: https://www.mongodb.com/docs/manual/reference/mongodb-extended-json/
	bson.UnmarshalExtJSON([]byte(d.request.Body), true, &d.input)

	// Validate input
	if err := validateInput(d.input); err != nil {
		errors.As(err, &stdError)
	}

	// Remove PDF file
	if stdError == nil {
		fileObject := d.getFileObject()
		if err = services.DeleteS3Object(fileObject, d.Cfg.AwsRegion, d.Cfg.S3Bucket); err != nil {
			stdError = &lerrors.StdError{
				Caller:     "handlers.Delete.process",
				Code:       lerrors.CodeApplicationError,
				Err:        err,
				Msg:        err.Error(),
				StatusCode: 400,
			}
		}
	}

	// Process any error
	if stdError != nil {
		rb.Code = stdError.Code
		rb.Message = stdError.Msg
		statusCode = stdError.StatusCode
		logError(stdError)
	} else {
		rb.Code = lerrors.CodeSuccess
		rb.Message = "File deleted"
	}

	// Create the response object
	body, _ = json.Marshal(&rb)
	d.response = events.APIGatewayProxyResponse{
		Body:       string(body),
		Headers:    headers,
		StatusCode: statusCode,
	}
}

func (d *Delete) getFileObject() (fileObj string) {

	var fileNm string

	switch *d.input.RequestType {
	case "estimate":
		fileNm = fmt.Sprintf("%s-%d.pdf", "est", *d.input.EstimateNumber)
	case "invoice":
		fileNm = fmt.Sprintf("%s-%d.pdf", "inv", *d.input.EstimateNumber)
	}

	return path.Join(*d.input.RequestType, fileNm)
}
