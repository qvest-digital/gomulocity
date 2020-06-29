package alarm

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

/*
See: https://cumulocity.com/guides/reference/alarms/#delete-delete-an-alarm-collection
 */
type AlarmsFilter struct {
	Status				[]Status	// Comma separated alarm statuses, for example ACTIVE,CLEARED.
	SourceId 			string	// Source device id.
	WithSourceAssets	bool	// When set to true also alarms for related source assets will be removed.
								// When this parameter is provided also source must be defined.
	WithSourceDevices	bool	// When set to true also alarms for related source devices will be removed.
								// When this parameter is provided also source must be defined.
	Resolved			string	// When set to true only resolved alarms will be removed (the one with status CLEARED),
								// false means alarms with status ACTIVE or ACKNOWLEDGED.
	Severity			Severity	// Alarm severity, for example MINOR.
	DateFrom			time.Time	// Start date or date and time of alarm occurrence.
	DateTo				time.Time	// End date or date and time of alarm occurrence.
	Type				string	// Alarm type.
}

/*
https://cumulocity.com/guides/reference/alarms/#put-bulk-update-of-alarm-collection
 */
type UpdateAlarmsFilter struct {
	Status				Status
	SourceId 			string
	Resolved			string
	Severity			Severity
	DateFrom			time.Time
	DateTo				time.Time
}

func (daf AlarmsFilter) appendFilter(r *http.Request) error {
	q := r.URL.Query()

	if len(daf.Status) > 0 {
		var statusesAsString []string
		for _, status := range daf.Status {
			fmt.Println(status)
			statusesAsString = append(statusesAsString, string(status))
		}
		q.Set("status", strings.Join(statusesAsString, ","))
	}

	if len(daf.SourceId) > 0 {
		q.Set("source", daf.SourceId)
	}

	if daf.WithSourceAssets {
		if len(daf.SourceId) == 0 {
			return fmt.Errorf("failed to build filter: when 'WithSourceAssets' parameter is defined also SourceID must be set.")
		}
		q.Set("withSourceAssets", "true")
	}

	if daf.WithSourceDevices {
		if len(daf.SourceId) == 0 {
			return fmt.Errorf("failed to build filter: when 'WithSourceDevices' parameter is defined also SourceID must be set.")
		}
		q.Set("withSourceDevices", "true")
	}

	if len(daf.Resolved) > 0 {
		resolved, err := strconv.ParseBool(daf.Resolved)
		if err != nil {
			return fmt.Errorf("failed to build filter: if 'Resolved' parameter is set, only 'true' and 'false' values are accepted.")
		}
		q.Set("resolved", strconv.FormatBool(resolved))
	}

	if len(daf.Severity) > 0 {
		q.Set("severity", fmt.Sprintf("%s", daf.Severity))
	}

	if !daf.DateFrom.IsZero() {
		q.Set("dateFrom", daf.DateFrom.Format(time.RFC3339))
	}

	if !daf.DateTo.IsZero() {
		q.Set("dateTo", daf.DateTo.Format(time.RFC3339))
	}

	if len(daf.Type) > 0 {
		q.Set("type", daf.Type)
	}

	r.URL.RawQuery = q.Encode()
	return nil
}

func (gaf UpdateAlarmsFilter) appendFilter(r *http.Request) error {
	q := r.URL.Query()

	if len(gaf.Status) > 0 {
		q.Set("status", fmt.Sprintf("%s", gaf.Status))
	}

	if len(gaf.SourceId) > 0 {
		q.Set("source", gaf.SourceId)
	}

	if len(gaf.Resolved) > 0 {
		resolved, err := strconv.ParseBool(gaf.Resolved)
		if err != nil {
			return fmt.Errorf("failed to build filter: if 'Resolved' parameter is set, only 'true' and 'false' values are accepted.")
		}
		q.Set("resolved", strconv.FormatBool(resolved))
	}

	if len(gaf.Severity) > 0 {
		q.Set("severity", fmt.Sprintf("%s", gaf.Severity))
	}

	if !gaf.DateFrom.IsZero() {
		q.Set("dateFrom", gaf.DateFrom.Format(time.RFC3339))
	}

	if !gaf.DateTo.IsZero() {
		q.Set("dateTo", gaf.DateTo.Format(time.RFC3339))
	}

	r.URL.RawQuery = q.Encode()
	return nil
}
