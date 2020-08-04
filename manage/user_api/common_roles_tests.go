package user_api

import "github.com/tarent/gomulocity/generic"

var RolesCollection = RoleCollection{
	Self: "https://t200588189.cumulocity.com/user/roles?pageSize=5&currentPage=1",
	Next: "https://t200588189.cumulocity.com/user/roles?pageSize=5&currentPage=2",
	Prev: "",
	Roles: []Role{
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
	},
	Statistics: &generic.PagingStatistics{
		TotalRecords: 0,
		TotalPages:   10,
		PageSize:     5,
		CurrentPage:  1,
	},
}

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
