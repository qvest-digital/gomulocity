package user_api

import (
	"github.com/tarent/gomulocity/generic"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestUsers_GetCurrentUser(t *testing.T) {
	var expectedUrl = "/user/currentUser"
	var capturedUrl string

	// given: A test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = r.URL.String()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(currentUserJSON))
	}))
	defer ts.Close()
	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// when
	currentUser, err := api.GetCurrentUser()

	// then
	if err != nil {
		t.Errorf("received an unexpected error: %s", err)
	}

	if !reflect.DeepEqual(testCurrentUser, currentUser) {
		t.Errorf("CreateUser() want: %v, actual: %v", testCurrentUser, currentUser)
	}

	if capturedUrl != expectedUrl {
		t.Errorf("unexpected request url: expected: %v, actual: %v", expectedUrl, capturedUrl)
	}
}

func TestUsers_GetCurrentUser_invalid_status(t *testing.T) {
	// given: A test server
	ts := buildHttpServer(http.StatusInternalServerError, currentUserErroneousJSON)
	defer ts.Close()

	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// when
	_, err := api.GetCurrentUser()

	// then
	expectedErr := currentUser_ErroneousResponse(http.StatusInternalServerError).Error()
	if err.Error() != expectedErr {
		t.Errorf("received an unexpected error: expected: %v, actual: %v", expectedErr, err)
	}
}

func TestUsers_GetCurrentUser_invalid_json(t *testing.T) {
	// given: A test server
	ts := buildHttpServer(http.StatusOK, `<>`)
	defer ts.Close()

	// and: the api as system under test
	api := buildUserApi(ts.URL)

	// when
	_, err := api.GetCurrentUser()

	// then
	expectedErr := generic.ClientError("Error while unmarshalling request body: invalid character '<' looking for beginning of value", "GetCurrentUser")
	if err.Error() != expectedErr.Error() {
		t.Errorf("received an unexpected error: expected: %v, actual: %v", expectedErr, err)
	}
}
