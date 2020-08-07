package user_api

import (
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserApi_UnassignRoleFromGroup(t *testing.T) {
	var expectedUrl = fmt.Sprintf("/user/%v/groups/%v/roles/%v", tenantID, groupID, roleName)
	var capturedUrl string

	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = r.URL.String()
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte(roleReferenceJSON))
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// when
	err := api.UnassignRoleFromGroup(tenantID, groupID, roleName)

	// then
	if err != nil {
		t.Fatalf("Received an unexpected error: %s", err)
	}

	if expectedUrl != capturedUrl {
		t.Errorf("unexpected request url: expected: %v, actual: %v", expectedUrl, capturedUrl)
	}
}

func TestUserApi_UnassignRoleFromGroup_roleName_missing(t *testing.T) {
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(erroneousResponseJSON))
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// when
	err := api.UnassignRoleFromGroup(tenantID, groupID, "")

	// then
	expectedErr := generic.ClientError("Unassign role from group without tenantID, groupID or roleName is not allowed","UnassignRoleFromGroup")
	if err.Error() != expectedErr.Error() {
		t.Fatalf("Received an unexpected error: %s", err)
	}
}

func TestUserApi_UnassignRoleFromGroup_invalid_status(t *testing.T) {
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(erroneousResponseJSON))
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// when
	err := api.UnassignRoleFromGroup(tenantID, username, roleName)

	// then
	expectedErr := createErroneousResponse(http.StatusInternalServerError)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("Received an unexpected error: %s", err)
	}
}

