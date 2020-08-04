package user_api

import (
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"github.com/tarent/gomulocity/manage/user_api/roles"
	"net/url"
	"time"
)

const (
	USER_CONTENT_TYPE = "application/vnd.com.nsn.cumulocity.user+json;ver=0.9"
	USER_ACCEPT       = "application/vnd.com.nsn.cumulocity.user+json;ver=0.9"
)

type User struct {
	ID       string `json:"id"`
	Self     string `json:"self"`
	Username string `json:"userName"`
	//Password          string      `json:"password"`
	FirstName         string       `json:"firstName"`
	LastName          string       `json:"lastName"`
	Phone             string       `json:"phone"`
	Email             string       `json:"email"`
	Enabled           bool         `json:"enabled"`
	Groups            []Group      `json:"groups"`
	Roles             []roles.Role `json:"roles"`
	DevicePermissions interface{}  `json:"devicePermissions"`
}

type CreateUser struct {
	Username          string      `json:"userName"`
	Password          string      `json:"password"`
	FirstName         string      `json:"firstName"`
	LastName          string      `json:"lastName"`
	Phone             string      `json:"phone"`
	Email             string      `json:"email"`
	Enabled           bool        `json:"enabled"`
	DevicePermissions interface{} `json:"devicePermissions"`
}

type UserCollection struct {
	Self       string                    `json:"self"`
	Users      []User                    `json:"userApi"`
	Statistics *generic.PagingStatistics `json:"statistics"`
	Prev       string                    `json:"prev"`
	Next       string                    `json:"next"`
}

type CurrentUser struct {
	ID                  string       `json:"id"`
	Self                string       `json:"self"`
	Username            string       `json:"userName"`
	FirstName           string       `json:"firstName"`
	LastName            string       `json:"lastName"`
	Phone               string       `json:"phone"`
	Email               string       `json:"email"`
	Enabled             bool         `json:"enabled"`
	DevicePermissions   interface{}  `json:"devicePermissions"`
	EffectiveRoles      []roles.Role `json:"effectiveRoles"`
	ShouldResetPassword bool         `json:"shouldResetPassword"`
	LastPasswordChange  time.Time    `json:"lastPasswordChange"`
}

type Group struct {
	ID                string       `json:"id"`
	Self              string       `json:"self"`
	Name              string       `json:"name"`
	Roles             []roles.Role `json:"roles"`
	DevicePermissions struct{}     `json:"devicePermissions"`
}

type QueryFilter struct {
	Username          string
	Groups            []Group
	Owner             string
	OnlyDevices       bool
	WithSubUsersCount bool
}

func (q QueryFilter) QueryParams(params *url.Values) error {
	if params == nil {
		return fmt.Errorf("The provided parameter values must not be nil!")
	}

	if len(q.Username) > 0 {
		params.Add("username", q.Username)
	}

	if len(q.Owner) > 0 {
		params.Add("owner", q.Owner)
	}

	if q.OnlyDevices {
		params.Add("onlyDevices", "true")
	}

	if q.WithSubUsersCount {
		params.Add("withSubusersCount", "true")
	}
	return nil
}

func (q QueryFilter) addGroups(query *url.Values) string {
	var groups string

	if len(q.Groups) > 0 {
		prefix := ""
		if len(*query) != 0 {
			prefix = "&"
		}
		groups = fmt.Sprintf("%vgroups=", prefix)
		for _, group := range q.Groups {
			groups += fmt.Sprintf("%v,", group.Name)
		}
	}
	return query.Encode() + groups[:len(groups)-1]
}

func (u User) HasDevicePermissions() bool {
	return u.DevicePermissions != nil
}
