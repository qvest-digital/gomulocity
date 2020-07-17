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
	self                 string
	externalId           string
	externalIdOfGlobalId string
}

type ExternalIDCollection struct {
	self        string
	externalIds []ExternalID
	prev        string
	next        string
}

type ExternalID struct {
	self          string
	externalId    string
	typ           string
	managedObject deviceinformation.ManagedObject
}

type IdentityAPI interface {
	GetIdentity() (Identity, *generic.Error)
	GetExternalIDCollection() //TODO
	GetExternalID(externalID ExternalID) (ExternalID, *generic.Error)
	CreateExternalID(ID ExternalID) (ExternalID, *generic.Error)
	DeleteExternalID(externalID ExternalID) *generic.Error
}

type identityAPI struct {
	basePath string
	client   *generic.Client
}

func NewIdentityAPI(client *generic.Client, basePath string) identityAPI {
	return identityAPI{
		client:   client,
		basePath: basePath,
	}
}

func (i identityAPI) GetIdentity(identity Identity) (Identity, *generic.Error) {
	body, status, err := i.client.Get(fmt.Sprintf("%s/%s", i.basePath, url.QueryEscape("identity")), generic.AcceptHeader(IDENTITY_TYPE))

	if err != nil {
		return Identity{}, generic.ClientError(fmt.Sprintf("Error while getting the Identity Ressource: %s", err.Error()), "Get")
	}
	if status == http.StatusNotFound {
		return Identity{}, nil
	}
	if status != http.StatusOK {
		return Identity{}, generic.CreateErrorFromResponse(body, status)
	}
	result := Identity{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return Identity{}, generic.ClientError(fmt.Sprintf("Error while parsing response JSON: %s", err.Error()), "ResponseParser")
	}

	return result, nil
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
