package inventory

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"github.com/tarent/gomulocity/models"
	"net/http"
	"strings"
)

const (
	manageObjectPath = "/inventory/managedObjects"

	manageObjectContentType = "application/vnd.com.nsn.cumulocity.managedObject+json"
)

func (c Client) GetManageObjects(fragmentFilter string) ([]models.ManagedObject, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%v%v", c.BaseURL, manageObjectPath),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize rest request: %w", err)
	}

	req.SetBasicAuth(c.Username, c.Password)

	//TODO: Add pagination
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute rest request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return []models.ManagedObject{}, generic.BadCredentialsErr
		case http.StatusForbidden:
			return []models.ManagedObject{}, generic.AccessDeniedErr
		default:
			return []models.ManagedObject{}, fmt.Errorf("received an unexpected status code: %v", resp.StatusCode)
		}
	}

	temporaryObjectData := struct {
		ManageObjects    []models.ManagedObject   `json:"manageObjects"`
		PagingStatistics generic.PagingStatistics `json:"statistics"`
	}{}

	if err = json.NewDecoder(resp.Body).Decode(&temporaryObjectData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	if len(fragmentFilter) == 0 {
		fragmentsFilter(fragmentFilter, &temporaryObjectData.ManageObjects)
	}
	return temporaryObjectData.ManageObjects, nil
}

func fragmentsFilter(filter string, objects *[]models.ManagedObject) {
	if objects != nil {
		for _, object := range *objects {
			var datapoints []models.Datapoints

			if object.C8yDashboard.Children.ID.Config.Datapoints != nil {
				dps := object.C8yDashboard.Children.ID.Config.Datapoints

				for _, datapoint := range dps {
					if strings.ToLower(datapoint.Fragment) == strings.ToLower(filter) {
						datapoints = append(datapoints, datapoint)
					}
				}
			}
			object.C8yDashboard.Children.ID.Config.Datapoints = datapoints
		}
	}
}

func (c Client) CreateManagedObject(name, binarySwitch string) (models.NewManagedObject, error) {
	data := struct {
		Name         string `json:"name"`
		BinarySwitch string `json:"com_cumulocity_model_BinarySwitch"`
	}{
		Name:         name,
		BinarySwitch: binarySwitch,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return models.NewManagedObject{}, fmt.Errorf("failed to marshal request body: %w", err)
	}
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%v%v", c.BaseURL, manageObjectPath),
		bytes.NewReader(body),
	)
	if err != nil {
		return models.NewManagedObject{}, fmt.Errorf("failed to initialize rest request: %w", err)
	}

	req.Header.Add("Content-Type", manageObjectContentType)
	req.Header.Add("Accept", manageObjectContentType)
	req.SetBasicAuth(c.Username, c.Password)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return models.NewManagedObject{}, fmt.Errorf("an error occurred while processing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return models.NewManagedObject{}, generic.BadCredentialsErr
		case http.StatusForbidden:
			return models.NewManagedObject{}, generic.AccessDeniedErr
		default:
			return models.NewManagedObject{}, fmt.Errorf("received an unexpected status code: %v", resp.StatusCode)
		}
	}

	managedObject := models.NewManagedObject{}
	if err = json.NewDecoder(resp.Body).Decode(&managedObject); err != nil {
		return models.NewManagedObject{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	return managedObject, nil
}
