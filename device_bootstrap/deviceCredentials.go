package devicecontrol

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"net/http"
	"strings"
)

/*
DeviceCredentials represent cumulocity's 'application/vnd.com.nsn.cumulocity.deviceCredentials+json'.
See: https://cumulocity.com/guides/reference/device-credentials/#devicecredentials-application-vnd-com-nsn-cumulocity-devicecredentials-json
*/
type DeviceCredentials struct {
	ID       string `json:"id"`
	TenantID string `json:"tenantId,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Self     string `json:"self,omitempty"`
}

var deviceCredentialsContentType = "application/vnd.com.nsn.cumulocity.deviceCredentials+json"

/*
CreateDeviceCredentials creates 'DeviceCredentials' and returns the generated credentials

Return created 'DeviceCredentials' on success.
Can return the following errors:
- generic.BadCredentialsErr (invalid username / password / host combination)
- generic.AccessDeniedErr (missing user rights)
- generic.Error (generic cloud error)
- error (unexpected)

See: https://cumulocity.com/guides/reference/device-credentials/#post-creates-a-device-credentials-request
*/
func (c Client) CreateDeviceCredentials(deviceID string) (DeviceCredentials, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(DeviceCredentials{
		ID: deviceID,
	})
	if err != nil {
		return DeviceCredentials{}, fmt.Errorf("failed to encode request body: %w", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/devicecontrol/deviceCredentials", c.BaseURL),
		&buf,
	)
	if err != nil {
		return DeviceCredentials{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(c.Username, c.Password)

	h := req.Header
	h.Add("Content-Type", deviceCredentialsContentType)
	h.Add("Accept", deviceCredentialsContentType)
	req.Header = h

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return DeviceCredentials{}, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return DeviceCredentials{}, generic.BadCredentialsErr
		case http.StatusForbidden:
			return DeviceCredentials{}, generic.AccessDeniedErr
		default:
			if strings.HasPrefix(resp.Header.Get("Content-Type"), generic.ErrorContentType) {
				var errResp generic.Error
				err := json.NewDecoder(resp.Body).Decode(&errResp)
				if err != nil {
					return DeviceCredentials{}, fmt.Errorf("failed to decode error response body: %s", err)
				}
				return DeviceCredentials{}, fmt.Errorf("failed to create device-credentials (%d): %w", resp.StatusCode, errResp)
			}
			return DeviceCredentials{}, fmt.Errorf("failed to create device-credentials with status code %d", resp.StatusCode)
		}
	}

	var respBody DeviceCredentials
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return DeviceCredentials{}, fmt.Errorf("failed to decode response body: %w", err)
	}
	return respBody, nil
}
