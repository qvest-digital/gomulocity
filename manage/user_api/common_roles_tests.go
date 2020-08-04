package user_api

import (
	"github.com/tarent/gomulocity/generic"
)

var roleCollection = &RoleCollection{
	Self:  "https://t200588189.cumulocity.com/user/roles?pageSize=5&currentPage=1",
	Next:  "https://t200588189.cumulocity.com/user/roles?pageSize=5&currentPage=2",
	Prev:  "",
	Roles: roles,
	Statistics: &generic.PagingStatistics{
		TotalRecords: 0,
		TotalPages:   10,
		PageSize:     5,
		CurrentPage:  1,
	},
}

var roleID = "ROLE_ACCOUNT_ADMIN"

var roles = []Role{
	{
		ID:   "ROLE_ACCOUNT_ADMIN",
		Name: "ROLE_ACCOUNT_ADMIN",
		Self: "https://t200588189.cumulocity.com/user/roles/ROLE_ACCOUNT_ADMIN",
	},
	{
		ID:   "ROLE_ALARM_ADMIN",
		Name: "ROLE_ALARM_ADMIN",
		Self: "https://t200588189.cumulocity.com/user/roles/ROLE_ALARM_ADMIN",
	},
	{
		ID:   "ROLE_ALARM_READ",
		Name: "ROLE_ALARM_READ",
		Self: "https://t200588189.cumulocity.com/user/roles/ROLE_ALARM_READ",
	},
}

var roleJSON = `
 {
	"id": "ROLE_ACCOUNT_ADMIN",
	"name": "ROLE_ACCOUNT_ADMIN",
	"self": "https://t200588189.cumulocity.com/user/roles/ROLE_ACCOUNT_ADMIN"
}
`

var roleCollectionJSON = `{
    "next": "https://t200588189.cumulocity.com/user/roles?pageSize=5&currentPage=2",
    "self": "https://t200588189.cumulocity.com/user/roles?pageSize=5&currentPage=1",
    "roles": [
        {
            "id": "ROLE_ACCOUNT_ADMIN",
            "name": "ROLE_ACCOUNT_ADMIN",
            "self": "https://t200588189.cumulocity.com/user/roles/ROLE_ACCOUNT_ADMIN"
        },
        {
            "id": "ROLE_ALARM_ADMIN",
            "name": "ROLE_ALARM_ADMIN",
            "self": "https://t200588189.cumulocity.com/user/roles/ROLE_ALARM_ADMIN"
        },
        {
            "id": "ROLE_ALARM_READ",
            "name": "ROLE_ALARM_READ",
            "self": "https://t200588189.cumulocity.com/user/roles/ROLE_ALARM_READ"
        }
    ],
    "statistics": {
        "pageSize": 5,
        "totalPages": 10,
        "currentPage": 1
    }
}`

var roleReferenceCollection = &RoleReferenceCollection{
	Self: "https://t200588189.cumulocity.com/user/" + tenantID + "/users/" + username + "/roles?pageSize=5&currentPage=2",
	Next: "https://t200588189.cumulocity.com/user/" + tenantID + "/users/" + username + "/roles?pageSize=5&currentPage=3",
	Prev: "https://t200588189.cumulocity.com/user/" + tenantID + "/users/" + username + "/roles?pageSize=5&currentPage=1",
	References: []RoleReference{
		{
			Role: Role{
				ID:   "ROLE_ACCOUNT_ADMIN",
				Name: "ROLE_ACCOUNT_ADMIN",
				Self: "https://t200588189.cumulocity.com/user/roles/ROLE_ACCOUNT_ADMIN",
			},
		},
		{
			Role: Role{
				ID:   "ROLE_ALARM_ADMIN",
				Name: "ROLE_ALARM_ADMIN",
				Self: "https://t200588189.cumulocity.com/user/roles/ROLE_ALARM_ADMIN",
			},
		},
	},
	Statistics: &generic.PagingStatistics{
		TotalRecords: 0,
		TotalPages:   10,
		PageSize:     5,
		CurrentPage:  1,
	},
}

var roleReferenceCollectionJSON = `{
    "self": "https://t200588189.cumulocity.com/user/1111111/users/msmith/roles?pageSize=5&currentPage=2",
    "next": "https://t200588189.cumulocity.com/user/1111111/users/msmith/roles?pageSize=5&currentPage=3",
    "prev": "https://t200588189.cumulocity.com/user/1111111/users/msmith/roles?pageSize=5&currentPage=1",
    "references": [
        {
            "role": {
                "id": "ROLE_ACCOUNT_ADMIN",
                "name": "ROLE_ACCOUNT_ADMIN",
                "self": "https://t200588189.cumulocity.com/user/roles/ROLE_ACCOUNT_ADMIN"
            }
        },
        {
            "role": {
                "id": "ROLE_ALARM_ADMIN",
                "name": "ROLE_ALARM_ADMIN",
                "self": "https://t200588189.cumulocity.com/user/roles/ROLE_ALARM_ADMIN"
            }
        }
    ],
    "statistics": {
        "pageSize": 5,
        "totalPages": 10,
        "currentPage": 1
    }
}`

var roleReference = &RoleReference{
	Self: "",
	Role: Role{
		ID:   "ROLE_ACCOUNT_ADMIN",
		Name: "ROLE_ACCOUNT_ADMIN",
		Self: "https://t200588189.cumulocity.com/user/roles/ROLE_ACCOUNT_ADMIN",
	},
}

var roleReferenceJSON = `{
	"self":"",
	"role": {
		"id":"ROLE_ACCOUNT_ADMIN",
		"name":"ROLE_ACCOUNT_ADMIN",
		"self":"https://t200588189.cumulocity.com/user/roles/ROLE_ACCOUNT_ADMIN"
	}
}
`

var roleCollectionTemplate = `{
    "next": "https://t200588189.cumulocity.com/user/roles?pageSize=5&currentPage=2",
    "self": "https://t200588189.cumulocity.com/user/roles?pageSize=5&currentPage=1",
    "roles": [%s],
    "statistics": {
        "currentPage": 1,
        "pageSize": 5
    }
}`
