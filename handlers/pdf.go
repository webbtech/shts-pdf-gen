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

var (
	Stage             string
	ValidRequestTypes = []string{"estimate", "invoice"}
)

// Pdf struct
type Pdf struct {
	Db       model.DbHandler
	input    *model.PdfRequest
	Request  events.APIGatewayProxyRequest
	response events.APIGatewayProxyResponse
}

const (
	ERR_INVALID_TYPE         = "Invalid request type in input"
	ERR_MISSING_NUMBER       = "Missing request number in input"
	ERR_MISSING_REQUEST_BODY = "Missing request body"
	ERR_MISSING_TYPE         = "Missing request type in input"
)

// ========================== Public Methods =============================== //

func (p *Pdf) Response(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	p.process()
	return p.response, nil
}

// ========================== Private Methods ============================== //

// TODO: maybe we should return errors???
func (p *Pdf) process() {

	rb := responseBody{}
	var body []byte
	var stdError *lerrors.StdError
	var statusCode int = 201

	json.Unmarshal([]byte(p.Request.Body), &p.input)

	// Validate input
	if err := p.validateInput(); err != nil {
		errors.As(err, &stdError)
	}

	// Fetch DB Record
	if stdError == nil {
		estimateRecord, err := p.Db.FetchEstimate(*p.input.EstimateNumber)
		fmt.Printf("estimateRecord: %+v\n", estimateRecord)
		fmt.Printf("err: %+v\n", err)
	}

	// Generate PDF file
	if stdError == nil {

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

func (p *Pdf) validateInput() (err *lerrors.StdError) {
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

	if p.input.RequestType == nil {
		inputErrs = append(inputErrs, ERR_MISSING_TYPE)
	}

	if p.input.RequestType != nil {
		_, found := findString(ValidRequestTypes, *p.input.RequestType)
		if !found {
			errStr := fmt.Sprintf(ERR_INVALID_TYPE+": \"%s\"", *p.input.RequestType)
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
