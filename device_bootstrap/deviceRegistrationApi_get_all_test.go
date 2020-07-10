package device_bootstrap

import (
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestDeviceRegistrationApi_CommonPropertiesOnGetAll(t *testing.T) {
	var expectedContentType = "application/vnd.com.nsn.cumulocity.newDeviceRequestCollection+json"
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
		_, err = res.Write([]byte(`{}`))
		if err != nil {
			t.Fatalf("failed to write resp body: %s", err)
		}
	}))
	defer testServer.Close()

	deviceRegistrationApi := buildDeviceRegistrationApi(testServer)
	deviceRegistrationApi.GetAll(11)

	if reqAccept != expectedContentType {
		t.Errorf("unexpected request accept header. Expected %q. Given: %q", expectedContentType, reqAccept)
	}
	if len(reqContentType) != 0 {
		t.Errorf("unexpected request content-type header. Expected none. Given: %q", reqContentType)
	}

	if reqBasicAuthUsername != USER {
		t.Errorf("unexpected c8y request basic-auth username. Expected %q. Given: %q", USER, reqBasicAuthUsername)
	}
	if reqBasicAuthPassword != PASSWORD {
		t.Errorf("unexpected c8y request basic-auth password. Expected %q. Given: %q", PASSWORD, reqBasicAuthPassword)
	}

	var expectedC8YRequestURL = "/devicecontrol/newDeviceRequests?pageSize=11"
	if reqURL != expectedC8YRequestURL {
		t.Errorf("unexpected c8y request url. Expected %q. Given: %q", expectedC8YRequestURL, reqURL)
	}
}

func TestDeviceRegistrationApi_GetAll(t *testing.T) {
	tests := []struct {
		name                        string
		c8yRespCode                 int
		c8yRespContentType          string
		c8yRespBody                 string
		expectedDeviceRegistrations *DeviceRegistrationCollection
		expectedErr                 *generic.Error
	}{
		{
			name:        "happy case",
			c8yRespCode: http.StatusOK,
			c8yRespBody: `{
				"self": "selfURL", 
				"newDeviceRequests":[{
					"id": "4711", 
					"status": "PENDING_ACCEPTANCE", 
					"self": "https://myFancyCloudInstance.com/devicecontrol/newDeviceRequests/4711",
					"owner": "me@company.com",
					"customProperties": {},
					"creationTime": "2020-07-03T10:16:35.870+02:00",
					"tenantId": "myCloud"}], 
				"statistics": {
					"pageSize":11, 
					"currentPage":1
				}, 
				"next":"nextURL"
			}`,
			expectedDeviceRegistrations: &DeviceRegistrationCollection{
				Self: "selfURL",
				DeviceRegistrations: []DeviceRegistration{{
					Id:               "4711",
					Status:           PENDING_ACCEPTANCE,
					Self:             "https://myFancyCloudInstance.com/devicecontrol/newDeviceRequests/4711",
					Owner:            "me@company.com",
					CustomProperties: map[string]interface{}{},
					CreationTime:     &deviceRegistrationTime,
					TenantId:         "myCloud",
				}},
				Statistics: &generic.PagingStatistics{
					PageSize:    11,
					CurrentPage: 1,
				},
				Next: "nextURL",
			},
			expectedErr: nil,
		}, {
			name:               "bad credentials",
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
			c8yRespCode: http.StatusInternalServerError,
			c8yRespBody: `#`,
			expectedErr: &generic.Error{
				ErrorType: "500: ClientError",
				Message:   "Error while parsing response JSON [#]: invalid character '#' looking for beginning of value",
				Info:      "CreateErrorFromResponse",
			},
		}, {
			name:        "invalid json response",
			c8yRespCode: http.StatusOK,
			c8yRespBody: `#`,
			expectedErr: &generic.Error{
				ErrorType: "ClientError",
				Message:   "Error while parsing response JSON: invalid character '#' looking for beginning of value",
				Info:      "GetDeviceRegistrationCollection",
			},
		}, {
			name:        "empty json response",
			c8yRespCode: http.StatusOK,
			expectedErr: &generic.Error{
				ErrorType: "ClientError",
				Message:   "Response body was empty",
				Info:      "GetDeviceRegistrationCollection",
			},
		}, {
			name:        "post error",
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
			deviceRegistration, err := deviceRegistrationApi.GetAll(11)

			setDynamicUrl(tt.expectedErr, testServer.URL)
			if fmt.Sprint(err) != fmt.Sprint(tt.expectedErr) {
				t.Fatalf("respond with unexpected error. \nExpected: %s\nGiven:    %s", tt.expectedErr, err)
			}

			if !reflect.DeepEqual(deviceRegistration, tt.expectedDeviceRegistrations) {
				t.Errorf("respond with unexpected deviceRegistration. \nExpected: %#v. \nGiven: %#v", tt.expectedDeviceRegistrations, deviceRegistration)
			}
		})
	}
}

func TestDeviceRegistrationApi_GetAll_PageSize(t *testing.T) {
	tests := []struct {
		name        string
		reqPageSize int
		expectedErr *generic.Error
	}{
		{
			name:        "Negative",
			reqPageSize: -1,
			expectedErr: &generic.Error{
				ErrorType: "ClientError",
				Message:   "Error while building pageSize parameter to fetch deviceRegistrations: The page size must be between 1 and 2000. Was -1",
				Info:      "GetAllDeviceRegistrations",
			},
		}, {
			name:        "Zero",
			reqPageSize: 0,
			expectedErr: &generic.Error{
				ErrorType: "ClientError",
				Message:   "Error while building pageSize parameter to fetch deviceRegistrations: The page size must be between 1 and 2000. Was 0",
				Info:      "GetAllDeviceRegistrations",
			},
		}, {
			name:        "too large",
			reqPageSize: 2001,
			expectedErr: &generic.Error{
				ErrorType: "ClientError",
				Message:   "Error while building pageSize parameter to fetch deviceRegistrations: The page size must be between 1 and 2000. Was 2001",
				Info:      "GetAllDeviceRegistrations",
			},
		}, {
			name:        "Min",
			reqPageSize: 1,
			expectedErr: nil,
		}, {
			name:        "Max",
			reqPageSize: 2000,
			expectedErr: nil,
		}, {
			name:        "in range",
			reqPageSize: 1000,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given: A test server
			var capturedUrl string
			testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				capturedUrl = req.URL.String()
				_, _ = res.Write([]byte("{}"))
			}))
			defer testServer.Close()

			deviceRegistrationApi := buildDeviceRegistrationApi(testServer)
			_, err := deviceRegistrationApi.GetAll(tt.reqPageSize)

			if fmt.Sprint(err) != fmt.Sprint(tt.expectedErr) {
				t.Fatalf("respond with unexpected error. \nExpected: %s\nGiven:    %s", tt.expectedErr, err)
			}

			if tt.expectedErr == nil && !strings.Contains(capturedUrl, fmt.Sprintf("pageSize=%d", tt.reqPageSize)) {
				t.Errorf("GetAll() expected pageSize '%d' in url. '%s' given", tt.reqPageSize, capturedUrl)
			}
		})
	}
}
