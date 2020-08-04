package user_api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserApi_FindRoleCollection(t *testing.T) {
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
	collection, err := api.RoleCollection(5)

}
