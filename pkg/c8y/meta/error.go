package meta

import "errors"

var BadCredentialsErr = errors.New("bad credentials")
var AccessDeniedErr = errors.New("access denied")

/*
ErrorBody represent cumulocity's 'application/vnd.com.nsn.cumulocity.error+json' without 'Error details'.
See: https://cumulocity.com/guides/reference/rest-implementation/#error-application-vnd-com-nsn-cumulocity-error-json
*/
type ErrorBody struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Info    string `json:"info"`
}
