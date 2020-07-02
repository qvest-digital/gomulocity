package managedObjects

import (
	"errors"
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_CreateManagedObject(t *testing.T) {
	tests := []struct {
		testName    string
		c8yRespCode int
		c8yRespBody string
		expectedErr error

		name         string
		state        string
		responseData NewManagedObject
	}{
		{
			testName:    "create managed object - happy",
			c8yRespCode: http.StatusCreated,
			c8yRespBody: `{"self":"<self>", "id":"<id>", "lastUpdated":"<lastUpdated>", "name":"<name>", "com_cumulocity_model_BinarySwitch":{"state":"<state>"}}`,

			name:  "<name>",
			state: "<state>",
			responseData: NewManagedObject{
				Self:        "<self>",
				ID:          "<id>",
				LastUpdated: "<lastUpdated>",
				Name:        "<name>",
				BinarySwitch: struct {
					State string `json:"state"`
				}{
					State: "<state>",
				},
			},
		},
		{
			testName:    "create managed object - bad credentials",
			c8yRespCode: http.StatusUnauthorized,
			c8yRespBody: `{}`,
			expectedErr: generic.BadCredentialsErr,

			name:  "<name>",
			state: "<state>",
		},
		{
			testName:    "create managed object - access denied",
			c8yRespCode: http.StatusForbidden,
			c8yRespBody: `{}`,
			expectedErr: generic.AccessDeniedErr,

			name:  "<name>",
			state: "<state>",
		},
		{
			testName:    "create managed object - unexpected error",
			c8yRespCode: http.StatusInternalServerError,
			c8yRespBody: `{}`,
			expectedErr: errors.New("received an unexpected status code: 500"),

			name:  "<name>",
			state: "<state>",
		},
		{
			testName:    "create managed object - invalid response body",
			c8yRespCode: http.StatusCreated,
			c8yRespBody: `<invalid response body>`,
			expectedErr: errors.New("failed to unmarshal response body: invalid character '<' looking for beginning of value"),

			name:  "<name>",
			state: "<state>",
		},
	}

	for _, tt := range tests {
		var username, password string
		t.Run(tt.name, func(t *testing.T) {
			testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				username, password, _ = req.BasicAuth()
				res.WriteHeader(tt.c8yRespCode)
				_, err := res.Write([]byte(tt.c8yRespBody))
				if err != nil {
					t.Fatalf("failed to write resp body: %s", err)
				}

			}))
			defer testServer.Close()

			u := "<username>"
			p := "<password>"
			c := ManagedObjectApi{
				Client: &generic.Client{
					HTTPClient: testServer.Client(),
					BaseURL:    testServer.URL,
					Username:   u,
					Password:   p,
				},
				ManagedObjectsPath: managedObjectPath,
			}

			newManagedObject, err := c.CreateManagedObject(tt.name, tt.state)
			if err != nil {
				if fmt.Sprint(err) != fmt.Sprint(tt.expectedErr) {
					t.Errorf("received an unexpected error: %v expected: %v", err, tt.expectedErr)
				}
			}

			if username != u {
				t.Errorf("expected c8y auth-username: %v expected: %v", u, username)
			}

			if password != p {
				t.Errorf("expected c8y auth-password: %v expected: %v", u, username)
			}

			if newManagedObject.ID != tt.responseData.ID {
				t.Errorf("ID is incorrect: %v expected: %v", newManagedObject.ID, tt.responseData.ID)
			}
			if newManagedObject.Name != tt.responseData.Name {
				t.Errorf("name is incorrect: %v expected: %v", newManagedObject.Name, tt.responseData.Name)
			}
			if newManagedObject.Self != tt.responseData.Self {
				t.Errorf("self is incorrect: %v expected: %v", newManagedObject.Self, tt.responseData.Self)
			}
			if newManagedObject.LastUpdated != tt.responseData.LastUpdated {
				t.Errorf("lastUpdated is incorrect: %v expected: %v", newManagedObject.LastUpdated, tt.responseData.LastUpdated)
			}
			if newManagedObject.BinarySwitch != tt.responseData.BinarySwitch {
				t.Errorf("BinarySwitch is incorrect: %v expected: %v", newManagedObject.BinarySwitch, tt.responseData.BinarySwitch)
			}
		})
	}
}
