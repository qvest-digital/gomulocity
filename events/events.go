package gomulocity_event

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"time"
)

func NewEventsApi(client Client) Events {
	return &events{client, "/event/events"}
}

type Events interface {
	CreateEvent(event CreateEvent)
	UpdateEvent(event UpdateEvent)
	DeleteEvent(eventId string)

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

func (e *events) CreateEvent(event CreateEvent) {
}
func (e *events) UpdateEvent(event UpdateEvent) {
}
func (e *events) DeleteEvent(eventId string) {
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
