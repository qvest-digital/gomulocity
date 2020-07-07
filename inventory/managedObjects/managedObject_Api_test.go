package managedObjects

import "testing"

func TestManagedObjectCollection_NextPage(t *testing.T) {
	collection := ManagedObjectCollection{
		Next: "https://nextpage.com/examplePath?exampleQuery=example",
	}

	nextPage, err := collection.NextPage()
	if err != nil {
		t.Errorf("failed to get next page: %v", err)
	}

	if nextPage != "/examplePath?exampleQuery=example"{
		t.Errorf("unexpected nextPage. Expected: %v, Actual: %v", "/examplePath?exampleQuery=example", nextPage)
	}
}

func TestNewManagedObjectCollection_HasNextPage(t *testing.T) {
	collection := ManagedObjectCollection{
		Next: "<nextPage>",
	}

	result := collection.hasNextPage()

	if !result {
		t.Error("collection hasn't a next page")
	}
}
