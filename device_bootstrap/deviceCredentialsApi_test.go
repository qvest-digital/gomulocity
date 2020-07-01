package device_bootstrap

import (
	"fmt"
	jsoncompare "github.com/orasik/gocomparejson"
	"github.com/tarent/gomulocity/generic"
	"strings"

	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestDeviceCredentialsApi_CommonPropertiesOnCreate(t *testing.T) {
	var expectedContentType = "application/vnd.com.nsn.cumulocity.deviceCredentials+json"
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
		res.WriteHeader(http.StatusCreated)
		_, err = res.Write([]byte(`{"id": "4711"}`))
		if err != nil {
			t.Fatalf("failed to write resp body: %s", err)
		}
	}))
	defer testServer.Close()

	u := "<username>"
	p := "<password>"
	c := &generic.Client{
		HTTPClient: testServer.Client(),
		BaseURL:    testServer.URL,
		Username:   u,
		Password:   p,
	}

	deviceCredentialsApi := NewDeviceCredentialsApi(c)
	deviceCredentialsApi.Create("4711")

	if reqAccept != expectedContentType {
		t.Errorf("unexpected request accept header. Expected %q. Given: %q", expectedContentType, reqAccept)
	}
	if reqContentType != expectedContentType {
		t.Errorf("unexpected request content-type header. Expected %q. Given: %q", expectedContentType, reqContentType)
	}

	if reqBasicAuthUsername != u {
		t.Errorf("unexpected c8y request basic-auth username. Expected %q. Given: %q", u, reqBasicAuthUsername)
	}
	if reqBasicAuthPassword != p {
		t.Errorf("unexpected c8y request basic-auth password. Expected %q. Given: %q", p, reqBasicAuthPassword)
	}

	var expectedC8YRequestURL = "/devicecontrol/deviceCredentials"
	if reqURL != expectedC8YRequestURL {
		t.Errorf("unexpected c8y request url. Expected %q. Given: %q", expectedC8YRequestURL, reqURL)
	}
}

func TestDeviceCredentialsApi_Create(t *testing.T) {
	tests := []struct {
		name                      string
		deviceID                  string
		c8yRespCode               int
		c8yRespContentType        string
		c8yRespBody               string
		c8yExpectedRequestBody    string
		expectedDeviceCredentials *DeviceCredentials
		expectedErr               *generic.Error
	}{
		{
			name:                   "happy case",
			deviceID:               "4711",
			c8yRespCode:            http.StatusCreated,
			c8yRespBody:            `{"id": "4711", "tenantId" : "test", "username" : "device_4711", "password" : "3rasfst4swfa", "self": "https://myFancyCloudInstance.com/devicecontrol/deviceCredentials"}`,
			c8yExpectedRequestBody: `{"id": "4711"}`,
			expectedDeviceCredentials: &DeviceCredentials{
				ID:       "4711",
				TenantID: "test",
				Username: "device_4711",
				Password: "3rasfst4swfa",
				Self:     "https://myFancyCloudInstance.com/devicecontrol/deviceCredentials",
			},
			expectedErr: nil,
		}, {
			name:               "bad credentials",
			deviceID:           "401",
			c8yRespCode:        http.StatusUnauthorized,
			c8yRespContentType: "application/vnd.com.nsn.cumulocity.error+json",
			c8yRespBody: `{
				"error": "security/Unauthorized",
				"message": "Invalid credentials! : Bad credentials",
				"info": "https://www.cumulocity.com/guides/reference-guide/#error_reporting"
			}`,
			expectedErr: &generic.Error{
				ErrorType: "security/Unauthorized",
				Message:   "Invalid credentials! : Bad credentials",
				Info:      "https://www.cumulocity.com/guides/reference-guide/#error_reporting",
			},
			c8yExpectedRequestBody: `{"id": "401"}`,
		}, {
			name:        "access denied",
			deviceID:    "403",
			c8yRespCode: http.StatusForbidden,
			c8yRespBody: `{    
					"error": "security/Forbidden",
					"message": "Access is denied",
					"info": "https://www.cumulocity.com/guides/reference-guide/#error_reporting"
				}`,
			expectedErr: &generic.Error{
				ErrorType: "security/Forbidden",
				Message:   "Access is denied",
				Info:      "https://www.cumulocity.com/guides/reference-guide/#error_reporting",
			},
			c8yExpectedRequestBody: `{"id": "403"}`,
		}, {
			name:        "without deviceId",
			c8yRespCode: http.StatusNotFound,
			c8yRespBody: `{
					"error": "devicecontrol/Not Found",
					"message": "There is no newDeviceRequest for device id .",
					"info": "https://www.cumulocity.com/guides/reference-guide/#error_reporting"
				}`,
			expectedErr: &generic.Error{
				ErrorType: "devicecontrol/Not Found",
				Message:   "There is no newDeviceRequest for device id .",
				Info:      "https://www.cumulocity.com/guides/reference-guide/#error_reporting",
			},
			c8yExpectedRequestBody: `{"id": ""}`,
		}, {
			name:        "invalid json response",
			deviceID:    "4711",
			c8yRespCode: http.StatusCreated,
			c8yRespBody: `#`,
			expectedErr: &generic.Error{
				ErrorType: "ClientError",
				Message:   "Error while parsing response JSON: invalid character '#' looking for beginning of value",
				Info:      "ResponseParser",
			},
			c8yExpectedRequestBody: `{"id": "4711"}`,
		}, {
			name:        "empty json response",
			deviceID:    "4711",
			c8yRespCode: http.StatusCreated,
			expectedErr: &generic.Error{
				ErrorType: "ClientError",
				Message:   "Response body was empty",
				Info:      "GetDeviceCredentials",
			},
			c8yExpectedRequestBody: `{"id": "4711"}`,
		}, {
			name:     "post error",
			deviceID: "4711",
			expectedErr: &generic.Error{
				ErrorType: "ClientError",
				Message:   "Error while posting new device credentials: Post <dynamic-URL>/devicecontrol/deviceCredentials: EOF",
				Info:      "CreateDeviceCredentials",
			},
			c8yExpectedRequestBody: `{"id": "4711"}`,
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

			u := "<username>"
			p := "<password>"
			c := &generic.Client{
				HTTPClient: testServer.Client(),
				BaseURL:    testServer.URL,
				Username:   u,
				Password:   p,
			}

			deviceCredentialsApi := NewDeviceCredentialsApi(c)

			deviceCredentials, err := deviceCredentialsApi.Create(tt.deviceID)

			if equal, _ := jsoncompare.CompareJSON(reqBody, tt.c8yExpectedRequestBody); !equal {
				t.Errorf("unexpected c8y request body. Expected %q. Given: %q", tt.c8yExpectedRequestBody, reqBody)
			}

			setDynamicUrl(tt.expectedErr, testServer.URL)
			if fmt.Sprint(err) != fmt.Sprint(tt.expectedErr) {
				t.Fatalf("respond with unexpected error. \nExpected: %s\nGiven:    %s", tt.expectedErr, err)
			}

			if !reflect.DeepEqual(deviceCredentials, tt.expectedDeviceCredentials) {
				t.Errorf("respond with unexpected deviceCredentials. \nExpected: %#v. \nGiven: %#v", tt.expectedDeviceCredentials, deviceCredentials)
			}
		})
	}
}

func setDynamicUrl(err *generic.Error, url string) {
	if err != nil {
		err.Message = strings.ReplaceAll(err.Message, "<dynamic-URL>", url)
	}
}
