package inventory

import (
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"reflect"
	"testing"
)

func TestInventoryApi_ParseManagedObjectResponse(t *testing.T) {
	tests := []struct {
		name                  string
		givenResponseBody     string
		expectedManagedObject *ManagedObject
		expectedErr           *generic.Error
	}{
		{
			name:                  "parsable body",
			givenResponseBody:     givenResponseBody,
			expectedManagedObject: expectedManagedObject,
			expectedErr:           nil,
		}, {
			name:                  "unparsable body",
			givenResponseBody:     `{`,
			expectedManagedObject: nil,
			expectedErr: &generic.Error{
				ErrorType: "ClientError",
				Message:   "Error while parsing response JSON: Error while unmarshaling json: unexpected end of JSON input",
				Info:      "ResponseParser",
			},
		}, {
			name:                  "without body",
			givenResponseBody:     ``,
			expectedManagedObject: nil,
			expectedErr: &generic.Error{
				ErrorType: "ClientError",
				Message:   "Response body was empty",
				Info:      "ResponseParser",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			managedObject, err := parseManagedObjectResponse([]byte(tt.givenResponseBody))

			if fmt.Sprint(err) != fmt.Sprint(tt.expectedErr) {
				t.Fatalf("Unexpected error was returned: %s\nExpected: %s", err, tt.expectedErr)
			}

			if !reflect.DeepEqual(managedObject, tt.expectedManagedObject) {
				t.Errorf("Unexpected managedObject was returned: %#v\nExpected: %#v", managedObject, tt.expectedManagedObject)
			}
		})
	}
}
