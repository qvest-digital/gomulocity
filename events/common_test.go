package events

import (
	"encoding/json"
	"github.com/tarent/gomulocity/generic"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"
)

var eventTime, _ = time.Parse(time.RFC3339, "2020-06-26T10:43:25.130Z")
var responseEvent = &Event{
	Id:               "",
	Type:             "TestEvent",
	Time:             eventTime,
	CreationTime:     eventTime,
	Text:             "This is my test event",
	Source:           Source{Id: "4711"},
	Self:             "https://t0815.cumulocity.com/event/events/1337",
	AdditionalFields: map[string]interface{}{},
}

func buildEventsApi(url string) Events {
	httpClient := http.DefaultClient
	client := generic.Client{
		HTTPClient: httpClient,
		BaseURL:    url,
		Username:   "foo",
		Password:   "bar",
	}
	return NewEventsApi(client)
}

func buildHttpServer(status int, body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
}

func updateEventHttpServer(status int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)

		var event UpdateEvent
		_ = json.Unmarshal(body, &event)
		updateEventCapture = &event
		requestCapture = r

		w.WriteHeader(status)
		response, _ := json.Marshal(responseEvent)
		_, _ = w.Write(response)
	}))
}

func createEventHttpServer(status int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)

		var event CreateEvent
		_ = json.Unmarshal(body, &event)
		createEventCapture = &event
		requestCapture = r

		w.WriteHeader(status)
		response, _ := json.Marshal(responseEvent)
		_, _ = w.Write(response)
	}))
}

var requestCapture *http.Request
var createEventCapture *CreateEvent
var updateEventCapture *UpdateEvent
var updateUrlCapture string

var deviceId = "1111111"
var eventId = "2222222"
var event = `{
	"creationTime": "2020-01-01T01:00:10.000Z",
	"source": {
		"name": "test-device",
		"self": "https://t0815.cumulocity.com/inventory/managedObjects/1111111",
		"id": "1111111"
	},
	"type": "threshold",
	"self": "https://t0815.cumulocity.com/event/events/2222222",
	"time": "2020-01-01T01:00:00.000Z",
	"text": "over 21°C",
	"id": "2222222",
	"custom1": "Hello",
	"custom2": [
		"Foo", "Bar"
	]
}`
var eventCollectionTemplate = `{
    "next": "https://t0815.cumulocity.com/event/events?source=1111111&pageSize=5&currentPage=2",
    "self": "https://t0815.cumulocity.com/event/events?source=1111111&pageSize=5&currentPage=1",
    "events": [%s],
    "statistics": {
        "currentPage": 1,
        "pageSize": 5
    }
}`
