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
	Find(query EventQuery) *EventCollection
	NextPage(c *EventCollection) *EventCollection
	PrevPage(c *EventCollection) *EventCollection
}

type EventQuery struct {
	DateFrom     time.Time
	DateTo       time.Time
	FragmentType string
	Type         string
	Source       string
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
func (e *events) Find(query EventQuery) *EventCollection {
	return nil
}
func (e *events) NextPage(c *EventCollection) *EventCollection {
	return nil
}
func (e *events) PrevPage(c *EventCollection) *EventCollection {
	return nil
}

func request(c *EventCollection) *EventCollection {
	return nil
}
