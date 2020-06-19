package gomulocity_event

import (
	"github.com/tarent/gomulocity/generic"
	"net/url"
	"time"
)

var contentType = "application/vnd.com.nsn.cumulocity.event+json"

// application/vnd.com.nsn.cumulocity.event+json
type CreateEvent struct {
	Type   string    `json:"type"`
	Time   time.Time `json:"time"`
	Text   string    `json:"test"`
	Source struct {
		Id string `json:"id"`
	} `json:"source"`
}

type UpdateEvent struct {
	Text string `json:"test"`
}

// application/vnd.com.nsn.cumulocity.event+json
type Event struct {
	Id           string    `json:"id"`
	Type         string    `json:"type"`
	Time         time.Time `json:"time"`
	CreationTime time.Time `json:"creationTime"`
	Text         string    `json:"test"`
	Source       struct {
		Id   string  `json:"id"`
		Self url.URL `json:"self"`
	} `json:"source"`
	Self url.URL `json:"self"`
}

// application/vnd.com.nsn.cumulocity.eventCollection+json
type EventCollection struct {
	Next       *url.URL                  `json:"next"`
	Self       *url.URL                  `json:"self"`
	Prev       *url.URL                  `json:"prev"`
	Events     []Event                   `json:"events"`
	Statistics *generic.PagingStatistics `json:"statistics"` // ToDo: Check for dependencies vs. module singularity!
}

func (c *EventCollection) CurrentPage() int {
	return c.Statistics.CurrentPage
}
