package gomulocity_event

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEvents_Delete_Event_Success(t *testing.T) {
	var capturedUrl string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		capturedUrl = r.URL.String()
		w.WriteHeader(http.StatusNoContent)
	}))

	// given: A test server
	defer ts.Close()

	// and: the api as system under test
	api := buildEventsApi(ts.URL)

	err := api.DeleteEvent(eventId)

	if err != nil {
		t.Fatalf("DeleteEvent() got an unexpected error: %s", err.Error())
	}

	if strings.Contains(capturedUrl, eventId) == false {
		t.Errorf("DeleteEvent() The target URL does not contains the event Id: url: %s - expected eventId %s", capturedUrl, eventId)
	}
}

func TestEvents_Delete_Event_NotFound(t *testing.T) {
	// given: A test server
	ts := buildHttpServer(http.StatusNotFound, "")
	defer ts.Close()

	// and: the api as system under test
	api := buildEventsApi(ts.URL)

	err := api.DeleteEvent(eventId)

	if err == nil {
		t.Errorf("DeleteEvent() expected error on 404 - not found")
		return
	}
}
