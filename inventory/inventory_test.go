package inventory

import (
	"errors"
	"fmt"
	"github.com/tarent/gomulocity/deviceinformation"
	"github.com/tarent/gomulocity/generic"
	"github.com/tarent/gomulocity/models"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestClient_GetManageObjects(t *testing.T) {
	tests := []struct {
		name           string
		req            func(*http.Request)
		c8yRespCode    int
		c8yRespBody    string
		PagingStatics  ManagedObjectsPagingStatics
		FragmentFilter string
		expectedErr    error
	}{
		{
			name:        "get manageObjects - happy",
			req:         generic.Page(2),
			c8yRespCode: http.StatusOK,
			c8yRespBody: deviceinformation.ResponseBodyDeviceInformation,
			PagingStatics: ManagedObjectsPagingStatics{
				Statistics: generic.PagingStatistics{
					TotalRecords: 10,
					CurrentPage:  1,
					PageSize:     5,
				},
			},
			FragmentFilter: "<filter>",
		},
		{
			name:        "inventory collection request - bad credentials",
			c8yRespCode: http.StatusUnauthorized,
			c8yRespBody: `{}`,
			expectedErr: generic.BadCredentialsErr,
		},
		{
			name:        "managed object - access denied",
			c8yRespCode: http.StatusForbidden,
			c8yRespBody: `{}`,
			expectedErr: generic.AccessDeniedErr,
		},
		{
			name:        "inventory collection request - unexpected error",
			c8yRespCode: http.StatusInternalServerError,
			c8yRespBody: `{}`,
			expectedErr: errors.New("received an unexpected status code: 500"),
		},
		{
			name:        "inventory collection request - invalid response body",
			c8yRespCode: http.StatusOK,
			c8yRespBody: `<invalid response body>`,
			expectedErr: errors.New("failed to unmarshal response body: invalid character '<' looking for beginning of value"),
		},
	}

	for _, tt := range tests {
		var username, password, reqURL string
		t.Run(tt.name, func(t *testing.T) {
			testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				reqURL = req.URL.String()
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
			c := Client{
				HTTPClient: testServer.Client(),
				BaseURL:    testServer.URL,
				Username:   u,
				Password:   p,
			}

			manageObjectPagingStatics = func(c Client) (ManagedObjectsPagingStatics, error) {
				return tt.PagingStatics, nil
			}

			managedObjects, err := c.GetManagedObjects(tt.FragmentFilter)
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

			expectedRequestUrl := fmt.Sprintf("/inventory/managedObjects?currentPage=%v", tt.PagingStatics.Statistics.PageSize)
			if len(tt.FragmentFilter) > 0 {
				expectedRequestUrl += fmt.Sprintf("&fragmentType=%v", tt.FragmentFilter)
			}
			reqURL, err = url.QueryUnescape(reqURL)
			if err != nil {
				t.Error(err)
			}
			if reqURL != expectedRequestUrl {
				t.Errorf("unexpected request url: %v expected: %v", reqURL, expectedRequestUrl)
			}

			if tt.c8yRespCode == http.StatusOK && len(managedObjects) == 0 {
				t.Log("no managedobject found")
			}
		})
	}
}

func TestClient_NewManagedObjectPagingStatics(t *testing.T) {
	tests := []struct {
		name          string
		req           func(*http.Request)
		c8yRespCode   int
		c8yRespBody   string
		PagingStatics ManagedObjectsPagingStatics
		expectedErr   error
	}{
		{
			name:        "managed object paging statics - happy",
			req:         generic.Page(1),
			c8yRespCode: http.StatusOK,
			c8yRespBody: `{"statistics":{"totalPages":180,"currentPage":8,"pageSize":5}}`,
			PagingStatics: ManagedObjectsPagingStatics{
				Statistics: generic.PagingStatistics{
					TotalPages:  180,
					CurrentPage: 8,
					PageSize:    5,
				},
			},
		},
		{
			name:        "managed object - bad credentials",
			c8yRespCode: http.StatusUnauthorized,
			c8yRespBody: `{}`,
			expectedErr: generic.BadCredentialsErr,
		},
		{
			name:        "managed object - access denied",
			c8yRespCode: http.StatusForbidden,
			c8yRespBody: `{}`,
			expectedErr: generic.AccessDeniedErr,
		},
		{
			name:        "inventory collection request - unexpected error",
			c8yRespCode: http.StatusInternalServerError,
			c8yRespBody: `{}`,
			expectedErr: errors.New("received an unexpected status code: 500"),
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
			c := Client{
				HTTPClient: testServer.Client(),
				BaseURL:    testServer.URL,
				Username:   u,
				Password:   p,
			}

			pagingStatics, err := newManagedObjectPagingStatics(c)
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

			if pagingStatics.Statistics.PageSize != tt.PagingStatics.Statistics.PageSize {
				t.Errorf("pageSize is incorrect: %v expected: %v", pagingStatics.Statistics.PageSize, tt.PagingStatics.Statistics.PageSize)
			}
			if pagingStatics.Statistics.CurrentPage != tt.PagingStatics.Statistics.CurrentPage {
				t.Errorf("currentPage is incorrect: %v expected: %v", pagingStatics.Statistics.CurrentPage, tt.PagingStatics.Statistics.CurrentPage)
			}
			if pagingStatics.Statistics.TotalPages != tt.PagingStatics.Statistics.TotalPages {
				t.Errorf("totalPages are incorrect: %v expected: %v", pagingStatics.Statistics.TotalPages, tt.PagingStatics.Statistics.TotalPages)
			}
		})
	}
}

func TestClient_CreateManagedObject(t *testing.T) {
	tests := []struct {
		testName    string
		c8yRespCode int
		c8yRespBody string
		expectedErr error

		name         string
		state        string
		responseData models.NewManagedObject
	}{
		{
			testName:    "create managed object - happy",
			c8yRespCode: http.StatusCreated,
			c8yRespBody: `{"self":"<self>", "id":"<id>", "lastUpdated":"<lastUpdated>", "name":"<name>", "com_cumulocity_model_BinarySwitch":{"state":"<state>"}}`,

			name:  "<name>",
			state: "<state>",
			responseData: models.NewManagedObject{
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
			testName: "create managed object - invalid response body",
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
			c := Client{
				HTTPClient: testServer.Client(),
				BaseURL:    testServer.URL,
				Username:   u,
				Password:   p,
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
				t.Errorf("self is incorrect: %v expected: %v", newManagedObject.Self, tt.responseData.Self,)
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
