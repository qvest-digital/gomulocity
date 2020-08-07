package audit

import (
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"net/http"
	"net/http/httptest"
	"time"
)

func buildAuditApi(url string) AuditApi {
	httpClient := http.DefaultClient
	client := &generic.Client{
		HTTPClient: httpClient,
		BaseURL:    url,
		Username:   "foo",
		Password:   "bar",
	}
	return NewAuditApi(client)
}

func buildHttpServer(status int, body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
}

var t, _ = time.Parse(time.RFC3339, "2020-07-21T09:47:21.092Z")
var creationTime, _ = time.Parse(time.RFC3339, "2020-07-21T09:47:22.166Z")

var auditID = "100197821"

var testAuditRecord = &AuditRecord{
	Severity:     "MAJOR",
	Activity:     "Alarm created",
	CreationTime: creationTime,
	Source: struct {
		Self string `json:"self"`
		ID   string `json:"id"`
	}{
		Self: "https://t200588189.cumulocity.com/inventory/managedObjects/10013234",
		ID:   "10013234",
	},
	Type:        "Alarm",
	Self:        "https://t200588189.cumulocity.com/audit/auditRecords/" + auditID,
	Time:        t,
	Text:        "Device name: 'TestDevice123', alarm text: 'Alarm 1.3 occured'",
	ID:          auditID,
	User:        "gfa-agent",
	Application: "Omniscape",
	Changes: []Changes{
		{
			NewValue:      "WARNING",
			Attribute:     "severity",
			Type:          "com.cumulocity.model.event.CumulocitySeverities",
			PreviousValue: "MINOR",
		},
	},
	AuditSourceDevice: struct {
		ID string `json:"id"`
	}{
		ID: "3249962",
	},
}

var testAuditRecordJSON = `
{
    "id": "` + auditID + `",
    "type": "Alarm",
    "self": "https://t200588189.cumulocity.com/audit/auditRecords/` + auditID + `",
    "time": "2020-07-21T09:47:21.092Z",
    "text": "Device name: 'TestDevice123', alarm text: 'Alarm 1.3 occured'",
    "user": "gfa-agent",
    "source": {
        "id": "10013234",
        "self": "https://t200588189.cumulocity.com/inventory/managedObjects/10013234"
    },
    "changes": [
        {
            "type": "com.cumulocity.model.event.CumulocitySeverities",
            "newValue": "WARNING",
            "attribute": "severity",
            "previousValue": "MINOR"
        }
    ],
    "severity": "MAJOR",
    "activity": "Alarm created",
    "application": "Omniscape",
    "creationTime": "2020-07-21T09:47:22.166Z",
    "com_cumulocity_model_event_AuditSourceDevice": {
        "id": "3249962"
    }
}
`

var testAuditRecords = &AuditRecordCollection{
	Self:         "https://t200588189.cumulocity.com/audit/auditRecords?pageSize=5&currentPage=2",
	Next:         "https://t200588189.cumulocity.com/audit/auditRecords?pageSize=5&currentPage=3",
	Prev:         "https://t200588189.cumulocity.com/audit/auditRecords?pageSize=5&currentPage=1",
	AuditRecords: []AuditRecord{*testAuditRecord},
	Statistics: &generic.PagingStatistics{
		TotalPages:  1946,
		PageSize:    5,
		CurrentPage: 1,
	},
}

var testAuditRecordsJSON = `
{
	"self":"https://t200588189.cumulocity.com/audit/auditRecords?pageSize=5&currentPage=2",
	"next":"https://t200588189.cumulocity.com/audit/auditRecords?pageSize=5&currentPage=3",
	"prev":"https://t200588189.cumulocity.com/audit/auditRecords?pageSize=5&currentPage=1",
	"auditRecords":[` + testAuditRecordJSON + `],
	"statistics":{
		"totalPages":1946,
		"pageSize":5,
		"currentPage":1
	}
}
`

var testCreateAuditRecord = &AuditRecord{
	Activity: "Alarm created",
	Type:     "Alarm",
	Time:     t,
	Text:     "Device name: 'TestDevice123', alarm text: 'Alarm 1.3 occured'",
}

var testCreateAuditRecordJSON = `
{
	"activity":"Alarm created",
	"type":"alarm",
	"time":"2020-07-21T09:47:21.092Z",
	"text":"Device name: 'TestDevice123', alarm text: 'Alarm 1.3 occured'"
}
`

var erroneousResponseJSON = `
{
"error": "general/internalError",
"info": "https://www.cumulocity.com/guides/reference-guide/#error_reporting"
}
`

var erroneousResponseCreateAuditJSON = `
{
    "error": "undefined/validationError",
    "message": "Following mandatory fields should be included: activity,type,time,text",
    "info": "https://www.cumulocity.com/guides/reference-guide/#error_reporting"
}
`

var createErroneousResponse = func(status int, message string) *generic.Error {
	return &generic.Error{
		ErrorType: fmt.Sprintf("%v: general/internalError", status),
		Message:   message,
		Info:      "https://www.cumulocity.com/guides/reference-guide/#error_reporting",
	}
}

var createErroneousResponseCreateAudit = func(status int) *generic.Error {
	return &generic.Error{
		ErrorType: fmt.Sprintf("%v: undefined/validationError", status),
		Message:   "Following mandatory fields should be included: activity,type,time,text",
		Info:      "https://www.cumulocity.com/guides/reference-guide/#error_reporting",
	}
}

var auditCollectionTemplate = `{
    "next": "https://t200588189.cumulocity.com/audit/auditRecords?pageSize=5&currentPage=2",
    "self": "https://t200588189.cumulocity.com/audit/auditRecords?pageSize=5&currentPage=1",
    "auditRecords": [%s],
    "statistics": {
        "currentPage": 1,
        "pageSize": 5
	}
}`
