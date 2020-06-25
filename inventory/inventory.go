package inventory

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"github.com/tarent/gomulocity/models"
	"net/http"
)

const (
	manageObjectPath = "/inventory/managedObjects"

	manageObjectContentType = "application/vnd.com.nsn.cumulocity.managedObject+json"
)

//monkey patch to test GetManagedObjects method
var manageObjectPagingStatics = newManagedObjectPagingStatics

type ManagedObjectsPagingStatics struct {
	Statistics generic.PagingStatistics `json:"statistics"`
}

func newManagedObjectPagingStatics(c Client) (ManagedObjectsPagingStatics, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%v/inventory/managedObjects", c.BaseURL),
		nil,
	)
	if err != nil {
		return ManagedObjectsPagingStatics{}, fmt.Errorf("failed to initialize inventoryCollection request: %w", err)
	}

	req.SetBasicAuth(c.Username, c.Password)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return ManagedObjectsPagingStatics{}, fmt.Errorf("failed to execute inventoryCollection request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return ManagedObjectsPagingStatics{}, generic.BadCredentialsErr
		case http.StatusForbidden:
			return ManagedObjectsPagingStatics{}, generic.AccessDeniedErr
		default:
			return ManagedObjectsPagingStatics{}, fmt.Errorf("received an unexpected status code: %v", resp.StatusCode)
		}
	}

	inventoryReqCollection := ManagedObjectsPagingStatics{}
	if err := json.NewDecoder(resp.Body).Decode(&inventoryReqCollection); err != nil {
		return ManagedObjectsPagingStatics{}, fmt.Errorf("failed to unmarshal inventoryRequestCollection: %w", err)
	}
	return inventoryReqCollection, nil
}

func (c Client) GetManagedObjects(fragmentFilter string) ([]models.ManagedObject, error) {
	InventoryReqCollection, err := manageObjectPagingStatics(c)
	if err != nil {
		return []models.ManagedObject{}, err
	}
	var aggregatedManagedObjects []models.ManagedObject

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%v%v", c.BaseURL, manageObjectPath), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize rest request: %w", err)
	}
	if len(fragmentFilter) > 0 {
		query := request.URL.Query()
		query.Add("fragmentType", fragmentFilter)
		request.URL.RawQuery = query.Encode()
	}
	request.SetBasicAuth(c.Username, c.Password)

	for InventoryReqCollection.Statistics.PageSize >= InventoryReqCollection.Statistics.CurrentPage {
		req := generic.Page(InventoryReqCollection.Statistics.CurrentPage)
		req(request)

		resp, err := c.HTTPClient.Do(request)
		if err != nil {
			return nil, fmt.Errorf("failed to execute rest request: %w", err)
		}

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

		tempData := struct {
			inventoryData models.InventoryStructure
		}{}
		if err = json.NewDecoder(resp.Body).Decode(&tempData.inventoryData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
		}
		resp.Body.Close()

		for _, tempObject := range tempData.inventoryData.ManagedObject {
			aggregatedManagedObjects = append(aggregatedManagedObjects, tempObject)
		}
		InventoryReqCollection.Statistics.CurrentPage++
	}

	return aggregatedManagedObjects, nil
}

func (c Client) CreateManagedObject(name, state string) (models.NewManagedObject, error) {
	data := struct {
		Name         string `json:"name"`
		BinarySwitch struct {
			State string `json:"state"`
		} `json:"com_cumulocity_model_BinarySwitch"`
	}{
		Name: name,
		BinarySwitch: struct {
			State string `json:"state"`
		}{
			State: state,
		},
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
