package alarm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"net/http"
	"net/url"
	"strings"
)

//var SourceIdNotExistsErr = errors.New("'source' with ID already exists")

const (
	alarmType           = "application/vnd.com.nsn.cumulocity.alarm+json"
	alarmApiPath        = "/alarm/alarms"
	alarmCollectionType = "application/vnd.com.nsn.cumulocity.alarmCollection+json"
)

type Status string
const (
	ACTIVE Status = "ACTIVE"
	ACKNOWLEDGED Status = "ACKNOWLEDGED"
	CLEARED Status = "CLEARED"
)

type Severity string
const (
	CRITICAL Severity = "CRITICAL"
	MAJOR Severity = "MAJOR"
	MINOR Severity = "MINOR"
	WARNING Severity = "WARNING"
)

/*
Represents cumulocity's alarm 'application/vnd.com.nsn.cumulocity.alarm+json'.
See: https://cumulocity.com/guides/reference/alarms/#alarm
*/
type NewAlarm struct {
	Type    	string `json:"type"`
	Time 		string `json:"time"`
	Text   		string `json:"text"`
	Source		struct {
		ID 	string `json:"id"`
	} `json:"source"`
	Status		Status `json:"status"`
	Severity	Severity `json:"severity"`
	// TODO: object - 0..n additional properties of the alarm.
	//Other map[string]interface{}
}

type Alarm struct {
	ID     				string `json:"id,omitempty"`
	Self   				string `json:"self,omitempty"`
	CreationTime		string `json:"creationTime,omitempty"`

	Type    			string `json:"type,omitempty"`
	Time 				string `json:"time,omitempty"`
	Text   				string `json:"text,omitempty"`
	Source				struct {
		ID 	string `json:"id,omitempty"`
	} `json:"source,omitempty"`
	Status				Status `json:"status,omitempty"`
	Severity			Severity `json:"severity,omitempty"`

	Count  				int `json:"count,omitempty"`
	FirstOccurrenceTime	string `json:"firstOccurrenceTime,omitempty"`

	// TODO: object - 0..n additional properties of the alarm.
}

type AlarmUpdate struct {
	Text   		string `json:"text"`
	Status		Status `json:"status"`
	Severity	Severity `json:"severity"`
	// TODO: object - 0..n additional properties of the alarm.
}


/*
/*
AlarmCollection represent cumulocity's 'application/vnd.com.nsn.cumulocity.alarmCollection+json'.
See: https://cumulocity.com/guides/reference/alarms/#alarm-collection
*/
type AlarmCollection struct {
	Self              string                   `json:"self"`
	Alarms 			  []Alarm       		   `json:"alarms"`
	Statistics        generic.PagingStatistics `json:"statistics,omitempty"`
	Prev              string                   `json:"prev,omitempty"`
	Next              string                   `json:"next,omitempty"`
}

/*
Creates an alarm for an existing device.

Return created 'Alarm' on success.
Can return the following errors:
- generic.BadCredentialsErr (invalid username / password / host combination)
- generic.AccessDeniedErr (missing user rights)
- generic.Error (generic cloud error)
- error (unexpected)

See: https://cumulocity.com/guides/reference/alarms/#post-create-a-new-alarm
*/
func (c Client) CreateAlarm(newAlarm NewAlarm) (Alarm, error) {

	body, err := json.Marshal(newAlarm)
	if err != nil {
		return Alarm{}, fmt.Errorf("failed to marshal request body: %w", err)
	}
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s%s", c.BaseURL, alarmApiPath),
		bytes.NewReader(body),
	)
	if err != nil {
		return Alarm{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(c.Username, c.Password)

	h := req.Header
	h.Add("Content-Type", alarmType)
	h.Add("Accept", alarmType)
	req.Header = h

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return Alarm{}, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return Alarm{}, generic.BadCredentialsErr
		case http.StatusForbidden:
			return Alarm{}, generic.AccessDeniedErr
		case http.StatusUnprocessableEntity:
			var errResp generic.Error
			err := json.NewDecoder(resp.Body).Decode(&errResp)
			if err != nil {
				return Alarm{}, fmt.Errorf("failed to decode error response body: %s", err)
			}
			return Alarm{}, errResp
		default:
			if strings.HasPrefix(resp.Header.Get("Content-Type"), generic.ErrorContentType) {
				var errResp generic.Error
				err := json.NewDecoder(resp.Body).Decode(&errResp)
				if err != nil {
					return Alarm{}, fmt.Errorf("failed to decode error response body: %s", err)
				}
				return Alarm{}, fmt.Errorf("failed to create new-device-reuest (%d): %w", resp.StatusCode, errResp)
			}
			return Alarm{}, fmt.Errorf("failed to create alarm with status code %d", resp.StatusCode)
		}
	}

	var respBody Alarm
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return Alarm{}, fmt.Errorf("failed to decode response body: %w", err)
	}

	return respBody, nil
}

/*
Get alarms.
Note: for paging use generic.Page(int) as reqOpts.

Return created 'AlarmCollection' on success.
Can return the following errors:
- generic.BadCredentialsErr (invalid username / password / host combination)
- generic.AccessDeniedErr (missing user rights)
- generic.Error (generic cloud error)
- error (unexpected)

See: https://cumulocity.com/guides/reference/alarms/#alarm-collection
*/
func (c Client) GetAlarms(alarmsFilter AlarmsFilter, reqOpts ...func(*http.Request)) (AlarmCollection, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s%s", c.BaseURL, alarmApiPath),
		nil,
	)
	if err != nil {
		return AlarmCollection{}, fmt.Errorf("failed to create request: %w", err)
	}

	for _, opt := range reqOpts {
		if opt != nil {
			opt(req)
		}
	}

	alarmsFilter.appendFilter(req)
	fmt.Printf("send GET request: %s\n", req.URL.String())

	req.SetBasicAuth(c.Username, c.Password)

	h := req.Header
	h.Add("Accept", alarmCollectionType)
	req.Header = h

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return AlarmCollection{}, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return AlarmCollection{}, generic.BadCredentialsErr
		case http.StatusForbidden:
			return AlarmCollection{}, generic.AccessDeniedErr
		default:
			if strings.HasPrefix(resp.Header.Get("Content-Type"), generic.ErrorContentType) {
				var errResp generic.Error
				err := json.NewDecoder(resp.Body).Decode(&errResp)
				if err != nil {
					return AlarmCollection{}, fmt.Errorf("failed to decode error response body: %s", err)
				}
				return AlarmCollection{}, fmt.Errorf("failed to find-all alarms (%d): %w", resp.StatusCode, errResp)
			}
			return AlarmCollection{}, fmt.Errorf("failed to find-all alarms with status code %d", resp.StatusCode)
		}
	}

	var respBody AlarmCollection
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return AlarmCollection{}, fmt.Errorf("failed to decode response body: %w", err)
	}
	return respBody, nil
}

/*
Get alarm by ID.

Return alarm on success.
Can return the following errors:
- generic.BadCredentialsErr (invalid username / password / host combination)
- generic.AccessDeniedErr (missing user rights)
- generic.Error (generic cloud error)
- error (unexpected)

See: https://cumulocity.com/guides/reference/alarms/#alarm-collection
*/
func (c Client) GetAlarm(id string) (Alarm, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s%s/%s", c.BaseURL, alarmApiPath, url.QueryEscape(id)),
		nil,
	)
	if err != nil {
		return Alarm{}, fmt.Errorf("failed to create request: %w", err)
	}

	fmt.Printf("send GET request: %s\n", req.URL.String())

	req.SetBasicAuth(c.Username, c.Password)

	h := req.Header
	h.Add("Accept", alarmType)
	req.Header = h

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return Alarm{}, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return Alarm{}, generic.BadCredentialsErr
		case http.StatusForbidden:
			return Alarm{}, generic.AccessDeniedErr
		default:
			if strings.HasPrefix(resp.Header.Get("Content-Type"), generic.ErrorContentType) {
				var errResp generic.Error
				err := json.NewDecoder(resp.Body).Decode(&errResp)
				if err != nil {
					return Alarm{}, fmt.Errorf("failed to decode error response body: %s", err)
				}
				return Alarm{}, fmt.Errorf("failed to find-all alarms (%d): %w", resp.StatusCode, errResp)
			}
			return Alarm{}, fmt.Errorf("failed to find-all alarms with status code %d", resp.StatusCode)
		}
	}

	var respBody Alarm
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return Alarm{}, fmt.Errorf("failed to decode response body: %w", err)
	}
	return respBody, nil
}

/*
Updates the alarm with given ID.

Can return the following errors:
- generic.BadCredentialsErr (invalid username / password / host combination)
- generic.AccessDeniedErr (missing user rights)
- generic.Error (generic cloud error)
- error (unexpected)

See: https://cumulocity.com/guides/reference/alarms/#update-an-alarm
*/
func (c Client) UpdateAlarm(alarm AlarmUpdate, id string) (Alarm, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(alarm)
	if err != nil {
		return Alarm{}, fmt.Errorf("failed to encode request body: %w", err)
	}

	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s%s/%s", c.BaseURL, alarmApiPath, url.QueryEscape(id)),
		&buf,
	)
	if err != nil {
		return Alarm{}, fmt.Errorf("failed to create request: %w", err)
	}

	//fmt.Printf("Request: %w", buf.String())

	req.SetBasicAuth(c.Username, c.Password)

	h := req.Header
	h.Add("Content-Type", alarmType)
	h.Add("Accept", alarmType)
	req.Header = h

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return Alarm{}, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return Alarm{}, generic.BadCredentialsErr
		case http.StatusForbidden:
			return Alarm{}, generic.AccessDeniedErr
		default:
			if strings.HasPrefix(resp.Header.Get("Content-Type"), generic.ErrorContentType) {
				var errResp generic.Error
				err := json.NewDecoder(resp.Body).Decode(&errResp)
				if err != nil {
					return Alarm{}, fmt.Errorf("failed to decode error response body: %s", err)
				}
				//fmt.Printf("failed to update alarm (%d): %w", resp.StatusCode, errResp)
				return Alarm{}, fmt.Errorf("failed to update alarm (%d): %w", resp.StatusCode, errResp)
			}
			return Alarm{}, fmt.Errorf("failed to update alarm with status code %d", resp.StatusCode)
		}
	}

	var respBody Alarm
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return Alarm{}, fmt.Errorf("failed to decode response body: %w", err)
	}
	return respBody, nil
}

/*
Updates status of alarms by filter.

Can return the following errors:
- generic.BadCredentialsErr (invalid username / password / host combination)
- generic.AccessDeniedErr (missing user rights)
- generic.Error (generic cloud error)
- error (unexpected)

See: https://cumulocity.com/guides/reference/alarms/#put-bulk-update-of-alarm-collection
*/
func (c Client) UpdateAlarms(updateAlarmsFilter UpdateAlarmsFilter, newStatus Status) error {
	alarmStatus := struct {
		Status	Status `json:"status"`
	}{
		Status: newStatus,
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(alarmStatus)
	if err != nil {
		return fmt.Errorf("failed to encode request body: %w", err)
	}

	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s%s", c.BaseURL, alarmApiPath),
		&buf,
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	updateAlarmsFilter.appendFilter(req)
	fmt.Printf("Request: %s\n", buf.String())
	fmt.Printf("send UPDATE request: %s\n", req.URL.String())

	req.SetBasicAuth(c.Username, c.Password)

	h := req.Header
	h.Add("Content-Type", alarmType)
	h.Add("Accept", alarmType)
	req.Header = h

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return generic.BadCredentialsErr
		case http.StatusForbidden:
			return generic.AccessDeniedErr
		default:
			if strings.HasPrefix(resp.Header.Get("Content-Type"), generic.ErrorContentType) {
				var errResp generic.Error
				err := json.NewDecoder(resp.Body).Decode(&errResp)
				if err != nil {
					return fmt.Errorf("failed to decode error response body: %s", err)
				}
				//fmt.Printf("failed to update alarm (%d): %w", resp.StatusCode, errResp)
				return fmt.Errorf("failed to update alarms (%d): %w", resp.StatusCode, errResp)
			}
			return fmt.Errorf("failed to update alarms with status code %d", resp.StatusCode)
		}
	}

	return nil
}

/*
Deletes alarms by filter.

Can return the following errors:
- generic.BadCredentialsErr (invalid username / password / host combination)
- generic.AccessDeniedErr (missing user rights)
- generic.Error (generic cloud error)
- error (unexpected)

See: https://cumulocity.com/guides/reference/alarms/#delete-delete-an-alarm-collection
*/
func (c Client) DeleteAlarms(alarmsFilter AlarmsFilter) error {
	req, err := http.NewRequest(
		http.MethodDelete,
		//http.MethodDelete,
		fmt.Sprintf("%s%s", c.BaseURL, alarmApiPath),
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	alarmsFilter.appendFilter(req)

	fmt.Printf("send DELETE request: %s\n", req.URL.String())

	req.SetBasicAuth(c.Username, c.Password)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return generic.BadCredentialsErr
		case http.StatusForbidden:
			return generic.AccessDeniedErr
		default:
			if strings.HasPrefix(resp.Header.Get("Content-Type"), generic.ErrorContentType) {
				var errResp generic.Error
				err := json.NewDecoder(resp.Body).Decode(&errResp)
				if err != nil {
					return fmt.Errorf("failed to decode error response body: %s", err)
				}
				return fmt.Errorf("failed to delete alarm (%d): %w", resp.StatusCode, errResp)
			}
			return fmt.Errorf("failed to delete alarm with status code %d", resp.StatusCode)
		}
	}

	return nil
}
