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
var inventoryRoleID = 2

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

var testInventoryRoleCollection = &InventoryRolesCollection{
	Self:       "https://t200588189.cumulocity.com/user/inventoryroles?pageSize=5&currentPage=2",
	Next:       "https://t200588189.cumulocity.com/user/inventoryroles?pageSize=5&currentPage=1",
	Prev:       "",
	Roles:      testInventoryRoles,
	Statistics: nil,
}

var testInventoryRoles = []InventoryRole{
	{
		Name:        "Reader",
		Description: "Can read all data of the asset and manage all inventory data, but cannot perform operations. Can also acknowledge and clear alarms. Can create and updates dashboards.",
		Self:        "https://t200588189.cumulocity.com/user/inventoryroles/2",
		ID:          2,
		Permissions: []Permission{
			{
				ID:         2,
				Permission: "READ",
				Type:       "*",
				Scope:      "*",
			},
			{
				ID:         3,
				Permission: "ADMIN",
				Type:       "*",
				Scope:      "MANAGED_OBJECT",
			},
		},
	},
	{
		ID:          2,
		Name:        "Manager",
		Description: "Can read all data of the asset and manage all inventory data, but cannot perform operations. Can also acknowledge and clear alarms. Can create and updates dashboards.",
		Self:        "https://t200588189.cumulocity.com/user/inventoryroles/2",
		Permissions: []Permission{
			{
				ID:         2,
				Permission: "READ",
				Type:       "*",
				Scope:      "*",
			},
			{
				ID:         3,
				Permission: "ADMIN",
				Type:       "*",
				Scope:      "MANAGED_OBJECT",
			},
		},
	},
}

var testInventoryRoleCollectionJSON = `
{
    "self": "https://t200588189.cumulocity.com/user/inventoryroles?pageSize=5&currentPage=2",
    "next": "https://t200588189.cumulocity.com/user/inventoryroles?pageSize=5&currentPage=1",
    "prev": "",
    "roles": [
        {
            "id": 2,
            "name": "Reader",
            "self": "https://t200588189.cumulocity.com/user/inventoryroles/2",
            "description": "Can read all data of the asset and manage all inventory data, but cannot perform operations. Can also acknowledge and clear alarms. Can create and updates dashboards.",
            "permissions": [
                {
                    "id": 2,
                    "type": "*",
                    "scope": "*",
                    "permission": "READ"
                },
                {
                    "id": 3,
                    "type": "*",
                    "scope": "MANAGED_OBJECT",
                    "permission": "ADMIN"
                }
            ]
        },
        {
            "id": 2,
            "name": "Manager",
            "self": "https://t200588189.cumulocity.com/user/inventoryroles/2",
            "description": "Can read all data of the asset and manage all inventory data, but cannot perform operations. Can also acknowledge and clear alarms. Can create and updates dashboards.",
            "permissions": [
                {
                    "id": 2,
                    "type": "*",
                    "scope": "*",
                    "permission": "READ"
                },
                {
                    "id": 3,
                    "type": "*",
                    "scope": "MANAGED_OBJECT",
                    "permission": "ADMIN"
                }
            ]
        }
    ],
    "statistics": null
}
`

var testInventoryRole = func(name string) *InventoryRole {
	return &InventoryRole{
		ID:          2,
		Name:        name,
		Self:        "https://t200588189.cumulocity.com/user/inventoryroles/2",
		Description: "Can read all data of the asset and manage all inventory data, but cannot perform operations. Can also acknowledge and clear alarms. Can create and updates dashboards.",
		Permissions: []Permission{
			{
				ID:         2,
				Permission: "READ",
				Type:       "*",
				Scope:      "*",
			},
			{
				ID:         3,
				Permission: "ADMIN",
				Type:       "*",
				Scope:      "MANAGED_OBJECT",
			},
		},
	}
}

var testInventoryRoleJSON = func(name string) string {
	return `
{
        "id": 2,
        "name": "` + name + `",
        "self": "https://t200588189.cumulocity.com/user/inventoryroles/2",
        "description": "Can read all data of the asset and manage all inventory data, but cannot perform operations. Can also acknowledge and clear alarms. Can create and updates dashboards.",
        "permissions": [
            {
                "id": 2,
                "type": "*",
                "scope": "*",
                "permission": "READ"
			},
			{
				"id": 3,
				"type": "*",
				"scope": "MANAGED_OBJECT",
				"permission": "ADMIN"
			}
		]
}
`
}

var TestInventoryRolesJSON = `
[
    {
        "id": 2,
        "name": "Reader",
        "self": "https://t200588189.cumulocity.com/user/inventoryroles/2",
        "description": "Can read all data of the asset and manage all inventory data, but cannot perform operations. Can also acknowledge and clear alarms. Can create and updates dashboards.",
        "permissions": [
            {
                "id": 2,
                "type": "*",
                "scope": "*",
                "permission": "READ"
            },
            {
                "id": 3,
                "type": "*",
                "scope": "MANAGED_OBJECT",
                "permission": "ADMIN"
            }
        ]
    },
    {
        "id": 2,
        "name": "Manager",
        "self": "https://t200588189.cumulocity.com/user/inventoryroles/2",
        "description": "Can read all data of the asset and manage all inventory data, but cannot perform operations. Can also acknowledge and clear alarms. Can create and updates dashboards.",
        "permissions": [
            {
                "id": 2,
                "type": "*",
                "scope": "*",
                "permission": "READ"
            },
            {
                "id": 3,
                "type": "*",
                "scope": "MANAGED_OBJECT",
                "permission": "ADMIN"
            }
        ]
    }
]
`
