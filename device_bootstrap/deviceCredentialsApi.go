package device_bootstrap

import (
	"encoding/json"
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"net/http"
)

type DeviceCredentialsApi interface {
	// Creates new device credentials for a given id
	Create(deviceId string) (*DeviceCredentials, *generic.Error)
}

type deviceCredentialsApi struct {
	client   *generic.Client
	basePath string
}

// Creates a new device credentials api object
// client - Must be a gomulocity client.
// returns - The `device credentials`-api object
func NewDeviceCredentialsApi(client *generic.Client) DeviceCredentialsApi {
	return &deviceCredentialsApi{client, DEVICE_CREDENTIALS_API_PATH}
}

/*
CreateDeviceCredentials creates 'DeviceCredentials'

Return generated 'DeviceCredentials' on success.

See: https://cumulocity.com/guides/reference/device-credentials/#post-creates-a-device-credentials-request
*/
func (deviceCredentialsApi *deviceCredentialsApi) Create(deviceId string) (*DeviceCredentials, *generic.Error) {
	bytes, err := json.Marshal(DeviceCredentials{ID: deviceId})
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while marshalling the device credentials request: %s", err.Error()), "CreateDeviceCredentials")
	}
	headers := generic.AcceptHeader(DEVICE_CREDENTIALS_TYPE)
	contentType := generic.ContentTypeHeader(DEVICE_CREDENTIALS_TYPE)
	for k, v := range contentType {
		headers[k] = v
	}

	body, status, err := deviceCredentialsApi.client.Post(deviceCredentialsApi.basePath, bytes, headers)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while posting new device credentials: %s", err.Error()), "CreateDeviceCredentials")
	}
	if status != http.StatusCreated {
		return nil, generic.CreateErrorFromResponse(body)
	}

	return parseDeviceCredentialsResponse(body)
}


// -- internal

func parseDeviceCredentialsResponse(body []byte) (*DeviceCredentials, *generic.Error) {
	var result DeviceCredentials
	if len(body) > 0 {
		err := json.Unmarshal(body, &result)
		if err != nil {
			return nil, generic.ClientError(fmt.Sprintf("Error while parsing response JSON: %s", err.Error()), "ResponseParser")
		}
	} else {
		return nil, generic.ClientError("Response body was empty", "GetDeviceCredentials")
	}

	return &result, nil
}


