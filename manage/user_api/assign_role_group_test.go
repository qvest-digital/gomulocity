package user_api

import (
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"reflect"

	//"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserApi_AssignRoleToGroup(t *testing.T) {
	var expectedUrl = fmt.Sprintf("/user/%v/groups/%v/roles", tenantID, groupID)
	var capturedUrl string

	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = r.URL.String()
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(roleReferenceJSON))
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// and a reference
	reference := &RoleReference{
		Role: Role{
			Self: "https://t200588189.cumulocity.com/user/roles/ROLE_ACCOUNT_ADMIN",
		},
	}

	// when
	ref, err := api.AssignRoleToGroup(tenantID, groupID, reference)

	// then
	if err != nil {
		t.Fatalf("Received an unexpected error: %s", err)
	}

	if ref != nil {
		if ref.Role.Self != reference.Role.Self || ref.Role.ID != "ROLE_ACCOUNT_ADMIN" || ref.Role.Name != "ROLE_ACCOUNT_ADMIN" {
			t.Errorf("AssignRoleToUser() received an invalid reference response. Expected: %v, actual: %v", roleReference, ref)
		}
	} else {
		t.Error("AssignRoleToGroup() response reference must not be nil")
	}

	if expectedUrl != capturedUrl {
		t.Errorf("unexpected request url: expected: %v, actual: %v", expectedUrl, capturedUrl)
	}

	if !reflect.DeepEqual(roleReference, ref) {
		t.Errorf("AssignRoleToGroup() want: %v, actual: %v", roleReference, ref)
	}
}

func TestUserApi_AssignRoleToGroup_error_unmarshalling_response(t *testing.T) {
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("<>"))
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// and a reference
	reference := &RoleReference{
		Role: Role{
			Self: "https://t200588189.cumulocity.com/user/roles/ROLE_ACCOUNT_ADMIN",
		},
	}

	// when
	_, err := api.AssignRoleToGroup(tenantID, groupID, reference)

	// then
	expectedErr := generic.ClientError("Error while unmarshalling response body: invalid character '<' looking for beginning of value", "AssignRoleToGroup")
	if err.Error() != expectedErr.Error() {
		t.Fatalf("Received an unexpected error: %s", err)
	}
}

func TestUserApi_AssignRoleToGroup_tenantID_and_groupID_missing(t *testing.T) {
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("<>"))
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// and a reference
	reference := &RoleReference{
		Role: Role{
			Self: "https://t200588189.cumulocity.com/user/roles/ROLE_ACCOUNT_ADMIN",
		},
	}

	// when
	_, err := api.AssignRoleToGroup("", "", reference)

	// then
	expectedErr := generic.ClientError("Assigning role to group without tenantID or groupID is not allowed", "AssignRoleToGroup")
	if err.Error() != expectedErr.Error() {
		t.Fatalf("Received an unexpected error: %s", err)
	}
}

func TestUserApi_AssignRoleToGroup_invalid_status(t *testing.T) {
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(erroneousResponseJSON))
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// and a reference
	reference := &RoleReference{
		Role: Role{
			Self: "https://t200588189.cumulocity.com/user/roles/ROLE_ACCOUNT_ADMIN",
		},
	}

	// when
	_, err := api.AssignRoleToGroup(tenantID, groupID, reference)

	// then
	expectedErr := createErroneousResponse(http.StatusInternalServerError)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("Received an unexpected error: %s", err)
	}
}
