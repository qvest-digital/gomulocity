package gomulocity_event

import (
	"fmt"
	"net/url"
	"testing"
)

func TestEventCollection_CurrentPage(t *testing.T) {
	tests := []struct {
		name         string
		givenPage    string
		expectedPage int
	}{
		{"Correct one digit page number", "5", 5},
		{"Correct multiple digit page number", "112", 112},
		{"No digit", "foobar", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			selfUrl, _ := url.Parse(fmt.Sprintf("http://0815.cumulocity.com/event/events?pageSize=100&currentPage=%s", tt.givenPage))
			c := &EventCollection{
				Next:       nil,
				Self:       selfUrl,
				Prev:       nil,
				Events:     []Event{},
				Statistics: nil,
			}
			if got := c.CurrentPage(); got != tt.expectedPage {
				t.Errorf("CurrentPage() = %v, want %v", got, tt.expectedPage)
			}
		})
	}
}
