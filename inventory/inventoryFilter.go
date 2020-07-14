package inventory

import (
	"fmt"
	"net/url"
	"strings"
)

type InventoryFilter struct {
	Type         string
	FragmentType string
	Ids          []string
	Text         string
}

// Appends the filter query parameters to the provided parameter values for a request.
// When provided values is nil an error will be created
func (inventoryFilter InventoryFilter) QueryParams(params *url.Values) error {
	if params == nil {
		return fmt.Errorf("The provided parameter values must not be nil!")
	}

	if len(inventoryFilter.Type) > 0 {
		params.Add("type", inventoryFilter.Type)
	}

	if len(inventoryFilter.FragmentType) > 0 {
		params.Add("fragmentType", inventoryFilter.FragmentType)
	}

	if len(inventoryFilter.Ids) > 0 {
		params.Add("ids", strings.Join(inventoryFilter.Ids, ","))
	}

	if len(inventoryFilter.Text) > 0 {
		params.Add("text", inventoryFilter.Text)
	}

	return nil
}
