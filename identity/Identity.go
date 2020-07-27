package identity

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/tarent/gomulocity/deviceinformation"
	"github.com/tarent/gomulocity/generic"
)

const (
	IDENTITY_TYPE               = "application/vnd.com.nsn.cumulocity.identityApi+json"
	EXTERNAL_ID_COLLECTION_TYPE = "application/vnd.com.nsn.cumulocity.exteralIdCollection+json"
	EXTERNAL_ID_TYPE            = "application/vnd.com.nsn.cumulocity.externalId+json"
)

type Identity struct {
	Self                 string `json:"self"`
	ExternalId           string `json:"externalId"`
	ExternalIdOfGlobalId string `json:"externalIdOfGlobalId"`
}

type ExternalIDCollection struct {
	Self        string
	ExternalIds []ExternalID
	Prev        string
	Next        string
}

type ExternalID struct {
	Self          string
	ExternalId    string
	Type          string
	ManagedObject deviceinformation.ManagedObject
}

type IdentityAPI interface {
	GetIdentity() (*Identity, *generic.Error)
	GetExternalID(externalIDType, externalID string) (*ExternalID, *generic.Error)
	CreateExternalID(ID ExternalID) (ExternalID, *generic.Error)
	DeleteExternalID(externalIDType, externalID string) *generic.Error
}

type identityAPI struct {
	basePath string
	client   generic.Client
}

func NewIdentityAPI(client generic.Client) identityAPI {
	return identityAPI{
		client:   client,
		basePath: "/identity",
	}
}

func (i identityAPI) GetIdentity() (*Identity, *generic.Error) {
	body, status, err := i.client.Get(i.basePath, generic.AcceptHeader(IDENTITY_TYPE))

	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while getting the Identity Ressource: %s", err.Error()), "Get")
	}
	if status == http.StatusNotFound {
		return nil, nil
	}
	if status != http.StatusOK {
		return nil, generic.CreateErrorFromResponse(body, status)
	}
	var result Identity
	if len(body) > 0 {
		err := generic.ObjectFromJson(body, &result)
		if err != nil {
			return nil, generic.ClientError(fmt.Sprintf("Error while parsing response JSON: %s", err.Error()), "ResponseParser")
		}
	} else {
		return nil, generic.ClientError("Response body was empty", "GetEvent")
	}

	return &result, nil
}

func (i identityAPI) CreateExternalID(externalId ExternalID) (ExternalID, *generic.Error) {
	bytes, err := json.Marshal(externalId)
	if err != nil {
		return ExternalID{}, generic.ClientError(fmt.Sprintf("Error while marshalling the event: %s", err.Error()), "CreateExternalID")
	}

	body, status, err := i.client.Post(i.basePath, bytes, generic.AcceptAndContentTypeHeader(EXTERNAL_ID_TYPE, EXTERNAL_ID_TYPE))
	if err != nil {
		return ExternalID{}, generic.ClientError(fmt.Sprintf("Error while posting a new event: %s", err.Error()), "CreateEvent")
	}
	if status != http.StatusCreated {
		return ExternalID{}, generic.CreateErrorFromResponse(body, status)
	}
	result := ExternalID{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return ExternalID{}, generic.ClientError(fmt.Sprintf("Error while parsing response JSON: %s", err.Error()), "ResponseParser")
	}

	return result, nil
}

func (i identityAPI) GetExternalID(externalIDtype string, externalID string) (*ExternalID, *generic.Error) {
	body, status, err := i.client.Get(fmt.Sprintf("%s/%s/%s/%s", i.basePath, url.QueryEscape("extrenalIds"), url.QueryEscape(externalIDtype), url.QueryEscape(externalID)), generic.AcceptHeader(IDENTITY_TYPE))

	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while getting an externalID: %s", err.Error()), "get")
	}
	if status == http.StatusNotFound {
		return nil, nil
	}
	if status != http.StatusOK {
		return nil, generic.CreateErrorFromResponse(body, status)
	}
	result := ExternalID{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while parsing response JSON: %s", err.Error()), "ResponseParser")
	}

	return &result, nil
}

func (i identityAPI) DeleteExternalID(externalIDType, externalID string) *generic.Error {
	if len(externalIDType) == 0 || len(externalID) == 0 {
		return generic.ClientError("Deleting deviceRegistrations without an id is not allowed", "DeleteDeviceRegistration")
	}

	path := fmt.Sprintf("%s/%s/%s/%s", i.basePath, "externalIds", url.QueryEscape(externalIDType), url.QueryEscape(externalID))
	body, status, err := i.client.Delete(path, generic.EmptyHeader())
	if err != nil {
		return generic.ClientError(fmt.Sprintf("Error while deleting an ExternalID with id %s", err.Error()), "Delete ExternalID")
	}

	if status != http.StatusNoContent {
		return generic.CreateErrorFromResponse(body, status)
	}

	return nil
}
