package devicecontrol

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tarent/gomulocity/pkg/c8y/meta"
	"net/http"
)

var NewDeviceRequestAlreadyExistsErr = errors.New("'newDeviceRequest' with ID already exists")

/*
NewDeviceRequest represent cumulocity's 'application/vnd.com.nsn.cumulocity.NewDeviceRequest+json'.
See: https://cumulocity.com/guides/reference/device-credentials/#newdevicerequest-application-vnd-com-nsn-cumulocity-newdevicerequest-json
*/
type NewDeviceRequest struct {
	ID     string `json:"id"`
	Status string `json:"status,omitempty"`
	Self   string `json:"self,omitempty"`
}

var newDeviceRequestType = "application/vnd.com.nsn.cumulocity.NewDeviceRequest+json"

/*
CreateNewDeviceRequest creates a 'newDeviceRequest' with the given id into configured c8y instance.
Return created 'newDeviceRequest' on success.

Can return the following errors:
- meta.BadCredentialsErr (invalid username / password / host combination)
- meta.AccessDeniedErr (missing user rights)
- NewDeviceRequestAlreadyExistsErr ('newDeviceRequest' with given id already exists)

See: https://cumulocity.com/guides/reference/device-credentials/#post-create-a-new-device-request
*/
func (c Client) CreateNewDeviceRequest(id string) (NewDeviceRequest, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(NewDeviceRequest{
		ID: id,
	})
	if err != nil {
		return NewDeviceRequest{}, fmt.Errorf("failed to encode new-device-request body: %w", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/devicecontrol/newDeviceRequests", c.BaseURL),
		&buf,
	)
	if err != nil {
		return NewDeviceRequest{}, fmt.Errorf("failed to create new-device-request request: %w", err)
	}

	req.SetBasicAuth(c.Username, c.Password)

	h := req.Header
	h.Add("Content-Type", newDeviceRequestType)
	h.Add("Accept", newDeviceRequestType)
	req.Header = h

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return NewDeviceRequest{}, fmt.Errorf("failed to execute new-device-request request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		var errRespBody meta.ErrorBody
		err := json.NewDecoder(resp.Body).Decode(&errRespBody)
		if err != nil {
			return NewDeviceRequest{}, fmt.Errorf("failed to decode new-device-request error response body: %s", err)
		}

		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return NewDeviceRequest{}, meta.BadCredentialsErr
		case http.StatusForbidden:
			return NewDeviceRequest{}, meta.AccessDeniedErr
		case http.StatusUnprocessableEntity:
			if errRespBody.Error == "devicecontrol/Non Unique Result" {
				fmt.Printf("## %#v ###\n", errRespBody)
				return NewDeviceRequest{}, NewDeviceRequestAlreadyExistsErr
			}
			return NewDeviceRequest{}, fmt.Errorf("failed to create new-device-request. Status: %d: %q %s See: %s",
				resp.StatusCode, errRespBody.Error, errRespBody.Message, errRespBody.Info)
		default:
			return NewDeviceRequest{}, fmt.Errorf("failed to create new-device-request. Status: %d: %q %s See: %s",
				resp.StatusCode, errRespBody.Error, errRespBody.Message, errRespBody.Info)
		}
	}

	var respBody NewDeviceRequest
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return NewDeviceRequest{}, fmt.Errorf("failed to decode new-device-request response body: %w", err)
	}
	return respBody, nil
}
