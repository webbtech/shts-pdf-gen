package handlers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-playground/validator/v10"
	lerrors "github.com/webbtech/shts-pdf-gen/errors"
	"github.com/webbtech/shts-pdf-gen/model"
)

var headers map[string]string = map[string]string{"Content-Type": "application/json"}

type Handler interface {
	Response(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
	process()
}

type response struct {
	Body       string
	Headers    map[string]string
	StatusCode int
}

type responseBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func validateInput(input *model.DocRequest) (err *lerrors.StdError) {

	var inputErrs []string
	validate := validator.New()

	// validate struct based on tags
	// see https://github.com/go-playground/validator
	valErr := validate.Struct(input)
	if valErr != nil {
		// for more on usage, see: https://github.com/go-playground/validator/blob/master/_examples/simple/main.go
		for _, err := range valErr.(validator.ValidationErrors) {
			inputErrs = append(inputErrs, fmt.Sprintf("Field validation for '%s' failed on the '%s' tag", err.Field(), err.Tag()))
		}
	}

	if input.RequestType != nil {
		_, found := findString(ValidRequestTypes, *input.RequestType)
		if !found {
			errStr := fmt.Sprintf(ERR_INVALID_TYPE+": \"%s\"", *input.RequestType)
			inputErrs = append(inputErrs, errStr)
		}
	}

	if len(inputErrs) > 0 {
		err = &lerrors.StdError{
			Caller:     "handlers.validateInput",
			Code:       lerrors.CodeBadInput,
			Err:        errors.New("Failed request input validation"),
			Msg:        strings.Join(inputErrs, "\n"),
			StatusCode: 400,
		}
		return err
	}

	return nil
}
