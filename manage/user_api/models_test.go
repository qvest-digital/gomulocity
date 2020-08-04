package user_api

import (
	"reflect"
	"testing"
)

func TestQueryFilter_QueryParams(t *testing.T) {
	filter := QueryFilter{
		Username: "username",
		Groups: []Group{
			{
				Name: "group1",
			},
			{
				Name: "group2",
			},
		},
		Owner:             "owner",
		OnlyDevices:       true,
		WithSubUsersCount: true,
	}

	query := filter.QueryParams()

	expectedQuery := "onlyDevices=true&owner=owner&username=username&withSubusersCount=true&groups=group1,group2"
	if query != expectedQuery {
		t.Errorf("Unexpected query: expected: %v actual: %v", expectedQuery, query)
	}
}

func TestQueryFilter_QueryParams_OnlyGroup(t *testing.T) {
	filter := QueryFilter{
		Groups: []Group{
			{
				Name: "group1",
			},
			{
				Name: "group2",
			},
		},
	}

	query := filter.QueryParams()

	expectedQuery := "groups=group1,group2"
	if query != expectedQuery {
		t.Errorf("Unexpected query: expected: %v actual: %v", expectedQuery, query)
	}
}

func TestUser_HasDevicePermissions_happy(t *testing.T) {
	user := User{
		DevicePermissions: reflect.Interface,
	}

	if !user.HasDevicePermissions() {
		t.Error("user does not have permissions")
	}
}