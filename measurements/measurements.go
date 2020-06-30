package measurements

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tarent/gomulocity/generic"
)

const (
	CONTENT_TYPE_MEASUREMENT            = "application/vnd.com.nsn.cumulocity.measurement+json;charset=UTF-8;ver=0.9"
	ACCEPT_MEASUREMENT                  = "application/vnd.com.nsn.cumulocity.measurement+json;charset=UTF-8;ver=0.9"
	CONTENT_TYPE_MEASUREMENT_COLLECTION = "application/vnd.com.nsn.cumulocity.measurementCollection+json;charset=UTF-8;ver=0.9"
	ACCEPT_MEASUREMENT_COLLECTION       = "application/vnd.com.nsn.cumulocity.measurementCollection+json;charset=UTF-8;ver=0.9"

	MEASUREMENTS_API = "measurement/measurements/"
)

type MeasurementsApiClient struct {
	Client generic.Client
}

func NewMeasurementsApiClient (httpclient generic.Client) MeasurementsApiClient{
	return MeasurementsApiClient{
		Client : httpclient
	}
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
	req.Header.Add("Accept", ACCEPT_MEASUREMENT)
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
	req.Header.Add("Accept", ACCEPT_MEASUREMENT)
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

func (c Client) GetMeasurements(resultSize int) (Measurement, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%v%v", c.BaseURL, MEASUREMENTS_API),
		nil,
	)
	if err != nil {
		return Measurement{}, fmt.Errorf("failed to initialize rest request: %w", err)
	}
	req.Header.Add("Accept", ACCEPT_MEASUREMENT)
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
	req.Header.Add("Accept", ACCEPT_MEASUREMENT)
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
