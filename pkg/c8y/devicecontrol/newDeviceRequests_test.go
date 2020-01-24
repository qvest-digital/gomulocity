package devicecontrol

import (
	"errors"
	"fmt"
	jsoncompare "github.com/orasik/gocomparejson"
	"github.com/tarent/gomulocity/pkg/c8y/meta"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestClient_CreateNewDeviceRequest(t *testing.T) {
	tests := []struct {
		name                     string
		deviceID                 string
		c8yRespCode              int
		c8yRespBody              string
		c8yExpectedRequestBody   string
		expectedNewDeviceRequest NewDeviceRequest
		expectedErr              error
	}{
		{
			name:                   "happy case",
			deviceID:               "4711",
			c8yRespCode:            http.StatusCreated,
			c8yRespBody:            `{"id": "4711", "status": "PENDING_ACCEPTANCE", "self": "https://myFancyCloudInstance.com/devicecontrol/deviceCredentials"}`,
			c8yExpectedRequestBody: `{"id": "4711"}`,
			expectedNewDeviceRequest: NewDeviceRequest{
				ID:     "4711",
				Status: "PENDING_ACCEPTANCE",
				Self:   "https://myFancyCloudInstance.com/devicecontrol/deviceCredentials",
			},
			expectedErr: nil,
		}, {
			name:                   "bad credentials",
			deviceID:               "nope 401",
			c8yRespCode:            http.StatusUnauthorized,
			c8yRespBody:            `{}`,
			expectedErr:            meta.BadCredentialsErr,
			c8yExpectedRequestBody: `{"id": "nope 401"}`,
		}, {
			name:                   "access denied",
			deviceID:               "nope 403",
			c8yRespCode:            http.StatusForbidden,
			c8yRespBody:            `{}`,
			expectedErr:            meta.AccessDeniedErr,
			c8yExpectedRequestBody: `{"id": "nope 403"}`,
		}, {
			name:                   "device already exists",
			deviceID:               "nope 422",
			c8yRespCode:            http.StatusUnprocessableEntity,
			c8yRespBody:            `{"error": "devicecontrol/Non Unique Result", "message": "That thing already exists", "info": "https://cumulocity.com/guides/reference/rest-implementation/#a-name-error-reporting-a-error-reporting"}`,
			expectedErr:            NewDeviceRequestAlreadyExistsErr,
			c8yExpectedRequestBody: `{"id": "nope 422"}`,
		}, {
			name:                   "unexpected error",
			deviceID:               "nope 500",
			c8yRespCode:            http.StatusInternalServerError,
			c8yRespBody:            `{"error": "myCustomError", "message": "something goes wrong.", "info": "my link"}`,
			expectedErr:            errors.New("failed to create new-device-request. Status: 500: \"myCustomError\" something goes wrong. See: my link"),
			c8yExpectedRequestBody: `{"id": "nope 500"}`,
		}, {
			name:                   "invalid json error response",
			deviceID:               "nope 500 1",
			c8yRespCode:            http.StatusInternalServerError,
			c8yRespBody:            `#`,
			expectedErr:            errors.New("failed to decode error response body: invalid character '#' looking for beginning of value"),
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

			newDeviceRequest, err := c.CreateNewDeviceRequest(tt.deviceID)
			if fmt.Sprint(err) != fmt.Sprint(tt.expectedErr) {
				t.Fatalf("respond with unexpected error. Expected: %s. Given: %s", tt.expectedErr, err)
			}

			if !reflect.DeepEqual(newDeviceRequest, tt.expectedNewDeviceRequest) {
				t.Errorf("respond with unexpected newDeviceRequest. \nExpected: %#v. \nGiven: %#v", tt.expectedNewDeviceRequest, newDeviceRequest)
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

			var expectedC8YRequestURL = "/devicecontrol/newDeviceRequests"
			if reqURL != expectedC8YRequestURL {
				t.Errorf("unexpected c8y request url. Expected %q. Given: %q", expectedC8YRequestURL, reqURL)
			}
		})
	}
}

func TestClient_NewDeviceRequests(t *testing.T) {
	tests := []struct {
		name                               string
		reqOpts                            func(*http.Request)
		c8yRespCode                        int
		c8yRespBody                        string
		expectedRequestURL                 string
		expectedNewDeviceRequestCollection NewDeviceRequestCollection
		expectedErr                        error
	}{
		{
			name:               "happy case",
			reqOpts:            meta.Page(4),
			c8yRespCode:        http.StatusOK,
			c8yRespBody:        `{"self": "selfURL", "newDeviceRequests":[{"id": "rID1", "status": "ACCEPTED", "self":"rSelf1"}], "statistics": {"totalRecords":0, "pageSize":5, "currentPage":1}, "prev":"prevURL", "next":"nextURL"}`,
			expectedRequestURL: "/devicecontrol/newDeviceRequests?currentPage=4",
			expectedNewDeviceRequestCollection: NewDeviceRequestCollection{
				Self: "selfURL",
				NewDeviceRequests: []NewDeviceRequest{
					{ID: "rID1", Status: "ACCEPTED", Self: "rSelf1"},
				},
				Statistics: meta.PagingStatistics{
					TotalRecords: 0,
					PageSize:     5,
					CurrentPage:  1,
				},
				Prev: "prevURL",
				Next: "nextURL",
			},
			expectedErr: nil,
		}, {
			name:               "bad credentials",
			reqOpts:            meta.Page(1),
			c8yRespCode:        http.StatusUnauthorized,
			c8yRespBody:        `{}`,
			expectedRequestURL: "/devicecontrol/newDeviceRequests?currentPage=1",
			expectedErr:        meta.BadCredentialsErr,
		}, {
			name:               "access denied",
			c8yRespCode:        http.StatusForbidden,
			c8yRespBody:        `{}`,
			expectedRequestURL: "/devicecontrol/newDeviceRequests",
			expectedErr:        meta.AccessDeniedErr,
		}, {
			name:               "unexpected error",
			reqOpts:            meta.Page(9999999),
			c8yRespCode:        http.StatusInternalServerError,
			c8yRespBody:        `{"error": "myCustomError", "message": "something goes wrong.", "info": "my link"}`,
			expectedRequestURL: "/devicecontrol/newDeviceRequests?currentPage=9999999",
			expectedErr:        errors.New("failed to find-all new-device-requests. Status: 500: \"myCustomError\" something goes wrong. See: my link"),
		}, {
			name:               "invalid json error response",
			c8yRespCode:        http.StatusOK,
			c8yRespBody:        `%`,
			expectedRequestURL: "/devicecontrol/newDeviceRequests",
			expectedErr:        errors.New("failed to decode response body: invalid character '%' looking for beginning of value"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqBasicAuthUsername, reqBasicAuthPassword, reqURL string

			testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				reqURL = req.URL.String()
				reqBasicAuthUsername, reqBasicAuthPassword, _ = req.BasicAuth()
				res.WriteHeader(tt.c8yRespCode)
				_, err := res.Write([]byte(tt.c8yRespBody))
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

			newDeviceRequest, err := c.NewDeviceRequests(tt.reqOpts)
			if fmt.Sprint(err) != fmt.Sprint(tt.expectedErr) {
				t.Fatalf("respond with unexpected error. Expected: %s. Given: %s", tt.expectedErr, err)
			}

			if !reflect.DeepEqual(newDeviceRequest, tt.expectedNewDeviceRequestCollection) {
				t.Errorf("respond with unexpected newDeviceRequestCollection. \nExpected: %#v. \nGiven: %#v", tt.expectedNewDeviceRequestCollection, newDeviceRequest)
			}

			if reqBasicAuthUsername != u {
				t.Errorf("unexpected c8y request basic-auth username. Expected %q. Given: %q", u, reqBasicAuthUsername)
			}
			if reqBasicAuthPassword != p {
				t.Errorf("unexpected c8y request basic-auth password. Expected %q. Given: %q", p, reqBasicAuthPassword)
			}
			if reqURL != tt.expectedRequestURL {
				t.Errorf("unexpected c8y request URL. Expected %q. Given: %q", tt.expectedRequestURL, reqURL)
			}
		})
	}
}
