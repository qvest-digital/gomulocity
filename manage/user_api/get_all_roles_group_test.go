package user_api

import (
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserApi_GetAllRolesOfAGroup(t *testing.T) {
	var expectedUrl = fmt.Sprintf("/user/%v/groups/%v/roles?pageSize=%v", tenantID, groupID, 5)
	var capturedUrl string

	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = r.URL.String()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(roleReferenceCollectionJSON))
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// when
	ref, err := api.GetAllRolesOfAGroup(tenantID, groupID, 5)

	// then
	if err != nil {
		t.Fatalf("Received an unexpected error: %s", err)
	}

	if ref != nil {
		if len(ref.References) == 0 {
			t.Error("role reference collection does not contain any references. It should contain two references!")
		}
	} else {
		t.Error("GetAllRolesOfAGroup() role reference must not be nil")
	}

	if expectedUrl != capturedUrl {
		t.Errorf("unexpected request url: expected: %v, actual: %v", expectedUrl, capturedUrl)
	}
}

func TestUserApi_GetAllRolesOfAGroup_without_groupID(t *testing.T) {
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(erroneousResponseJSON))
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// when
	_, err := api.GetAllRolesOfAGroup(tenantID, "", 5)

	// then
	expectedErr := generic.ClientError("Getting role reference collection without username or groupID is not allowed", "FindRoleReferenceCollection")
	if err.Error() != expectedErr.Error() {
		t.Fatalf("Received an unexpected error: %s", err)
	}
}

func TestUserApi_GetAllRolesOfAGroup_invalid_status(t *testing.T) {
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(erroneousResponseJSON))
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// when
	_, err := api.GetAllRolesOfAGroup(tenantID, groupID, 5)

	// then
	expectedErr := createErroneousResponse(http.StatusInternalServerError)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("Received an unexpected error: %s", err)
	}
}

