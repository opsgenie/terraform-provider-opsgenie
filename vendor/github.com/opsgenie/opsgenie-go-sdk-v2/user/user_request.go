package user

import (
	"net/http"
	"strconv"

	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/pkg/errors"
)

type UserRoleRequest struct {
	RoleName string `json:"name"`
}

type UserAddressRequest struct {
	Country string `json:"country,omitempty"`
	State   string `json:"state,omitempty"`
	City    string `json:"city,omitempty"`
	Line    string `json:"line,omitempty"`
	ZipCode string `json:"zipCode,omitempty"`
}

type Order string
type SortField string
type Identifier uint32

const (
	Asc  Order = "asc"
	Desc Order = "desc"

	Username    SortField = "username"
	Role        SortField = "role"
	FullName    SortField = "fullName"
	FullNameRaw SortField = "fullName.raw"
	Verified    SortField = "verified"
	Blocked     SortField = "blocked"
	CreatedAt   SortField = "createdAt"

	Name Identifier = iota
	Id
)

type CreateRequest struct {
	client.BaseRequest
	Username           string              `json:"username"`
	FullName           string              `json:"fullName"`
	Role               *UserRoleRequest    `json:"role"`
	SkypeUsername      string              `json:"skypeUsername,omitempty"`
	UserAddressRequest *UserAddressRequest `json:"userAddress,omitempty"`
	Tags               []string            `json:"tags,omitempty"`
	Details            map[string][]string `json:"details,omitempty"`
	TimeZone           string              `json:"timeZone,omitempty"`
	Locale             string              `json:"locale,omitempty"`
	InvitationDisabled string              `json:"invitationDisabled,omitempty"`
}

func (r *CreateRequest) Validate() error {
	if r.Username == "" {
		return errors.New("Username can not be empty")
	}
	if r.FullName == "" {
		return errors.New("FullName can not be empty")
	}
	if r.Role == nil || r.Role.RoleName == "" {
		return errors.New("User Role can not be empty")
	}

	return nil
}

func (r *CreateRequest) ResourcePath() string {

	return "/v2/users"
}

func (r *CreateRequest) Method() string {
	return http.MethodPost
}

type GetRequest struct {
	client.BaseRequest
	Identifier string
	Expand     string
}

func (r *GetRequest) Validate() error {
	if r.Identifier == "" {
		return errors.New("Identifier can not be empty")
	}
	return nil
}

func (r *GetRequest) ResourcePath() string {

	return "/v2/users/" + r.Identifier
}

func (r *GetRequest) Method() string {
	return http.MethodGet
}

func (r *GetRequest) RequestParams() map[string]string {
	params := make(map[string]string)

	if r.Expand != "" {
		params["expand"] = r.Expand
	}

	return params
}

type UpdateRequest struct {
	client.BaseRequest
	Identifier         string
	Username           string              `json:"username,omitempty"`
	FullName           string              `json:"fullName,omitempty"`
	Role               *UserRoleRequest    `json:"role,omitempty"`
	SkypeUsername      string              `json:"skypeUsername,omitempty"`
	UserAddressRequest *UserAddressRequest `json:"userAddress,omitempty"`
	Tags               []string            `json:"tags,omitempty"`
	Details            map[string][]string `json:"details,omitempty"`
	TimeZone           string              `json:"timeZone,omitempty"`
	Locale             string              `json:"locale,omitempty"`
	InvitationDisabled string              `json:"invitationDisabled,omitempty"`
}

func (r *UpdateRequest) Validate() error {
	if r.Identifier == "" {
		return errors.New("Identifier can not be empty")
	}
	if r.Role != nil && r.Role.RoleName == "" {
		return errors.New("User role name can not be empty")
	}

	return nil
}

func (r *UpdateRequest) ResourcePath() string {

	return "/v2/users/" + r.Identifier
}

func (r *UpdateRequest) Method() string {
	return http.MethodPatch
}

type DeleteRequest struct {
	client.BaseRequest
	Identifier string
}

func (r *DeleteRequest) Validate() error {
	if r.Identifier == "" {
		return errors.New("Identifier can not be empty")
	}
	return nil
}

func (r *DeleteRequest) ResourcePath() string {

	return "/v2/users/" + r.Identifier
}

func (r *DeleteRequest) Method() string {
	return http.MethodDelete
}

type ListRequest struct {
	client.BaseRequest
	Limit  int
	Offset int
	Sort   SortField
	Order  Order
	Query  string
}

func (r *ListRequest) Validate() error {
	return nil
}

func (r *ListRequest) ResourcePath() string {

	return "/v2/users/"
}

func (r *ListRequest) Method() string {
	return http.MethodGet
}

func (r *ListRequest) RequestParams() map[string]string {

	params := make(map[string]string)

	if r.Limit != 0 {
		params["limit"] = strconv.Itoa(r.Limit)
	}

	if r.Sort != "" {
		params["sort"] = string(r.Sort)
	}

	if r.Offset != 0 {
		params["offset"] = strconv.Itoa(r.Offset)
	}

	if r.Query != "" {
		params["query"] = r.Query
	}

	if r.Order != "" {
		params["order"] = string(r.Order)
	}

	return params
}

type ListUserEscalationsRequest struct {
	client.BaseRequest
	Identifier string
}

func (r *ListUserEscalationsRequest) Validate() error {
	if r.Identifier == "" {
		return errors.New("Identifier can not be empty")
	}
	return nil
}

func (r *ListUserEscalationsRequest) ResourcePath() string {

	return "/v2/users/" + r.Identifier + "/escalations"
}

func (r *ListUserEscalationsRequest) Method() string {
	return http.MethodGet
}

type ListUserTeamsRequest struct {
	client.BaseRequest
	Identifier string
}

func (r *ListUserTeamsRequest) Validate() error {
	if r.Identifier == "" {
		return errors.New("Identifier can not be empty")
	}
	return nil
}

func (r *ListUserTeamsRequest) ResourcePath() string {

	return "/v2/users/" + r.Identifier + "/teams"
}

func (r *ListUserTeamsRequest) Method() string {
	return http.MethodGet
}

type ListUserForwardingRulesRequest struct {
	client.BaseRequest
	Identifier string
}

func (r *ListUserForwardingRulesRequest) Validate() error {
	if r.Identifier == "" {
		return errors.New("Identifier can not be empty")
	}
	return nil
}

func (r *ListUserForwardingRulesRequest) ResourcePath() string {

	return "/v2/users/" + r.Identifier + "/forwarding-rules"
}

func (r *ListUserForwardingRulesRequest) Method() string {
	return http.MethodGet
}

type ListUserSchedulesRequest struct {
	client.BaseRequest
	Identifier string
}

func (r *ListUserSchedulesRequest) Validate() error {
	if r.Identifier == "" {
		return errors.New("Identifier can not be empty")
	}
	return nil
}

func (r *ListUserSchedulesRequest) ResourcePath() string {

	return "/v2/users/" + r.Identifier + "/schedules"
}

func (r *ListUserSchedulesRequest) Method() string {
	return http.MethodGet
}

type GetSavedSearchRequest struct {
	client.BaseRequest
	Identifier     string
	IdentifierType Identifier
}

func (r *GetSavedSearchRequest) Validate() error {
	if r.Identifier == "" {
		return errors.New("Identifier can not be empty")
	}
	return nil
}

func (r *GetSavedSearchRequest) ResourcePath() string {

	return "/v2/users/saved-searches/" + r.Identifier
}

func (r *GetSavedSearchRequest) Method() string {
	return http.MethodGet
}

func (r *GetSavedSearchRequest) RequestParams() map[string]string {
	params := make(map[string]string)

	if r.IdentifierType == Name {
		params["identifierType"] = "name"
	} else {
		params["identifierType"] = "id"
	}

	return params
}

type ListSavedSearchesRequest struct {
	client.BaseRequest
}

func (r *ListSavedSearchesRequest) Validate() error {
	return nil
}

func (r *ListSavedSearchesRequest) ResourcePath() string {

	return "/v2/users/saved-searches"
}

func (r *ListSavedSearchesRequest) Method() string {
	return http.MethodGet
}

type DeleteSavedSearchRequest struct {
	client.BaseRequest
	Identifier     string
	IdentifierType Identifier
}

func (r *DeleteSavedSearchRequest) Validate() error {
	if r.Identifier == "" {
		return errors.New("Identifier can not be empty")
	}
	return nil
}

func (r *DeleteSavedSearchRequest) ResourcePath() string {

	return "/v2/users/saved-searches/" + r.Identifier
}

func (r *DeleteSavedSearchRequest) Method() string {
	return http.MethodDelete
}

func (r *DeleteSavedSearchRequest) RequestParams() map[string]string {
	params := make(map[string]string)

	if r.IdentifierType == Name {
		params["identifierType"] = "name"
	} else {
		params["identifierType"] = "id"
	}

	return params
}
