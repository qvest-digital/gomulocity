package alarm

import (
	"github.com/tarent/gomulocity/generic"
	"time"
)

const (
	ALARM_API_PATH        = "/alarm/alarms"
	ALARM_TYPE            = "application/vnd.com.nsn.cumulocity.alarm+json"
	ALARM_COLLECTION_TYPE = "application/vnd.com.nsn.cumulocity.alarmCollection+json"
)

type Status string

const (
	ACTIVE       Status = "ACTIVE"
	ACKNOWLEDGED Status = "ACKNOWLEDGED"
	CLEARED      Status = "CLEARED"
)

type Severity string

const (
	CRITICAL Severity = "CRITICAL"
	MAJOR    Severity = "MAJOR"
	MINOR    Severity = "MINOR"
	WARNING  Severity = "WARNING"
)

type Source struct {
	Id   string `json:"id"`
	Self string `json:"self,omitempty"`
	Name string `json:"name,omitempty"`
}

/*
Represents cumulocity's alarm structure for creation purposes.
See: https://cumulocity.com/guides/reference/alarms/#post-create-a-new-alarm
*/
type NewAlarm struct {
	Type             string                 `json:"type"`
	Time             time.Time              `json:"time"`
	Text             string                 `json:"text"`
	Source           Source                 `json:"source"`
	Status           Status                 `json:"status"`
	Severity         Severity               `json:"severity"`
	AdditionalFields map[string]interface{} `jsonc:"flat"`
}

/*
Represents cumulocity's alarm 'application/vnd.com.nsn.cumulocity.alarm+json'.
See: https://cumulocity.com/guides/reference/alarms/#alarm
https://cumulocity.com/guides/reference/alarms/#post-create-a-new-alarm
*/
type Alarm struct {
	Id           string    `json:"id,omitempty"`
	Self         string    `json:"self,omitempty"`
	CreationTime *time.Time `json:"creationTime,omitempty"`

	Type     string    `json:"type,omitempty"`
	Time     *time.Time `json:"time,omitempty"`
	Text     string    `json:"text,omitempty"`
	Source   Source    `json:"source,omitempty"`
	Status   Status    `json:"status,omitempty"`
	Severity Severity  `json:"severity,omitempty"`

	Count               int       `json:"count,omitempty"`
	FirstOccurrenceTime *time.Time `json:"firstOccurrenceTime,omitempty"`

	// TODO: object - 0..n additional properties of the alarm.
}

/*
Represents cumulocity's alarm structure for update purposes.
See: https://cumulocity.com/guides/reference/alarms/#update-an-alarm
*/
type UpdateAlarm struct {
	Text             string                 `json:"text,omitempty"`
	Status           Status                 `json:"status,omitempty"`
	Severity         Severity               `json:"severity,omitempty"`
	AdditionalFields map[string]interface{} `jsonc:"flat"`
}

/*
AlarmCollection represent cumulocity's 'application/vnd.com.nsn.cumulocity.alarmCollection+json'.
See: https://cumulocity.com/guides/reference/alarms/#alarm-collection
*/
type AlarmCollection struct {
	Self       string                    `json:"self"`
	Alarms     []Alarm                   `json:"alarms"`
	Statistics *generic.PagingStatistics `json:"statistics,omitempty"`
	Prev       string                    `json:"prev,omitempty"`
	Next       string                    `json:"next,omitempty"`
}
