package events

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

func TestEvents_Get_CustomElements(t *testing.T) {
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

	if len(event.AdditionalFields) != 2 {
		t.Fatalf("Get() AdditionalFields length = %d, want %d", len(event.AdditionalFields), 2)
	}

	custom1, ok1 := event.AdditionalFields["custom1"].(string)
	custom2, ok2 := event.AdditionalFields["custom2"].([]interface{})

	if !(ok1 && custom1 == "Hello") {
		t.Errorf("Get() custom1 = %v, want %v", custom1, "Hello")
	}
	if !(ok2 && custom2[0] == "Foo" && custom2[1] == "Bar") {
		t.Errorf("Get() custom2 = [%v, %v], want [%v, %v]", custom2[0], custom2[1], "Foo", "Bar")
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
