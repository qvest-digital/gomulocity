package alarm

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

/*
See: https://cumulocity.com/guides/reference/alarms/#delete-delete-an-alarm-collection
 */
type AlarmFilter struct {
	Status				[]Status	// Comma separated alarm statuses, for example ACTIVE,CLEARED.
									// Please note: when resolved parameter is set then status parameter will be ignored
	SourceId 			string		// Source device id.
	WithSourceAssets	bool		// When set to true also alarms for related source assets will be removed.
									// When this parameter is provided also source must be defined.
	WithSourceDevices	bool		// When set to true also alarms for related source devices will be removed.
									// When this parameter is provided also source must be defined.
	Resolved			string		// When this parameter is provided then status parameter will be ignored.
									// When set to true only resolved alarms will be removed (the one with status CLEARED),
									// false means alarms with status ACTIVE or ACKNOWLEDGED.
	Severity			Severity	// Alarm severity, for example MINOR.
	DateFrom			*time.Time	// Start date or date and time of alarm occurrence.
	DateTo				*time.Time	// End date or date and time of alarm occurrence.
	Type				string		// Alarm type.
}

/*
https://cumulocity.com/guides/reference/alarms/#put-bulk-update-of-alarm-collection
 */
type UpdateAlarmsFilter struct {
	Status				Status
	SourceId 			string
	Resolved			string
	Severity			Severity
	DateFrom			*time.Time
	DateTo				*time.Time
}

func (alarmFilter AlarmFilter) QueryParams() (string, error) {
	params := url.Values{}

	if len(alarmFilter.Status) > 0 {
		var statusesAsString []string
		for _, status := range alarmFilter.Status {
			statusesAsString = append(statusesAsString, string(status))
		}
		params.Add("status", strings.Join(statusesAsString, ","))
	}

	if len(alarmFilter.SourceId) > 0 {
		params.Add("source", alarmFilter.SourceId)
	}

	if alarmFilter.WithSourceAssets {
		if len(alarmFilter.SourceId) == 0 {
			return "", fmt.Errorf("failed to build filter: when 'WithSourceAssets' parameter is defined also SourceID must be set.")
		}
		params.Add("withSourceAssets", "true")
	}

	if alarmFilter.WithSourceDevices {
		if len(alarmFilter.SourceId) == 0 {
			return "", fmt.Errorf("failed to build filter: when 'WithSourceDevices' parameter is defined also SourceID must be set.")
		}
		params.Add("withSourceDevices", "true")
	}

	if len(alarmFilter.Resolved) > 0 {
		resolved, err := strconv.ParseBool(alarmFilter.Resolved)
		if err != nil {
			return "", fmt.Errorf("failed to build filter: if 'Resolved' parameter is set, only 'true' and 'false' values are accepted.")
		}
		params.Add("resolved", strconv.FormatBool(resolved))
	}

	if len(alarmFilter.Severity) > 0 {
		params.Add("severity", fmt.Sprintf("%s", alarmFilter.Severity))
	}

	if alarmFilter.DateFrom != nil {
		params.Add("dateFrom", alarmFilter.DateFrom.Format(time.RFC3339))
	}

	if alarmFilter.DateTo != nil {
		params.Add("dateTo", alarmFilter.DateTo.Format(time.RFC3339))
	}

	if len(alarmFilter.Type) > 0 {
		params.Add("type", alarmFilter.Type)
	}

	return params.Encode(), nil
}

func (updateAlarmsFilter UpdateAlarmsFilter) QueryParams() (string, error) {
	params := url.Values{}

	if len(updateAlarmsFilter.Status) > 0 {
		params.Add("status", fmt.Sprintf("%s", updateAlarmsFilter.Status))
	}

	if len(updateAlarmsFilter.SourceId) > 0 {
		params.Add("source", updateAlarmsFilter.SourceId)
	}

	if len(updateAlarmsFilter.Resolved) > 0 {
		resolved, err := strconv.ParseBool(updateAlarmsFilter.Resolved)
		if err != nil {
			return "", fmt.Errorf("failed to build filter: if 'Resolved' parameter is set, only 'true' and 'false' values are accepted.")
		}
		params.Add("resolved", strconv.FormatBool(resolved))
	}

	if len(updateAlarmsFilter.Severity) > 0 {
		params.Add("severity", fmt.Sprintf("%s", updateAlarmsFilter.Severity))
	}

	if updateAlarmsFilter.DateFrom != nil {
		params.Add("dateFrom", updateAlarmsFilter.DateFrom.Format(time.RFC3339))
	}

	if updateAlarmsFilter.DateTo != nil {
		params.Add("dateTo", updateAlarmsFilter.DateTo.Format(time.RFC3339))
	}

	return params.Encode(), nil
}
