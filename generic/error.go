package generic

import (
	"errors"
	"fmt"
)

var BadCredentialsErr = errors.New("bad credentials")
var AccessDeniedErr = errors.New("access denied")

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
	return fmt.Sprintf("request failed: %q %s See: %s", e.ErrorType, e.Message, e.Info)
}

var ErrorContentType = "application/vnd.com.nsn.cumulocity.error+json"
