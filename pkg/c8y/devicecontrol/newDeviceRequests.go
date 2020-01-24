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
NewDeviceRequest represent cumulocity's 'application/vnd.com.nsn.cumulocity.NewDeviceRequest+json' without 'PagingStatistics'.
See: https://cumulocity.com/guides/reference/device-credentials/#newdevicerequest-application-vnd-com-nsn-cumulocity-newdevicerequest-json
*/
type NewDeviceRequest struct {
	ID     string `json:"id"`
	Status string `json:"status,omitempty"`
	Self   string `json:"self,omitempty"`
}

var newDeviceRequestType = "application/vnd.com.nsn.cumulocity.NewDeviceRequest+json"

/*
/*
NewDeviceRequestCollection represent cumulocity's 'application/vnd.com.nsn.cumulocity.newDeviceRequestCollection+json'.
See: https://cumulocity.com/guides/reference/device-credentials/#newdevicerequestcollection-application-vnd-com-nsn-cumulocity-newdevicerequestcollection-json
*/
type NewDeviceRequestCollection struct {
	Self              string                `json:"self"`
	NewDeviceRequests []NewDeviceRequest    `json:"newDeviceRequests"`
	Statistics        meta.PagingStatistics `json:"statistics"`
	Prev              string                `json:"prev"`
	Next              string                `json:"next"`
}

var newDeviceRequestCollectionType = "application/vnd.com.nsn.cumulocity.newDeviceRequestCollection+json"

/*
CreateNewDeviceRequest creates a 'newDeviceRequest' with the given id.

Return created 'newDeviceRequest' on success.
Can return the following errors:
- meta.BadCredentialsErr (invalid username / password / host combination)
- meta.AccessDeniedErr (missing user rights)
- NewDeviceRequestAlreadyExistsErr ('newDeviceRequest' with given id already exists)
- error (unexpected)

See: https://cumulocity.com/guides/reference/device-credentials/#post-create-a-new-device-request
*/
func (c Client) CreateNewDeviceRequest(id string) (NewDeviceRequest, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(NewDeviceRequest{
		ID: id,
	})
	if err != nil {
		return NewDeviceRequest{}, fmt.Errorf("failed to encode request body: %w", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/devicecontrol/newDeviceRequests", c.BaseURL),
		&buf,
	)
	if err != nil {
		return NewDeviceRequest{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(c.Username, c.Password)

	h := req.Header
	h.Add("Content-Type", newDeviceRequestType)
	h.Add("Accept", newDeviceRequestType)
	req.Header = h

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return NewDeviceRequest{}, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return NewDeviceRequest{}, meta.BadCredentialsErr
		case http.StatusForbidden:
			return NewDeviceRequest{}, meta.AccessDeniedErr
		case http.StatusUnprocessableEntity:
			var errRespBody meta.ErrorBody
			err := json.NewDecoder(resp.Body).Decode(&errRespBody)
			if err != nil {
				return NewDeviceRequest{}, fmt.Errorf("failed to decode error response body: %s", err)
			}
			if errRespBody.Error == "devicecontrol/Non Unique Result" {
				return NewDeviceRequest{}, NewDeviceRequestAlreadyExistsErr
			}
			return NewDeviceRequest{}, fmt.Errorf("failed to create new-device-request. Status: %d: %q %s See: %s",
				resp.StatusCode, errRespBody.Error, errRespBody.Message, errRespBody.Info)
		default:
			var errRespBody meta.ErrorBody
			err := json.NewDecoder(resp.Body).Decode(&errRespBody)
			if err != nil {
				return NewDeviceRequest{}, fmt.Errorf("failed to decode error response body: %s", err)
			}
			return NewDeviceRequest{}, fmt.Errorf("failed to create new-device-request. Status: %d: %q %s See: %s",
				resp.StatusCode, errRespBody.Error, errRespBody.Message, errRespBody.Info)
		}
	}

	var respBody NewDeviceRequest
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return NewDeviceRequest{}, fmt.Errorf("failed to decode response body: %w", err)
	}
	return respBody, nil
}

/*
NewDeviceRequests find all 'newDeviceRequest's.
Note: for paging use meta.Page(int) as reqOpts.

Return created 'NewDeviceRequestCollection' on success.
Can return the following errors:
- meta.BadCredentialsErr (invalid username / password / host combination)
- meta.AccessDeniedErr (missing user rights)
- error (unexpected)

See: https://cumulocity.com/guides/reference/device-credentials/#get-returns-all-new-device-requests
*/
func (c Client) NewDeviceRequests(reqOpts ...func(*http.Request)) (NewDeviceRequestCollection, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/devicecontrol/newDeviceRequests", c.BaseURL),
		nil,
	)
	if err != nil {
		return NewDeviceRequestCollection{}, fmt.Errorf("failed to create request: %w", err)
	}

	for _, opt := range reqOpts {
		if opt != nil {
			opt(req)
		}
	}

	req.SetBasicAuth(c.Username, c.Password)

	h := req.Header
	h.Add("Accept", newDeviceRequestCollectionType)
	req.Header = h

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return NewDeviceRequestCollection{}, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return NewDeviceRequestCollection{}, meta.BadCredentialsErr
		case http.StatusForbidden:
			return NewDeviceRequestCollection{}, meta.AccessDeniedErr
		default:
			var errRespBody meta.ErrorBody
			err := json.NewDecoder(resp.Body).Decode(&errRespBody)
			if err != nil {
				return NewDeviceRequestCollection{}, fmt.Errorf("failed to decode error response body: %s", err)
			}
			return NewDeviceRequestCollection{}, fmt.Errorf("failed to find-all new-device-requests. Status: %d: %q %s See: %s",
				resp.StatusCode, errRespBody.Error, errRespBody.Message, errRespBody.Info)
		}
	}

	var respBody NewDeviceRequestCollection
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return NewDeviceRequestCollection{}, fmt.Errorf("failed to decode response body: %w", err)
	}
	return respBody, nil
}
