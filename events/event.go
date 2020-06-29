package events

import (
	"github.com/tarent/gomulocity/generic"
	"time"
)

type Source struct {
	Id   string `json:"id"`
	Self string `json:"self,omitempty"`
}

// application/vnd.com.nsn.cumulocity.event+json
type CreateEvent struct {
	Type   string    `json:"type"`
	Time   time.Time `json:"time"`
	Text   string    `json:"test"`
	Source Source    `json:"source"`
}

type UpdateEvent struct {
	Text string `json:"test"`
}

// ---- Event
// application/vnd.com.nsn.cumulocity.event+json
type Event struct {
	Id           string    `json:"id"`
	Type         string    `json:"type"`
	Time         time.Time `json:"time"`
	CreationTime time.Time `json:"creationTime"`
	Text         string    `json:"test"`
	Source       Source    `json:"source"`
	Self         string    `json:"self"`
}

// application/vnd.com.nsn.cumulocity.eventCollection+json
// ---- EventCollection
type EventCollection struct {
	Next       string                    `json:"next"`
	Self       string                    `json:"self"`
	Prev       string                    `json:"prev"`
	Events     []Event                   `json:"events"`
	Statistics *generic.PagingStatistics `json:"statistics"` // ToDo: Check for dependencies vs. module singularity!
}

func (c *EventCollection) CurrentPage() int {
	return c.Statistics.CurrentPage
}
