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

func TestInventoryReferenceApi_CommonPropertiesOnGetMany(t *testing.T) {
	var expectedContentType = "application/vnd.com.nsn.cumulocity.managedObjectReferenceCollection+json"
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
		_, err = res.Write([]byte(givenManagedObjectReferenceCollectionResponse))
		if err != nil {
			t.Fatalf("failed to write resp body: %s", err)
		}
	}))
	defer testServer.Close()

	inventoryReferenceApi := buildInventoryReferenceApi(testServer)
	managedObjectReferenceCollection, err := inventoryReferenceApi.GetMany(managedObjectId, CHILD_DEVICES, 5)
	if err != nil {
		t.Errorf("received an unexpected error: %v", err)
	}

	if !reflect.DeepEqual(managedObjectReferenceCollection, expectedManagedObjectReferenceCollection) {
		t.Errorf("received an unexpected ManagedObjectReference: %#v. \nExpected: %#v", managedObjectReferenceCollection, expectedManagedObjectReferenceCollection)
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

	var expectedC8YRequestURL = fmt.Sprintf("/inventory/managedObjects/%s/childDevices?pageSize=5", managedObjectId)
	if reqURL != expectedC8YRequestURL {
		t.Errorf("unexpected c8y request url. Expected %q. Given: %q", expectedC8YRequestURL, reqURL)
	}
}

func TestInventoryReferenceApi_GetManyManagedObjectReferences(t *testing.T) {
	tests := []struct {
		name            string
		managedObjectId string
		pageSize        int
		c8yRespCode     int
		c8yRespBody     string

		expectedErr *generic.Error
	}{
		{
			name:            "managedObjectId is empty",
			managedObjectId: "",
			pageSize:        5,
			c8yRespCode:     http.StatusOK,
			c8yRespBody:     givenManagedObjectReferenceCollectionResponse,
			expectedErr: &generic.Error{
				ErrorType: "ClientError",
				Message:   "managedObjectId must not be empty",
				Info:      "GetManyManagedObjectReferences",
			},
		}, {
			name:            "invalid pageSize",
			managedObjectId: "9963944",
			pageSize:        -1,
			c8yRespCode:     http.StatusOK,
			c8yRespBody:     givenManagedObjectReferenceCollectionResponse,
			expectedErr: &generic.Error{
				ErrorType: "ClientError",
				Message:   "Error while building pageSize parameter to fetch managedObjectReferences: The page size must be between 1 and 2000. Was -1",
				Info:      "GetManyManagedObjectReferences",
			},
		}, {
			name:        "empty response body",
			managedObjectId: "9963944",
			pageSize:        5,
			c8yRespCode: http.StatusOK,
			c8yRespBody: "",
			expectedErr: &generic.Error{
				ErrorType: "ClientError",
				Message:   "Response body was empty",
				Info:      "GetManagedObjectReferenceCollection",
			},
		}, {
			name:        "unparsable response body",
			managedObjectId: "9963944",
			pageSize:        5,
			c8yRespCode: http.StatusOK,
			c8yRespBody: "{",
			expectedErr: &generic.Error{
				ErrorType: "ClientError",
				Message:   "Error while parsing response JSON: unexpected end of JSON input",
				Info:      "GetManagedObjectReferenceCollection",
			},
		}, {
			name:            "bad credentials",
			managedObjectId: "9963944",
			pageSize:        5,
			c8yRespCode:     http.StatusUnauthorized,
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
			name:            "error without status code on GET",
			managedObjectId: "9963944",
			pageSize:        5,
			c8yRespCode:     0,
			expectedErr: &generic.Error{
				ErrorType: "ClientError",
				Message:   "Error while getting managedObjectReferences: Get.*",
				Info:      "GetManagedObjectReferenceCollection",
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

			inventoryReferenceApi := buildInventoryReferenceApi(testServer)
			_, err := inventoryReferenceApi.GetMany(tt.managedObjectId, CHILD_DEVICES, tt.pageSize)

			if matched, _ := regexp.MatchString(fmt.Sprint(tt.expectedErr), fmt.Sprint(err)); !matched {
				t.Fatalf("received an unexpected error: %s\nExpected: %s", err, tt.expectedErr)
			}
		})
	}
}
