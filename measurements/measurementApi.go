package measurements

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/tarent/gomulocity/models"

	"github.com/tarent/gomulocity/generic"
)

const (
	MEASUREMENTS_API = "measurement/measurements/"

	MEASUREMENT_TYPE            = "application/vnd.com.nsn.cumulocity.measurement+json;charset=UTF-8;ver=0.9"
	MEASUREMENT_COLLECTION_TYPE = "application/vnd.com.nsn.cumulocity.measurementCollection+json;charset=UTF-8;ver=0.9"
)

type MeasurementApi interface {
	// Create a new measurement and returns the created entity with id and creation time
	Create(measurement *Measurement) (*Measurement, *generic.Error)

	CreateMany(measurement *MeasurementCollection) (*MeasurementCollection, *generic.Error)

	// Gets an exiting measurement by its id. If the id does not exists, nil is returned.
	Get(measurementId string) (*Measurement, *generic.Error)

	// Deletion by measurement id. If error is nil, measurement was deleted successfully.
	Delete(measurementId string) *generic.Error

	// Deletes measurements by filter. If error is nil, measurements were deleted successfully.
	DeleteMany(measurementQuery *MeasurementQuery) *generic.Error

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

func (measurementApi *measurementApi) Create(measurement *Measurement) (*Measurement, *generic.Error) {
	bytes, err := json.Marshal(measurement)
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

func (measurementApi *measurementApi) CreateMany(measurement *MeasurementCollection) (*MeasurementCollection, *generic.Error) {
	bytes, err := json.Marshal(measurement)
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



func (c Client) createMeasurement(measurement Measurement) (Measurement, error) {
	body, err := json.Marshal(measurement)
	if err != nil {
		return Measurement{}, fmt.Errorf("failed to marshal request body: %w", err)
	}
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%v%v", c.BaseURL, MEASUREMENTS_API),
		bytes.NewReader(body),
	)
	if err != nil {
		return Measurement{}, fmt.Errorf("failed to initialize rest request: %w", err)
	}
	req.Header.Add("Accept", MEASUREMENT_TYPE)
	req.Header.Add("Content-Type", CONTENT_TYPE_MEASUREMENT)
	req.SetBasicAuth(c.Username, c.Password)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return Measurement{}, fmt.Errorf("failed to execute rest request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return Measurement{}, generic.BadCredentialsErr
		case http.StatusForbidden:
			return Measurement{}, generic.AccessDeniedErr
		default:
			return Measurement{}, fmt.Errorf("received an unexpected status code: %v", resp.StatusCode)
		}
	}

	measurementFromAPI := Measurement{}

	if err = json.NewDecoder(resp.Body).Decode(&measurementFromAPI); err != nil {
		return Measurement{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return measurementFromAPI, nil
}

func (c Client) getMeasurement(id string) (Measurement, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%v%v", c.BaseURL, MEASUREMENTS_API),
		nil,
	)
	if err != nil {
		return Measurement{}, fmt.Errorf("failed to initialize rest request: %w", err)
	}
	req.Header.Add("Accept", MEASUREMENT_TYPE)
	req.SetBasicAuth(c.Username, c.Password)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return Measurement{}, fmt.Errorf("failed to execute rest request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return Measurement{}, generic.BadCredentialsErr
		case http.StatusForbidden:
			return Measurement{}, generic.AccessDeniedErr
		default:
			return Measurement{}, fmt.Errorf("received an unexpected status code: %v", resp.StatusCode)
		}
	}

	measurementFromAPI := Measurement{}

	if err = json.NewDecoder(resp.Body).Decode(&measurementFromAPI); err != nil {
		return Measurement{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return measurementFromAPI, nil
}

func (c Client) GetMeasurements(resultSize int, query MeasurementQuery) (Measurement, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%v%v%v", c.BaseURL, MEASUREMENTS_API, query.QueryParams()),
		nil,
	)
	if err != nil {
		return Measurement{}, fmt.Errorf("failed to initialize rest request: %w", err)
	}
	req.Header.Add("Accept", MEASUREMENT_COLLECTION_TYPE)
	req.SetBasicAuth(c.Username, c.Password)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return Measurement{}, fmt.Errorf("failed to execute rest request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return Measurement{}, generic.BadCredentialsErr
		case http.StatusForbidden:
			return Measurement{}, generic.AccessDeniedErr
		default:
			return Measurement{}, fmt.Errorf("received an unexpected status code: %v", resp.StatusCode)
		}
	}

	measurementFromAPI := Measurement{}

	if err = json.NewDecoder(resp.Body).Decode(&measurementFromAPI); err != nil {
		return Measurement{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return measurementFromAPI, nil
}
func (c Client) deleteMeasurement(id string) (Measurement, error) {
	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%v%v", c.BaseURL, MEASUREMENTS_API),
		nil,
	)
	if err != nil {
		return Measurement{}, fmt.Errorf("failed to initialize rest request: %w", err)
	}
	req.Header.Add("Accept", MEASUREMENT_TYPE)
	req.SetBasicAuth(c.Username, c.Password)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return Measurement{}, fmt.Errorf("failed to execute rest request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return Measurement{}, generic.BadCredentialsErr
		case http.StatusForbidden:
			return Measurement{}, generic.AccessDeniedErr
		default:
			return Measurement{}, fmt.Errorf("received an unexpected status code: %v", resp.StatusCode)
		}
	}

	measurementFromAPI := Measurement{}

	if err = json.NewDecoder(resp.Body).Decode(&measurementFromAPI); err != nil {
		return Measurement{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return measurementFromAPI, nil
}

func parseMeasurementResponse(body []byte) (*Measurement, *generic.Error) {
	var result Measurement
	if len(body) > 0 {
		err := json.Unmarshal(body, &result)
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
		err := json.Unmarshal(body, &result)
		if err != nil {
			return nil, generic.ClientError(fmt.Sprintf("Error while parsing response JSON: %s", err.Error()), "CollectionResponseParser")
		}
	} else {
		return nil, generic.ClientError("Response body was empty", "CollectionResponseParser")
	}

	return &result, nil
}
