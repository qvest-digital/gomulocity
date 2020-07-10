package events

import (
	"reflect"
	"testing"
	"time"
)

// given: A create event
var createEvent = &CreateEvent{
	Type:             "TestEvent",
	Time:             time.Time{},
	Text:             "This is my test event",
	Source:           Source{Id: "4711"},
	AdditionalFields: map[string]interface{}{},
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
		t.Errorf("CreateEvent() \nevent = %v \nwant %v", createEventCapture, createEvent)
	}

	header := requestCapture.Header.Get("Accept")
	want := "application/vnd.com.nsn.cumulocity.eventApi+json"
	if header != want {
		t.Errorf("CreateEvent() accent header = %v, want %v", header, want)
	}
}

func TestEvents_Create_Event_CustomFields(t *testing.T) {
	// given: A test server
	ts := createEventHttpServer(201)
	defer ts.Close()

	// and: the api as system under test
	api := buildEventsApi(ts.URL)
	createEvent := &CreateEvent{
		Type:   "TestEvent",
		Time:   time.Time{},
		Text:   "This is my test event",
		Source: Source{Id: "4711"},
		AdditionalFields: map[string]interface{}{
			"Custom1": 4711,
			"Custom2": "Hello World",
		},
	}
	_, err := api.CreateEvent(createEvent)

	if err != nil {
		t.Fatalf("CreateEvent() got an unexpected error: %s", err.Error())
	}

	if createEventCapture == nil {
		t.Fatalf("CreateEvent() Captured event is nil.")
	}

	if createEvent.Type != createEventCapture.Type ||
		createEvent.Time != createEventCapture.Time ||
		createEvent.Text != createEventCapture.Text ||
		createEvent.Source != createEventCapture.Source {

		t.Errorf("CreateEvent() basic fields - \nevent = %v \nwant %v", createEventCapture, createEvent)
	}

	custom1, _ := createEventCapture.AdditionalFields["Custom1"].(float64)
	custom2, _ := createEventCapture.AdditionalFields["Custom2"].(string)
	if custom1 != 4711 ||
		custom2 != "Hello World" {
		t.Errorf("CreateEvent() additional fields - \nevent = %v \nwant %v", createEventCapture, createEvent)
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
