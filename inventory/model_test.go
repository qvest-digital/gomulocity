package inventory

import (
	"testing"
)

func TestManagedObject_FilterForField(t *testing.T) {
	object := ManagedObject{
		AdditionalFields: map[string]interface{}{
			"custom2": struct {
				ID   int
				Name string
			}{
				ID:   1,
				Name: "foo",
			},
		},
	}

	result, _, err := object.FilterAdditionalFieldByName("custom2")
	if err != nil {
		t.Errorf("received an unexpected error: %s", err)
	}

	if _, ok := result.(struct {
		ID   int
		Name string
	}); !ok {
		t.Errorf("received an unexpected type")
	}
}

func TestManagedObject_FilterForField_Error(t *testing.T) {
	object := ManagedObject{
		AdditionalFields: map[string]interface{}{},
	}

	_, _, err := object.FilterAdditionalFieldByName("custom2")
	if err == nil {
		t.Errorf("received an unexpected error: %s", err)
	}
}
