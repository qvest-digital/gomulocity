package alarm

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)



func TestAlarmApi_Delete_All_Success(t *testing.T) {
	var expectedUrl = "alarm/alarms"
	var capturedUrl string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUrl = r.URL.String()
		w.WriteHeader(http.StatusNoContent)
	}))

	// given: A test server
	defer ts.Close()

	// and: the api as system under test
	api := buildAlarmApi(ts.URL)

	err := api.DeleteAll()

	if err != nil {
		t.Fatalf("DeleteAll() got an unexpected error: %s", err.Error())
	}

	if strings.HasSuffix(capturedUrl, expectedUrl) == false {
		t.Errorf("DeleteAll(): Wrong target URL: %s - expected %s", capturedUrl, expectedUrl)
	}
}

func TestAlarmApi_Delete_All_BadRequest(t *testing.T) {
	// given: A test server
	ts := buildHttpServer(http.StatusBadRequest, "")
	defer ts.Close()

	// and: the api as system under test
	api := buildAlarmApi(ts.URL)

	err := api.DeleteAll()

	if err == nil {
		t.Errorf("DeleteAll() expected error on 404 - not found")
		return
	}

	if !strings.Contains(err.ErrorType, "400") {
		t.Errorf("DeleteAll() expected error on 400 - bad request. Got: %s", err.ErrorType)
		return
	}
}
