package models

import (
	"time"
)

type Measurement struct {
	id              string        `json:"id"`
	time            time.Time     `json:"time"`
	measurementType string        `json:"type"`
	source          ManagedObject `json:"managedObject"`
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
