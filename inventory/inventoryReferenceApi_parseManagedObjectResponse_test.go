package inventory

import (
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"reflect"
	"testing"
)

func TestInventoryReferenceApi_ParseManagedObjectResponse(t *testing.T) {
	tests := []struct {
		name                           string
		givenResponseBody              string
		expectedManagedObjectReference *ManagedObjectReference
		expectedErr                    *generic.Error
	}{
		{
			name:                           "parsable body",
			givenResponseBody:              givenReferenceResponseBody,
			expectedManagedObjectReference: expectedManagedObjectReference,
			expectedErr:                    nil,
		}, {
			name:                           "unparsable body",
			givenResponseBody:              `{`,
			expectedManagedObjectReference: nil,
			expectedErr: &generic.Error{
				ErrorType: "ClientError",
				Message:   "Error while parsing response JSON: unexpected end of JSON input",
				Info:      "ResponseParser",
			},
		}, {
			name:                           "without body",
			givenResponseBody:              ``,
			expectedManagedObjectReference: nil,
			expectedErr: &generic.Error{
				ErrorType: "ClientError",
				Message:   "Response body was empty",
				Info:      "ResponseParser",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			managedObject, err := parseManagedObjectReferenceResponse([]byte(tt.givenResponseBody))

			if fmt.Sprint(err) != fmt.Sprint(tt.expectedErr) {
				t.Fatalf("Unexpected error was returned: %s\nExpected: %s", err, tt.expectedErr)
			}

			if !reflect.DeepEqual(managedObject, tt.expectedManagedObjectReference) {
				t.Errorf("Unexpected managedObject was returned: %#v\nExpected: %#v", managedObject, tt.expectedManagedObjectReference)
			}
		})
	}
}
