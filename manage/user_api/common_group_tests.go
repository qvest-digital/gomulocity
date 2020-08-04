package user_api

import (
	"fmt"
	"github.com/tarent/gomulocity/generic"
)

var testGroup = &Group{
	ID:                groupID,
	Self:              "https://t200588189.cumulocity.com/user/management/groups/" + groupID,
	Name:              groupName,
	Roles:             roles,
	DevicePermissions: nil,
}

var groupJSON = `
{
	"id":"12",
	"self":"https://t200588189.cumulocity.com/user/management/groups/12",
	"name":"TEST_GROUP",
	"roles":[
		{
			"id":   "ROLE_ACCOUNT_ADMIN",
			"name": "ROLE_ACCOUNT_ADMIN",
			"self": "https://t200588189.cumulocity.com/user/roles/ROLE_ACCOUNT_ADMIN"
		},
		{
			"id":   "ROLE_ALARM_ADMIN",
			"name": "ROLE_ALARM_ADMIN",
			"self": "https://t200588189.cumulocity.com/user/roles/ROLE_ALARM_ADMIN"
		},
		{
			"id":   "ROLE_ALARM_READ",
			"name": "ROLE_ALARM_READ",
			"self": "https://t200588189.cumulocity.com/user/roles/ROLE_ALARM_READ"
		}
	],
	"devicePermissions":null
}
`

var groupJSON_WithGroupName = func(groupName string) string {
	return `
{
	"id":"12",
	"self":"https://t200588189.cumulocity.com/user/management/groups/12",
	"name":"` + groupName + `",
	"roles":[
		{
			"id":   "ROLE_ACCOUNT_ADMIN",
			"name": "ROLE_ACCOUNT_ADMIN",
			"self": "https://t200588189.cumulocity.com/user/roles/ROLE_ACCOUNT_ADMIN"
		},
		{
			"id":   "ROLE_ALARM_ADMIN",
			"name": "ROLE_ALARM_ADMIN",
			"self": "https://t200588189.cumulocity.com/user/roles/ROLE_ALARM_ADMIN"
		},
		{
			"id":   "ROLE_ALARM_READ",
			"name": "ROLE_ALARM_READ",
			"self": "https://t200588189.cumulocity.com/user/roles/ROLE_ALARM_READ"
		}
	],
	"devicePermissions":null
}
`
}

var testGroupReferenceCollection = &GroupReferenceCollection{
	Self:   fmt.Sprintf("https://t200588189.cumulocity.com/user/%v/users/%v/groups?pageSize=%v&currentPage=%v", tenantID, username, 5, 2),
	Next:   fmt.Sprintf("https://t200588189.cumulocity.com/user/%v/users/%v/groups?pageSize=%v&currentPage=%v", tenantID, username, 5, 3),
	Prev:   fmt.Sprintf("https://t200588189.cumulocity.com/user/%v/users/%v/groups?pageSize=%v&currentPage=%v", tenantID, username, 5, 1),
	Groups: groups,
	Statistics: &generic.PagingStatistics{
		TotalPages:  3,
		PageSize:    5,
		CurrentPage: 2,
	},
}

var groups = []Group{
	{
		ID:                "12",
		Self:              "https://t200588189.cumulocity.com/user/management/groups/12",
		Name:              "TEST_GROUP_12",
		Roles:             roles,
		DevicePermissions: nil,
	},
	{
		ID:                "13",
		Self:              "https://t200588189.cumulocity.com/user/management/groups/13",
		Name:              "TEST_GROUP_13",
		Roles:             roles,
		DevicePermissions: nil,
	},
	{
		ID:                "14",
		Self:              "https://t200588189.cumulocity.com/user/management/groups/14",
		Name:              "TEST_GROUP_14",
		Roles:             roles,
		DevicePermissions: nil,
	},
}

var groupReferenceCollectionJSON = `
{
    "self": "https://t200588189.cumulocity.com/user/1111111/users/msmith/groups?pageSize=5&currentPage=2",
    "next": "https://t200588189.cumulocity.com/user/1111111/users/msmith/groups?pageSize=5&currentPage=3",
    "prev": "https://t200588189.cumulocity.com/user/1111111/users/msmith/groups?pageSize=5&currentPage=1",
    "groups": [
        {
            "id": "12",
            "self": "https://t200588189.cumulocity.com/user/management/groups/12",
            "name": "TEST_GROUP_12",
            "roles": [
                {
                    "id": "ROLE_ACCOUNT_ADMIN",
                    "self": "https://t200588189.cumulocity.com/user/roles/ROLE_ACCOUNT_ADMIN",
                    "name": "ROLE_ACCOUNT_ADMIN"
                },
                {
                    "id": "ROLE_ALARM_ADMIN",
                    "self": "https://t200588189.cumulocity.com/user/roles/ROLE_ALARM_ADMIN",
                    "name": "ROLE_ALARM_ADMIN"
                },
                {
                    "id": "ROLE_ALARM_READ",
                    "self": "https://t200588189.cumulocity.com/user/roles/ROLE_ALARM_READ",
                    "name": "ROLE_ALARM_READ"
                }
            ],
            "devicePermissions": null
        },
        {
            "id": "13",
            "self": "https://t200588189.cumulocity.com/user/management/groups/13",
            "name": "TEST_GROUP_13",
            "roles": [
                {
                    "id": "ROLE_ACCOUNT_ADMIN",
                    "self": "https://t200588189.cumulocity.com/user/roles/ROLE_ACCOUNT_ADMIN",
                    "name": "ROLE_ACCOUNT_ADMIN"
                },
                {
                    "id": "ROLE_ALARM_ADMIN",
                    "self": "https://t200588189.cumulocity.com/user/roles/ROLE_ALARM_ADMIN",
                    "name": "ROLE_ALARM_ADMIN"
                },
                {
                    "id": "ROLE_ALARM_READ",
                    "self": "https://t200588189.cumulocity.com/user/roles/ROLE_ALARM_READ",
                    "name": "ROLE_ALARM_READ"
                }
            ],
            "devicePermissions": null
        },
        {
            "id": "14",
            "self": "https://t200588189.cumulocity.com/user/management/groups/14",
            "name": "TEST_GROUP_14",
            "roles": [
                {
                    "id": "ROLE_ACCOUNT_ADMIN",
                    "self": "https://t200588189.cumulocity.com/user/roles/ROLE_ACCOUNT_ADMIN",
                    "name": "ROLE_ACCOUNT_ADMIN"
                },
                {
                    "id": "ROLE_ALARM_ADMIN",
                    "self": "https://t200588189.cumulocity.com/user/roles/ROLE_ALARM_ADMIN",
                    "name": "ROLE_ALARM_ADMIN"
                },
                {
                    "id": "ROLE_ALARM_READ",
                    "self": "https://t200588189.cumulocity.com/user/roles/ROLE_ALARM_READ",
                    "name": "ROLE_ALARM_READ"
                }
            ],
            "devicePermissions": null
        }
    ],
    "statistics": {
        "pageSize": 5,
        "totalPages": 3,
        "currentPage": 2
    }
}
`

var groupCollectionTemplate = `{
    "next": "https://t200588189.cumulocity.com/user/1111111/users/msmith/groups?pageSize=5&currentPage=2",
    "self": "https://t200588189.cumulocity.com/user/1111111/users/msmith/groups?pageSize=5&currentPage=1",
    "groups": [%s],
    "statistics": {
        "currentPage": 1,
        "pageSize": 5
    }
}`
