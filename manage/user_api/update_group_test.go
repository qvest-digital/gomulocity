package user_api

import (
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserApi_UpdateGroup(t *testing.T) {
	var expectedUrl = fmt.Sprintf("/user/%v/groups/%v", tenantID, groupID)
	var capturedUrl string
	var newGroupName = "newGroup"

	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = r.URL.String()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(groupJSON_WithGroupName(newGroupName)))
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// and a group
	group := &Group{
		Name: newGroupName,
	}

	// when
	group, err := api.UpdateGroup(tenantID, groupID, group)

	// then
	if err != nil {
		t.Fatalf("Received an unexpected error: %s", err)
	}

	if group != nil {
		if group.Name != newGroupName {
			t.Errorf("GroupName does not match! Expected ID: %v, actual: %v", newGroupName, group.Name)
		}
		if len(group.Roles) == 0 {
			t.Error("Group does not contain any roles. It should contain three roles!")
		}
	} else {
		t.Error("UpdateGroup() response group must not be nil")
	}

	if expectedUrl != capturedUrl {
		t.Errorf("unexpected request url: expected: %v, actual: %v", expectedUrl, capturedUrl)
	}
}

func TestUserApi_UpdateGroup_without_groupID(t *testing.T) {
	var newGroupName = "newGroup"

	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(erroneousResponseJSON))
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// and a group
	group := &Group{
		Name: newGroupName,
	}

	// when
	_, err := api.UpdateGroup(tenantID, "", group)

	// then
	expectedErr := generic.ClientError("Updating a group without tenantID and groupID is not allowed", "UpdateGroup")
	if err.Error() != expectedErr.Error() {
		t.Fatalf("Received an unexpected error: %s", err)
	}
}

func TestUserApi_UpdateGroup_error_unmarshalling(t *testing.T) {
	var newGroupName = "newGroup"

	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<>"))
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// and a group
	group := &Group{
		Name: newGroupName,
	}

	// when
	_, err := api.UpdateGroup(tenantID, groupID, group)

	// then
	expectedErr := generic.ClientError("Error while unmarshalling response body: invalid character '<' looking for beginning of value", "UpdateGroup")
	if err.Error() != expectedErr.Error() {
		t.Fatalf("Received an unexpected error: %s", err)
	}
}

func TestUserApi_UpdateGroup_invalid_status(t *testing.T) {
	var newGroupName = "newGroup"

	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(erroneousResponseJSON))
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// and a group
	group := &Group{
		Name: newGroupName,
	}

	// when
	_, err := api.UpdateGroup(tenantID, groupID, group)

	// then
	expectedErr := createErroneousResponse(http.StatusInternalServerError)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("Received an unexpected error: %s", err)
	}
}
