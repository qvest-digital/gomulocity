package user_api

import (
	"github.com/tarent/gomulocity/generic"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestUserApi_AssignNewInventoryRole(t *testing.T) {
	var expectedUrl = "/user/inventoryroles"
	var capturedUrl string

	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = r.URL.String()
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(testInventoryRoleJSON("Reader")))
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// when
	role, err := api.AssignNewInventoryRole(&testInventoryRoles[0])

	// then
	if err != nil {
		t.Fatalf("Received an unexpected error: %s", err)
	}

	if !reflect.DeepEqual(testInventoryRole("Reader"), role) {
		t.Errorf("AssignNewInventoryRole() want: %v, actual: %v", testInventoryRole("Reader"), role)
	}

	if expectedUrl != capturedUrl {
		t.Errorf("unexpected request url: expected: %v, actual: %v", expectedUrl, capturedUrl)
	}
}

func TestUserApi_AssignNewInventoryRole_error_unmarshalling(t *testing.T) {
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("<>"))
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// when
	_, err := api.AssignNewInventoryRole(&testInventoryRoles[0])

	// then
	expectedErr := generic.ClientError("Error while unmarshalling response body: invalid character '<' looking for beginning of value", "AssignNewInventoryRole")
	if err.Error() != expectedErr.Error() {
		t.Fatalf("Received an unexpected error: %s", err)
	}
}

func TestUserApi_AssignNewInventoryRole_invalid_status(t *testing.T) {
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(erroneousResponseJSON))
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// when
	_, err := api.AssignNewInventoryRole(&testInventoryRoles[0])

	// then
	expectedErr := createErroneousResponse(http.StatusInternalServerError)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("Received an unexpected error: %s", err)
	}
}
