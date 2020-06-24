package alarm

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

/*
See: https://cumulocity.com/guides/reference/alarms/#delete-delete-an-alarm-collection
 */
type DeleteAlarmsFilter struct {
	Status				Status	// Comma separated alarm statuses, for example ACTIVE,CLEARED.
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
https://cumulocity.com/guides/reference/alarms/#alarm-api
 */
type GetAlarmsFilter struct {
	Status				Status
	SourceId 			string
	WithSourceAssets	bool
	WithSourceDevices	bool
	Resolved			string
	Severity			Severity
	DateFrom			time.Time
	DateTo				time.Time
	Type				string
}

func (daf DeleteAlarmsFilter) buildFilter() (string, error) {
	filter := ""

	if len(daf.Status) > 0 {
		addFilter(filter, "status", fmt.Sprintf("%s", daf.Status))
	}

	if len(daf.SourceId) > 0 {
		addFilter(filter, "source", daf.SourceId)
	}

	if daf.WithSourceAssets {
		if len(daf.SourceId) == 0 {
			return filter, fmt.Errorf("failed to build filter: when 'WithSourceAssets' parameter is defined also SourceID must be set.")
		}
		addFilter(filter, "withSourceAssets", "true")
	}

	if daf.WithSourceDevices {
		if len(daf.SourceId) == 0 {
			return filter, fmt.Errorf("failed to build filter: when 'WithSourceDevices' parameter is defined also SourceID must be set.")
		}
		addFilter(filter, "withSourceDevices", "true")
	}

	if len(daf.Resolved) > 0 {
		resolved, err := strconv.ParseBool(daf.Resolved)
		if err != nil {
			return filter, fmt.Errorf("failed to build filter: if 'Resolved' parameter is set, only 'true' and 'false' values are accepted.")
		}
		addFilter(filter, "resolved", strconv.FormatBool(resolved))
	}

	if len(daf.Severity) > 0 {
		addFilter(filter, "severity", fmt.Sprintf("%s", daf.Severity))
	}

	if !daf.DateFrom.IsZero() {
		addFilter(filter, "dateFrom", daf.DateFrom.Format(time.RFC3339))
	}

	if !daf.DateTo.IsZero() {
		addFilter(filter, "dateTo", daf.DateTo.Format(time.RFC3339))
	}

	if len(daf.Type) > 0 {
		addFilter(filter, "type", daf.Type)
	}

	return filter, nil
}

func (daf DeleteAlarmsFilter) appendFilter(r *http.Request) error {
	q := r.URL.Query()

	if len(daf.Status) > 0 {
		q.Set("status", fmt.Sprintf("%s", daf.Status))
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

func (gaf GetAlarmsFilter) appendFilter(r *http.Request) error {
	q := r.URL.Query()

	if len(gaf.Status) > 0 {
		q.Set("status", fmt.Sprintf("%s", gaf.Status))
	}

	if len(gaf.SourceId) > 0 {
		q.Set("source", gaf.SourceId)
	}

	if gaf.WithSourceAssets {
		if len(gaf.SourceId) == 0 {
			return fmt.Errorf("failed to build filter: when 'WithSourceAssets' parameter is defined also SourceID must be set.")
		}
		q.Set("withSourceAssets", "true")
	}

	if gaf.WithSourceDevices {
		if len(gaf.SourceId) == 0 {
			return fmt.Errorf("failed to build filter: when 'WithSourceDevices' parameter is defined also SourceID must be set.")
		}
		q.Set("withSourceDevices", "true")
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

	if len(gaf.Type) > 0 {
		q.Set("type", gaf.Type)
	}

	r.URL.RawQuery = q.Encode()
	return nil
}

func addFilter(filter string, filterName string, filterValue string) {
	if len(filter) > 0 {
		filter = filter + "&"
	} else {
		filter = "?"
	}

	filter = filter + filterName + "=" + filterValue
}