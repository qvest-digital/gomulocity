package gomulocity_event

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func createEventHttpServer(status int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)

		var event CreateEvent
		_ = json.Unmarshal(body, &event)
		eventCapture = &event

		w.WriteHeader(status)
	}))
}

var eventCapture *CreateEvent

// given: A create event
var createEvent = &CreateEvent{
	Type:   "TestEvent",
	Time:   time.Time{},
	Text:   "This is my test event",
	Source: Source{Id: "4711"},
}

func TestEvents_Create_Event_Success(t *testing.T) {
	// given: A test server
	ts := createEventHttpServer(201)
	defer ts.Close()

	// and: the api as system under test
	api := buildEventsApi(ts.URL)

	t.Run("create successful", func(t *testing.T) {
		err := api.CreateEvent(createEvent)

		if err != nil {
			t.Fatalf("CreateEvent() got an unexpected error: %s", err.Error())
		}

		if eventCapture == nil {
			t.Fatalf("CreateEvent() Captured event is nil.")
		}

		if !reflect.DeepEqual(createEvent, eventCapture) {
			t.Errorf("CreateEvent() event = %v, want %v", createEvent, eventCapture)
		}
	})
}

func TestEvents_Create_Event_BadRequest(t *testing.T) {
	// given: A test server
	ts := createEventHttpServer(400)
	defer ts.Close()

	// and: the api as system under test
	api := buildEventsApi(ts.URL)

	t.Run("non existing device id", func(t *testing.T) {
		err := api.CreateEvent(createEvent)

		if err == nil {
			t.Errorf("CreateEvent() expected error on 400 - bad request")
			return
		}
	})
}
