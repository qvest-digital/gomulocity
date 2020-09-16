package inventory

import (
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"testing"
)

func TestInventoryApi_CommonPropertiesOnFindByQuery(t *testing.T) {
	var expectedContentType = "application/vnd.com.nsn.cumulocity.managedObjectCollection+json"
	var reqBasicAuthUsername, reqBasicAuthPassword, reqURL, reqAccept, reqContentType, reqBody string

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		reqURL = req.URL.String()
		reqBodyBytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			t.Fatalf("failed to read c8y request body: %s", err)
		}
		defer req.Body.Close()
		reqBody = string(reqBodyBytes)

		reqBasicAuthUsername, reqBasicAuthPassword, _ = req.BasicAuth()
		reqAccept = req.Header.Get("Accept")
		reqContentType = req.Header.Get("Content-Type")
		res.WriteHeader(http.StatusOK)
		_, err = res.Write([]byte(givenManagedObjectCollectionResponse))
		if err != nil {
			t.Fatalf("failed to write resp body: %s", err)
		}
	}))
	defer testServer.Close()

	inventoryApi := buildInventoryApi(testServer)
	managedObjects, err := inventoryApi.FindByQuery(query, 5)
	if err != nil {
		t.Errorf("received an unexpected error: %v", err)
	}

	if !reflect.DeepEqual(managedObjects, expectedManagedObjectCollection) {
		t.Errorf("received an unexpected ManagedObject: %#v. \nExpected: %#v", managedObjects, expectedManagedObjectCollection)
	}

	if len(reqBody) > 0 {
		t.Errorf("processed an unexpected c8y request body %q", reqBody)
	}

	if reqAccept != expectedContentType {
		t.Errorf("unexpected request accept header. Expected %q. Given: %q", expectedContentType, reqAccept)
	}
	if len(reqContentType) > 0 {
		t.Errorf("unexpected request content-type header:%q", reqContentType)
	}

	if reqBasicAuthUsername != USER {
		t.Errorf("unexpected c8y request basic-auth username. Expected %q. Given: %q", USER, reqBasicAuthUsername)
	}
	if reqBasicAuthPassword != PASSWORD {
		t.Errorf("unexpected c8y request basic-auth password. Expected %q. Given: %q", PASSWORD, reqBasicAuthPassword)
	}

	var expectedC8YRequestURL = "/inventory/managedObjects?pageSize=5&query=%24filter%3Dname+eq+%27%2ATest%2A%27+%24orderby%3Did+desc"
	if reqURL != expectedC8YRequestURL {
		t.Errorf("unexpected c8y request url. Expected %q. Given: %q", expectedC8YRequestURL, reqURL)
	}
}

func TestInventoryApi_FindManagedObjectByQuery(t *testing.T) {
	tests := []struct {
		name        string
		query       string
		pageSize    int
		c8yRespCode int
		c8yRespBody string

		expectedErr *generic.Error
	}{
		{
			name:        "query filter is empty",
			query:       "",
			pageSize:    5,
			c8yRespCode: http.StatusOK,
			c8yRespBody: givenManagedObjectCollectionResponse,
			expectedErr: nil,
		}, {
			name:        "invalid pageSize",
			query:       query,
			pageSize:    -1,
			c8yRespCode: http.StatusOK,
			c8yRespBody: givenManagedObjectCollectionResponse,
			expectedErr: &generic.Error{
				ErrorType: "ClientError",
				Message:   "Error while building pageSize parameter to fetch managedObjects: The page size must be between 1 and 2000. Was -1",
				Info:      "FindManagedObject",
			},
		}, {
			name:        "empty response body",
			query:       query,
			pageSize:    1,
			c8yRespCode: http.StatusOK,
			c8yRespBody: "",
			expectedErr: &generic.Error{
				ErrorType: "ClientError",
				Message:   "Response body was empty",
				Info:      "GetManagedObjectCollection",
			},
		}, {
			name:        "unparsable response body",
			query:       query,
			pageSize:    1,
			c8yRespCode: http.StatusOK,
			c8yRespBody: "{",
			expectedErr: &generic.Error{
				ErrorType: "ClientError",
				Message:   "Error while parsing response JSON: Error while unmarshaling json: unexpected end of JSON input",
				Info:      "GetManagedObjectCollection",
			},
		}, {
			name:        "bad credentials",
			query:       query,
			pageSize:    5,
			c8yRespCode: http.StatusUnauthorized,
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
			name:        "error without status code on GET",
			query:       query,
			pageSize:    5,
			c8yRespCode: 0,
			expectedErr: &generic.Error{
				ErrorType: "ClientError",
				Message:   "Error while getting managedObjects: Get.*",
				Info:      "GetManagedObjectCollection",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				res.WriteHeader(tt.c8yRespCode)
				_, err := res.Write([]byte(tt.c8yRespBody))
				if err != nil {
					t.Fatalf("failed to write resp body: %s", err)
				}
			}))
			defer testServer.Close()

			inventoryApi := buildInventoryApi(testServer)
			_, err := inventoryApi.FindByQuery(query, tt.pageSize)

			if matched, _ := regexp.MatchString(fmt.Sprint(tt.expectedErr), fmt.Sprint(err)); !matched {
				t.Fatalf("received an unexpected error: %s\nExpected: %s", err, tt.expectedErr)
			}
		})
	}
}
