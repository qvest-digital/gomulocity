package user_api

import (
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestUsers_CreateUser_Happy(t *testing.T) {
	var expectedUrl = fmt.Sprintf("/user/%v/userApi", tenantID)
	var capturedUrl string

	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = r.URL.String()
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(testUserJSON))
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// when
	user, err := api.CreateUser(tenantID, createUser)

	// then
	if err != nil {
		t.Errorf("received an unexpected error: %s", err)
	}

	if !reflect.DeepEqual(testUser, user) {
		t.Errorf("CreateUser() want: %v, actual: %v", testUser, user)
	}

	if capturedUrl != expectedUrl {
		t.Errorf("unexpected request url: expected: %v, actual: %v", expectedUrl, capturedUrl)
	}

}

func TestUsers_CreateUser_invalid_status(t *testing.T) {
	// given: A test server
	ts := buildHttpServer(http.StatusInternalServerError, createErroneousResponseJSON)
	defer ts.Close()

	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// when
	_, err := api.CreateUser(tenantID, createUser)

	// then
	expectedErr := createErroneousResponse(http.StatusInternalServerError).Error()
	if err.Error() != expectedErr {
		t.Errorf("received an unexpected error: expected: %v, actual: %v", expectedErr, err)
	}
}

func TestUsers_CreateUser_empty_tenantID(t *testing.T) {
	// given: A test server
	ts := buildHttpServer(http.StatusInternalServerError, createErroneousResponseJSON)
	defer ts.Close()

	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// when
	_, err := api.CreateUser("", createUser)

	// then
	expectedErr := generic.ClientError("Creating user without a tenantID is not allowed", "CreateUser").Error()
	if err.Error() != expectedErr {
		t.Errorf("received an unexpected error: expected: %v, actual: %v", expectedErr, err)
	}
}


