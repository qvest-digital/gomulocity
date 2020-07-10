package device_bootstrap

import (
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeviceRegistrationApi_CommonPropertiesOnDelete(t *testing.T) {
	var reqBasicAuthUsername, reqBasicAuthPassword, reqURL, reqAccept, reqContentType string

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		reqURL = req.URL.String()
		_, err := ioutil.ReadAll(req.Body)
		if err != nil {
			t.Fatalf("failed to read c8y request body: %s", err)
		}
		defer req.Body.Close()

		reqBasicAuthUsername, reqBasicAuthPassword, _ = req.BasicAuth()
		reqAccept = req.Header.Get("Accept")
		reqContentType = req.Header.Get("Content-Type")
		res.WriteHeader(http.StatusCreated)
		_, err = res.Write(nil)
		if err != nil {
			t.Fatalf("failed to write resp body: %s", err)
		}
	}))
	defer testServer.Close()

	deviceRegistrationApi := buildDeviceRegistrationApi(testServer)
	deviceRegistrationApi.Delete("4711")

	if reqAccept != "" {
		t.Errorf("unexpected request accept header. Expected none. Given: %q", reqAccept)
	}
	if reqContentType != "" {
		t.Errorf("unexpected request content-type header. Expected none. Given: %q", reqContentType)
	}

	if reqBasicAuthUsername != USER {
		t.Errorf("unexpected c8y request basic-auth username. Expected %q. Given: %q", USER, reqBasicAuthUsername)
	}
	if reqBasicAuthPassword != PASSWORD {
		t.Errorf("unexpected c8y request basic-auth password. Expected %q. Given: %q", PASSWORD, reqBasicAuthPassword)
	}

	var expectedC8YRequestURL = "/devicecontrol/newDeviceRequests/4711"
	if reqURL != expectedC8YRequestURL {
		t.Errorf("unexpected c8y request url. Expected %q. Given: %q", expectedC8YRequestURL, reqURL)
	}
}

func TestDeviceRegistrationApi_Delete(t *testing.T) {
	tests := []struct {
		name               string
		deviceId           string
		c8yRespCode        int
		c8yRespContentType string
		c8yRespBody        string
		expectedErr        *generic.Error
	}{
		{
			name:        "happy case",
			deviceId:    "4711",
			c8yRespCode: http.StatusNoContent,
			expectedErr: nil,
		}, {
			name:               "bad credentials",
			deviceId:           "401",
			c8yRespCode:        http.StatusUnauthorized,
			c8yRespContentType: "application/vnd.com.nsn.cumulocity.error+json",
			c8yRespBody: `{
				"error": "security/Unauthorized",
				"message": "Invalid credentials! : Bad credentials",
				"info": "https://www.cumulocity.com/guides/reference-guide/#error_reporting"
			}`,
			expectedErr: &generic.Error{
				ErrorType: "401: security/Unauthorized",
				Message:   "Invalid credentials! : Bad credentials",
				Info:      "https://www.cumulocity.com/guides/reference-guide/#error_reporting",
			},
		}, {
			name:        "invalid json error response",
			deviceId:    "4711",
			c8yRespCode: http.StatusInternalServerError,
			c8yRespBody: `#`,
			expectedErr: &generic.Error{
				ErrorType: "500: ClientError",
				Message:   "Error while parsing response JSON [#]: invalid character '#' looking for beginning of value",
				Info:      "CreateErrorFromResponse",
			},
		}, {
			name:        "without deviceId",
			c8yRespCode: http.StatusMethodNotAllowed,
			expectedErr: &generic.Error{
				ErrorType: "ClientError",
				Message:   "Deleting deviceRegistrations without an id is not allowed",
				Info:      "DeleteDeviceRegistration",
			},
		}, {
			name:        "invalid json response",
			deviceId:    "4711",
			c8yRespCode: http.StatusInternalServerError,
			c8yRespBody: `#`,
			expectedErr: &generic.Error{
				ErrorType: "500: ClientError",
				Message:   "Error while parsing response JSON [#]: invalid character '#' looking for beginning of value",
				Info:      "CreateErrorFromResponse",
			},
		}, {
			name:     "post error",
			deviceId: "4711",
			expectedErr: &generic.Error{
				ErrorType: "ClientError",
				Message:   "Error while deleting a deviceRegistration with id 4711: Delete <dynamic-URL>/devicecontrol/newDeviceRequests/4711: EOF",
				Info:      "DeleteDeviceRegistration",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqBody string

			testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				res.Header().Set("Content-Type", tt.c8yRespContentType)
				reqBodyBytes, err := ioutil.ReadAll(req.Body)
				if err != nil {
					t.Fatalf("failed to read c8y request body: %s", err)
				}
				defer req.Body.Close()
				reqBody = string(reqBodyBytes)

				res.WriteHeader(tt.c8yRespCode)
				_, err = res.Write([]byte(tt.c8yRespBody))
				if err != nil {
					t.Fatalf("failed to write resp body: %s", err)
				}
			}))
			defer testServer.Close()

			deviceRegistrationApi := buildDeviceRegistrationApi(testServer)
			err := deviceRegistrationApi.Delete(tt.deviceId)

			if len(reqBody) > 0 {
				t.Errorf("unexpected c8y request body. Expected none. Given: %q", reqBody)
			}

			setDynamicUrl(tt.expectedErr, testServer.URL)
			if fmt.Sprint(err) != fmt.Sprint(tt.expectedErr) {
				t.Fatalf("respond with unexpected error. \nExpected: %s\nGiven:    %s", tt.expectedErr, err)
			}
		})
	}
}
