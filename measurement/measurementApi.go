package measurement

import (
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"log"
	"net/http"
	"net/url"
)

const (
	MEASUREMENTS_API = "/measurement/measurements"

	MEASUREMENT_TYPE            = "application/vnd.com.nsn.cumulocity.measurement+json;charset=UTF-8;ver=0.9"
	MEASUREMENT_COLLECTION_TYPE = "application/vnd.com.nsn.cumulocity.measurementCollection+json;charset=UTF-8;ver=0.9"
)

type MeasurementApi interface {
	// Create a new measurement and returns the created entity with id and creation time
	Create(measurement *NewMeasurement) (*Measurement, *generic.Error)

	CreateMany(measurement *NewMeasurements) (*MeasurementCollection, *generic.Error)

	// Gets an exiting measurement by its id. If the id does not exists, nil is returned.
	Get(measurementId string) (*Measurement, *generic.Error)

	// Deletion by measurement id. If error is nil, measurement was deleted successfully.
	Delete(measurementId string) *generic.Error

	// Deletes measurements by filter. If error is nil, measurements were deleted successfully.
	// ATTENTION: at least one filter should be set otherwise an error will be thrown.
	// Use DeleteAll() (with caution!) instead if you want delete all measurements!
	DeleteMany(measurementQuery *MeasurementQuery) *generic.Error

	// Deletes all measurements. If error is nil, measurements were deleted successfully.
	// ATTENTION: use it with caution!
	DeleteAll() *generic.Error

	// Gets a measurement collection by a source (aka managed object id).
	GetForDevice(sourceId string, pageSize int) (*MeasurementCollection, *generic.Error)

	// Returns an measurement collection, found by the given measurement query parameters.
	// All query parameters are AND concatenated.
	Find(measurementQuery *MeasurementQuery, pageSize int) (*MeasurementCollection, *generic.Error)

	// Gets the next page from an existing measurement collection.
	// If there is no next page, nil is returned.
	NextPage(c *MeasurementCollection) (*MeasurementCollection, *generic.Error)

	// Gets the previous page from an existing measurement collection.
	// If there is no previous page, nil is returned.
	PreviousPage(c *MeasurementCollection) (*MeasurementCollection, *generic.Error)
}

type measurementApi struct {
	client   *generic.Client
	basePath string
}

// Creates a new measurement api object
// client - Must be a gomulocity client.
// returns - The `measurement`-api object
func NewMeasurementApi(client *generic.Client) MeasurementApi {
	return &measurementApi{client, MEASUREMENTS_API}
}

/*
Creates a measurement for an existing device.

Returns created 'Measurement' on success, otherwise an error.
*/
func (measurementApi *measurementApi) Create(measurement *NewMeasurement) (*Measurement, *generic.Error) {
	bytes, err := generic.JsonFromObject(measurement)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while marhalling the measurement: %s", err.Error()), "CreateMeasurement")
	}
	headers := generic.AcceptAndContentTypeHeader(MEASUREMENT_TYPE, MEASUREMENT_TYPE)

	body, status, err := measurementApi.client.Post(measurementApi.basePath, bytes, headers)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while posting a new measurement: %s", err.Error()), "CreateMeasurement")
	}
	if status != http.StatusCreated {
		return nil, generic.CreateErrorFromResponse(body, status)
	}

	return parseMeasurementResponse(body)
}

/*
Creates many measurements at once for an existing device.

Returns a 'Measurement' collection on success, otherwise an error.
*/
func (measurementApi *measurementApi) CreateMany(measurements *NewMeasurements) (*MeasurementCollection, *generic.Error) {
	bytes, err := generic.JsonFromObject(measurements)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while marhalling the measurements: %s", err.Error()), "CreateManyMeasurement")
	}
	headers := generic.AcceptAndContentTypeHeader(MEASUREMENT_COLLECTION_TYPE, MEASUREMENT_COLLECTION_TYPE)

	body, status, err := measurementApi.client.Post(measurementApi.basePath, bytes, headers)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while posting new measurements: %s", err.Error()), "CreateManyMeasurement")
	}
	if status != http.StatusCreated {
		return nil, generic.CreateErrorFromResponse(body, status)
	}

	return parseMeasurementCollectionResponse(body)
}

/*
Gets a measurement for a given Id.

Returns 'Measurement' on success or nil if the id does not exist.
*/
func (measurementApi *measurementApi) Get(measurementId string) (*Measurement, *generic.Error) {
	if len(measurementId) == 0 {
		return nil, generic.ClientError("Getting measurement without an id is not allowed", "GetMeasurement")
	}

	path := fmt.Sprintf("%s/%s", measurementApi.basePath, url.QueryEscape(measurementId))
	body, status, err := measurementApi.client.Get(path, generic.AcceptHeader(MEASUREMENT_TYPE))

	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while getting a measurement: %s", err.Error()), "GetMeasurement")
	}
	if status == http.StatusNotFound {
		return nil, nil
	}
	if status != http.StatusOK {
		return nil, generic.CreateErrorFromResponse(body, status)
	}

	return parseMeasurementResponse(body)
}

/*
Deletes measurement by id.
*/
func (measurementApi *measurementApi) Delete(measurementId string) *generic.Error {
	if len(measurementId) == 0 {
		return generic.ClientError("Deleting measurement without an id will lead into deletion of all measurements "+
			"which is not allowed by this function. Therefore use `DeleteAll()` instead.", "DeleteMeasurement")
	}

	body, status, err := measurementApi.client.Delete(fmt.Sprintf("%s/%s", measurementApi.basePath, url.QueryEscape(measurementId)), generic.EmptyHeader())
	if err != nil {
		return generic.ClientError(fmt.Sprintf("Error while deleting measurement with id [%s]: %s", measurementId, err.Error()), "DeleteMeasurement")
	}

	if status != http.StatusNoContent {
		return generic.CreateErrorFromResponse(body, status)
	}

	return nil
}

/*
Deletes measurements by filter.
*/
func (measurementApi *measurementApi) DeleteMany(measurementQuery *MeasurementQuery) *generic.Error {
	if measurementQuery == nil {
		return generic.ClientError("No filter set. At least one filter has to be set to avoid accident deletion of all measurements. Use `DeleteAll()` if you really want to remove them all", "DeleteManyMeasurements")
	}
	queryParamsValues := &url.Values{}
	err := measurementQuery.QueryParams(queryParamsValues)
	if err != nil {
		return generic.ClientError(fmt.Sprintf("Error while building query parameters for deletion of measurements: %s", err.Error()), "DeleteManyMeasurements")
	}
	if len(*queryParamsValues) == 0 {
		return generic.ClientError("No filter set. At least one filter has to be set to avoid accident deletion of all measurements. Use `DeleteAll()` if you really want to remove them all", "DeleteManyMeasurements")
	}

	body, status, err := measurementApi.client.Delete(fmt.Sprintf("%s?%s", measurementApi.basePath, queryParamsValues.Encode()), generic.EmptyHeader())
	if err != nil {
		return generic.ClientError(fmt.Sprintf("Error while deleting measurements: %s", err.Error()), "DeleteManyMeasurements")
	}

	if status != http.StatusNoContent {
		return generic.CreateErrorFromResponse(body, status)
	}

	return nil
}

/*
ATTENTION: This function deletes all measurements
*/
func (measurementApi *measurementApi) DeleteAll() *generic.Error {
	body, status, err := measurementApi.client.Delete(fmt.Sprintf("%s", measurementApi.basePath), generic.EmptyHeader())
	if err != nil {
		return generic.ClientError(fmt.Sprintf("Error while deleting measurements: %s", err.Error()), "DeleteAllMeasurements")
	}

	if status != http.StatusNoContent {
		return generic.CreateErrorFromResponse(body, status)
	}
	log.Println("WARNING: all measurements of the tenant were deleted!")

	return nil
}

func (measurementApi *measurementApi) GetForDevice(sourceId string, pageSize int) (*MeasurementCollection, *generic.Error) {
	return measurementApi.Find(&MeasurementQuery{SourceId: sourceId}, pageSize)
}

func (measurementApi *measurementApi) Find(measurementQuery *MeasurementQuery, pageSize int) (*MeasurementCollection, *generic.Error) {
	queryParamsValues := &url.Values{}
	err := measurementQuery.QueryParams(queryParamsValues)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while building query parameters to search for measurements: %s", err.Error()), "FindMeasurements")
	}

	err = generic.PageSizeParameter(pageSize, queryParamsValues)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while building pageSize parameter to fetch measurements: %s", err.Error()), "FindMeasurements")
	}

	return measurementApi.getCommon(fmt.Sprintf("%s?%s", measurementApi.basePath, queryParamsValues.Encode()))
}

func (measurementApi *measurementApi) NextPage(c *MeasurementCollection) (*MeasurementCollection, *generic.Error) {
	return measurementApi.getPage(c.Next)
}

func (measurementApi *measurementApi) PreviousPage(c *MeasurementCollection) (*MeasurementCollection, *generic.Error) {
	return measurementApi.getPage(c.Prev)
}

// -- internal

func (measurementApi *measurementApi) getPage(reference string) (*MeasurementCollection, *generic.Error) {
	if reference == "" {
		log.Print("No page reference given. Returning nil.")
		return nil, nil
	}

	nextUrl, err := url.Parse(reference)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Unparsable URL given for page reference: '%s'", reference), "GetPage")
	}

	collection, genErr := measurementApi.getCommon(fmt.Sprintf("%s?%s", nextUrl.Path, nextUrl.RawQuery))
	if genErr != nil {
		return nil, genErr
	}

	if len(collection.Measurements) == 0 {
		log.Print("Returned collection is empty. Returning nil.")
		return nil, nil
	}

	return collection, nil
}

func (measurementApi *measurementApi) getCommon(path string) (*MeasurementCollection, *generic.Error) {
	body, status, err := measurementApi.client.Get(path, generic.AcceptHeader(MEASUREMENT_COLLECTION_TYPE))
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while getting measurements: %s", err.Error()), "GetMeasurementCollection")
	}

	if status != http.StatusOK {
		return nil, generic.CreateErrorFromResponse(body, status)
	}

	return parseMeasurementCollectionResponse(body)
}

func parseMeasurementResponse(body []byte) (*Measurement, *generic.Error) {
	var result Measurement
	if len(body) > 0 {
		err := generic.ObjectFromJson(body, &result)
		if err != nil {
			return nil, generic.ClientError(fmt.Sprintf("Error while parsing response JSON: %s", err.Error()), "ResponseParser")
		}
	} else {
		return nil, generic.ClientError("Response body was empty", "GetMeasurement")
	}

	return &result, nil
}

func parseMeasurementCollectionResponse(body []byte) (*MeasurementCollection, *generic.Error) {
	var result MeasurementCollection
	if len(body) > 0 {
		err := generic.ObjectFromJson(body, &result)
		if err != nil {
			return nil, generic.ClientError(fmt.Sprintf("Error while parsing response JSON: %s", err.Error()), "CollectionResponseParser")
		}
	} else {
		return nil, generic.ClientError("Response body was empty", "CollectionResponseParser")
	}

	return &result, nil
}
