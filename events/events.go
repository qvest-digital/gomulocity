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

func NewEventsApi(client Client) Events {
	return &events{client, "/event/events"}
}

type Events interface {
	CreateEvent(event *CreateEvent) error
	UpdateEvent(eventId string, event *UpdateEvent) error
	DeleteEvent(eventId string) error

	Get(eventId string) (*Event, error)
	GetForDevice(source string) (*EventCollection, error)
	Find(query EventQuery) (*EventCollection, error)
	NextPage(c *EventCollection) *EventCollection
	PrevPage(c *EventCollection) *EventCollection
}

type EventQuery struct {
	DateFrom     *time.Time
	DateTo       *time.Time
	FragmentType string
	Type         string
	Source       string
}

func (q EventQuery) QueryParams() string {
	params := url.Values{}

	if q.DateFrom != nil {
		params.Add("dateFrom", q.DateFrom.Format(time.RFC3339))
	}

	if q.DateTo != nil {
		params.Add("dateTo", q.DateTo.Format(time.RFC3339))
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

	return params.Encode()
}

type events struct {
	client   Client
	basePath string
}

func (e *events) CreateEvent(event *CreateEvent) error {
	bytes, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error while marhalling the event: %s", err.Error())
	}

	response, status, err := e.client.post(e.basePath, bytes)
	if err != nil {
		log.Printf("Error while posting a new event: %s", err.Error())
		return err
	}
	if status != http.StatusCreated {
		var msg map[string]interface{}
		_ = json.Unmarshal(response, &msg)
		return errors.New(fmt.Sprintf("Event creation failed. Server returns error: %s", msg["error"]))
	}

	return nil
}

func (e *events) UpdateEvent(eventId string, event *UpdateEvent) error {
	bytes, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error while marhalling the update event: %s", err.Error())
	}
	path := fmt.Sprintf("%s/%s", e.basePath, url.QueryEscape(eventId))

	body, status, err := e.client.put(path, bytes)

	if status != http.StatusOK {
		var msg map[string]interface{}
		_ = json.Unmarshal(body, &msg)
		return errors.New(fmt.Sprintf("Event update failed. Server returns error: %s", msg["error"]))
	}

	return err
}

func (e *events) DeleteEvent(eventId string) error {
	body, status, err := e.client.delete(fmt.Sprintf("%s/%s", e.basePath, url.QueryEscape(eventId)))

	if status != http.StatusNoContent {
		var msg map[string]interface{}
		_ = json.Unmarshal(body, &msg)
		return errors.New(fmt.Sprintf("Event creation failed. Server returns error: %s", msg["error"]))
	}

	return err
}

func (e *events) Get(eventId string) (*Event, error) {
	response, status, err := e.client.get(fmt.Sprintf("%s/%s", e.basePath, url.QueryEscape(eventId)))

	if status != 200 {
		log.Printf("Event with id %s was not found", eventId)
		return nil, nil
	}

	var result Event
	if len(response) > 0 {
		err = json.Unmarshal(response, &result)
		if err != nil {
			log.Printf("Error while parsing response JSON: %s", err.Error())
			return nil, err
		}
	} else {
		return nil, errors.New("GetEvent: response body was empty")
	}

	return &result, nil
}

func (e *events) GetForDevice(source string) (*EventCollection, error) {
	params := url.Values{}
	params.Add("source", source)
	response, _, err := e.client.get(fmt.Sprintf("%s?%s", e.basePath, params.Encode()))

	var result EventCollection
	if len(response) > 0 {
		err = json.Unmarshal(response, &result)
		if err != nil {
			log.Printf("Error while parsing response JSON: %s", err.Error())
			return nil, err
		}
	} else {
		return nil, errors.New("GetForDevice: response body was empty")
	}

	return &result, nil
}
func (e *events) Find(query EventQuery) (*EventCollection, error) {
	body, status, err := e.client.get(fmt.Sprintf("%s?%s", e.basePath, query.QueryParams()))

	if status != http.StatusOK {
		var msg map[string]interface{}
		_ = json.Unmarshal(body, &msg)
		return nil, errors.New(fmt.Sprintf("Query failed. Server returns error: %s", msg["error"]))
	}

	var result EventCollection
	if len(body) > 0 {
		err = json.Unmarshal(body, &result)
		if err != nil {
			log.Printf("Error while parsing response JSON: %s", err.Error())
			return nil, err
		}
	} else {
		return nil, errors.New("Find: response body was empty")
	}

	return &result, nil
}
func (e *events) NextPage(c *EventCollection) *EventCollection {
	return nil
}
func (e *events) PrevPage(c *EventCollection) *EventCollection {
	return nil
}
