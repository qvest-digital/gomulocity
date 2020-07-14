package inventory

import (
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestManagedObjectApi_CreateManagedObject(t *testing.T) {
	creationDate, _ := time.Parse(time.RFC3339, "2019-11-07T20:48:43.472Z")
	lastUpdated, _ := time.Parse(time.RFC3339, "2020-04-24T09:36:05.112Z")

	tests := []struct {
		name        string
		c8yRespCode int
		c8yRespBody string
		expectedErr error

		requestData  *NewManagedObject
		responseData NewManagedObject
	}{
		{
			name: "create managed object - happy",

			requestData: &NewManagedObject{
				Type:         "<type>",
				Name:         "<name>",
				CreationDate: creationDate,
			},

			c8yRespCode: http.StatusCreated,
			c8yRespBody: `{"com_cumulocity_model_BinarySwitch": {"state": "<state>"},"id": "<id>","self": "<self>","type": "<type>","name": "<name>","creationDate":"2019-11-07T20:48:43.472Z","lastUpdated": "2020-04-24T09:36:05.112Z"}`,

			responseData: NewManagedObject{
				ID:           "<id>",
				Self:         "<self>",
				Type:         "<type>",
				Name:         "<name>",
				CreationDate: creationDate,
				LastUpdated:  lastUpdated,
				BinarySwitch: struct {
					State string `json:"state"`
				}{
					State: "<state>",
				},
			},
		},
		{
			name: "create managed object - status is not StatusCreated",
			requestData: &NewManagedObject{
				Type:         "type",
				Name:         "name",
				CreationDate: time.Time{},
			},
			c8yRespCode: http.StatusInternalServerError,
			c8yRespBody: `{"error":"someErr", "message":"someMessage", "info":"someInfo"}`,
			expectedErr: generic.Error{
				ErrorType: "someErr",
				Message:   "someMessage",
				Info:      "someInfo",
			},
		},
		{
			name: "create managed object - failed to unmarshal result",
			requestData: &NewManagedObject{
				Type:         "type",
				Name:         "name",
				CreationDate: time.Time{},
			},
			c8yRespCode: http.StatusCreated,
			c8yRespBody: `<invalid response body>`,
			expectedErr: clientError("Error while unmarshalling response: invalid character '<' looking for beginning of value", "NewManagedObject"),
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
			c := inventoryApi{
				Client: &generic.Client{
					HTTPClient: testServer.Client(),
					BaseURL:    testServer.URL,
					Username:   u,
					Password:   p,
				},
				ManagedObjectsPath: managedObjectPath,
			}

			newManagedObject, err := c.CreateManagedObject(tt.requestData)
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

			if err == nil {
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
			}
		})
	}
}
