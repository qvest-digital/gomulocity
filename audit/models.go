package audit

import (
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"net/url"
	"time"
)

const (
	AUDIT_RECORD_ACCEPT = "application/vnd.com.nsn.cumulocity.auditrecord+json"
	AUDIT_CONTENT_TYPE  = "application/vnd.com.nsn.cumulocity.auditrecord+json"
)

type (
	AuditRecord struct {
		Severity     string    `json:"severity"`
		Activity     string    `json:"activity"`
		CreationTime time.Time `json:"creationTime"`
		Source       struct {
			Self string `json:"self"`
			ID   string `json:"id"`
		} `json:"source"`
		Type              string    `json:"type"`
		Self              string    `json:"self"`
		Time              time.Time `json:"time"`
		Text              string    `json:"text"`
		ID                string    `json:"id"`
		User              string    `json:"user"`
		Application       string    `json:"application"`
		Changes           []Changes `json:"changes"`
		AuditSourceDevice struct {
			ID string `json:"id"`
		} `json:"com_cumulocity_model_event_AuditSourceDevice"`
	}

	Changes struct {
		NewValue      string `json:"newValue"`
		Attribute     string `json:"attribute"`
		Type          string `json:"type"`
		PreviousValue string `json:"previousValue"`
	}
)

type AuditRecordCollection struct {
	Self         string                    `json:"self"`
	Next         string                    `json:"next"`
	Prev         string                    `json:"prev"`
	AuditRecords []AuditRecord             `json:"auditRecords"`
	Statistics   *generic.PagingStatistics `json:"statistics"`
}

type AuditQuery struct {
	Revert bool // In case of executing range queries on audit logs API, like query by dateFrom and dateTo,
	// audits are returned by default in order from the newest to the oldest.
	// It is possible to change the order by adding query parameter “revert=false” to the request URL
	// https://cumulocity.com/guides/reference/auditing/#audit-record-collection

	Type        string
	Application string
	User        string
}

func (a AuditQuery) QueryParams(params *url.Values) error {
	if params == nil {
		return fmt.Errorf("The provided parameter values must not be nil!")
	}

	if len(a.Type) > 0 {
		params.Add("type", a.Type)
	}
	if len(a.User) > 0 {
		params.Add("user", a.User)
	}
	if len(a.Application) > 0 {
		params.Add("application", a.Application)
	}
	if a.Revert {
		params.Add("revert", "true")
	} else {
		params.Add("revert", "false")
	}
	return nil
}
