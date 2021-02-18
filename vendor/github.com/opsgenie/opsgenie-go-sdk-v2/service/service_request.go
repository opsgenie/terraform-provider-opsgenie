package service

import (
	"net/http"
	"strconv"

	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/pkg/errors"
)

type CreateRequest struct {
	client.BaseRequest
	Name        string     `json:"name"`
	TeamId      string     `json:"teamId"`
	Description string     `json:"description,omitempty"`
	Visibility  Visibility `json:"visibility,omitempty"`
	Tags        []string   `json:"tags,omitempty"`
}

func (r *CreateRequest) Validate() error {
	if r.Name == "" {
		return errors.New("Name field cannot be empty.")
	}
	if r.TeamId == "" {
		return errors.New("Team ID field cannot be empty.")
	}
	err := validateVisibility(r.Visibility)
	if err != nil {
		return err
	}
	return nil
}

func (r *CreateRequest) ResourcePath() string {
	return "/v1/services"
}

func (r *CreateRequest) Method() string {
	return http.MethodPost
}

type UpdateRequest struct {
	client.BaseRequest
	Id          string
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	Visibility  Visibility `json:"visibility,omitempty"`
}

func (r *UpdateRequest) Validate() error {
	if r.Id == "" {
		return errors.New("Service ID cannot be blank.")
	}
	err := validateVisibility(r.Visibility)
	if err != nil {
		return err
	}
	return nil
}

func (r *UpdateRequest) ResourcePath() string {
	return "/v1/services/" + r.Id
}

func (r *UpdateRequest) Method() string {
	return http.MethodPatch
}

type DeleteRequest struct {
	client.BaseRequest
	Id string
}

func (r *DeleteRequest) Validate() error {
	if r.Id == "" {
		return errors.New("Service ID cannot be blank.")
	}
	return nil
}

func (r *DeleteRequest) ResourcePath() string {
	return "/v1/services/" + r.Id
}

func (r *DeleteRequest) Method() string {
	return http.MethodDelete
}

type GetRequest struct {
	client.BaseRequest
	Id string
}

func (r *GetRequest) Validate() error {
	if r.Id == "" {
		return errors.New("Service ID cannot be blank.")
	}
	return nil
}

func (r *GetRequest) ResourcePath() string {
	return "/v1/services/" + r.Id
}

func (r *GetRequest) Method() string {
	return http.MethodGet
}

type ListRequest struct {
	client.BaseRequest
	Limit  int
	Offset int
}

func (r *ListRequest) Validate() error {
	return nil
}

func (r *ListRequest) ResourcePath() string {
	return "/v1/services"
}

func (r *ListRequest) Method() string {
	return http.MethodGet
}

func (r *ListRequest) RequestParams() map[string]string {
	params := map[string]string{}
	if r.Limit != 0 {
		params["limit"] = strconv.Itoa(r.Limit)
	}
	if r.Offset != 0 {
		params["offset"] = strconv.Itoa(r.Offset)
	}
	return params
}

type Visibility string

const (
	TeamMembers   Visibility = "TEAM_MEMBERS"
	OpsgenieUsers Visibility = "OPSGENIE_USERS"
)

func validateVisibility(visibility Visibility) error {
	switch visibility {
	case TeamMembers, OpsgenieUsers, "":
		return nil
	}
	return errors.New("Visibility should be one of these: " +
		"'TeamMembers', 'OpsgenieUsers' or empty.")
}
