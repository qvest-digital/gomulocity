package managedObjects

import (
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestClient_GetDeviceInformation(t *testing.T) {

	tests := []struct {
		name                     string
		expectedObjectCollection ManagedObjectCollection
		respCode                 int
		respBody                 string
		expectedRespBody         string
		expectedErr              error

		queryValues map[string][]string
	}{
		{
			name: "Happy",
			expectedObjectCollection: ManagedObjectCollection{
				Next: "<next>",
				Self: "<self>",
				ManagedObjects: []ManagedObjects{
					{
						C8YIsDevice: reflect.Interface,
						ID:          "<ID-1>",
						Name:        "<Name-1>",
					},
				},
			},
			queryValues: map[string][]string{
				"pageSize": {"10"},
			},
			respCode: http.StatusOK,
			respBody: NewManagedObjectCollection_ResponseBody(""),
		},
		{
			name:        "Unhappy - unauthorized",
			respCode:    http.StatusUnauthorized,
			respBody:    NewManagedObjectCollection_ResponseBody(""),
			expectedErr: generic.BadCredentialsErr,
		},
		{
			name:        "Unhappy - forbidden",
			respCode:    http.StatusForbidden,
			respBody:    NewManagedObjectCollection_ResponseBody(""),
			expectedErr: generic.AccessDeniedErr,
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

			c := ManagedObjectApi{
				Client: &generic.Client{
					HTTPClient: testserver.Client(),
					BaseURL:    testserver.URL,
					Username:   "<username>",
					Password:   "<password>",
				},
			}

			managedObjectCollection, err := c.ManagedObjectCollection(tt.queryValues)
			if err != tt.expectedErr {
				t.Errorf("received an unexpected error: expected: %v, actual: %v", tt.expectedErr, err)
			}

			if username != u {
				t.Fatalf("unexpected c8y auth-username: %v expected: %v", u, username)
			}
			if password != p {
				t.Fatalf("unexpected c8y auth-password: %v expected: %v", p, password)
			}

			expectedReqURL := fmt.Sprintf("%v", managedObjectPath)
			if reqURL != expectedReqURL {
				t.Fatalf("unexpected request url: %v expected: %v", reqURL, managedObjectPath)
			}

			if err == nil {
				if len(managedObjectCollection.ManagedObjects) != 3 {
					t.Errorf("unexpected count of managedObject. Expected: %v, actual: %v", 3, len(managedObjectCollection.ManagedObjects))
				}
			}
		})
	}
}

func Test_BuildURL_Happy(t *testing.T) {
	// given
	next := "https://t200588189.cumulocity.com/inventory/managedObjects?query=has(c8y_IsDevice)&pageSize=3&currentPage=1"

	// when
	result, err := buildURL(next)

	// then
	if err != nil {
		t.Error(err)
	}
	if result != "/inventory/managedObjects?query=has(c8y_IsDevice)&pageSize=3&currentPage=1" {
		t.Errorf("failed to build url")
	}
}
