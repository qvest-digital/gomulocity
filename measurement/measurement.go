package measurement

import (
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"net/url"
	"strconv"
	"time"
)

type MeasurementCollection struct {
	Measurements []Measurement             `json:"measurements"`
	Self         string                    `json:"self,omitempty"`
	Statistics   *generic.PagingStatistics `json:"statistics,omitempty"`
	Prev         string                    `json:"prev,omitempty"`
	Next         string                    `json:"next,omitempty"`
}

type NewMeasurements struct {
	Measurements []NewMeasurement `json:"measurements"`
}

type Source struct {
	Id   string `json:"id"`
	Self string `json:"self,omitempty"`
}

type NewMeasurement struct {
	Time            *time.Time             `json:"time"`
	MeasurementType string                 `json:"type"`
	Source          Source                 `json:"source"`
	Metrics         map[string]interface{} `jsonc:"flat"`
}

type Measurement struct {
	Id              string                 `json:"id"`
	Self            string                 `json:"self"`
	Time            *time.Time             `json:"time"`
	MeasurementType string                 `json:"type"`
	Source          Source                 `json:"source"`
	Metrics         map[string]interface{} `jsonc:"flat"`
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
	SourceId            string
	Revert              bool // It's not a filter. It's the sort order. As per default the measurements will be delivered in ascending sort order.
	// That means, the oldest measurements are returned first. Setting to true is only valid with DateFrom and DateTo filters. In that case
	// the latest measurement of the given time period will be at the first place.
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

	if q.Revert {
		if q.DateFrom != nil && q.DateTo != nil {
			params.Add("revert", strconv.FormatBool(q.Revert))
		} else {
			return fmt.Errorf("failed to build filter: if 'Revert' parameter is set to true, 'DateFrom' and 'DateTo' should be set as well.")
		}
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

	if len(q.SourceId) > 0 {
		params.Add("source", q.SourceId)
	}
	return nil
}
