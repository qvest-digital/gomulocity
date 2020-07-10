package device_bootstrap

import (
	"fmt"
	jsoncompare "github.com/orasik/gocomparejson"
	"github.com/tarent/gomulocity/generic"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestDeviceRegistrationApi_CommonPropertiesOnUpdate(t *testing.T) {
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
	deviceRegistrationApi.Update("4711", ACCEPTED)

	if reqAccept != expectedContentType {
		t.Errorf("unexpected request accept header. Expected %q. Given: %q", expectedContentType, reqAccept)
	}
	if reqContentType != expectedContentType {
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

func TestDeviceRegistrationApi_Update(t *testing.T) {
	tests := []struct {
		name                       string
		deviceId                   string
		deviceRegistrationStatus   Status
		c8yRespCode                int
		c8yRespContentType         string
		c8yRespBody                string
		c8yExpectedRequestBody     string
		expectedDeviceRegistration *DeviceRegistration
		expectedErr                *generic.Error
	}{
		{
			name:                     "happy case",
			deviceId:                 "4711",
			deviceRegistrationStatus: ACCEPTED,
			c8yRespCode:              http.StatusOK,
			c8yRespBody: `{
				"id": "4711", 
				"status": "ACCEPTED", 
				"self": "https://myFancyCloudInstance.com/devicecontrol/deviceCredentials/4711",
				"owner": "me@company.com",
				"customProperties": {},
				"creationTime": "2020-07-03T10:16:35.870+02:00",
				"tenantId": "myCloud"
			}`,
			c8yExpectedRequestBody: `{"status": "ACCEPTED"}`,
			expectedDeviceRegistration: &DeviceRegistration{
				Id:               "4711",
				Status:           ACCEPTED,
				Self:             "https://myFancyCloudInstance.com/devicecontrol/deviceCredentials/4711",
				Owner:            "me@company.com",
				CustomProperties: map[string]interface{}{},
				CreationTime:     &deviceRegistrationTime,
				TenantId:         "myCloud",
			},
			expectedErr: nil,
		}, {
			name:                     "bad credentials",
			deviceId:                 "401",
			deviceRegistrationStatus: ACCEPTED,
			c8yRespCode:              http.StatusUnauthorized,
			c8yRespContentType:       "application/vnd.com.nsn.cumulocity.error+json",
			c8yExpectedRequestBody:   `{"status": "ACCEPTED"}`,
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
			name:                     "invalid json error response",
			deviceId:                 "4711",
			deviceRegistrationStatus: ACCEPTED,
			c8yRespCode:              http.StatusInternalServerError,
			c8yRespBody:              `#`,
			expectedErr: &generic.Error{
				ErrorType: "500: ClientError",
				Message:   "Error while parsing response JSON [#]: invalid character '#' looking for beginning of value",
				Info:      "CreateErrorFromResponse",
			},
			c8yExpectedRequestBody: `{"status": "ACCEPTED"}`,
		}, {
			name:                     "without deviceId",
			deviceRegistrationStatus: ACCEPTED,
			c8yRespCode:              http.StatusUnprocessableEntity,
			expectedErr: &generic.Error{
				ErrorType: "ClientError",
				Message:   "Updating a deviceRegistration without an id is not allowed",
				Info:      "UpdateDeviceRegistration",
			},
			c8yExpectedRequestBody: `{}`,
		}, {
			name:     "without status",
			deviceId: "4711",
			//deviceRegistrationStatus: ACCEPTED,
			c8yRespCode: http.StatusUnprocessableEntity,
			c8yRespBody: `{
				"error": "undefined/validationError",
				"message": "Following mandatory fields should be included: status",
				"info": "https://www.cumulocity.com/guides/reference-guide/#error_reporting"
			}`,
			expectedErr: &generic.Error{
				ErrorType: "422: undefined/validationError",
				Message:   "Following mandatory fields should be included: status",
				Info:      "https://www.cumulocity.com/guides/reference-guide/#error_reporting",
			},
			c8yExpectedRequestBody: `{}`,
		}, {
			name:                     "invalid json response",
			deviceId:                 "4711",
			deviceRegistrationStatus: ACCEPTED,
			c8yRespCode:              http.StatusOK,
			c8yRespBody:              `#`,
			expectedErr: &generic.Error{
				ErrorType: "ClientError",
				Message:   "Error while parsing response JSON: invalid character '#' looking for beginning of value",
				Info:      "ResponseParser",
			},
			c8yExpectedRequestBody: `{"status": "ACCEPTED"}`,
		}, {
			name:                     "empty json response",
			deviceId:                 "4711",
			deviceRegistrationStatus: ACCEPTED,
			c8yRespCode:              http.StatusOK,
			expectedErr: &generic.Error{
				ErrorType: "ClientError",
				Message:   "Response body was empty",
				Info:      "GetDeviceRegistration",
			},
			c8yExpectedRequestBody: `{"status": "ACCEPTED"}`,
		}, {
			name:                     "post error",
			deviceId:                 "4711",
			deviceRegistrationStatus: ACCEPTED,
			c8yRespCode:              http.StatusInternalServerError,
			expectedErr: &generic.Error{
				ErrorType: "500: ClientError",
				Message:   "Error while parsing response JSON []: unexpected end of JSON input",
				Info:      "CreateErrorFromResponse",
			},
			c8yExpectedRequestBody: `{"status": "ACCEPTED"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqBody = `{}`

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
			deviceRegistration, err := deviceRegistrationApi.Update(tt.deviceId, tt.deviceRegistrationStatus)

			if equal, _ := jsoncompare.CompareJSON(reqBody, tt.c8yExpectedRequestBody); !equal {
				t.Errorf("unexpected c8y request body. Expected %q. Given: %q", tt.c8yExpectedRequestBody, reqBody)
			}

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
