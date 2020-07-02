package managedObjects

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"net/http"
	"net/url"
)

type ManagedObjectApi struct {
	Client             *generic.Client
	ManagedObjectsPath string
}

func NewManagedObjectApi(client *generic.Client) ManagedObjectApi {
	return ManagedObjectApi{
		Client:             client,
		ManagedObjectsPath: managedObjectPath,
	}
}

func (m *ManagedObjectApi) ManagedObjectCollection(filter map[string][]string) (ManagedObjectCollection, error) {
	var tempCollection ManagedObjectCollection

	url := managedObjectPath

breakOuterLoop:
	for {
		result, statusCode, err := m.Client.Get(url, nil)
		if err != nil {
			return ManagedObjectCollection{}, fmt.Errorf("failed to execute rest request: %w", err)
		}

		if statusCode != http.StatusOK {
			switch statusCode {
			case http.StatusUnauthorized:
				return ManagedObjectCollection{}, generic.BadCredentialsErr
			case http.StatusForbidden:
				return ManagedObjectCollection{}, generic.AccessDeniedErr
			default:
				return ManagedObjectCollection{}, fmt.Errorf("received an unexpected status code: %v", statusCode)
			}
		}

		objectCollection := ManagedObjectCollection{}
		if err := json.NewDecoder(bytes.NewReader(result)).Decode(&objectCollection); err != nil {
			return ManagedObjectCollection{}, fmt.Errorf("failed to unmarshal response body: %w", err)
		}

		for _, object := range objectCollection.ManagedObjects {
			if object.C8YIsDevice != nil {
				tempCollection.ManagedObjects = append(tempCollection.ManagedObjects, object)
			}
		}
		fmt.Println(objectCollection.Next)
		if len(objectCollection.Next) == 0 {
			break breakOuterLoop
		}
		url, err = buildURL(objectCollection.Next)
		if err != nil {
			return ManagedObjectCollection{}, err
		}
	}
	return tempCollection, nil
}

func buildURL(next string) (string, error) {
	url, err := url.Parse(next)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v?%v", url.Path, url.RawQuery), nil
}

func (d ManagedObjectCollection) PrintToConsole() {
	for _, managedObject := range d.ManagedObjects {
		fmt.Println(fmt.Sprintf("Device ID: %v Device name: %v", managedObject.ID, managedObject.Name))
	}
	fmt.Printf("Amount of devices: %v", len(d.ManagedObjects))
}
