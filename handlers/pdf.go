package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	log "github.com/sirupsen/logrus"

	"github.com/webbtech/shts-pdf-gen/config"
	lerrors "github.com/webbtech/shts-pdf-gen/errors"
	"github.com/webbtech/shts-pdf-gen/model"
	"github.com/webbtech/shts-pdf-gen/pdf"
)

const (
	ERR_INVALID_TYPE         = "Invalid request type in input"
	ERR_MISSING_NUMBER       = "Missing request number in input"
	ERR_MISSING_REQUEST_BODY = "Missing request body"
	ERR_MISSING_TYPE         = "Missing request type in input"
)

var (
	Stage             string
	ValidRequestTypes = []string{"estimate", "invoice"}
)

// Pdf struct
type Pdf struct {
	Cfg      *config.Config
	Db       model.DbHandler
	input    *model.PdfRequest
	request  events.APIGatewayProxyRequest
	response events.APIGatewayProxyResponse
}

// ========================== Public Methods =============================== //
// NOTE: why is request named here?
func (p *Pdf) Response(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	p.request = request
	p.process()
	return p.response, nil
}

// ========================== Private Methods ============================== //

func (p *Pdf) process() {

	rb := responseBody{}
	var body []byte
	var err error
	var estimateRecord *model.Estimate
	var statusCode int = 201
	var stdError *lerrors.StdError

	// Validate input
	json.Unmarshal([]byte(p.request.Body), &p.input)
	if err := p.validateInput(); err != nil {
		errors.As(err, &stdError)
	}

	// Fetch DB Record
	if stdError == nil {
		estimateRecord, err = p.Db.FetchEstimate(*p.input.EstimateNumber)
		if err != nil {
			stdError = &lerrors.StdError{
				Caller:     "handlers.Pdf.process",
				Code:       lerrors.CodeApplicationError,
				Err:        err,
				Msg:        err.Error(),
				StatusCode: 400,
			}
		}
	}

	// Generate PDF file
	if stdError == nil {
		pdf, err := pdf.New(p.Cfg, *p.input.RequestType, estimateRecord)
		if err != nil {
			stdError = &lerrors.StdError{
				Caller:     "handlers.Pdf.process",
				Code:       lerrors.CodeApplicationError,
				Err:        err,
				Msg:        err.Error(),
				StatusCode: 400,
			}
		}

		l, err := pdf.SaveToS3()
		if err != nil {
			stdError = &lerrors.StdError{
				Caller:     "handlers.Pdf.process",
				Code:       lerrors.CodeApplicationError,
				Err:        err,
				Msg:        err.Error(),
				StatusCode: 400,
			}
		} else {
			log.Infof("Saved pdf to: %s", l)
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

	// Create the response object
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
