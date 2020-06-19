package gomulocity_event

import (
	"github.com/tarent/gomulocity/generic"
	"log"
	"net/url"
	"regexp"
	"strconv"
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
	Statistics *generic.PagingStatistics `json:"statistics"`
}

func (c *EventCollection) CurrentPage() int {
	parameterPattern := regexp.MustCompile("^.*pageSize=(\\d+).*currentPage=(\\d+)$")
	match := parameterPattern.FindStringSubmatch(c.Self.RequestURI())

	if len(match) != 3 {
		log.Printf("Could not extract the current page number from self URL: %s", c.Self.String())
		return 0
	}

	val, err := strconv.Atoi(match[2])
	if err != nil {
		log.Printf("Could not extract the current page number from self URL: %s", c.Self.String())
		return 0
	}

	return val
}
