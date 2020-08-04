package user_api

import (
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestUserApi_GroupByName(t *testing.T) {
	var expectedUrl = fmt.Sprintf("/user/%v/groupByName/%v", tenantID, groupName)
	var capturedUrl string

	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = r.URL.String()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(groupJSON))
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// when
	group, err := api.GroupByName(tenantID, groupName)

	// then
	if err != nil {
		t.Fatalf("Received an unexpected error: %s", err)
	}

	if group != nil {
		if group.Name != groupName {
			t.Errorf("GroupName does not match! Expected ID: %v, actual: %v", groupName, group.Name)
		}
		if len(group.Roles) == 0 {
			t.Error("Group does not contain any roles. It should contain three roles!")
		}
	} else {
		t.Error("GroupByName() response group must not be nil")
	}

	if expectedUrl != capturedUrl {
		t.Errorf("unexpected request url: expected: %v, actual: %v", expectedUrl, capturedUrl)
	}

	if !reflect.DeepEqual(testGroup, group) {
		t.Errorf("GroupDetails() want: %v, actual: %v", testGroup, group)
	}
}

func TestUserApi_GroupByName_without_groupName(t *testing.T) {
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(erroneousResponseJSON))
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// when
	_, err := api.GroupByName(tenantID, "")

	// then
	expectedErr := generic.ClientError("Getting group without tenantID or group name is not allowed", "GroupByName")
	if err.Error() != expectedErr.Error() {
		t.Fatalf("Received an unexpected error: %s", err)
	}
}

func TestUserApi_GroupByName_error_unmarshalling(t *testing.T) {
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<>"))
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// when
	_, err := api.GroupByName(tenantID, groupName)

	// then
	expectedErr := generic.ClientError("Error while unmarshalling response body: invalid character '<' looking for beginning of value", "GroupByName")
	if err.Error() != expectedErr.Error() {
		t.Fatalf("Received an unexpected error: %s", err)
	}
}

func TestUserApi_GroupByName_invalid_status(t *testing.T) {
	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(erroneousResponseJSON))
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// when
	_, err := api.GroupByName(tenantID, groupName)

	// then
	expectedErr := createErroneousResponse(http.StatusInternalServerError)
	if err.Error() != expectedErr.Error() {
		t.Fatalf("Received an unexpected error: %s", err)
	}
}
