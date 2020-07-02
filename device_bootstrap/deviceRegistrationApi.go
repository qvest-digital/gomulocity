package device_bootstrap

import (
	"encoding/json"
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"log"
	"net/http"
	"net/url"
)

//var NewDeviceRequestAlreadyExistsErr = errors.New("'newDeviceRequest' with Id already exists")

type DeviceRegistrationApi interface {
	// Creates a new deviceRegistration and returns the created entity with status
	Create(deviceId string) (*DeviceRegistration, *generic.Error)

	// Gets an exiting deviceRegistration by device id. If the id does not exists, nil is returned.
	Get(deviceId string) (*DeviceRegistration, *generic.Error)

	// Updates an exiting deviceRegistration and returns the updated deviceRegistration entity.
	Update(deviceId string, newStatus string) (*DeviceRegistration, *generic.Error)

	// Deletes deviceRegistrations by device id. If error is nil, deviceRegistrations were deleted successfully.
	Delete(deviceId string) *generic.Error

	// Returns page by page all deviceRegistrations.
	GetAll(pageSize int) (*DeviceRegistrationCollection, *generic.Error)

	// Gets the next page from an existing deviceRegistration collection.
	// If there is no next page, nil is returned.
	NextPage(c *DeviceRegistrationCollection) (*DeviceRegistrationCollection, *generic.Error)

	// Gets the previous page from an existing deviceRegistration collection.
	// If there is no previous page, nil is returned.
	PreviousPage(c *DeviceRegistrationCollection) (*DeviceRegistrationCollection, *generic.Error)
}

type deviceRegistrationApi struct {
	client   *generic.Client
	basePath string
}

// Creates a new deviceRegistration api object
// client - Must be a gomulocity client.
// returns - The `deviceRegistration`-api object
func NewDeviceRegistrationApi(client *generic.Client) DeviceRegistrationApi {
	return &deviceRegistrationApi{client, DEVICE_REGISTRATION_API_PATH}
}

/*
Creates a 'DeviceRegistration' with the given id.

Returns created 'DeviceRegistration' on success, otherwise an error.
See: https://cumulocity.com/guides/reference/device-credentials/#post-create-a-new-device-request
*/
func (deviceRegistrationApi *deviceRegistrationApi) Create(deviceId string) (*DeviceRegistration, *generic.Error) {
	bytes, err := json.Marshal(DeviceRegistration{Id: deviceId})
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while marshalling the deviceRegistration: %s", err.Error()), "CreateDeviceRegistration")
	}
	headers := generic.AcceptAndContentTypeHeader(DEVICE_CREDENTIALS_TYPE, DEVICE_CREDENTIALS_TYPE)

	body, status, err := deviceRegistrationApi.client.Post(deviceRegistrationApi.basePath, bytes, headers)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while posting a new deviceRegistration: %s", err.Error()), "CreateDeviceRegistration")
	}
	if status != http.StatusCreated {
		return nil, generic.CreateErrorFromResponse(body)
	}

	return parseDeviceRegistrationResponse(body)
}


/*
Gets a device registration for a given Id.

Returns 'DeviceRegistration' on success or nil if the id does not exist.
See: https://cumulocity.com/guides/reference/device-credentials/#get-returns-a-new-device-request
*/
func (deviceRegistrationApi *deviceRegistrationApi) Get(deviceId string) (*DeviceRegistration, *generic.Error) {
	body, status, err := deviceRegistrationApi.client.Get(fmt.Sprintf("%s/%s", deviceRegistrationApi.basePath, url.QueryEscape(deviceId)), generic.AcceptHeader(DEVICE_REGISTRATION_TYPE))

	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while getting an deviceRegistration: %s", err.Error()), "Get")
	}
	if status != http.StatusOK {
		return nil, nil
	}

	return parseDeviceRegistrationResponse(body)
}

/*
Delivers all 'DeviceRegistration's page by page.

Returns created 'DeviceRegistrationCollection' on success.
See: https://cumulocity.com/guides/reference/device-credentials/#get-returns-all-new-device-requests
*/
func (deviceRegistrationApi *deviceRegistrationApi) GetAll(pageSize int) (*DeviceRegistrationCollection, *generic.Error) {
	pageSizeParams, err := generic.PageSizeParameter(pageSize, nil)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while building pageSize parameter to fetch deviceRegistrations: %s", err.Error()), "GetAllDeviceRegistrations")
	}

	return deviceRegistrationApi.getCommon(fmt.Sprintf("%s?%s", deviceRegistrationApi.basePath, pageSizeParams))
}


/*
Updates status of the deviceRegistration with given Id.

See: https://cumulocity.com/guides/reference/device-credentials/#put-updates-a-new-device-request
*/
func (deviceRegistrationApi *deviceRegistrationApi) Update(deviceId string, newStatus string) (*DeviceRegistration, *generic.Error) {
	bytes, err := json.Marshal(DeviceRegistration{Status: newStatus})
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while marshalling the update deviceRegistration: %s", err.Error()), "UpdateDeviceRegistration")
	}

	path := fmt.Sprintf("%s/%s", deviceRegistrationApi.basePath, url.QueryEscape(deviceId))
	headers := generic.AcceptAndContentTypeHeader(DEVICE_REGISTRATION_TYPE, DEVICE_REGISTRATION_TYPE)

	body, status, err := deviceRegistrationApi.client.Put(path, bytes, headers)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while updating an deviceRegistration: %s", err.Error()), "UpdateDeviceRegistration")
	}
	if status != http.StatusOK {
		return nil, generic.CreateErrorFromResponse(body)
	}

	return parseDeviceRegistrationResponse(body)
}

/*
Deletes deviceRegistration with given Id.

See: https://cumulocity.com/guides/reference/device-credentials/#delete-deletes-a-new-device-request
*/
func (deviceRegistrationApi *deviceRegistrationApi) Delete(deviceId string) *generic.Error {
	path := fmt.Sprintf("%s/%s", deviceRegistrationApi.basePath, url.QueryEscape(deviceId))
	body, status, err := deviceRegistrationApi.client.Delete(path, generic.EmptyHeader())
	if err != nil {
		return generic.ClientError(fmt.Sprintf("Error while deleting deviceRegistrations: %s", err.Error()), "DeleteDeviceRegistrations")
	}

	if status != http.StatusNoContent {
		return generic.CreateErrorFromResponse(body)
	}

	return nil
}

func (deviceRegistrationApi *deviceRegistrationApi) NextPage(c *DeviceRegistrationCollection) (*DeviceRegistrationCollection, *generic.Error) {
	return deviceRegistrationApi.getPage(c.Next)
}

func (deviceRegistrationApi *deviceRegistrationApi) PreviousPage(c *DeviceRegistrationCollection) (*DeviceRegistrationCollection, *generic.Error) {
	return deviceRegistrationApi.getPage(c.Prev)
}

// -- internal

func parseDeviceRegistrationResponse(body []byte) (*DeviceRegistration, *generic.Error) {
	var result DeviceRegistration
	if len(body) > 0 {
		err := json.Unmarshal(body, &result)
		if err != nil {
			return nil, generic.ClientError(fmt.Sprintf("Error while parsing response JSON: %s", err.Error()), "ResponseParser")
		}
	} else {
		return nil, generic.ClientError("Response body was empty", "GetDeviceRegistration")
	}

	return &result, nil
}

func (deviceRegistrationApi *deviceRegistrationApi) getPage(reference string) (*DeviceRegistrationCollection, *generic.Error) {
	if reference == "" {
		log.Print("No page reference given. Returning nil.")
		return nil, nil
	}

	nextUrl, err := url.Parse(reference)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Unparsable URL given for page reference: '%s'", reference), "GetPage")
	}

	collection, err2 := deviceRegistrationApi.getCommon(fmt.Sprintf("%s?%s", nextUrl.Path, nextUrl.RawQuery))
	if err2 != nil {
		return nil, err2
	}

	if len(collection.DeviceRegistrations) == 0 {
		log.Print("Returned collection is empty. Returning nil.")
		return nil, nil
	}

	return collection, nil
}

func (deviceRegistrationApi *deviceRegistrationApi) getCommon(path string) (*DeviceRegistrationCollection, *generic.Error) {
	body, status, err := deviceRegistrationApi.client.Get(path, generic.AcceptHeader(DEVICE_REGISTRATION_COLLECTION_TYPE))

	if status != http.StatusOK {
		return nil, generic.CreateErrorFromResponse(body)
	}

	var result DeviceRegistrationCollection
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
