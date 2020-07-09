package generic

import (
	"encoding/json"
	"errors"
	"fmt"
)

var BadCredentialsErr = errors.New("bad credentials")

/*
Error represent cumulocity's 'application/vnd.com.nsn.cumulocity.error+json' without 'Error details'.
See: https://cumulocity.com/guides/reference/rest-implementation/#error-application-vnd-com-nsn-cumulocity-error-json
*/
type Error struct {
	ErrorType string `json:"error"`
	Message   string `json:"message"`
	Info      string `json:"info"`
}

func (e Error) Error() string {
	return fmt.Sprintf("request failed: %q %s. See: %s", e.ErrorType, e.Message, e.Info)
}

var ErrorContentType = "application/vnd.com.nsn.cumulocity.error+json"

func ClientError(message string, info string) *Error {
	return &Error{
		ErrorType: "ClientError",
		Message:   message,
		Info:      info,
	}
}

func CreateErrorFromResponse(responseBody []byte, status int) *Error {
	var error Error
	err := json.Unmarshal(responseBody, &error)
	if err != nil {
		error = *ClientError(fmt.Sprintf("Error while parsing response JSON [%s]: %s", responseBody, err.Error()), "CreateErrorFromResponse")
	}

	error.ErrorType = fmt.Sprintf("%d: %s", status, error.ErrorType)

	return &error
}
