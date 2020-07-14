package inventory

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type InventoryFilter struct {
	Type         string
	FragmentType string
	Ids          []int
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
		var idsAsString []string
		for _, id := range inventoryFilter.Ids {
			idsAsString = append(idsAsString, strconv.Itoa(id))
		}
		params.Add("ids", strings.Join(idsAsString, ","))
	}

	if len(inventoryFilter.Text) > 0 {
		params.Add("text", inventoryFilter.Text)
	}

	return nil
}
