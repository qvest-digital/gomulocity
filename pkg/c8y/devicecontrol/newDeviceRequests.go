package devicecontrol

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tarent/gomulocity/pkg/c8y/meta"
	"net/http"
	"strings"
)

var NewDeviceRequestAlreadyExistsErr = errors.New("'newDeviceRequest' with ID already exists")

/*
NewDeviceRequest represent cumulocity's 'application/vnd.com.nsn.cumulocity.NewDeviceRequest+json'.
See: https://cumulocity.com/guides/reference/device-credentials/#newdevicerequest-application-vnd-com-nsn-cumulocity-newdevicerequest-json
*/
type NewDeviceRequest struct {
	ID     string `json:"id,omitempty"`
	Status string `json:"status,omitempty"`
	Self   string `json:"self,omitempty"`
}

var newDeviceRequestContentType = "application/vnd.com.nsn.cumulocity.NewDeviceRequest+json"

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

var newDeviceRequestCollectionContentType = "application/vnd.com.nsn.cumulocity.newDeviceRequestCollection+json"

/*
CreateNewDeviceRequest creates a 'newDeviceRequest' with the given id.

Return created 'newDeviceRequest' on success.
Can return the following errors:
- meta.BadCredentialsErr (invalid username / password / host combination)
- meta.AccessDeniedErr (missing user rights)
- meta.Error (generic cloud error)
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
	h.Add("Content-Type", newDeviceRequestContentType)
	h.Add("Accept", newDeviceRequestContentType)
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
			var errResp meta.Error
			err := json.NewDecoder(resp.Body).Decode(&errResp)
			if err != nil {
				return NewDeviceRequest{}, fmt.Errorf("failed to decode error response body: %s", err)
			}
			if errResp.ErrorType == "devicecontrol/Non Unique Result" {
				return NewDeviceRequest{}, NewDeviceRequestAlreadyExistsErr
			}
			return NewDeviceRequest{}, errResp
		default:
			if strings.HasPrefix(resp.Header.Get("Content-Type"), meta.ErrorContentType) {
				var errResp meta.Error
				err := json.NewDecoder(resp.Body).Decode(&errResp)
				if err != nil {
					return NewDeviceRequest{}, fmt.Errorf("failed to decode error response body: %s", err)
				}
				return NewDeviceRequest{}, fmt.Errorf("failed to create new-device-reuest (%d): %w", resp.StatusCode, errResp)
			}
			return NewDeviceRequest{}, fmt.Errorf("failed to create new-device-request with status code %d", resp.StatusCode)
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
- meta.Error (generic cloud error)
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
	h.Add("Accept", newDeviceRequestCollectionContentType)
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
			if strings.HasPrefix(resp.Header.Get("Content-Type"), meta.ErrorContentType) {
				var errResp meta.Error
				err := json.NewDecoder(resp.Body).Decode(&errResp)
				if err != nil {
					return NewDeviceRequestCollection{}, fmt.Errorf("failed to decode error response body: %s", err)
				}
				return NewDeviceRequestCollection{}, fmt.Errorf("failed to find-all new-device-requests (%d): %w", resp.StatusCode, errResp)
			}
			return NewDeviceRequestCollection{}, fmt.Errorf("failed to find-all new-device-requests with status code %d", resp.StatusCode)
		}
	}

	var respBody NewDeviceRequestCollection
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return NewDeviceRequestCollection{}, fmt.Errorf("failed to decode response body: %w", err)
	}
	return respBody, nil
}

/*
UpdateNewDeviceRequest updates status of 'newDeviceRequest' with given ID.

Can return the following errors:
- meta.BadCredentialsErr (invalid username / password / host combination)
- meta.AccessDeniedErr (missing user rights)
- meta.Error (generic cloud error)
- error (unexpected)

See: https://cumulocity.com/guides/reference/device-credentials/#put-updates-a-new-device-request
*/
func (c Client) UpdateNewDeviceRequest(id, status string) (NewDeviceRequest, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(NewDeviceRequest{
		Status: status,
	})
	if err != nil {
		return NewDeviceRequest{}, fmt.Errorf("failed to encode request body: %w", err)
	}

	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/devicecontrol/newDeviceRequests/%s", c.BaseURL, id),
		&buf,
	)
	if err != nil {
		return NewDeviceRequest{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(c.Username, c.Password)

	h := req.Header
	h.Add("Content-Type", newDeviceRequestContentType)
	h.Add("Accept", newDeviceRequestContentType)
	req.Header = h

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return NewDeviceRequest{}, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return NewDeviceRequest{}, meta.BadCredentialsErr
		case http.StatusForbidden:
			return NewDeviceRequest{}, meta.AccessDeniedErr
		default:
			if strings.HasPrefix(resp.Header.Get("Content-Type"), meta.ErrorContentType) {
				var errResp meta.Error
				err := json.NewDecoder(resp.Body).Decode(&errResp)
				if err != nil {
					return NewDeviceRequest{}, fmt.Errorf("failed to decode error response body: %s", err)
				}
				return NewDeviceRequest{}, fmt.Errorf("failed to update new-device-request (%d): %w", resp.StatusCode, errResp)
			}
			return NewDeviceRequest{}, fmt.Errorf("failed to update new-device-request with status code %d", resp.StatusCode)
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
DeleteNewDeviceRequest deletes 'newDeviceRequest' with given ID.

Can return the following errors:
- meta.BadCredentialsErr (invalid username / password / host combination)
- meta.AccessDeniedErr (missing user rights)
- meta.Error (generic cloud error)
- error (unexpected)

See: https://cumulocity.com/guides/reference/device-credentials/#delete-deletes-a-new-device-request
*/
func (c Client) DeleteNewDeviceRequest(id string) error {
	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/devicecontrol/newDeviceRequests/%s", c.BaseURL, id),
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(c.Username, c.Password)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return meta.BadCredentialsErr
		case http.StatusForbidden:
			return meta.AccessDeniedErr
		default:
			if strings.HasPrefix(resp.Header.Get("Content-Type"), meta.ErrorContentType) {
				var errResp meta.Error
				err := json.NewDecoder(resp.Body).Decode(&errResp)
				if err != nil {
					return fmt.Errorf("failed to decode error response body: %s", err)
				}
				return fmt.Errorf("failed to delete new-device-request (%d): %w", resp.StatusCode, errResp)
			}
			return fmt.Errorf("failed to delete new-device-request with status code %d", resp.StatusCode)
		}
	}

	return nil
}
