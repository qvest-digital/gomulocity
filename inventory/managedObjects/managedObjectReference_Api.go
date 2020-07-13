package managedObjects

import (
	"encoding/json"
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"log"
	"net/http"
	"net/url"
)

const (
	MANAGED_OBJECT_REFERENCE_TYPE            = "application/vnd.com.nsn.cumulocity.managedObjectReference+json"
	MANAGED_OBJECT_REFERENCE_COLLECTION_TYPE = "application/vnd.com.nsn.cumulocity.managedObjectReferenceCollection+json"

	MANAGED_OBJECT_REFERENCE_API_PATH = "/inventory/managedObjects"
)

type ManagedObjectReferenceApi interface {
	// Create a new managed object reference and returns the created entity with id, creation time and other properties
	Create(managedObjectId string, referenceType ReferenceType) (*ManagedObjectReference, *generic.Error)

	// Gets an exiting managed object reference by its id. If the id does not exists, nil is returned.
	Get(managedObjectId string, referenceType ReferenceType, referenceID string) (*ManagedObjectReference, *generic.Error)

	GetMany(managedObjectId string, referenceType ReferenceType, pageSize int) (*ManagedObjectReferenceCollection, *generic.Error)

	// Deletion by managedObjectReference id. If error is nil, managed object reference was deleted successfully.
	Delete(managedObjectId string, referenceType ReferenceType, referenceID string) *generic.Error

	// Gets the next page from an existing managed object reference collection.
	// If there is no next page, nil is returned.
	NextPage(c *ManagedObjectReferenceCollection) (*ManagedObjectReferenceCollection, *generic.Error)

	// Gets the previous page from an existing managed object reference collection.
	// If there is no previous page, nil is returned.
	PreviousPage(c *ManagedObjectReferenceCollection) (*ManagedObjectReferenceCollection, *generic.Error)
}

type managedObjectReferenceApi struct {
	client   *generic.Client
	basePath string
}

// Creates a new managed object reference api object
//
// client - Must be a gomulocity client.
// returns - The `ManagedObjectReferenceApi` object
func NewManagedObjectReferenceApi(client *generic.Client) ManagedObjectReferenceApi {
	return &managedObjectReferenceApi{client, MANAGED_OBJECT_REFERENCE_API_PATH}
}

/*
Creates a new managed object reference based on the given variables.

See: https://cumulocity.com/guides/reference/inventory/#post-create-a-new-managedobject
*/
func (managedObjectReferenceApi *managedObjectReferenceApi) Create(managedObjectId string, referenceType ReferenceType) (*ManagedObjectReference, *generic.Error) {
	if len(managedObjectId) == 0 {
		return nil, generic.ClientError("managedObjectId must not be empty", "GetManagedObjectReference")
	}

	newManagedObjectReference := NewManagedObjectReference{Source{Id: managedObjectId}}
	bytes, err := json.Marshal(newManagedObjectReference)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while marshalling the managedObjectReference: %s", err.Error()), "CreateManagedObjectReference")
	}
	headers := generic.AcceptAndContentTypeHeader(MANAGED_OBJECT_REFERENCE_TYPE, MANAGED_OBJECT_REFERENCE_TYPE)

	path := fmt.Sprintf("%s/%s/%s", managedObjectReferenceApi.basePath, url.QueryEscape(managedObjectId), url.QueryEscape(string(referenceType)))
	body, status, err := managedObjectReferenceApi.client.Post(path, bytes, headers)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while posting a new managedObjectReference: %s", err.Error()), "CreateManagedObjectReference")
	}
	if status != http.StatusCreated {
		return nil, generic.CreateErrorFromResponse(body, status)
	}

	return parseManagedObjectReferenceResponse(body)
}

/*
Gets a managedObjectReference for a given Id.

Returns 'ManagedObjectReference' on success or nil if the id does not exist.
*/
func (managedObjectReferenceApi *managedObjectReferenceApi) Get(managedObjectId string, referenceType ReferenceType, referenceID string) (*ManagedObjectReference, *generic.Error) {
	if len(managedObjectId) == 0 {
		return nil, generic.ClientError("managedObjectId must not be empty", "GetManagedObjectReference")
	}
	if len(referenceID) == 0 {
		return nil, generic.ClientError("referenceID must not be empty", "GetManagedObjectReference")
	}

	path := fmt.Sprintf("%s/%s/%s/%s", managedObjectReferenceApi.basePath, url.QueryEscape(managedObjectId), url.QueryEscape(string(referenceType)), url.QueryEscape(referenceID))
	body, status, err := managedObjectReferenceApi.client.Get(path, generic.AcceptHeader(MANAGED_OBJECT_REFERENCE_TYPE))

	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while getting a managedObjectReference: %s", err.Error()), "GetManagedObjectReference")
	}
	if status == http.StatusNotFound {
		return nil, nil
	}
	if status != http.StatusOK {
		return nil, generic.CreateErrorFromResponse(body, status)
	}

	return parseManagedObjectReferenceResponse(body)
}

/*
   Returns a collection of managed object references on success or nil if the id does not exist.
*/
func (managedObjectReferenceApi *managedObjectReferenceApi) GetMany(managedObjectId string, referenceType ReferenceType, pageSize int) (*ManagedObjectReferenceCollection, *generic.Error) {
	if len(managedObjectId) == 0 {
		return nil, generic.ClientError("managedObjectId must not be empty", "GetManagedObjectReference")
	}
	queryParamsValues := &url.Values{}
	err := generic.PageSizeParameter(pageSize, queryParamsValues)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while building pageSize parameter to fetch managedObjectReferences: %s", err.Error()), "FindManagedObjectReferences")
	}

	path := fmt.Sprintf("%s/%s/%s?%s", managedObjectReferenceApi.basePath, url.QueryEscape(managedObjectId), url.QueryEscape(string(referenceType)), queryParamsValues.Encode())

	return managedObjectReferenceApi.getCommon(path)
}

/*
Deletes managedObjectReference by id.
*/
func (managedObjectReferenceApi *managedObjectReferenceApi) Delete(managedObjectId string, referenceType ReferenceType, referenceID string) *generic.Error {
	if len(managedObjectId) == 0 {
		return generic.ClientError("Deleting managedObjectReference without an id is not allowed", "DeleteManagedObjectReference")
	}
	if len(referenceID) == 0 {
		return generic.ClientError("referenceID must not be empty", "DeleteManagedObjectReference")
	}

	path := fmt.Sprintf("%s/%s/%s/%s", managedObjectReferenceApi.basePath, url.QueryEscape(managedObjectId), url.QueryEscape(string(referenceType)), url.QueryEscape(referenceID))

	body, status, err := managedObjectReferenceApi.client.Delete(path, generic.EmptyHeader())
	if err != nil {
		return generic.ClientError(fmt.Sprintf("Error while deleting managedObjectReference with id [%s]: %s", managedObjectId, err.Error()), "DeleteManagedObjectReference")
	}

	if status != http.StatusNoContent {
		return generic.CreateErrorFromResponse(body, status)
	}

	return nil
}

func (managedObjectReferenceApi *managedObjectReferenceApi) NextPage(c *ManagedObjectReferenceCollection) (*ManagedObjectReferenceCollection, *generic.Error) {
	return managedObjectReferenceApi.getPage(c.Next)
}

func (managedObjectReferenceApi *managedObjectReferenceApi) PreviousPage(c *ManagedObjectReferenceCollection) (*ManagedObjectReferenceCollection, *generic.Error) {
	return managedObjectReferenceApi.getPage(c.Prev)
}

// -- internal

func (managedObjectReferenceApi *managedObjectReferenceApi) getPage(reference string) (*ManagedObjectReferenceCollection, *generic.Error) {
	if reference == "" {
		log.Print("No page reference given. Returning nil.")
		return nil, nil
	}

	nextUrl, err := url.Parse(reference)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Unparsable URL given for page reference: '%s'", reference), "GetPage")
	}

	collection, genErr := managedObjectReferenceApi.getCommon(fmt.Sprintf("%s?%s", nextUrl.Path, nextUrl.RawQuery))
	if genErr != nil {
		return nil, genErr
	}

	if len(collection.References) == 0 {
		log.Print("Returned collection is empty. Returning nil.")
		return nil, nil
	}

	return collection, nil
}

func (managedObjectReferenceApi *managedObjectReferenceApi) getCommon(path string) (*ManagedObjectReferenceCollection, *generic.Error) {
	body, status, err := managedObjectReferenceApi.client.Get(path, generic.AcceptHeader(MANAGED_OBJECT_REFERENCE_COLLECTION_TYPE))
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while getting managedObjectReferences: %s", err.Error()), "GetManagedObjectReferenceCollection")
	}

	if status == http.StatusNotFound {
		return nil, nil
	}
	if status != http.StatusOK {
		return nil, generic.CreateErrorFromResponse(body, status)
	}

	var result ManagedObjectReferenceCollection
	if len(body) > 0 {
		err = json.Unmarshal(body, &result)
		if err != nil {
			return nil, generic.ClientError(fmt.Sprintf("Error while parsing response JSON: %s", err.Error()), "GetManagedObjectReferenceCollection")
		}
	} else {
		return nil, generic.ClientError("Response body was empty", "GetManagedObjectReferenceCollection")
	}

	return &result, nil
}

func parseManagedObjectReferenceResponse(body []byte) (*ManagedObjectReference, *generic.Error) {
	var result ManagedObjectReference
	if len(body) > 0 {
		err := json.Unmarshal(body, &result)
		if err != nil {
			return nil, generic.ClientError(fmt.Sprintf("Error while parsing response JSON: %s", err.Error()), "ResponseParser")
		}
	} else {
		return nil, generic.ClientError("Response body was empty", "GetManagedObjectReference")
	}

	return &result, nil
}