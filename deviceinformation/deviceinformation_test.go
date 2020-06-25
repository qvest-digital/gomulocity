package deviceinformation

import (
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_GetDeviceInformation(t *testing.T) {
	tests := []struct {
		name                      string
		expectedDeviceCredentials DeviceCredentials
		respCode                  int
		respBody                  string
		expectedRespBody          string
	}{
		{
			name: "Happy",
			expectedDeviceCredentials: DeviceCredentials{
				ManagedObjects: []ManagedObject{
					{
						ID:       "1",
						Name:     "<devicename>",
						IsDevice: nil,
					},
				},
			},
			respCode: http.StatusOK,
			respBody: ResponseBodyDeviceInformation,
		},
		{
			name:              "unauthorized",
			expectedDeviceCredentials: DeviceCredentials{},
			respCode:          http.StatusUnauthorized,
			respBody:          `{"error":"security/Unauthorized","message":"Full authentication is required to access this resource","info":"https://www.cumulocity.com/guides/reference-guide/#error_reporting"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var username, password, reqURL string

			testserver := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
				reqURL = request.URL.String()

				username, password, _ = request.BasicAuth()
				response.WriteHeader(tt.respCode)
				_, err := response.Write([]byte(tt.respBody))
				if err != nil {
					t.Fatal("failed to write response body")
				}
			}))
			defer testserver.Close()

			u := "<username>"
			p := "<password>"

			c := Client{
				testserver.Client(),
				testserver.URL,
				"<username>",
				"<password>",
			}

			deviceCreds, err := c.GetDeviceInformation()
			if err != nil {
				if err != generic.BadCredentialsErr {
					t.Fatal(err)
				}
			}

			if username != u {
				t.Fatalf("unexpected c8y auth-username: %v expected: %v", u, username)
			}
			if password != p {
				t.Fatalf("unexpected c8y auth-password: %v expected: %v", p, password)
			}

			expectedReqURL := fmt.Sprintf("%v%v", deviceCredsPath, deviceCredsQuery)
			if reqURL != expectedReqURL {
				t.Fatalf("unexpected request url: %v expected: %v", reqURL, deviceCredsPath)
			}

			if len(deviceCreds.ManagedObjects) == 0 {
				t.Log("no deviceCreds found")
			}
		})
	}
}
