package inventory

import (
	"encoding/json"
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"log"
	"net/http"
	"net/url"
)

const (
	MANAGED_OBJECT_TYPE = "application/vnd.com.nsn.cumulocity.managedObject+json"
	MANAGED_OBJECT_COLLECTION_TYPE = "application/vnd.com.nsn.cumulocity.managedObjectCollection+json"

	MANAGED_OBJECT_API_PATH = "/inventory/managedObjects"
)

type ManagedObjectApi interface {
	// Create a new managed object and returns the created entity with id, creation time and other properties
	Create(newManagedObject *NewManagedObject) (*ManagedObject, *generic.Error)

	// Gets an exiting managed object by its id. If the id does not exists, nil is returned.
	Get(managedObjectId string) (*ManagedObject, *generic.Error)

	Update(managedObjectId string, managedObject *ManagedObjectUpdate) (*ManagedObject, *generic.Error)

	// Deletion by managedObject id. If error is nil, managed object was deleted successfully.
	Delete(managedObjectId string) *generic.Error

	// Returns a managed object collection, found by the given managed object filter parameters.
	// All query parameters are AND concatenated.
	Find(managedObjectFilter *ManagedObjectFilter, pageSize int) (*ManagedObjectCollection, *generic.Error)

	// Returns a managed object collection, found by the given managed object query.
	// See the query language: https://cumulocity.com/guides/reference/inventory/#query-language
	FindByQuery(query string, pageSize int) (*ManagedObjectCollection, *generic.Error)

	// Gets the next page from an existing managed object collection.
	// If there is no next page, nil is returned.
	NextPage(c *ManagedObjectCollection) (*ManagedObjectCollection, *generic.Error)

	// Gets the previous page from an existing managed object collection.
	// If there is no previous page, nil is returned.
	PreviousPage(c *ManagedObjectCollection) (*ManagedObjectCollection, *generic.Error)
}

type managedObjectApi struct {
	client   *generic.Client
	basePath string
}

// Creates a new managed object api object
//
// client - Must be a gomulocity client.
// returns - The `managed object`-api object
func NewManagedObjectApi(client *generic.Client) ManagedObjectApi {
	return &managedObjectApi{client, MANAGED_OBJECT_API_PATH}
}

/*
Creates a new managed object based on the given variables.

See: https://cumulocity.com/guides/reference/inventory/#post-create-a-new-managedobject
*/
func (managedObjectApi *managedObjectApi) Create(newManagedObject *NewManagedObject) (*ManagedObject, *generic.Error) {
	bytes, err := json.Marshal(newManagedObject)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while marshalling the managedObject: %s", err.Error()), "CreateManagedObject")
	}
	headers := generic.AcceptAndContentTypeHeader(MANAGED_OBJECT_TYPE, MANAGED_OBJECT_TYPE)

	body, status, err := managedObjectApi.client.Post(managedObjectApi.basePath, bytes, headers)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while posting a new managedObject: %s", err.Error()), "CreateManagedObject")
	}
	if status != http.StatusCreated {
		return nil, generic.CreateErrorFromResponse(body, status)
	}

	return parseManagedObjectResponse(body)
}

/*
Gets a managedObject for a given Id.

Returns 'ManagedObject' on success or nil if the id does not exist.
*/
func (managedObjectApi *managedObjectApi) Get(managedObjectId string) (*ManagedObject, *generic.Error) {
	if len(managedObjectId) == 0 {
		return nil, generic.ClientError("managedObjectId must not be empty", "GetManagedObject")
	}

	path := fmt.Sprintf("%s/%s", managedObjectApi.basePath, url.QueryEscape(managedObjectId))
	body, status, err := managedObjectApi.client.Get(path, generic.AcceptHeader(MANAGED_OBJECT_TYPE))

	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while getting a managedObject: %s", err.Error()), "GetManagedObject")
	}
	if status == http.StatusNotFound {
		return nil, nil
	}
	if status != http.StatusOK {
		return nil, generic.CreateErrorFromResponse(body, status)
	}

	return parseManagedObjectResponse(body)
}


/*
Updates the managedObject with given Id.

See: https://cumulocity.com/guides/reference/managedObjects/#update-an-managedObject
*/
func (managedObjectApi *managedObjectApi) Update(managedObjectId string, managedObject *ManagedObjectUpdate) (*ManagedObject, *generic.Error) {
	if len(managedObjectId) == 0 {
		return nil, generic.ClientError("Updating managedObject without an id is not allowed", "UpdateManagedObject")
	}
	bytes, err := json.Marshal(managedObject)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while marshalling the update managedObject: %s", err.Error()), "UpdateManagedObject")
	}

	path := fmt.Sprintf("%s/%s", managedObjectApi.basePath, url.QueryEscape(managedObjectId))
	headers := generic.AcceptAndContentTypeHeader(MANAGED_OBJECT_TYPE, MANAGED_OBJECT_TYPE)

	body, status, err := managedObjectApi.client.Put(path, bytes, headers)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while updating an managedObject: %s", err.Error()), "UpdateManagedObject")
	}
	if status != http.StatusOK {
		return nil, generic.CreateErrorFromResponse(body, status)
	}

	return parseManagedObjectResponse(body)
}

/*
Deletes managedObject by id.
*/
func (managedObjectApi *managedObjectApi) Delete(managedObjectId string) *generic.Error {
	if len(managedObjectId) == 0 {
		return generic.ClientError("Deleting managedObject without an id is not allowed", "DeleteManagedObject")
	}

	body, status, err := managedObjectApi.client.Delete(fmt.Sprintf("%s/%s", managedObjectApi.basePath, url.QueryEscape(managedObjectId)), generic.EmptyHeader())
	if err != nil {
		return generic.ClientError(fmt.Sprintf("Error while deleting managedObject with id [%s]: %s", managedObjectId, err.Error()), "DeleteManagedObject")
	}

	if status != http.StatusNoContent {
		return generic.CreateErrorFromResponse(body, status)
	}

	return nil
}

/*
   Returns a collection of managed objects.

   See: https://cumulocity.com/guides/reference/inventory/#managed-object-collection
*/
func (managedObjectApi *managedObjectApi) Find(managedObjectFilter *ManagedObjectFilter, pageSize int) (*ManagedObjectCollection, *generic.Error) {
	queryParamsValues := &url.Values{}
	err := managedObjectFilter.QueryParams(queryParamsValues)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while building query parameters to search for managedObjects: %s", err.Error()), "FindManagedObjects")
	}

	err = generic.PageSizeParameter(pageSize, queryParamsValues)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while building pageSize parameter to fetch managedObjects: %s", err.Error()), "FindManagedObjects")
	}

	return managedObjectApi.getCommon(fmt.Sprintf("%s?%s", managedObjectApi.basePath, queryParamsValues.Encode()))
}

func (managedObjectApi *managedObjectApi) FindByQuery(query string, pageSize int) (*ManagedObjectCollection, *generic.Error) {
	queryParamsValues := &url.Values{}
	if len(query) > 0 {
		queryParamsValues.Add("query", query)
	}

	err := generic.PageSizeParameter(pageSize, queryParamsValues)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while building pageSize parameter to fetch managedObjects: %s", err.Error()), "FindManagedObjectsByQuery")
	}

	return managedObjectApi.getCommon(fmt.Sprintf("%s?%s", managedObjectApi.basePath, queryParamsValues.Encode()))
}


func (managedObjectApi *managedObjectApi) NextPage(c *ManagedObjectCollection) (*ManagedObjectCollection, *generic.Error) {
	return managedObjectApi.getPage(c.Next)
}

func (managedObjectApi *managedObjectApi) PreviousPage(c *ManagedObjectCollection) (*ManagedObjectCollection, *generic.Error) {
	return managedObjectApi.getPage(c.Prev)
}



// -- internal

func (managedObjectApi *managedObjectApi) getPage(reference string) (*ManagedObjectCollection, *generic.Error) {
	if reference == "" {
		log.Print("No page reference given. Returning nil.")
		return nil, nil
	}

	nextUrl, err := url.Parse(reference)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Unparsable URL given for page reference: '%s'", reference), "GetPage")
	}

	collection, genErr := managedObjectApi.getCommon(fmt.Sprintf("%s?%s", nextUrl.Path, nextUrl.RawQuery))
	if genErr != nil {
		return nil, genErr
	}

	if len(collection.ManagedObjects) == 0 {
		log.Print("Returned collection is empty. Returning nil.")
		return nil, nil
	}

	return collection, nil
}

func (managedObjectApi *managedObjectApi) getCommon(path string) (*ManagedObjectCollection, *generic.Error) {
	body, status, err := managedObjectApi.client.Get(path, generic.AcceptHeader(MANAGED_OBJECT_COLLECTION_TYPE))
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while getting managedObjects: %s", err.Error()), "GetManagedObjectCollection")
	}

	if status != http.StatusOK {
		return nil, generic.CreateErrorFromResponse(body, status)
	}

	var result ManagedObjectCollection
	if len(body) > 0 {
		err = json.Unmarshal(body, &result)
		if err != nil {
			return nil, generic.ClientError(fmt.Sprintf("Error while parsing response JSON: %s", err.Error()), "GetManagedObjectCollection")
		}
	} else {
		return nil, generic.ClientError("Response body was empty", "GetManagedObjectCollection")
	}

	return &result, nil
}

func parseManagedObjectResponse(body []byte) (*ManagedObject, *generic.Error) {
	var result ManagedObject
	if len(body) > 0 {
		err := json.Unmarshal(body, &result)
		if err != nil {
			return nil, generic.ClientError(fmt.Sprintf("Error while parsing response JSON: %s", err.Error()), "ResponseParser")
		}
	} else {
		return nil, generic.ClientError("Response body was empty", "GetManagedObject")
	}

	return &result, nil
}
