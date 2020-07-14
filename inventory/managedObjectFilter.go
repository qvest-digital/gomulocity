package inventory

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type ManagedObjectFilter struct {
	Type         string
	FragmentType string
	Ids          []int
	Text         string
}

// Appends the filter query parameters to the provided parameter values for a request.
// When provided values is nil an error will be created
func (managedObjectFilter ManagedObjectFilter) QueryParams(params *url.Values) error {
	if params == nil {
		return fmt.Errorf("The provided parameter values must not be nil!")
	}

	if len(managedObjectFilter.Type) > 0 {
		params.Add("type", managedObjectFilter.Type)
	}

	if len(managedObjectFilter.FragmentType) > 0 {
		params.Add("fragmentType", managedObjectFilter.FragmentType)
	}

	if len(managedObjectFilter.Ids) > 0 {
		var idsAsString []string
		for _, id := range managedObjectFilter.Ids {
			idsAsString = append(idsAsString, strconv.Itoa(id))
		}
		params.Add("ids", strings.Join(idsAsString, ","))
	}

	if len(managedObjectFilter.Text) > 0 {
		params.Add("text", managedObjectFilter.Text)
	}

	return nil
}
