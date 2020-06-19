package gomulocity_event

import (
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"net/url"
	"testing"
)

func TestEventCollection_CurrentPage(t *testing.T) {
	tests := []struct {
		name         string
		givenPage    int
		expectedPage int
	}{
		{"Correct one digit page number", 5, 5},
		{"Correct multiple digit page number", 112, 112},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			selfUrl, _ := url.Parse(fmt.Sprintf("http://0815.cumulocity.com/event/events?pageSize=100&currentPage=%d", tt.givenPage))
			c := &EventCollection{
				Next:   nil,
				Self:   selfUrl,
				Prev:   nil,
				Events: []Event{},
				Statistics: &generic.PagingStatistics{
					PageSize:    5,
					CurrentPage: tt.givenPage,
				},
			}
			if got := c.CurrentPage(); got != tt.expectedPage {
				t.Errorf("CurrentPage() = %v, want %v", got, tt.expectedPage)
			}
		})
	}
}
