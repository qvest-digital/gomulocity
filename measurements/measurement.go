package measurements

import (
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"net/url"
	"time"
)

type MeasurementCollection struct {
	Measurements []Measurement             `json:"measurements"`
	Self         string                    `json:"self,omitempty"`
	Statistics   *generic.PagingStatistics `json:"statistics,omitempty"`
	Prev         string                    `json:"prev,omitempty"`
	Next         string                    `json:"next,omitempty"`
}

type Source struct {
	Id   string `json:"id"`
	Self string `json:"self,omitempty"`
}

type Measurement struct {
	Id              string    `json:"id,omitempty"`
	Self            string    `json:"self,omitempty"`
	Time            time.Time `json:"time"`
	MeasurementType string    `json:"type"`
	Source          Source    `json:"source"`
	Temperature     Temperature
}

type Temperature struct {
	Cellar      ValueFragment
	GroundFloor ValueFragment
}

type ValueFragment struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit,omitempty"`
}

type MeasurementQuery struct {
	DateFrom            *time.Time
	DateTo              *time.Time
	Type                string
	ValueFragmentType   string
	ValueFragmentSeries string
	sourceId            string
}

func (q MeasurementQuery) QueryParams(params *url.Values) error {
	if params == nil {
		return fmt.Errorf("The provided parameter values must not be nil!")
	}

	if q.DateFrom != nil {
		params.Add("dateFrom", q.DateFrom.Format(time.RFC3339))
	}

	if q.DateTo != nil {
		params.Add("dateTo", q.DateTo.Format(time.RFC3339))
	}

	if len(q.Type) > 0 {
		params.Add("type", q.Type)
	}

	if len(q.ValueFragmentType) > 0 {
		params.Add("valueFragmentType", q.ValueFragmentType)
	}

	if len(q.ValueFragmentSeries) > 0 {
		params.Add("valueFragmentSeries", q.ValueFragmentSeries)
	}

	if len(q.sourceId) > 0 {
		params.Add("source", q.sourceId)
	}
	return nil
}
