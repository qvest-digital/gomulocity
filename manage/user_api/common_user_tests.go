package user_api

import (
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"net/http"
	"net/http/httptest"
	"time"
)

func buildUserApi(url string) UserApi {
	httpClient := http.DefaultClient
	client := &generic.Client{
		HTTPClient: httpClient,
		BaseURL:    url,
		Username:   "foo",
		Password:   "bar",
	}
	return NewUserApi(client)
}

func buildHttpServer(status int, body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
}

var createErroneousResponse = func(status int) *generic.Error {
	return &generic.Error{
		ErrorType: fmt.Sprintf("%v: userManagement/Forbidden", status),
		Message:   "authenticated user's tenant different from the one in URL path",
		Info:      "https://www.cumulocity.com/guides/reference-guide/#error_reporting",
	}
}

var erroneousResponseJSON = `
{
    "error": "userManagement/Forbidden",
    "message": "authenticated user's tenant different from the one in URL path",
    "info": "https://www.cumulocity.com/guides/reference-guide/#error_reporting"
}
`

var testUserJSON = `
{
    "id":"msmith",
	"self":"https://t200588189.cumulocity.com/user/msmith",
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

var testUser = &User{
	ID:        "msmith",
	Self:      "https://t200588189.cumulocity.com/user/msmith",
	Username:  "msmith",
	FirstName: "Michael",
	LastName:  "Smith",
	Phone:     "+1234567890",
	Email:     "ms@abc.com",
	Enabled:   true,
	Groups: []Group{
		{
			ID:   "group1",
			Name: "group1",
			Roles: []Role{
				{
					ID:   "role1",
					Name: "role1",
				},
			},
		},
	},
	Roles: []Role{
		{
			ID:   "role1",
			Name: "role1",
		},
	},
	DevicePermissions: nil,
}

var (
	tenantID  = "1111111"
	userID    = "msmith"
	username  = "msmith"
	roleName  = "ROLE_ACCOUNT_ADMIN"
	groupID   = "12"
	groupName = "TEST_GROUP"
)

var createUser = &CreateUser{
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

var testCurrentUser = &CurrentUser{
	ID:                "msmith",
	Self:              "https://t200588189.cumulocity.com/user/currentUser",
	Username:          "msmith",
	FirstName:         "Michael",
	LastName:          "Smith",
	Phone:             "+1234567890",
	Email:             "ms@abc.com",
	Enabled:           true,
	DevicePermissions: nil,
	EffectiveRoles: []Role{
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

var currentUserJSON = `{
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
var currentUserErroneousJSON = `
{
	"error": "security/Unauthorized",
	"message": "Full authentication is required to access this resource",
	"info": "https://www.cumulocity.com/guides/reference-guide/#error_reporting"
}`

var currentUser_ErroneousResponse = func(status int) *generic.Error {
	return &generic.Error{
		ErrorType: fmt.Sprintf("%v: security/Unauthorized", status),
		Message:   "Full authentication is required to access this resource",
		Info:      "https://www.cumulocity.com/guides/reference-guide/#error_reporting",
	}
}

var testUserCollection = &UserCollection{
	Self: "https://t200588189.cumulocity.com/user/1111111/users?pageSize=2&username=mmark&groups=group1,group3&currentPage=1",
	Users: []User{
		{
			ID:        "msmith",
			Username:  "msmith",
			FirstName: "Michael",
			LastName:  "Smith",
			Phone:     "+1234567890",
			Email:     "ms@abc.com",
			Enabled:   true,
			Groups: []Group{
				{
					ID:   "group1",
					Name: "group1",
					Roles: []Role{
						{
							ID:   "role1",
							Name: "role1",
						},
					},
				},
			},
			Roles: []Role{
				{
					ID:   "role1",
					Name: "role1",
				},
			},
			DevicePermissions: nil,
		},
		{
			ID:        "mmark",
			Username:  "mmark",
			FirstName: "Manuel",
			LastName:  "Mark",
			Phone:     "+1234567890",
			Email:     "mm@abc.com",
			Enabled:   true,
			Groups: []Group{
				{
					ID:   "group3",
					Name: "group3",
					Roles: []Role{
						{
							ID:   "role2",
							Name: "role2",
						},
					},
				},
			},
			Roles: []Role{
				{
					ID:   "role4",
					Name: "role4",
				},
			},
			DevicePermissions: nil,
		},
	},
	Statistics: nil,
	Prev:       "",
	Next:       "https://t200588189.cumulocity.com/user/1111111/users?pageSize=2&username=mmark&groups=group1,group3&currentPage=2",
}

var userCollectionJSON = `
{
    "self": "https://t200588189.cumulocity.com/user/1111111/users?pageSize=2&username=mmark&groups=group1,group3&currentPage=1",
    "next": "https://t200588189.cumulocity.com/user/1111111/users?pageSize=2&username=mmark&groups=group1,group3&currentPage=2",
    "prev": "",
    "users": [
        {
            "id": "msmith",
			"userName": "msmith",
            "lastName": "Smith",
            "firstName": "Michael",
            "phone": "+1234567890",
            "email": "ms@abc.com",
			"enabled": true,
            "roles": [
                {
                    "id": "role1",
                    "name": "role1"
                }
            ],
            "groups": [
                {
                    "id": "group1",
                    "name": "group1",
                    "roles": [
                        {
                            "id": "role1",
                            "name": "role1"
                        }
                    ]
                }
            ],
            "devicePermissions": null
        },
        {
            "id": "mmark",
			"userName": "mmark",
            "firstName": "Manuel",
			"lastName": "Mark",
            "phone": "+1234567890",
            "email": "mm@abc.com",
			"enabled": true,
            "roles": [
                {
                    "id": "role4",
                    "name": "role4"
                }
            ],
            "groups": [
                {
                    "id": "group3",
                    "name": "group3",
                    "roles": [
                        {
                            "id": "role2",
                            "name": "role2"
                        }
                    ]
                }
            ],
            "devicePermissions": null
        }
    ]
}
`
var userCollectionTemplate = `{
    "next": "https://t200588189.cumulocity.com/user/1111111/users?username=mmark&pageSize=2&currentPage=2",
    "self": "https://t200588189.cumulocity.com/user/1111111/users?username=mmark&pageSize=2&currentPage=1",
    "users": [%s],
    "statistics": {
        "currentPage": 1,
        "pageSize": 5
    }
}`
