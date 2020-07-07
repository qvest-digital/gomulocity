package measurements

import (
	"net/url"
	"time"

	"github.com/tarent/gomulocity/models"
)

type MeasurementCollection struct {
	measurements []Measurement `json:"measurements"`
}

type Measurement struct {
	id              string               `json:"id"`
	time            time.Time            `json:"time"`
	measurementType string               `json:"type"`
	source          models.ManagedObject `json:"managedObject"`
}

func (m Measurement) getID() string {
	return m.id
}

func (m Measurement) getTime() time.Time {
	return m.time
}

func (m Measurement) getType() string {
	return m.measurementType
}


type MeasurementQuery struct {
	DateFrom            *time.Time
	DateTo              *time.Time
	Type                string
	ValueFragmentType   string
	ValueFragmentSeries string
	source              *models.ManagedObject
}

func (q MeasurementQuery) QueryParams() string {
	params := url.Values{}
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

	if len(q.source.ID) > 0 {
		params.Add("source", q.source.ID)
	}
	return params.Encode()
}
