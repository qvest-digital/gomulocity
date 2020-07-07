package managedObjects

import (
	"net/url"
)

type ManagedObjectCollectionFilter struct {
	Type          string
	Owner         string
	FragmentType  string
	QueryLanguage string //please note: If the 'QueryLanguage' has been set, other parameters like 'DeviceID' and 'Type' will be ignored.

}

func (m ManagedObjectCollectionFilter) QueryParams() string {
	params := url.Values{}

	if len(m.QueryLanguage) > 0 {
		return m.QueryLanguage
	}

	if len(m.Type) > 0 {
		params.Add("type", m.Type)
	}
	if len(m.Owner) > 0 {
		params.Add("owner", m.Type)
	}
	if len(m.FragmentType) > 0 {
		params.Add("fragmentType", m.FragmentType)
	}
	return params.Encode()
}
