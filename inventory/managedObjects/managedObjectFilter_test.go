package managedObjects

import (
	"strings"
	"testing"
)

func TestManagedObjectCollectionFilter_QueryParams_WithoutQueryLanguage(t *testing.T) {
	// given
	collectionFilter := ManagedObjectCollectionFilter{
		Type:          "Type",
		Owner:         "Owner",
		FragmentType:  "FragmentType",
		QueryLanguage: "",
	}

	// when
	query := collectionFilter.QueryParams()

	// then
	if len(query) == 0 {
		t.Error("query is empty")
	}

	if !strings.Contains(query, "type") {
		t.Error("query does not contain type")
	}
	if !strings.Contains(query, "owner") {
		t.Error("query does not contain type")
	}
	if !strings.Contains(query, "fragment") {
		t.Error("query does not contain type")
	}
}

func TestManagedObjectCollectionFilter_QueryParams_WithQueryLanguage(t *testing.T) {
	// given
	collectionFilter := ManagedObjectCollectionFilter{
		Type:          "Type",
		Owner:         "Owner",
		FragmentType:  "FragmentType",
		QueryLanguage: "query=has(c8y_IsDevice)",
	}

	// when
	query := collectionFilter.QueryParams()

	// then
	if len(query) == 0 {
		t.Error("query is empty")
	}

	if strings.Contains(query, "type") {
		t.Error("query does not contain type")
	}
	if strings.Contains(query, "owner") {
		t.Error("query does not contain type")
	}
	if strings.Contains(query, "fragment") {
		t.Error("query does not contain type")
	}
	if query != collectionFilter.QueryLanguage {
		t.Error("query does not contain the given query language")
	}
}
