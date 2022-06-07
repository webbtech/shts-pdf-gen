package handlers

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/aws/aws-lambda-go/events"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/webbtech/shts-pdf-gen/config"
	lerrors "github.com/webbtech/shts-pdf-gen/errors"
	"github.com/webbtech/shts-pdf-gen/model"
	"github.com/webbtech/shts-pdf-gen/pdf"
)

const (
	ERR_INVALID_TYPE         = "Invalid request type in input"
	ERR_MISSING_REQUEST_BODY = "Missing request body"
)

var (
	Stage             string
	ValidRequestTypes = []string{"estimate", "invoice"}
)

// Post struct
type Post struct {
	Cfg      *config.Config
	Db       model.DbHandler
	input    *model.DocRequest
	request  events.APIGatewayProxyRequest
	response events.APIGatewayProxyResponse
}

// ========================== Public Methods =============================== //
func (p *Post) Response(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	p.request = request
	p.process()
	return p.response, nil
}

// ========================== Private Methods ============================== //

func (p *Post) process() {

	rb := responseBody{}
	var body []byte
	var err error
	var estimateRecord *model.Estimate
	var statusCode int = 201
	var stdError *lerrors.StdError

	// we're getting ExtJSON from the Realm function: createPDF,
	// so unmarshaling must be done on an EJSON formatted doc
	// see: https://www.mongodb.com/docs/manual/reference/mongodb-extended-json/
	bson.UnmarshalExtJSON([]byte(p.request.Body), true, &p.input)

	// Validate input
	if err := validateInput(p.input); err != nil {
		errors.As(err, &stdError)
	}

	// Fetch DB Record
	if stdError == nil {
		estimateRecord, err = p.Db.FetchEstimate(*p.input.EstimateNumber)
		if err != nil {
			stdError = &lerrors.StdError{
				Caller:     "handlers.Post.process",
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

		err = pdf.SaveToS3()
		if err != nil {
			stdError = &lerrors.StdError{
				Caller:     "handlers.Pdf.process",
				Code:       lerrors.CodeApplicationError,
				Err:        err,
				Msg:        err.Error(),
				StatusCode: 400,
			}
		} else {
			log.Info("Saved pdf to s3")
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
