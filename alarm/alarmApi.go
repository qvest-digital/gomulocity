package alarm

import (
	"encoding/json"
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"log"
	"net/http"
	"net/url"
)

type AlarmApi interface {
	// Create a new alarm and returns the created entity with id and creation time
	Create(alarm *NewAlarm) (*Alarm, *generic.Error)

	// Gets an exiting alarm by its id. If the id does not exists, nil is returned.
	Get(alarmId string) (*Alarm, *generic.Error)

	// Updates an exiting alarm and returns the updated alarm entity.
	Update(alarmId string, alarm *UpdateAlarm) (*Alarm, *generic.Error)

	// Updates status of many alarms.
	BulkStatusUpdate(query *UpdateAlarmsFilter, newStatus Status) *generic.Error

	// Deletion by alarm id is not supported/allowed by cumulocity.
	// Deletes alarms by filter. If error is nil, alarms were deleted successfully.
	Delete(query *AlarmFilter) *generic.Error

	// Gets a alarm collection by a source (aka managed object id).
	GetForDevice(sourceId string, pageSize int) (*AlarmCollection, *generic.Error)

	// Returns an alarm collection, found by the given alarm query parameters.
	// All query parameters are AND concatenated.
	Find(query *AlarmFilter, pageSize int) (*AlarmCollection, *generic.Error)

	// Gets the next page from an existing alarm collection.
	// If there is no next page, nil is returned.
	NextPage(c *AlarmCollection) (*AlarmCollection, *generic.Error)

	// Gets the previous page from an existing alarm collection.
	// If there is no previous page, nil is returned.
	PreviousPage(c *AlarmCollection) (*AlarmCollection, *generic.Error)
}

type alarmApi struct {
	client   *generic.Client
	basePath string
}

// Creates a new alarm api object
// client - Must be a gomulocity client.
// returns - The `alarm`-api object
func NewAlarmApi(client *generic.Client) AlarmApi {
	return &alarmApi{client, ALARM_API_PATH}
}

/*
Creates an alarm for an existing device.

Returns created 'Alarm' on success, otherwise an error.
See: https://cumulocity.com/guides/reference/alarms/#post-create-a-new-alarm
*/
func (alarmApi *alarmApi) Create(newAlarm *NewAlarm) (*Alarm, *generic.Error) {
	bytes, err := generic.JsonFromObject(newAlarm)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while marshalling the alarm: %s", err.Error()), "CreateAlarm")
	}
	headers := generic.AcceptAndContentTypeHeader(ALARM_TYPE, ALARM_TYPE)

	body, status, err := alarmApi.client.Post(alarmApi.basePath, bytes, headers)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while posting a new alarm: %s", err.Error()), "CreateAlarm")
	}
	if status != http.StatusCreated {
		return nil, generic.CreateErrorFromResponse(body, status)
	}

	return parseAlarmResponse(body)
}

/*
Gets an alarm for a given Id.

Returns 'Alarm' on success or nil if the id does not exist.

See: https://cumulocity.com/guides/reference/alarms/#get-an-alarm
*/
func (alarmApi *alarmApi) Get(alarmId string) (*Alarm, *generic.Error) {
	body, status, err := alarmApi.client.Get(fmt.Sprintf("%s/%s", alarmApi.basePath, url.QueryEscape(alarmId)), generic.AcceptHeader(ALARM_TYPE))

	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while getting an alarm: %s", err.Error()), "Get")
	}
	if status == http.StatusNotFound {
		return nil, nil
	}
	if status != http.StatusOK {
		return nil, generic.CreateErrorFromResponse(body, status)
	}

	return parseAlarmResponse(body)
}

/*
Updates the alarm with given Id.

See: https://cumulocity.com/guides/reference/alarms/#update-an-alarm
*/
func (alarmApi *alarmApi) Update(alarmId string, alarm *UpdateAlarm) (*Alarm, *generic.Error) {
	bytes, err := json.Marshal(alarm)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while marshalling the update alarm: %s", err.Error()), "UpdateAlarm")
	}

	path := fmt.Sprintf("%s/%s", alarmApi.basePath, url.QueryEscape(alarmId))
	headers := generic.AcceptAndContentTypeHeader(ALARM_TYPE, ALARM_TYPE)

	body, status, err := alarmApi.client.Put(path, bytes, headers)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while updating an alarm: %s", err.Error()), "UpdateAlarm")
	}
	if status != http.StatusOK {
		return nil, generic.CreateErrorFromResponse(body, status)
	}

	return parseAlarmResponse(body)
}

/*
Updates the status of many alarms at once searching by filter.

See: https://cumulocity.com/guides/reference/alarms/#put-bulk-update-of-alarm-collection
*/
func (alarmApi *alarmApi) BulkStatusUpdate(updateAlarmsFilter *UpdateAlarmsFilter, newStatus Status) *generic.Error {
	alarmStatus := UpdateAlarm{Status: newStatus}

	bytes, err := json.Marshal(alarmStatus)
	if err != nil {
		return generic.ClientError(fmt.Sprintf("Error while marshalling the update of alarms: %s", err.Error()), "BulkStatusUpdate")
	}

	queryParamsValues := &url.Values{}
	err = updateAlarmsFilter.QueryParams(queryParamsValues)
	if err != nil {
		return generic.ClientError(fmt.Sprintf("Error while building query parameters for update of alarms: %s", err.Error()), "BulkStatusUpdate")
	}

	path := fmt.Sprintf("%s?%s", alarmApi.basePath, queryParamsValues.Encode())
	headers := generic.AcceptHeader(ALARM_TYPE)

	body, status, err := alarmApi.client.Put(path, bytes, headers)
	if err != nil {
		return generic.ClientError(fmt.Sprintf("Error while updating alarms: %s", err.Error()), "BulkStatusUpdate")
	}

	// Since this operations can take a lot of time, request returns after maximum 0.5 sec of processing,
	// and updating is continued as a background process in the platform.
	// Therefore following possible response statuses can be interpret as successful:
	//	200 - if the process has completed, all alarms have been updated
	//	202 - if process continues in background (maybe )
	if status != http.StatusOK && status != http.StatusAccepted {
		return generic.CreateErrorFromResponse(body, status)
	}

	return nil
}

/*
Deletes alarms by filter.

See: https://cumulocity.com/guides/reference/alarms/#delete-delete-an-alarm-collection
*/
func (alarmApi *alarmApi) Delete(alarmFilter *AlarmFilter) *generic.Error {
	queryParamsValues := &url.Values{}
	err := alarmFilter.QueryParams(queryParamsValues)
	if err != nil {
		return generic.ClientError(fmt.Sprintf("Error while building query parameters for deletion of alarms: %s", err.Error()), "DeleteAlarms")
	}

	body, status, err := alarmApi.client.Delete(fmt.Sprintf("%s?%s", alarmApi.basePath, queryParamsValues.Encode()), generic.EmptyHeader())
	if err != nil {
		return generic.ClientError(fmt.Sprintf("Error while deleting alarms: %s", err.Error()), "DeleteAlarms")
	}

	if status != http.StatusNoContent {
		return generic.CreateErrorFromResponse(body, status)
	}

	return nil
}

func (alarmApi *alarmApi) GetForDevice(sourceId string, pageSize int) (*AlarmCollection, *generic.Error) {
	return alarmApi.Find(&AlarmFilter{SourceId: sourceId}, pageSize)
}

func (alarmApi *alarmApi) Find(alarmFilter *AlarmFilter, pageSize int) (*AlarmCollection, *generic.Error) {
	queryParamsValues := &url.Values{}
	err := alarmFilter.QueryParams(queryParamsValues)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while building query parameters to search for alarms: %s", err.Error()), "FindAlarms")
	}

	err = generic.PageSizeParameter(pageSize, queryParamsValues)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while building pageSize parameter to fetch alarms: %s", err.Error()), "FindAlarms")
	}

	return alarmApi.getCommon(fmt.Sprintf("%s?%s", alarmApi.basePath, queryParamsValues.Encode()))
}

func (alarmApi *alarmApi) NextPage(c *AlarmCollection) (*AlarmCollection, *generic.Error) {
	return alarmApi.getPage(c.Next)
}

func (alarmApi *alarmApi) PreviousPage(c *AlarmCollection) (*AlarmCollection, *generic.Error) {
	return alarmApi.getPage(c.Prev)
}

// -- internal

func parseAlarmResponse(body []byte) (*Alarm, *generic.Error) {
	var result Alarm
	if len(body) > 0 {
		err := json.Unmarshal(body, &result)
		if err != nil {
			return nil, generic.ClientError(fmt.Sprintf("Error while parsing response JSON: %s", err.Error()), "ResponseParser")
		}
	} else {
		return nil, generic.ClientError("Response body was empty", "GetAlarm")
	}

	return &result, nil
}

func (alarmApi *alarmApi) getPage(reference string) (*AlarmCollection, *generic.Error) {
	if reference == "" {
		log.Print("No page reference given. Returning nil.")
		return nil, nil
	}

	nextUrl, err := url.Parse(reference)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Unparsable URL given for page reference: '%s'", reference), "GetPage")
	}

	collection, genErr := alarmApi.getCommon(fmt.Sprintf("%s?%s", nextUrl.Path, nextUrl.RawQuery))
	if genErr != nil {
		return nil, genErr
	}

	if len(collection.Alarms) == 0 {
		log.Print("Returned collection is empty. Returning nil.")
		return nil, nil
	}

	return collection, nil
}

func (alarmApi *alarmApi) getCommon(path string) (*AlarmCollection, *generic.Error) {
	body, status, err := alarmApi.client.Get(path, generic.AcceptHeader(ALARM_COLLECTION_TYPE))

	if status != http.StatusOK {
		return nil, generic.CreateErrorFromResponse(body, status)
	}

	var result AlarmCollection
	if len(body) > 0 {
		err = json.Unmarshal(body, &result)
		if err != nil {
			return nil, generic.ClientError(fmt.Sprintf("Error while parsing response JSON: %s", err.Error()), "GetCollection")
		}
	} else {
		return nil, generic.ClientError("Response body was empty", "GetCollection")
	}

	return &result, nil
}
