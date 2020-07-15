package identity

import (
	"github.com/tarent/gomulocity/deviceinformation"
	"github.com/tarent/gomulocity/generic"
)

type IdentityAPI struct{
	client *generic.Client
}

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

func NewIdentityAPI(client *generic.Client)IdentityAPI{
	return IdentityAPI{client: client}
}
