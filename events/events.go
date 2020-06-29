package gomulocity_event

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

const EVENT_ACCEPT_HEADER = "application/vnd.com.nsn.cumulocity.eventApi+json"

func NewEventsApi(client Client) Events {
	return &events{client, "/event/events"}
}

type Events interface {
	CreateEvent(event *CreateEvent) (*Event, error)
	UpdateEvent(eventId string, event *UpdateEvent) error
	DeleteEvent(eventId string) error

	Get(eventId string) (*Event, error)
	GetForDevice(source string, pageSize int) (*EventCollection, error)
	Find(query EventQuery) (*EventCollection, error)
	NextPage(c *EventCollection) (*EventCollection, error)
	PreviousPage(c *EventCollection) (*EventCollection, error)
}

type EventQuery struct {
	DateFrom     *time.Time
	DateTo       *time.Time
	FragmentType string
	Type         string
	Source       string
	PageSize     int
}

func (q EventQuery) QueryParams() (string, error) {
	if q.PageSize < 0 || q.PageSize > 2000 {
		return "", errors.New(fmt.Sprintf("The page size must be between 1 and 2000. Was %d", q.PageSize))
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
	client   Client
	basePath string
}

func (e *events) CreateEvent(event *CreateEvent) (*Event, error) {
	bytes, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error while marhalling the event: %s", err.Error())
	}

	body, status, err := e.client.post(e.basePath, bytes, AcceptHeader(EVENT_ACCEPT_HEADER))
	if err != nil {
		log.Printf("Error while posting a new event: %s", err.Error())
		return nil, err
	}
	if status != http.StatusCreated {
		return nil, createErrorFromResponse(body)
	}

	var result Event
	if len(body) > 0 {
		err = json.Unmarshal(body, &result)
		if err != nil {
			log.Printf("Error while parsing response JSON: %s", err.Error())
			return nil, err
		}
	} else {
		return nil, errors.New("GetEvent: response body was empty")
	}

	return &result, nil
}

func (e *events) UpdateEvent(eventId string, event *UpdateEvent) error {
	bytes, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error while marhalling the update event: %s", err.Error())
	}
	path := fmt.Sprintf("%s/%s", e.basePath, url.QueryEscape(eventId))

	body, status, err := e.client.put(path, bytes, EmptyHeader())
	if status != http.StatusNoContent {
		return createErrorFromResponse(body)
	}

	return err
}

func (e *events) DeleteEvent(eventId string) error {
	body, status, err := e.client.delete(fmt.Sprintf("%s/%s", e.basePath, url.QueryEscape(eventId)), EmptyHeader())

	if status != http.StatusNoContent {
		return createErrorFromResponse(body)
	}

	return err
}

func (e *events) Get(eventId string) (*Event, error) {
	body, status, err := e.client.get(fmt.Sprintf("%s/%s", e.basePath, url.QueryEscape(eventId)), EmptyHeader())

	if status != http.StatusOK {
		log.Printf("Event with id %s was not found", eventId)
		return nil, nil
	}

	var result Event
	if len(body) > 0 {
		err = json.Unmarshal(body, &result)
		if err != nil {
			log.Printf("Error while parsing response JSON: %s", err.Error())
			return nil, err
		}
	} else {
		return nil, errors.New("GetEvent: response body was empty")
	}

	return &result, nil
}

func (e *events) GetForDevice(source string, pageSize int) (*EventCollection, error) {
	return e.Find(EventQuery{Source: source, PageSize: pageSize})
}

func (e *events) Find(query EventQuery) (*EventCollection, error) {
	queryParams, err := query.QueryParams()
	if err != nil {
		return nil, err
	}

	return e.getCommon(fmt.Sprintf("%s?%s", e.basePath, queryParams))
}

func (e *events) NextPage(c *EventCollection) (*EventCollection, error) {
	return e.getPage(c.Next)
}

func (e *events) PreviousPage(c *EventCollection) (*EventCollection, error) {
	return e.getPage(c.Prev)
}

func (e *events) getPage(reference string) (*EventCollection, error) {
	if reference == "" {
		log.Print("No page reference given. Returning nil.")
		return nil, nil
	}

	nextUrl, err := url.Parse(reference)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unparsable URL given for page reference: '%s'", reference))
	}

	collection, err := e.getCommon(fmt.Sprintf("%s?%s", nextUrl.Path, nextUrl.RawQuery))
	if err != nil {
		return nil, err
	}

	if len(collection.Events) == 0 {
		log.Print("Returned collection is empty. Returning nil.")
		return nil, nil
	}

	return collection, nil
}

// -- internal

func (e *events) getCommon(path string) (*EventCollection, error) {
	body, status, err := e.client.get(path, EmptyHeader())

	if status != http.StatusOK {
		return nil, createErrorFromResponse(body)
	}

	var result EventCollection
	if len(body) > 0 {
		err = json.Unmarshal(body, &result)
		if err != nil {
			log.Printf("Error while parsing response JSON: %s", err.Error())
			return nil, err
		}
	} else {
		return nil, errors.New("GetCollection: response body was empty")
	}

	return &result, nil
}

func createErrorFromResponse(responseBody []byte) error {
	var msg map[string]interface{}
	_ = json.Unmarshal(responseBody, &msg)
	return errors.New(fmt.Sprintf("Request failed. Server returns error: {%s: %s}", msg["error"], msg["message"]))
}
