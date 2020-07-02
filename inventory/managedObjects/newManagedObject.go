package managedObjects

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"net/http"
)

const (
	AcceptHeader = "application/vnd.com.nsn.cumulocity.managedObject+json"
)

func (m ManagedObjectApi) CreateManagedObject(name, state string) (NewManagedObject, error) {
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
		return NewManagedObject{}, fmt.Errorf("failed to marshal request body: %w", err)
	}
	result, statusCode, err := m.Client.Post(managedObjectPath, body, generic.AcceptHeader(AcceptHeader))
	if err != nil {
		return NewManagedObject{}, fmt.Errorf("an error occurred while processing request: %w", err)
	}

	if statusCode != http.StatusCreated {
		switch statusCode {
		case http.StatusUnauthorized:
			return NewManagedObject{}, generic.BadCredentialsErr
		case http.StatusForbidden:
			return NewManagedObject{}, generic.AccessDeniedErr
		default:
			return NewManagedObject{}, fmt.Errorf("received an unexpected status code: %v", statusCode)
		}
	}

	managedObject := NewManagedObject{}
	if err = json.NewDecoder(bytes.NewReader(result)).Decode(&managedObject); err != nil {
		return NewManagedObject{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	return managedObject, nil
}
