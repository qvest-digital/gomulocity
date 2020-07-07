package managedObjects

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"net/http"
	"net/url"
)

type ManagedObjectApi interface {
	CreateManagedObject(object *CreateManagedObject) (*NewManagedObject, *generic.Error)
	ManagedObjectCollection(filter ManagedObjectCollectionFilter) (ManagedObjectCollection, *generic.Error)
	ManagedObjectByID(deviceID string) (ManagedObject, *generic.Error)
	ReferenceByID(deviceID, reference, referenceID string) (Reference, *generic.Error)
	ReferenceCollection(deviceID, reference string) (ReferenceCollection, *generic.Error)
	DeleteReference(deviceID, reference, referenceID string) *generic.Error
	AddReferenceToCollection(deviceID, reference string) (interface{}, *generic.Error)
	UpdateManagedObject(deviceID string, model *Update) (UpdateResponse, *generic.Error)
	DeleteManagedObject(deviceID string) *generic.Error
}

type managedObjectApi struct {
	Client             *generic.Client
	ManagedObjectsPath string
}

func NewManagedObjectApi(client *generic.Client) ManagedObjectApi {
	return managedObjectApi{
		Client:             client,
		ManagedObjectsPath: managedObjectPath,
	}
}

/*
Creates a new managed object based on the given variables.

See: https://cumulocity.com/guides/reference/inventory/#managed-object-collection
*/
func (m managedObjectApi) CreateManagedObject(model *CreateManagedObject) (*NewManagedObject, *generic.Error) {
	bytes, err := json.Marshal(model)
	if err != nil {
		return nil, clientError(fmt.Sprintf("Error while marshalling managed object: %s", err), "CreateManagedObject")
	}
	result, statusCode, err := m.Client.Post(fmt.Sprintf("%v", managedObjectPath), bytes, generic.AcceptAndContentTypeHeader(MANAGED_OBJECT_ACCEPT, MANAGED_OBJECT_CONTENT_TYPE))
	if err != nil {
		return nil, clientError(fmt.Sprintf("Error while creating a new managed object: %s", err), "CreateManagedObject")
	}

	if statusCode != http.StatusCreated {
		return nil, createErrorFromResponse(result)
	}

	newManagedObject := &NewManagedObject{}
	if err = json.Unmarshal(result, newManagedObject); err != nil {
		return nil, clientError(fmt.Sprintf("Error while unmarshalling response: %s", err), "CreateManagedObject")
	}
	return newManagedObject, nil
}

/*
Returns a collection of managed objects.

See: https://cumulocity.com/guides/reference/inventory/#managed-object-collection
*/
func (m managedObjectApi) ManagedObjectCollection(filter ManagedObjectCollectionFilter) (ManagedObjectCollection, *generic.Error) {
	var tempCollection ManagedObjectCollection

	url := fmt.Sprintf("%v?%v", m.ManagedObjectsPath, filter.QueryParams())
	for {
		result, statusCode, err := m.Client.Get(url, generic.EmptyHeader())
		if err != nil {
			return ManagedObjectCollection{}, clientError(fmt.Sprintf("failed to execute rest request: %s", err), "ManagedObjectCollection")
		}
		if statusCode != http.StatusOK {
			return ManagedObjectCollection{}, createErrorFromResponse(result)
		}

		objectCollection := ManagedObjectCollection{}
		if err := json.NewDecoder(bytes.NewReader(result)).Decode(&objectCollection); err != nil {
			return ManagedObjectCollection{}, clientError(fmt.Sprintf("failed to unmarshal response body: %s", err), "ManagedObjectCollection")
		}

		for _, collection := range objectCollection.ManagedObjects {
			tempCollection.ManagedObjects = append(tempCollection.ManagedObjects, collection)
		}

		if objectCollection.hasNextPage() {
			var genericErr *generic.Error
			url, genericErr = objectCollection.NextPage()
			if genericErr != nil {
				return ManagedObjectCollection{}, genericErr
			}
		} else {
			break
		}

	}
	return tempCollection, nil
}

/*
Returns a single managed object.
*/
func (m managedObjectApi) ManagedObjectByID(deviceID string) (ManagedObject, *generic.Error) {
	if len(deviceID) == 0 {
		return ManagedObject{}, clientError("given deviceID is empty", "ManagedObjectByID")
	}

	result, code, err := m.Client.Get(fmt.Sprintf("%v/%v", m.ManagedObjectsPath, url.QueryEscape(deviceID)), generic.EmptyHeader())
	if err != nil {
		return ManagedObject{}, clientError(fmt.Sprintf("error while getting managedObject: %v", err), "ManagedObjectByID")
	}

	if code != http.StatusOK {
		return ManagedObject{}, createErrorFromResponse(result)
	}

	managedObject := ManagedObject{}
	if err := json.Unmarshal(result, &managedObject); err != nil {
		return ManagedObject{}, clientError(fmt.Sprintf("error while unmarshalling managedObject: %v", err), "ManagedObjectByID")
	}
	return managedObject, nil
}

/*
See: https://cumulocity.com/guides/reference/inventory/#managed-object-reference
*/
func (m managedObjectApi) ReferenceByID(deviceID, reference, referenceID string) (Reference, *generic.Error) {
	if len(deviceID) == 0 || len(reference) == 0 || len(referenceID) == 0 {
		return Reference{}, clientError("given deviceID, reference or referenceID is empty", "ReferenceByID")
	}

	result, code, err := m.Client.Get(fmt.Sprintf("%v/%v/%v/%v", m.ManagedObjectsPath, url.QueryEscape(deviceID), url.QueryEscape(reference), url.QueryEscape(referenceID)), generic.EmptyHeader())
	if err != nil {
		return Reference{}, clientError(fmt.Sprintf("error while getting reference: %v with referenceID: %v for device: %v, %s", reference, referenceID, deviceID, err), "ReferenceByID")
	}

	if code != http.StatusOK {
		return Reference{}, createErrorFromResponse(result)
	}

	var referenceModel Reference
	if err := json.Unmarshal(result, &referenceModel); err != nil {
		return Reference{}, clientError(fmt.Sprintf("received an error while unmarshalling response: %s", err), "ReferenceByID")
	}
	return referenceModel, nil
}

/*
Returns a reference collection for a device and reference.

See: https://cumulocity.com/guides/reference/inventory/#managed-object-reference
*/
func (m managedObjectApi) ReferenceCollection(deviceID, reference string) (ReferenceCollection, *generic.Error) {
	if len(deviceID) == 0 || len(reference) == 0 {
		return ReferenceCollection{}, clientError("given deviceID or reference is empty", "ReferenceCollection")
	}

	result, code, err := m.Client.Get(fmt.Sprintf("%v/%v/%v", m.ManagedObjectsPath, url.QueryEscape(deviceID), url.QueryEscape(reference)), generic.EmptyHeader())
	if err != nil {
		return ReferenceCollection{}, clientError(fmt.Sprintf("error while getting reference collections: %s", err), "ReferenceCollection")
	}

	if code != http.StatusOK {
		if code == http.StatusNotFound {
			return ReferenceCollection{}, clientError(fmt.Sprintf("no reference collection found for reference: %v", reference), "ReferenceCollection")
		}
		return ReferenceCollection{}, createErrorFromResponse(result)
	}

	referenceCollection := ReferenceCollection{}
	if err := json.Unmarshal(result, &referenceCollection); err != nil {
		return ReferenceCollection{}, clientError(fmt.Sprintf("received an error while unmarshalling response: %s", err), "ReferenceCollection")
	}
	return referenceCollection, nil
}

func (m managedObjectApi) AddReferenceToCollection(deviceID, reference string) (interface{}, *generic.Error) {
	return nil, nil
}

func (m managedObjectApi) DeleteReference(deviceID, reference, referenceID string) *generic.Error {
	if len(deviceID) == 0 || len(reference) == 0 || len(referenceID) == 0 {
		return clientError("given deviceID, reference or referenceID is empty", "DeleteReference")
	}

	result, code, err := m.Client.Delete(fmt.Sprintf("%v/%v/%v/%v", m.ManagedObjectsPath, url.QueryEscape(deviceID), url.QueryEscape(reference), url.QueryEscape(referenceID)), generic.EmptyHeader())
	if err != nil {
		return clientError(fmt.Sprintf("received an error while deleting reference: %s", err), "DeleteReference")
	}

	if code != http.StatusNoContent {
		return createErrorFromResponse(result)
	}
	return nil
}

func (m managedObjectApi) UpdateManagedObject(deviceID string, model *Update) (UpdateResponse, *generic.Error) {
	if len(deviceID) == 0 {
		return UpdateResponse{}, clientError("given deviceID is empty", "UpdateManagedObject")
	}
	bytes, err := json.Marshal(model)
	if err != nil {
		return UpdateResponse{}, clientError(fmt.Sprintf("received an error while marshalling managedObject: %s", err), "UpdateManagedObject")
	}

	result, code, err := m.Client.Put(fmt.Sprintf("%v/%v", m.ManagedObjectsPath, url.QueryEscape(deviceID)), bytes, generic.ContentTypeHeader(MANAGED_OBJECT_CONTENT_TYPE))
	if err != nil {
		return UpdateResponse{}, clientError(fmt.Sprintf("received an error while updating managedObject: %s", err), "UpdateManagedObject")
	}
	if code != http.StatusOK {
		return UpdateResponse{}, createErrorFromResponse(result)
	}

	updateResponse := UpdateResponse{}
	if err = json.Unmarshal(result, &updateResponse); err != nil {
		return UpdateResponse{}, clientError(fmt.Sprintf("received an error while unmarshalling response: %s", err), "UpdateManagedObject")
	}
	return updateResponse, nil
}

func (m managedObjectApi) DeleteManagedObject(deviceID string) *generic.Error {
	if len(deviceID) == 0 {
		return clientError("given deviceID is empty", "DeleteManagedObject")
	}
	result, code, err := m.Client.Delete(fmt.Sprintf("%v/%v", m.ManagedObjectsPath, url.QueryEscape(deviceID)), generic.EmptyHeader())
	if err != nil {
		return clientError(fmt.Sprintf("received an error while deleting managed object: %s", err), "DeleteManagedObject")
	}
	if code != http.StatusNoContent {
		return createErrorFromResponse(result)
	}
	return nil
}

func (d ManagedObjectCollection) PrintToConsole() {
	for _, managedObject := range d.ManagedObjects {
		fmt.Println(fmt.Sprintf("Device ID: %v Device name: %v", managedObject.ID, managedObject.Name))
	}
	fmt.Printf("Amount of devices: %v", len(d.ManagedObjects))
}

func (d ManagedObjectCollection) NextPage() (string, *generic.Error) {
	return buildURL(d.Next)
}

func (d ManagedObjectCollection) hasNextPage() bool {
	return len(d.Next) > 0
}

func buildURL(next string) (string, *generic.Error) {
	url, err := url.Parse(next)
	if err != nil {
		return "", clientError(fmt.Sprintf("failed to parse url of the next managedObject page: %s", err), "buildURL")
	}
	return fmt.Sprintf("%v?%v", url.Path, url.RawQuery), nil
}

func clientError(message string, info string) *generic.Error {
	return &generic.Error{
		ErrorType: "ClientError",
		Message:   message,
		Info:      info,
	}
}

func createErrorFromResponse(responseBody []byte) *generic.Error {
	var err generic.Error
	_ = json.Unmarshal(responseBody, &err)
	return &err
}
