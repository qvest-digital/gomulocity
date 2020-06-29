package gomulocity_event

import (
	"testing"
)

func TestEvents_Get_ExistingId(t *testing.T) {
	// given: A test server
	ts := buildHttpServer(200, event)
	defer ts.Close()

	// and: the api as system under test
	api := buildEventsApi(ts.URL)

	event, err := api.Get(deviceId)

	if err != nil {
		t.Fatalf("Get() got an unexpected error: %s", err.Error())
	}

	if event == nil {
		t.Fatalf("Get() returns an unexpected nil event.")
	}

	if event.Id != eventId {
		t.Errorf("Get() event id = %v, want %v", event.Id, eventId)
	}
}

func TestEvents_Get_NotExistingId(t *testing.T) {
	// given: A test server
	ts := buildHttpServer(404, "")
	defer ts.Close()

	// and: the api as system under test
	api := buildEventsApi(ts.URL)

	event, err := api.Get(eventId)

	if err != nil {
		t.Fatalf("Get() got an unexpected error: %s", err.Error())
		return
	}

	if event != nil {
		t.Fatalf("Get() got an unexpected event. Should be nil.")
		return
	}
}

func TestEvents_Get_MalformedJson(t *testing.T) {
	// given: A test server
	ts := buildHttpServer(200, "{ foo ...")
	defer ts.Close()

	// and: the api as system under test
	api := buildEventsApi(ts.URL)

	_, err := api.Get(eventId)

	if err == nil {
		t.Error("Get() Error expected but nil was given")
		return
	}
}
