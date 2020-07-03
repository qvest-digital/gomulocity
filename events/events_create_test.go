package events

import (
	"reflect"
	"testing"
	"time"
)

var eventTime, _ = time.Parse(time.RFC3339, "2020-06-26T10:43:25.130Z")
var responseEvent = &Event{
	Id:           "",
	Type:         "TestEvent",
	Time:         eventTime,
	CreationTime: eventTime,
	Text:         "This is my test event",
	Source:       Source{Id: "4711"},
	Self:         "https://t0815.cumulocity.com/event/events/1337",
}

// given: A create event
var createEvent = &CreateEvent{
	Type:   "TestEvent",
	Time:   time.Time{},
	Text:   "This is my test event",
	Source: Source{Id: "4711"},
}

func TestEvents_Create_Event_Success_SendsData(t *testing.T) {
	// given: A test server
	ts := createEventHttpServer(201)
	defer ts.Close()

	// and: the api as system under test
	api := buildEventsApi(ts.URL)
	_, err := api.CreateEvent(createEvent)

	if err != nil {
		t.Fatalf("CreateEvent() got an unexpected error: %s", err.Error())
	}

	if createEventCapture == nil {
		t.Fatalf("CreateEvent() Captured event is nil.")
	}

	if !reflect.DeepEqual(createEvent, createEventCapture) {
		t.Errorf("CreateEvent() event = %v, want %v", createEvent, createEventCapture)
	}

	header := requestCapture.Header.Get("Accept")
	want := "application/vnd.com.nsn.cumulocity.eventApi+json"
	if header != want {
		t.Errorf("CreateEvent() accent header = %v, want %v", header, want)
	}
}

func TestEvents_Create_Event_Success_ReceivesData(t *testing.T) {
	// given: A test server
	ts := createEventHttpServer(201)
	defer ts.Close()

	// and: the api as system under test
	api := buildEventsApi(ts.URL)
	event, err := api.CreateEvent(createEvent)

	if err != nil {
		t.Fatalf("CreateEvent() got an unexpected error: %s", err.Error())
	}

	if !reflect.DeepEqual(event, responseEvent) {
		t.Errorf("CreateEvent() event = %v, want %v", createEvent, createEventCapture)
	}
}

func TestEvents_Create_Event_BadRequest(t *testing.T) {
	// given: A test server
	ts := createEventHttpServer(400)
	defer ts.Close()

	// and: the api as system under test
	api := buildEventsApi(ts.URL)
	_, err := api.CreateEvent(createEvent)

	if err == nil {
		t.Errorf("CreateEvent() expected error on 400 - bad request")
		return
	}
}
