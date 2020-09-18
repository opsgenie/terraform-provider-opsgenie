package custom_user_role

import (
	"net/http"

	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/pkg/errors"
)

type ExtendedRole string
type Identifier uint32

const (
	Name Identifier = iota
	Id

	ExtendedRoleUser        ExtendedRole = "user"
	ExtendedRoleObserver    ExtendedRole = "observer"
	ExtendedRoleStakeholder ExtendedRole = "stakeholder"
)

type CreateRequest struct {
	client.BaseRequest
	Name             string       `json:"name"`
	ExtendedRole     ExtendedRole `json:"extendedRole,omitempty"`
	GrantedRights    []string     `json:"grantedRights,omitempty"`
	DisallowedRights []string     `json:"disallowedRights,omitempty"`
}

func (r *CreateRequest) Validate() error {
	if r.Name == "" {
		return errors.New("Name can not be empty")
	}

	switch r.ExtendedRole {
	case ExtendedRoleStakeholder, ExtendedRoleUser, ExtendedRoleObserver, "":
		break
	default:
		return errors.New("ExtendedRole should be one of these: 'observer', 'user', 'stakeholder' or empty")
	}

	return nil
}

func (r *CreateRequest) ResourcePath() string {
	return "/v2/roles"
}

func (r *CreateRequest) Method() string {
	return http.MethodPost
}

type GetRequest struct {
	client.BaseRequest
	Identifier     string
	IdentifierType Identifier
}

func (r *GetRequest) Validate() error {
	if r.Identifier == "" {
		return errors.New("Identifier can not be empty")
	}
	return nil
}

func (r *GetRequest) ResourcePath() string {

	return "/v2/roles/" + r.Identifier
}

func (r *GetRequest) Method() string {
	return http.MethodGet
}

func (r *GetRequest) RequestParams() map[string]string {
	params := make(map[string]string)

	if r.IdentifierType == Name {
		params["identifierType"] = "name"
	} else {
		params["identifierType"] = "id"
	}

	return params
}

type UpdateRequest struct {
	client.BaseRequest
	Identifier       string
	IdentifierType   Identifier
	Name             string       `json:"name,omitempty"`
	ExtendedRole     ExtendedRole `json:"extendedRole,omitempty"`
	GrantedRights    []string     `json:"grantedRights,omitempty"`
	DisallowedRights []string     `json:"disallowedRights,omitempty"`
}

func (r *UpdateRequest) Validate() error {
	if r.Identifier == "" {
		return errors.New("Identifier can not be empty")
	}

	switch r.ExtendedRole {
	case ExtendedRoleStakeholder, ExtendedRoleUser, ExtendedRoleObserver, "":
		break
	default:
		return errors.New("ExtendedRole should be one of these: 'observer', 'user', 'stakeholder' or empty")
	}

	return nil
}

func (r *UpdateRequest) ResourcePath() string {
	return "/v2/roles/" + r.Identifier
}

func (r *UpdateRequest) Method() string {
	return http.MethodPut
}

func (r *UpdateRequest) RequestParams() map[string]string {
	params := make(map[string]string)

	if r.IdentifierType == Name {
		params["identifierType"] = "name"
	} else {
		params["identifierType"] = "id"
	}

	return params
}

type DeleteRequest struct {
	client.BaseRequest
	Identifier     string
	IdentifierType Identifier
}

func (r *DeleteRequest) Validate() error {
	if r.Identifier == "" {
		return errors.New("Identifier can not be empty")
	}
	return nil
}

func (r *DeleteRequest) ResourcePath() string {

	return "/v2/roles/" + r.Identifier
}

func (r *DeleteRequest) Method() string {
	return http.MethodDelete
}

func (r *DeleteRequest) RequestParams() map[string]string {
	params := make(map[string]string)

	if r.IdentifierType == Name {
		params["identifierType"] = "name"
	} else {
		params["identifierType"] = "id"
	}

	return params
}

type ListRequest struct {
	client.BaseRequest
}

func (r *ListRequest) Validate() error {
	return nil
}

func (r *ListRequest) ResourcePath() string {

	return "/v2/roles/"
}

func (r *ListRequest) Method() string {
	return http.MethodGet
}
