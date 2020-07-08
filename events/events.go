package events

import (
	"encoding/json"
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"log"
	"net/http"
	"net/url"
	"time"
)

const EVENT_ACCEPT_HEADER = "application/vnd.com.nsn.cumulocity.eventApi+json"

// Creates a new events api object
// client - Must be a gomulocity client.
// returns - The `Events`-api object
func NewEventsApi(client generic.Client) Events {
	return &events{client, "/event/events"}
}

type Events interface {
	// Create a new event and returns the created entity with
	// id and creation time
	CreateEvent(event *CreateEvent) (*Event, *generic.Error)

	// Updated an exiting event and returns the updated event entity.
	UpdateEvent(eventId string, event *UpdateEvent) (*Event, *generic.Error)

	// Deletes an exiting event. If error is nil, the event was deleted
	// successfully.
	DeleteEvent(eventId string) *generic.Error

	// Gets an exiting event by its id. If the id does not exists, nil is returned.
	Get(eventId string) (*Event, *generic.Error)

	// Gets a event collection by a source (aka managed object id).
	GetForDevice(source string, pageSize int) (*EventCollection, *generic.Error)

	// Returns an event collection, found by the given event query parameters.
	// all query parameters are AND concat.
	Find(query EventQuery) (*EventCollection, *generic.Error)

	// Gets the next page from an existing event collection.
	// If there is no next page, nil is returned.
	NextPage(c *EventCollection) (*EventCollection, *generic.Error)

	// Gets the previous page from an existing event collection.
	// If there is no previous page, nil is returned.
	PreviousPage(c *EventCollection) (*EventCollection, *generic.Error)
}

type EventQuery struct {
	DateFrom     *time.Time
	DateTo       *time.Time
	FragmentType string
	Type         string
	Source       string
	PageSize     int
}

func (q EventQuery) QueryParams() (string, *generic.Error) {
	if q.PageSize < 0 || q.PageSize > 2000 {
		return "", clientError(fmt.Sprintf("The page size must be between 1 and 2000. Was %d", q.PageSize), "QueryParams")
	}

	params := url.Values{}

	if q.DateFrom != nil {
		params.Add("dateFrom", q.DateFrom.Format(time.RFC3339))
	}

	if q.DateTo != nil {
		params.Add("dateTo", q.DateTo.Format(time.RFC3339))
	}

	if q.PageSize > 0 {
		params.Add("pageSize", fmt.Sprintf("%d", q.PageSize))
	}

	if len(q.FragmentType) > 0 {
		params.Add("fragmentType", q.FragmentType)
	}

	if len(q.Type) > 0 {
		params.Add("type", q.Type)
	}

	if len(q.Source) > 0 {
		params.Add("source", q.Source)
	}

	return params.Encode(), nil
}

type events struct {
	client   generic.Client
	basePath string
}

func (e *events) DeleteEvent(eventId string) *generic.Error {
	body, status, err := e.client.Delete(fmt.Sprintf("%s/%s", e.basePath, url.QueryEscape(eventId)), generic.EmptyHeader())

	if err != nil {
		return clientError(fmt.Sprintf("Error while deleting an event: %s", err.Error()), "DeleteEvent")
	}

	if status != http.StatusNoContent {
		return createErrorFromResponse(body)
	}

	return nil
}

func (e *events) CreateEvent(event *CreateEvent) (*Event, *generic.Error) {
	bytes, err := json.Marshal(event)
	if err != nil {
		return nil, clientError(fmt.Sprintf("Error while marhalling the event: %s", err.Error()), "CreateEvent")
	}

	body, status, err := e.client.Post(e.basePath, bytes, generic.AcceptHeader(EVENT_ACCEPT_HEADER))
	if err != nil {
		return nil, clientError(fmt.Sprintf("Error while posting a new event: %s", err.Error()), "CreateEvent")
	}
	if status != http.StatusCreated {
		return nil, createErrorFromResponse(body)
	}

	return parseEventResponse(body)
}

func (e *events) UpdateEvent(eventId string, event *UpdateEvent) (*Event, *generic.Error) {
	bytes, err := json.Marshal(event)
	if err != nil {
		return nil, clientError(fmt.Sprintf("Error while marhalling the update event: %s", err.Error()), "UpdateEvent")
	}

	path := fmt.Sprintf("%s/%s", e.basePath, url.QueryEscape(eventId))
	body, status, err := e.client.Put(path, bytes, generic.AcceptHeader(EVENT_ACCEPT_HEADER))
	if err != nil {
		return nil, clientError(fmt.Sprintf("Error while updating an event: %s", err.Error()), "UpdateEvent")
	}
	if status != http.StatusOK {
		return nil, createErrorFromResponse(body)
	}

	return parseEventResponse(body)
}

func (e *events) Get(eventId string) (*Event, *generic.Error) {
	body, status, err := e.client.Get(fmt.Sprintf("%s/%s", e.basePath, url.QueryEscape(eventId)), generic.EmptyHeader())

	if err != nil {
		return nil, clientError(fmt.Sprintf("Error while getting an event: %s", err.Error()), "Get")
	}
	if status != http.StatusOK {
		return nil, nil
	}

	return parseEventResponse(body)
}

func (e *events) GetForDevice(source string, pageSize int) (*EventCollection, *generic.Error) {
	return e.Find(EventQuery{Source: source, PageSize: pageSize})
}

func (e *events) Find(query EventQuery) (*EventCollection, *generic.Error) {
	queryParams, err := query.QueryParams()
	if err != nil {
		return nil, err
	}

	return e.getCommon(fmt.Sprintf("%s?%s", e.basePath, queryParams))
}

func (e *events) NextPage(c *EventCollection) (*EventCollection, *generic.Error) {
	return e.getPage(c.Next)
}

func (e *events) PreviousPage(c *EventCollection) (*EventCollection, *generic.Error) {
	return e.getPage(c.Prev)
}

// -- internal

func parseEventResponse(body []byte) (*Event, *generic.Error) {
	var result Event
	if len(body) > 0 {
		err := generic.ObjectFromJson(body, &result)
		if err != nil {
			return nil, clientError(fmt.Sprintf("Error while parsing response JSON: %s", err.Error()), "ResponseParser")
		}
	} else {
		return nil, clientError("Response body was empty", "GetEvent")
	}

	return &result, nil
}

func (e *events) getPage(reference string) (*EventCollection, *generic.Error) {
	if reference == "" {
		log.Print("No page reference given. Returning nil.")
		return nil, nil
	}

	nextUrl, err := url.Parse(reference)
	if err != nil {
		return nil, clientError(fmt.Sprintf("Unparsable URL given for page reference: '%s'", reference), "GetPage")
	}

	collection, err2 := e.getCommon(fmt.Sprintf("%s?%s", nextUrl.Path, nextUrl.RawQuery))
	if err2 != nil {
		return nil, err2
	}

	if len(collection.Events) == 0 {
		log.Print("Returned collection is empty. Returning nil.")
		return nil, nil
	}

	return collection, nil
}

func (e *events) getCommon(path string) (*EventCollection, *generic.Error) {
	body, status, err := e.client.Get(path, generic.EmptyHeader())

	if status != http.StatusOK {
		return nil, createErrorFromResponse(body)
	}

	var result EventCollection
	if len(body) > 0 {
		err = generic.ObjectFromJson(body, &result)
		if err != nil {
			return nil, clientError(fmt.Sprintf("Error while parsing response JSON: %s", err.Error()), "GetCollection")
		}
	} else {
		return nil, clientError("Response body was empty", "GetCollection")
	}

	return &result, nil
}

func clientError(message string, info string) *generic.Error {
	return &generic.Error{
		ErrorType: "ClientError",
		Message:   message,
		Info:      info,
	}
}

func createErrorFromResponse(responseBody []byte) *generic.Error {
	var err generic.Error
	_ = json.Unmarshal(responseBody, &err)
	return &err
}
