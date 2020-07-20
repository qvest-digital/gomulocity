package inventory

import (
	"net/url"
	"testing"
)

func TestInventoryFilter_QueryParams_Successful(t *testing.T) {
	// given
	collectionFilter := InventoryFilter{
		Type:         "Test-Type",
		FragmentType: "Test-FragmentType",
		Ids:          []string{"1","2","3"},
		Text:         "Test Text",
	}
	queryParamsValues := &url.Values{}
	queryParamsValues.Add("a", "b")

	// when
	err := collectionFilter.QueryParams(queryParamsValues)

	// then
	if err != nil  {
		t.Errorf("Unexpected error was returned: %s", err)
	}

	expectedQuery := "a=b&fragmentType=Test-FragmentType&ids=1%2C2%2C3&text=Test+Text&type=Test-Type"
	if queryParamsValues.Encode() != expectedQuery {
		t.Errorf("Unexpected query params were created: %s; expected: %s", queryParamsValues.Encode(), expectedQuery)
	}
}

func TestInventoryFilter_QueryParams_WithOneId(t *testing.T) {
	// given
	collectionFilter := InventoryFilter{
		Ids:          []string{"1"},
	}
	queryParamsValues := &url.Values{}

	// when
	err := collectionFilter.QueryParams(queryParamsValues)

	// then
	if err != nil  {
		t.Errorf("Unexpected error was returned: %s", err)
	}

	expectedQuery := "ids=1"
	if queryParamsValues.Encode() != expectedQuery {
		t.Errorf("Unexpected query params were created: %s; expected: %s", queryParamsValues.Encode(), expectedQuery)
	}
}

func TestInventoryFilter_QueryParams_WithEmptyFilter(t *testing.T) {
	// given
	collectionFilter := InventoryFilter{}
	queryParamsValues := &url.Values{}
	queryParamsValues.Add("a", "b")

	// when
	err := collectionFilter.QueryParams(queryParamsValues)

	// then
	if err != nil  {
		t.Errorf("Unexpected error was returned: %s", err)
	}

	expectedQuery := "a=b"
	if queryParamsValues.Encode() != expectedQuery {
		t.Errorf("Unexpected query params were created: %s; expected: %s", queryParamsValues.Encode(), expectedQuery)
	}
}

func TestInventoryFilter_QueryParams_WithoutQueryParamValues(t *testing.T) {
	// given
	collectionFilter := InventoryFilter{
		Type:         "Test-Type",
		FragmentType: "Test-FragmentType",
		Ids:          []string{"1","2","3"},
		Text:         "Test Text",
	}

	// when
	err := collectionFilter.QueryParams(nil)

	// then
	if err == nil  {
		t.Error("Expected an error but no one was returned")
	}

	expectedError := "The provided parameter values must not be nil!"
	if err.Error() != expectedError {
		t.Errorf("Unexpected error was returned: %s; expected: %s", err, expectedError)
	}
}
