package common_tests

import (
	"fmt"
	"github.com/tarent/gomulocity/generic"
	api "github.com/tarent/gomulocity/manage/user_api"
	"net/http"
	"net/http/httptest"
	"time"
)

func BuildUserApi(url string) api.UserApi {
	httpClient := http.DefaultClient
	client := &generic.Client{
		HTTPClient: httpClient,
		BaseURL:    url,
		Username:   "foo",
		Password:   "bar",
	}
	return api.NewUserApi(client)
}

func BuildHttpServer(status int, body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
}

var Create_ErroneousResponse = func(status int) *generic.Error {
	return &generic.Error{
		ErrorType: fmt.Sprintf("%v: userManagement/Forbidden", status),
		Message:   "authenticated user's tenant different from the one in URL path",
		Info:      "https://www.cumulocity.com/guides/reference-guide/#error_reporting",
	}
}

var Create_ErroneousResponseJSON = `
{
    "error": "userManagement/Forbidden",
    "message": "authenticated user's tenant different from the one in URL path",
    "info": "https://www.cumulocity.com/guides/reference-guide/#error_reporting"
}
`

var UserJSON = `
{
    "id":"msmith",
    "userName" : "msmith",
    "firstName" : "Michael",
    "lastName" : "Smith",
    "phone" : "+1234567890",
    "email" : "ms@abc.com",
    "enabled" : true,
    "groups":[
        {
            "id":"group1",
            "name":"group1",
            "roles":[
                {
                    "id":"role1",
                    "name":"role1"
                }
            ]
        }
    ],
     "roles":[
                {
                    "id":"role1",
                    "name":"role1"
                }
            ],
    "devicePermissions": null
  }
`

var TestUser = &api.User{
	ID:        "msmith",
	Username:  "msmith",
	FirstName: "Michael",
	LastName:  "Smith",
	Phone:     "+1234567890",
	Email:     "ms@abc.com",
	Enabled:   true,
	Groups: []api.Group{
		{
			ID:   "group1",
			Name: "group1",
			Roles: []api.Role{
				{
					ID:   "role1",
					Name: "role1",
				},
			},
		},
	},
	Roles: []api.Role{
		{
			ID:   "role1",
			Name: "role1",
		},
	},
	DevicePermissions: nil,
}

var TenantID = "1111111"

var CreateUser = &api.CreateUser{
	Username:  "username",
	Password:  "password",
	FirstName: "Michael",
	LastName:  "Smith",
	Phone:     "+1234567890",
	Email:     "ms@abc.com",
	Enabled:   true,
}

var lastPasswordChange = func(value string) time.Time {
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		fmt.Println(err)
		return time.Time{}
	}
	return t
}

var TestCurrentUser = api.CurrentUser{
	ID:                "msmith",
	Self:              "https://t200588189.cumulocity.com/user/currentUser",
	Username:          "msmith",
	FirstName:         "Michael",
	LastName:          "Smith",
	Phone:             "+1234567890",
	Email:             "ms@abc.com",
	Enabled:           true,
	DevicePermissions: nil,
	EffectiveRoles: []api.Role{
		{
			ID:   "role1",
			Name: "role1",
		},
		{
			ID:   "role2",
			Name: "role2",
		},
		{
			ID:   "role3",
			Name: "role3",
		},
	},
	ShouldResetPassword: false,
	LastPasswordChange:  lastPasswordChange("2020-05-25T11:52:56.999Z"),
}

var CurrentUserJSON = `{
    "id": "msmith",
    "self": "https://t200588189.cumulocity.com/user/currentUser",
    "phone": "+1234567890",
    "email": "ms@abc.com",
    "enabled": true,
    "userName": "msmith",
    "lastName": "Smith",
    "firstName": "Michael",
    "effectiveRoles": [
        {
            "id": "role1",
            "name": "role1"
        },
        {
            "id": "role2",
            "name": "role2"
        },
        {
            "id": "role3",
            "name": "role3"
        }
    ],
    "devicePermissions": null,
    "lastPasswordChange": "2020-05-25T11:52:56.999Z",
    "shouldResetPassword": false
}
`
var CurrentUserErroneousJSON = `
{
	"error": "security/Unauthorized",
	"message": "Full authentication is required to access this resource",
	"info": "https://www.cumulocity.com/guides/reference-guide/#error_reporting"
}`

var CurrentUser_ErroneousResponse = func(status int) *generic.Error {
	return &generic.Error{
		ErrorType: fmt.Sprintf("%v: security/Unauthorized", status),
		Message:   "Full authentication is required to access this resource",
		Info:      "https://www.cumulocity.com/guides/reference-guide/#error_reporting",
	}
}
