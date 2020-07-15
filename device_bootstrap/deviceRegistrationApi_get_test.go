package device_bootstrap

import (
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestDeviceRegistrationApi_CommonPropertiesOnGet(t *testing.T) {
	var expectedContentType = "application/vnd.com.nsn.cumulocity.NewDeviceRequest+json"
	var reqBasicAuthUsername, reqBasicAuthPassword, reqURL, reqAccept, reqContentType string

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		reqURL = req.URL.String()
		res.Header().Set("Content-Type", expectedContentType)
		_, err := ioutil.ReadAll(req.Body)
		if err != nil {
			t.Fatalf("failed to read c8y request body: %s", err)
		}
		defer req.Body.Close()

		reqBasicAuthUsername, reqBasicAuthPassword, _ = req.BasicAuth()
		reqAccept = req.Header.Get("Accept")
		reqContentType = req.Header.Get("Content-Type")
		res.WriteHeader(http.StatusOK)
		_, err = res.Write([]byte(`{"id": "4711"}`))
		if err != nil {
			t.Fatalf("failed to write resp body: %s", err)
		}
	}))
	defer testServer.Close()

	deviceRegistrationApi := buildDeviceRegistrationApi(testServer)
	deviceRegistrationApi.Get("4711")

	if reqAccept != expectedContentType {
		t.Errorf("unexpected request accept header. Expected %q. Given: %q", expectedContentType, reqAccept)
	}
	if len(reqContentType) != 0 {
		t.Errorf("unexpected request content-type header. Expected %q. Given: %q", expectedContentType, reqContentType)
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

func TestDeviceRegistrationApi_Get(t *testing.T) {
	tests := []struct {
		name                       string
		deviceId                   string
		c8yRespCode                int
		c8yRespContentType         string
		c8yRespBody                string
		expectedDeviceRegistration *DeviceRegistration
		expectedErr                *generic.Error
	}{
		{
			name:        "happy case",
			deviceId:    "4711",
			c8yRespCode: http.StatusOK,
			c8yRespBody: `{
				"id": "4711", 
				"status": "PENDING_ACCEPTANCE", 
				"self": "https://myFancyCloudInstance.com/devicecontrol/newDeviceRequests/4711",
				"owner": "me@company.com",
				"customProperties": {},
				"creationTime": "2020-07-03T10:16:35.870+02:00",
				"tenantId": "myCloud"
			}`,
			expectedDeviceRegistration: &DeviceRegistration{
				Id:               "4711",
				Status:           PENDING_ACCEPTANCE,
				Self:             "https://myFancyCloudInstance.com/devicecontrol/newDeviceRequests/4711",
				Owner:            "me@company.com",
				CustomProperties: map[string]interface{}{},
				CreationTime:     &deviceRegistrationTime,
				TenantId:         "myCloud",
			},
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
			c8yRespCode: http.StatusNotAcceptable,
			expectedErr: &generic.Error{
				ErrorType: "ClientError",
				Message:   "Getting deviceRegistration without an id is not allowed",
				Info:      "GetDeviceRegistration",
			},
		}, {
			name:        "invalid json response",
			deviceId:    "4711",
			c8yRespCode: http.StatusOK,
			c8yRespBody: `#`,
			expectedErr: &generic.Error{
				ErrorType: "ClientError",
				Message:   "Error while parsing response JSON: invalid character '#' looking for beginning of value",
				Info:      "ResponseParser",
			},
		}, {
			name:        "empty json response",
			deviceId:    "4711",
			c8yRespCode: http.StatusOK,
			expectedErr: &generic.Error{
				ErrorType: "ClientError",
				Message:   "Response body was empty",
				Info:      "GetDeviceRegistration",
			},
		}, {
			name:        "post error",
			deviceId:    "4711",
			c8yRespCode: http.StatusInternalServerError,
			expectedErr: &generic.Error{
				ErrorType: "500: ClientError",
				Message:   "Error while parsing response JSON []: unexpected end of JSON input",
				Info:      "CreateErrorFromResponse",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				res.Header().Set("Content-Type", tt.c8yRespContentType)
				_, err := ioutil.ReadAll(req.Body)
				if err != nil {
					t.Fatalf("failed to read c8y request body: %s", err)
				}
				defer req.Body.Close()

				res.WriteHeader(tt.c8yRespCode)
				_, err = res.Write([]byte(tt.c8yRespBody))
				if err != nil {
					t.Fatalf("failed to write resp body: %s", err)
				}
			}))
			defer testServer.Close()

			deviceRegistrationApi := buildDeviceRegistrationApi(testServer)
			deviceRegistration, err := deviceRegistrationApi.Get(tt.deviceId)

			setDynamicUrl(tt.expectedErr, testServer.URL)
			if fmt.Sprint(err) != fmt.Sprint(tt.expectedErr) {
				t.Fatalf("respond with unexpected error. \nExpected: %s\nGiven:    %s", tt.expectedErr, err)
			}

			if !reflect.DeepEqual(deviceRegistration, tt.expectedDeviceRegistration) {
				t.Errorf("respond with unexpected deviceRegistration. \nExpected: %#v. \nGiven: %#v", tt.expectedDeviceRegistration, deviceRegistration)
			}
		})
	}
}
