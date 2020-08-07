package user_api

import (
	"encoding/json"
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"log"
	"net/http"
	"net/url"
)

type UserApi interface {
	CreateUser(tenantID string, model *CreateUser) (*User, *generic.Error)
	UserCollection(filter *QueryFilter, pageSize int) (*UserCollection, *generic.Error)
	GetCurrentUser() (*CurrentUser, *generic.Error)
	UserByName(tenantID, username string) (*User, *generic.Error)
	FindUserCollection(userQuery *QueryFilter, pageSize int) (*UserCollection, *generic.Error)
	NextPageUserCollection(r *UserCollection) (*UserCollection, *generic.Error)
	PreviousPageUserCollection(r *UserCollection) (*UserCollection, *generic.Error)

	RoleCollection(pageSize int) (*RoleCollection, *generic.Error)
	FindRoleCollection(pageSize int) (*RoleCollection, *generic.Error)
	NextPageRoleCollection(r *RoleCollection) (*RoleCollection, *generic.Error)
	PreviousPageRoleCollection(r *RoleCollection) (*RoleCollection, *generic.Error)
	FindRoleReferenceCollection(tenantID, username, groupID string, pageSize int) (*RoleReferenceCollection, *generic.Error)
	AssignRoleToUser(tenantID, username string, reference *RoleReference) (*RoleReference, *generic.Error)
	AssignRoleToGroup(tenantID, groupID string, reference *RoleReference) (*RoleReference, *generic.Error)
	UnassignRoleFromUser(tenantID, username, roleName string) *generic.Error
	UnassignRoleFromGroup(tenantID, groupID, roleName string) *generic.Error
	GetAllRolesOfAUser(tenantID, username string, pageSize int) (*RoleReferenceCollection, *generic.Error)
	GetAllRolesOfAGroup(tenantID, groupID string, pageSize int) (*RoleReferenceCollection, *generic.Error)

	GroupDetails(groupID string) (*Group, *generic.Error)
	GroupByName(tenantID, groupName string) (*Group, *generic.Error)
	RemoveGroup(tenantID, groupID string) *generic.Error
	UpdateGroup(tenantID, groupID string, group *Group) (*Group, *generic.Error)
	GetAllGroupsOfUser(tenantID, username string, pageSize int) (*GroupReferenceCollection, *generic.Error)
	FindGroupReferenceCollection(tenantID, username string, pageSize int) (*GroupReferenceCollection, *generic.Error)
	NextPageGroupReferenceCollection(r *GroupReferenceCollection) (*GroupReferenceCollection, *generic.Error)
	PreviousPageGroupCollection(r *GroupReferenceCollection) (*GroupReferenceCollection, *generic.Error)
}

func NewUserApi(client *generic.Client) UserApi {
	return &userApi{
		client:   client,
		basePath: "/user",
	}
}

type userApi struct {
	client   *generic.Client
	basePath string
}

/*
TBD:
Get all available inventory roles
Assign a new inventory role
Retrieve an inventory role
Update an inventory role
Delete an inventory role
*/

func (u *userApi) CreateUser(tenantID string, model *CreateUser) (*User, *generic.Error) {
	if len(tenantID) == 0 {
		return nil, generic.ClientError("Creating user without a tenantID is not allowed", "CreateUser")
	}

	bytes, err := json.Marshal(model)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while marshalling user model: %s", err), "CreateUser")
	}

	body, status, err := u.client.Post(fmt.Sprintf("%v/%v/userApi", u.basePath, tenantID), bytes, generic.ContentTypeHeaderAndContentLength(USER_CONTENT_TYPE, len(bytes)))
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while posting a new user: %s", err), "CreateUser")
	}

	if status != http.StatusCreated {
		return nil, generic.CreateErrorFromResponse(body, status)
	}

	user := &User{}
	if err := json.Unmarshal(body, user); err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while unmarshalling response body: %s", err), "CreateUser")
	}
	return user, nil
}

func (u *userApi) UserCollection(filter *QueryFilter, pageSize int) (*UserCollection, *generic.Error) {
	return u.FindUserCollection(filter, pageSize)
}

func (u *userApi) GetCurrentUser() (*CurrentUser, *generic.Error) {
	body, status, err := u.client.Get(fmt.Sprintf("%v/currentUser", u.basePath), generic.AcceptHeader(USER_ACCEPT))
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while getting current user data: %s", err), "GetCurrentUser")
	}

	if status != http.StatusOK {
		return nil, generic.CreateErrorFromResponse(body, status)
	}

	user := &CurrentUser{}
	if err := json.Unmarshal(body, user); err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while unmarshalling response body: %s", err), "GetCurrentUser")
	}
	return user, nil
}

func (u *userApi) UserByName(tenantID, username string) (*User, *generic.Error) {
	if len(tenantID) == 0 || len(username) == 0 {
		return nil, generic.ClientError("Getting user without a tenantID or username is not allowed", "UserByName")
	}

	body, status, err := u.client.Get(fmt.Sprintf("%v/%v/userByName/%v", u.basePath, tenantID, username), generic.AcceptHeader(USER_ACCEPT))
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while getting user %v by name: %s", username, err), "UserByName")
	}

	if status != http.StatusOK {
		return nil, generic.CreateErrorFromResponse(body, status)
	}

	user := &User{}
	if err := json.Unmarshal(body, user); err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while unmarshalling response body: %s", err), "UserByName")
	}
	return user, nil
}

func (u *userApi) FindUserCollection(filter *QueryFilter, pageSize int) (*UserCollection, *generic.Error) {
	queryParamsValues := &url.Values{}

	if filter == nil {
		return nil, generic.ClientError("Given filter is empty", "FindUserCollection")
	}
	err := filter.QueryParams(queryParamsValues)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while building query parameters to search for measurements: %s", err.Error()), "FindUserCollection")
	}

	err = generic.PageSizeParameter(pageSize, queryParamsValues)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while building pageSize parameter to fetch measurements: %s", err.Error()), "FindUserCollection")
	}
	queryWithGroups := filter.addGroups(queryParamsValues)
	return u.getCommonUserCollection(fmt.Sprintf("%s?%s", u.basePath, queryWithGroups))
}

func (u *userApi) NextPageUserCollection(r *UserCollection) (*UserCollection, *generic.Error) {
	return u.getPageUserCollection(r.Next)
}

func (u *userApi) PreviousPageUserCollection(r *UserCollection) (*UserCollection, *generic.Error) {
	return u.getPageUserCollection(r.Prev)
}

func (u *userApi) getPageUserCollection(reference string) (*UserCollection, *generic.Error) {
	if reference == "" {
		log.Print("No page reference given. Returning nil.")
		return nil, nil
	}

	nextUrl, err := url.Parse(reference)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Unparsable URL given for page reference: '%s'", reference), "GetPage")
	}

	collection, genErr := u.getCommonUserCollection(fmt.Sprintf("%s?%s", nextUrl.Path, nextUrl.RawQuery))
	if genErr != nil {
		return nil, genErr
	}

	if len(collection.Users) == 0 {
		log.Print("Returned collection is empty. Returning nil.")
		return nil, nil
	}

	return collection, nil
}

func (u *userApi) getCommonUserCollection(path string) (*UserCollection, *generic.Error) {
	body, status, err := u.client.Get(path, generic.AcceptHeader(USER_ACCEPT))
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while getting measurements: %s", err.Error()), "GetMeasurementCollection")
	}

	if status != http.StatusOK {
		return nil, generic.CreateErrorFromResponse(body, status)
	}

	return parseUserCollectionResponse(body)
}

func parseUserCollectionResponse(body []byte) (*UserCollection, *generic.Error) {
	var result UserCollection
	if len(body) > 0 {
		err := generic.ObjectFromJson(body, &result)
		if err != nil {
			return nil, generic.ClientError(fmt.Sprintf("Error while parsing response JSON: %s", err.Error()), "CollectionResponseParser")
		}
	} else {
		return nil, generic.ClientError("Response body was empty", "CollectionResponseParser")
	}

	return &result, nil
}

// Roles

func (u *userApi) FindRoleCollection(pageSize int) (*RoleCollection, *generic.Error) {
	queryParamsValues := &url.Values{}
	err := generic.PageSizeParameter(pageSize, queryParamsValues)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while building pageSize parameter to fetch role collection: %s", err.Error()), "FindRoleCollection")
	}
	return u.getCommonRoleCollection(fmt.Sprintf("%s/roles?%v", u.basePath, queryParamsValues.Encode()))
}

func (u *userApi) FindRoleReferenceCollection(tenantID, username, groupID string, pageSize int) (*RoleReferenceCollection, *generic.Error) {
	queryParamsValues := &url.Values{}
	err := generic.PageSizeParameter(pageSize, queryParamsValues)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while building pageSize parameter to fetch reference collection: %s", err.Error()), "FindRoleReferenceCollection")
	}

	if len(username) > 0 {
		return u.getCommonRoleReferenceCollection(fmt.Sprintf("%s/%s/users/%s/roles?%s", u.basePath, tenantID, username, queryParamsValues.Encode()))
	} else if len(groupID) > 0 {
		return u.getCommonRoleReferenceCollection(fmt.Sprintf("%s/%s/groups/%v/roles?%s", u.basePath, tenantID, groupID, queryParamsValues.Encode()))
	} else {
		return nil, generic.ClientError("Getting role reference collection without username or groupID is not allowed", "FindRoleReferenceCollection")
	}
}

func (u *userApi) NextPageRoleCollection(r *RoleCollection) (*RoleCollection, *generic.Error) {
	return u.getPageRoleCollection(r.Next)
}

func (u *userApi) PreviousPageRoleCollection(r *RoleCollection) (*RoleCollection, *generic.Error) {
	return u.getPageRoleCollection(r.Prev)
}

func (u *userApi) getPageRoleCollection(reference string) (*RoleCollection, *generic.Error) {
	if reference == "" {
		log.Print("No page reference given. Returning nil.")
		return nil, nil
	}

	nextUrl, err := url.Parse(reference)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Unparsable URL given for page reference: '%s'", reference), "GetPage")
	}

	collection, genErr := u.getCommonRoleCollection(fmt.Sprintf("%s?%s", nextUrl.Path, nextUrl.RawQuery))
	if genErr != nil {
		return nil, genErr
	}

	if len(collection.Roles) == 0 {
		log.Print("Returned collection is empty. Returning nil.")
		return nil, nil
	}

	return collection, nil
}

func (u *userApi) getCommonRoleCollection(path string) (*RoleCollection, *generic.Error) {
	body, status, err := u.client.Get(path, generic.AcceptHeader(USER_ACCEPT))
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while getting measurements: %s", err.Error()), "GetMeasurementCollection")
	}

	if status != http.StatusOK {
		return nil, generic.CreateErrorFromResponse(body, status)
	}

	return parseRoleCollectionResponse(body)
}

func (u *userApi) getCommonRoleReferenceCollection(path string) (*RoleReferenceCollection, *generic.Error) {
	body, status, err := u.client.Get(path, generic.AcceptHeader(USER_ACCEPT))
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while getting measurements: %s", err.Error()), "GetMeasurementCollection")
	}

	if status != http.StatusOK {
		return nil, generic.CreateErrorFromResponse(body, status)
	}

	return parseRoleReferenceCollectionResponse(body)
}

func parseRoleCollectionResponse(body []byte) (*RoleCollection, *generic.Error) {
	var result RoleCollection
	if len(body) > 0 {
		err := generic.ObjectFromJson(body, &result)
		if err != nil {
			return nil, generic.ClientError(fmt.Sprintf("Error while parsing response JSON: %s", err.Error()), "CollectionResponseParser")
		}
	} else {
		return nil, generic.ClientError("Response body was empty", "CollectionResponseParser")
	}

	return &result, nil
}

func parseRoleReferenceCollectionResponse(body []byte) (*RoleReferenceCollection, *generic.Error) {
	var result RoleReferenceCollection
	if len(body) > 0 {
		err := generic.ObjectFromJson(body, &result)
		if err != nil {
			return nil, generic.ClientError(fmt.Sprintf("Error while parsing response JSON: %s", err.Error()), "CollectionResponseParser")
		}
	} else {
		return nil, generic.ClientError("Response body was empty", "CollectionResponseParser")
	}

	return &result, nil
}

func (u *userApi) RoleCollection(pageSize int) (*RoleCollection, *generic.Error) {
	return u.FindRoleCollection(pageSize)
}

func (u *userApi) AssignRoleToUser(tenantID, username string, reference *RoleReference) (*RoleReference, *generic.Error) {
	if len(tenantID) == 0 || len(username) == 0 {
		return nil, generic.ClientError("Assigning role to user without tenantID or username is not allowed", "AssignRoleToUser")
	}

	bytes, err := json.Marshal(reference)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while marshalling given role reference: %s", err), "AssignRoleToUser")
	}

	body, status, err := u.client.Post(fmt.Sprintf("%v/%v/users/%v/roles", u.basePath, tenantID, username), bytes, generic.ContentTypeHeader(USER_CONTENT_TYPE))
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while assignign role %v to user %v, %s", reference.Self, username, err), "AssignRoleToUser")
	}

	if status != http.StatusCreated {
		return nil, generic.CreateErrorFromResponse(body, status)
	}

	r := &RoleReference{}
	if err := json.Unmarshal(body, r); err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while unmarshalling response body: %s", err), "AssignRoleToUser")
	}
	return r, nil
}

func (u *userApi) AssignRoleToGroup(tenantID, groupID string, reference *RoleReference) (*RoleReference, *generic.Error) {
	if len(tenantID) == 0 || len(groupID) == 0 {
		return nil, generic.ClientError("Assigning role to group without tenantID or groupID is not allowed", "AssignRoleToGroup")
	}

	bytes, err := json.Marshal(reference)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while marshalling given role reference: %s", err), "AssignRoleToGroup")
	}

	body, status, err := u.client.Post(fmt.Sprintf("%v/%v/groups/%v/roles", u.basePath, tenantID, groupID), bytes, generic.ContentTypeHeader(USER_CONTENT_TYPE))
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while assignign role %v to group %v, %s", reference.Self, groupID, err), "AssignRoleToGroup")
	}

	if status != http.StatusCreated {
		return nil, generic.CreateErrorFromResponse(body, status)
	}

	r := &RoleReference{}
	if err := json.Unmarshal(body, r); err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while unmarshalling response body: %s", err), "AssignRoleToGroup")
	}
	return r, nil
}

func (u *userApi) UnassignRoleFromUser(tenantID, username, roleName string) *generic.Error {
	if len(tenantID) == 0 || len(username) == 0 || len(roleName) == 0 {
		return generic.ClientError("Unassign role from user without tenantID, username or roleName is not allowed", "UnassignRoleFromUser")
	}

	body, status, err := u.client.Delete(fmt.Sprintf("%v/%v/users/%v/roles/%v", u.basePath, tenantID, username, roleName), generic.EmptyHeader())
	if err != nil {
		return generic.ClientError(fmt.Sprintf("Error while unassign role %v from user %v: %s", roleName, username, err), "UnassignRoleFromUser")
	}

	if status != http.StatusNoContent {
		return generic.CreateErrorFromResponse(body, status)
	}
	return nil
}

func (u *userApi) UnassignRoleFromGroup(tenantID, groupID, roleName string) *generic.Error {
	if len(tenantID) == 0 || len(groupID) == 0 || len(roleName) == 0 {
		return generic.ClientError("Unassign role from group without tenantID, groupID or roleName is not allowed", "UnassignRoleFromGroup")
	}

	body, status, err := u.client.Delete(fmt.Sprintf("%v/%v/groups/%v/roles/%v", u.basePath, tenantID, groupID, roleName), generic.EmptyHeader())
	if err != nil {
		return generic.ClientError(fmt.Sprintf("Error while unassign role %v from group %v: %s", roleName, groupID, err), "UnassignRoleFromGroup")
	}

	if status != http.StatusNoContent {
		return generic.CreateErrorFromResponse(body, status)
	}
	return nil
}

func (u *userApi) GetAllRolesOfAUser(tenantID, username string, pageSize int) (*RoleReferenceCollection, *generic.Error) {
	return u.FindRoleReferenceCollection(tenantID, username, "", pageSize)
}

func (u *userApi) GetAllRolesOfAGroup(tenantID, groupID string, pageSize int) (*RoleReferenceCollection, *generic.Error) {
	return u.FindRoleReferenceCollection(tenantID, "", groupID, pageSize)
}

func (u *userApi) GroupDetails(groupID string) (*Group, *generic.Error) {
	if len(groupID) == 0 {
		return nil, generic.ClientError("Getting group details without groupID is not allowed", "GroupDetails")
	}

	body, status, err := u.client.Get(fmt.Sprintf("%v/management/groups/%v", u.basePath, groupID), generic.AcceptHeader(USER_ACCEPT))
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while getting details for group: %v, %s", groupID, err), "GroupDetails")
	}

	if status != http.StatusOK {
		return nil, generic.CreateErrorFromResponse(body, status)
	}

	group := &Group{}
	if err := json.Unmarshal(body, group); err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while unmarshalling response: %s", err), "GroupDetails")
	}
	return group, nil
}

func (u *userApi) GroupByName(tenantID, groupName string) (*Group, *generic.Error) {
	if len(tenantID) == 0 || len(groupName) == 0 {
		return nil, generic.ClientError("Getting group without tenantID or group name is not allowed", "GroupByName")
	}

	body, status, err := u.client.Get(fmt.Sprintf("%v/%v/groupByName/%v", u.basePath, tenantID, groupName), generic.AcceptHeader(USER_ACCEPT))
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while getting group %v by name: %s", groupName, err), "GroupByName")
	}

	if status != http.StatusOK {
		return nil, generic.CreateErrorFromResponse(body, status)
	}

	group := &Group{}
	if err := json.Unmarshal(body, group); err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while unmarshalling response body: %s", err), "GroupByName")
	}
	return group, nil
}

func (u *userApi) RemoveGroup(tenantID, groupID string) *generic.Error {
	if len(tenantID) == 0 || len(groupID) == 0 {
		return generic.ClientError("Removing a group without tenantID and groupID is not allowed", "RemoveGroup")
	}

	body, status, err := u.client.Delete(fmt.Sprintf("%v/%v/groups/%v", u.basePath, tenantID, groupID), generic.EmptyHeader())
	if err != nil {
		return generic.ClientError(fmt.Sprintf("Error while removing group: %s", err), "RemoveGroup")
	}

	if status != http.StatusNoContent {
		return generic.CreateErrorFromResponse(body, status)
	}
	return nil
}

func (u *userApi) UpdateGroup(tenantID, groupID string, group *Group) (*Group, *generic.Error) {
	if len(tenantID) == 0 || len(groupID) == 0 {
		return nil, generic.ClientError("Updating a group without tenantID and groupID is not allowed", "UpdateGroup")
	}

	bytes, err := json.Marshal(group)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while marshalling given group: %s", err), "UpdateGroup")
	}

	body, status, err := u.client.Put(fmt.Sprintf("%v/%v/groups/%v", u.basePath, tenantID, groupID), bytes, generic.EmptyHeader())
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while updating group: %s", err), "UpdateGroup")
	}

	if status != http.StatusOK {
		return nil, generic.CreateErrorFromResponse(body, status)
	}

	g := &Group{}
	if err := json.Unmarshal(body, g); err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while unmarshalling response body: %s", err), "UpdateGroup")
	}
	return g, nil
}

func (u *userApi) GetAllGroupsOfUser(tenantID, username string, pageSize int) (*GroupReferenceCollection, *generic.Error) {
	return u.FindGroupReferenceCollection(tenantID, username, pageSize)
}

func (u *userApi) FindGroupReferenceCollection(tenantID, username string, pageSize int) (*GroupReferenceCollection, *generic.Error) {
	if len(tenantID) == 0 || len(username) == 0 {
		return nil, generic.ClientError("Getting a group reference collection without tenantID and username is not allowed", "FindGroupReferenceCollection")
	}

	queryParamsValues := &url.Values{}
	err := generic.PageSizeParameter(pageSize, queryParamsValues)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while building pageSize parameter to fetch group references: %s", err.Error()), "FindGroupReferenceCollection")
	}
	return u.getCommonGroupReferenceCollection(fmt.Sprintf("%v/%v/users/%v/groups?%v", u.basePath, tenantID, username, queryParamsValues.Encode()))
}

func (u *userApi) NextPageGroupReferenceCollection(r *GroupReferenceCollection) (*GroupReferenceCollection, *generic.Error) {
	return u.getPageGroupReferenceCollection(r.Next)
}

func (u *userApi) PreviousPageGroupCollection(r *GroupReferenceCollection) (*GroupReferenceCollection, *generic.Error) {
	return u.getPageGroupReferenceCollection(r.Prev)
}

func (u *userApi) getPageGroupReferenceCollection(reference string) (*GroupReferenceCollection, *generic.Error) {
	if reference == "" {
		log.Print("No page reference given. Returning nil.")
		return nil, nil
	}

	nextUrl, err := url.Parse(reference)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Unparsable URL given for page reference: '%s'", reference), "GetPage")
	}

	collection, genErr := u.getCommonGroupReferenceCollection(fmt.Sprintf("%s?%s", nextUrl.Path, nextUrl.RawQuery))
	if genErr != nil {
		return nil, genErr
	}

	if len(collection.Groups) == 0 {
		log.Print("Returned collection is empty. Returning nil.")
		return nil, nil
	}

	return collection, nil
}

func (u *userApi) getCommonGroupReferenceCollection(path string) (*GroupReferenceCollection, *generic.Error) {
	body, status, err := u.client.Get(path, generic.AcceptHeader(USER_ACCEPT))
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while getting measurements: %s", err.Error()), "GetMeasurementCollection")
	}

	if status != http.StatusOK {
		return nil, generic.CreateErrorFromResponse(body, status)
	}

	return parseGroupReferenceCollectionResponse(body)
}

func parseGroupReferenceCollectionResponse(body []byte) (*GroupReferenceCollection, *generic.Error) {
	var result GroupReferenceCollection
	if len(body) > 0 {
		err := generic.ObjectFromJson(body, &result)
		if err != nil {
			return nil, generic.ClientError(fmt.Sprintf("Error while parsing response JSON: %s", err.Error()), "CollectionResponseParser")
		}
	} else {
		return nil, generic.ClientError("Response body was empty", "CollectionResponseParser")
	}

	return &result, nil
}
