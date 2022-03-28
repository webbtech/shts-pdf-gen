package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	log "github.com/sirupsen/logrus"

	lerrors "github.com/webbtech/shts-pdf-gen/errors"
	"github.com/webbtech/shts-pdf-gen/model"
)

// Stage variable
var Stage string

var ValidTypes = []string{"estimate", "invoice"}

// PDF struct
type PDF struct {
	request  events.APIGatewayProxyRequest
	response events.APIGatewayProxyResponse
	input    *model.PdfRequest
}

const (
	ERR_INVALID_TYPE         = "Invalid PDF type in input"
	ERR_MISSING_NUMBER       = "Missing Estimate number in input"
	ERR_MISSING_REQUEST_BODY = "Missing request body"
	ERR_MISSING_TYPE         = "Missing PDF type in input"
)

// ========================== Public Methods =============================== //

func (p *PDF) process() {

	rb := responseBody{}
	var body []byte
	var stdError *lerrors.StdError
	var statusCode int = 201

	json.Unmarshal([]byte(p.request.Body), &p.input)

	// validate input
	if stdError == nil {
		err := p.validateInput()
		if err != nil {
			errors.As(err, &stdError)
		}
	}

	// Process any error
	if stdError != nil {
		rb.Code = stdError.Code
		rb.Message = stdError.Msg
		statusCode = stdError.StatusCode
		logError(stdError)
	} else {
		rb.Code = "SUCCESS"
		rb.Message = "Success"
	}

	// Now creact the actual response object
	body, _ = json.Marshal(&rb)
	p.response = events.APIGatewayProxyResponse{
		Body:       string(body),
		Headers:    headers,
		StatusCode: statusCode,
	}
}

// ========================== Private Methods ============================== //

func (p *PDF) validateInput() (err *lerrors.StdError) {
	var inputErrs []string

	if p.input == nil {
		err = &lerrors.StdError{
			Caller: "handlers.validateInput",
			Code:   lerrors.CodeBadInput,
			Err:    errors.New(ERR_MISSING_REQUEST_BODY),
			Msg:    ERR_MISSING_REQUEST_BODY, StatusCode: 400,
		}
		return err
	}

	if p.input.EstimateNumber == nil {
		inputErrs = append(inputErrs, ERR_MISSING_NUMBER)
	}

	if p.input.FileType == nil {
		inputErrs = append(inputErrs, ERR_MISSING_TYPE)
	}

	if p.input.FileType != nil {
		_, found := findString(ValidTypes, *p.input.FileType)
		if !found {
			errStr := fmt.Sprintf(ERR_INVALID_TYPE+": \"%s\"", *p.input.FileType)
			inputErrs = append(inputErrs, errStr)
		}
	}

	if len(inputErrs) > 0 {
		error := errors.New(strings.Join(inputErrs, "\n"))
		err = &lerrors.StdError{Caller: "handlers.validateInput", Code: lerrors.CodeBadInput, Err: error, Msg: error.Error(), StatusCode: 400}
		return err
	}

	return nil
}

// NOTE: these could go into it's own package
func logError(err *lerrors.StdError) {
	if Stage == "" {
		Stage = os.Getenv("Stage")
	}

	if Stage != "test" {
		log.Error(err)
	}
}

func findString(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}

	return -1, false
}