package devicecontrol

import (
	"errors"
	"fmt"
	jsoncompare "github.com/orasik/gocomparejson"
	"github.com/tarent/gomulocity/generic"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestClient_CreateDeviceCredentials(t *testing.T) {
	tests := []struct {
		name                      string
		deviceID                  string
		c8yRespCode               int
		c8yRespContentType        string
		c8yRespBody               string
		c8yExpectedRequestBody    string
		expectedDeviceCredentials DeviceCredentials
		expectedErr               error
	}{
		{
			name:                   "happy case",
			deviceID:               "4711",
			c8yRespCode:            http.StatusCreated,
			c8yRespBody:            `{"id": "4711", "status": "PENDING_ACCEPTANCE", "self": "https://myFancyCloudInstance.com/devicecontrol/deviceCredentials"}`,
			c8yExpectedRequestBody: `{"id": "4711"}`,
			expectedDeviceCredentials: DeviceCredentials{
				ID:   "4711",
				Self: "https://myFancyCloudInstance.com/devicecontrol/deviceCredentials",
			},
			expectedErr: nil,
		}, {
			name:                   "bad credentials",
			deviceID:               "nope 401",
			c8yRespCode:            http.StatusUnauthorized,
			c8yRespBody:            `{}`,
			expectedErr:            generic.BadCredentialsErr,
			c8yExpectedRequestBody: `{"id": "nope 401"}`,
		}, {
			name:                   "access denied",
			deviceID:               "nope 403",
			c8yRespCode:            http.StatusForbidden,
			c8yRespBody:            `{}`,
			expectedErr:            generic.AccessDeniedErr,
			c8yExpectedRequestBody: `{"id": "nope 403"}`,
		}, {
			name:                   "unexpected error",
			deviceID:               "nope 500",
			c8yRespCode:            http.StatusInternalServerError,
			c8yRespBody:            `{"error": "myCustomError", "message": "something goes wrong.", "info": "my link"}`,
			c8yRespContentType:     "application/vnd.com.nsn.cumulocity.error+json;q=0.7,en;q=0.3",
			expectedErr:            errors.New("failed to create device-credentials (500): request failed: \"myCustomError\" something goes wrong. See: my link"),
			c8yExpectedRequestBody: `{"id": "nope 500"}`,
		}, {
			name:                   "invalid json error response",
			deviceID:               "nope 500 1",
			c8yRespCode:            http.StatusInternalServerError,
			c8yRespBody:            `#`,
			expectedErr:            errors.New("failed to create device-credentials with status code 500"),
			c8yExpectedRequestBody: `{"id": "nope 500 1"}`,
		}, {
			name:                   "invalid json response",
			deviceID:               "nope 201",
			c8yRespCode:            http.StatusCreated,
			c8yRespBody:            `#`,
			expectedErr:            errors.New("failed to decode response body: invalid character '#' looking for beginning of value"),
			c8yExpectedRequestBody: `{"id": "nope 201"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqBasicAuthUsername, reqBasicAuthPassword, reqBody, reqURL string

			testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				reqURL = req.URL.String()
				res.Header().Set("Content-Type", tt.c8yRespContentType)
				reqBodyBytes, err := ioutil.ReadAll(req.Body)
				if err != nil {
					t.Fatalf("failed to read c8y request body: %s", err)
				}
				defer req.Body.Close()
				reqBody = string(reqBodyBytes)

				reqBasicAuthUsername, reqBasicAuthPassword, _ = req.BasicAuth()
				res.WriteHeader(tt.c8yRespCode)
				_, err = res.Write([]byte(tt.c8yRespBody))
				if err != nil {
					t.Fatalf("failed to write resp body: %s", err)
				}
			}))
			defer testServer.Close()

			u := "<username>"
			p := "<password>"
			c := Client{
				HTTPClient: testServer.Client(),
				BaseURL:    testServer.URL,
				Username:   u,
				Password:   p,
			}

			newDeviceRequest, err := c.CreateDeviceCredentials(tt.deviceID)
			if fmt.Sprint(err) != fmt.Sprint(tt.expectedErr) {
				t.Fatalf("respond with unexpected error. \nExpected: %s\nGiven:    %s", tt.expectedErr, err)
			}

			if !reflect.DeepEqual(newDeviceRequest, tt.expectedDeviceCredentials) {
				t.Errorf("respond with unexpected newDeviceRequest. \nExpected: %#v. \nGiven: %#v", tt.expectedDeviceCredentials, newDeviceRequest)
			}

			if reqBasicAuthUsername != u {
				t.Errorf("unexpected c8y request basic-auth username. Expected %q. Given: %q", u, reqBasicAuthUsername)
			}
			if reqBasicAuthPassword != p {
				t.Errorf("unexpected c8y request basic-auth password. Expected %q. Given: %q", p, reqBasicAuthPassword)
			}

			if equal, _ := jsoncompare.CompareJSON(reqBody, tt.c8yExpectedRequestBody); !equal {
				t.Errorf("unexpected c8y request body. Expected %q. Given: %q", tt.c8yExpectedRequestBody, reqBody)
			}

			var expectedC8YRequestURL = "/devicecontrol/deviceCredentials"
			if reqURL != expectedC8YRequestURL {
				t.Errorf("unexpected c8y request url. Expected %q. Given: %q", expectedC8YRequestURL, reqURL)
			}
		})
	}
}
