package gomulocity_event

import "time"

func NewEventsApi() Events {
	return &events{}
}

type Events interface {
	CreateEvent(event CreateEvent)
	UpdateEvent(event UpdateEvent)
	DeleteEvent(eventId string)

	Get(eventId string) *Event
	GetForDevice(source string) *EventCollection
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

type events struct{}

func (e *events) CreateEvent(event CreateEvent) {
}
func (e *events) UpdateEvent(event UpdateEvent) {
}
func (e *events) DeleteEvent(eventId string) {
}
func (e *events) Get(eventId string) *Event {
	return nil
}
func (e *events) GetForDevice(source string) *EventCollection {
	return nil
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
