package identity

import (
	"github.com/tarent/gomulocity/deviceinformation"
	"github.com/tarent/gomulocity/generic"
)

type Identity struct {
	self string
	externalId string
	externalIdOfGlobalId string
}

type ExternalIDCollection{
	self string
	externalIds []ExternalID
	prev string
	next string
}

type ExternalID struct {
	self string
	externalId string
	typ string
	managedObject deviceinformation.ManagedObject
}

type IdentityAPI interface{
	GetIdentity()(Identity, *generic.Error)
	GetExternalIDCollection() //TODO
	GetExternalID(externalID ExternalID)(ExternalID, *generic.Error)
	CreateExternalID(ID ExternalID)(ExternalID, *generic.Error)
	DeleteExternalID(externalID ExternalID)(*generic.Error)
}
func NewIdentityAPI(client *generic.Client)IdentityAPI{
	return IdentityAPI{client: client}
}
