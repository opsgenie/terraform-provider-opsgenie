package heartbeat

import (
	"errors"
	"net/http"

	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/og"
)

type pingRequest struct {
	client.BaseRequest
	HeartbeatName string
}

func nameValidation(name string) error {
	if name == "" {
		return errors.New("HeartbeatName cannot be empty")
	}
	return nil
}

func (r pingRequest) Validate() error {
	return nameValidation(r.HeartbeatName)
}

func (r pingRequest) ResourcePath() string {
	return "/v2/heartbeats/" + r.HeartbeatName + "/ping"
}

func (r pingRequest) Method() string {
	return http.MethodGet
}

type getRequest struct {
	client.BaseRequest
	HeartbeatName string
}

func (r getRequest) Validate() error {
	return nameValidation(r.HeartbeatName)
}

func (r getRequest) ResourcePath() string {
	return "/v2/heartbeats/" + r.HeartbeatName
}

func (r getRequest) Method() string {
	return http.MethodGet
}

type listRequest struct {
	client.BaseRequest
}

func (r listRequest) Validate() error {
	return nil
}

func (r listRequest) ResourcePath() string {
	return "/v2/heartbeats"
}

func (r listRequest) Method() string {
	return http.MethodGet
}

type UpdateRequest struct {
	client.BaseRequest
	Name          string       `json:"name"`
	Description   string       `json:"description,omitempty"`
	Interval      int          `json:"interval"`
	IntervalUnit  Unit         `json:"intervalUnit"`
	Enabled       bool         `json:"enabled,omitempty"`
	OwnerTeam     og.OwnerTeam `json:"ownerTeam"`
	AlertMessage  string       `json:"alertMessage,omitempty"`
	AlertTag      []string     `json:"alertTags,omitempty"`
	AlertPriority string       `json:"alertPriority,omitempty"`
}

func (r UpdateRequest) Validate() error {
	if r.Name == "" {
		return errors.New("Invalid request. Name cannot be empty. ")
	}
	if &r.OwnerTeam == nil || (r.OwnerTeam.Id == "" && r.OwnerTeam.Name == "") {
		return errors.New("Invalid request. Owner team cannot be empty. ")
	}
	if r.Interval < 1 {
		return errors.New("Invalid request. Interval cannot be smaller than 1. ")
	}
	if r.IntervalUnit == "" {
		return errors.New("Invalid request. IntervalUnit cannot be empty. ")
	}
	return nil
}

func (r UpdateRequest) ResourcePath() string {
	return "/v2/heartbeats/" + r.Name
}

func (r UpdateRequest) Method() string {
	return http.MethodPatch
}

type AddRequest struct {
	client.BaseRequest
	Name          string       `json:"name"`
	Description   string       `json:"description,omitempty"`
	Interval      int          `json:"interval"`
	IntervalUnit  Unit         `json:"intervalUnit"`
	Enabled       bool         `json:"enabled"`
	OwnerTeam     og.OwnerTeam `json:"ownerTeam"`
	AlertMessage  string       `json:"alertMessage,omitempty"`
	AlertTag      []string     `json:"alertTags,omitempty"`
	AlertPriority string       `json:"alertPriority,omitempty"`
}

func (r AddRequest) Validate() error {
	if r.Name == "" {
		return errors.New("Invalid request. Name cannot be empty. ")
	}
	if &r.OwnerTeam == nil || (r.OwnerTeam.Id == "" && r.OwnerTeam.Name == "") {
		return errors.New("Invalid request. Owner team cannot be empty. ")
	}
	if r.Interval < 1 {
		return errors.New("Invalid request. Interval cannot be smaller than 1. ")
	}
	if r.IntervalUnit == "" {
		return errors.New("Invalid request. IntervalUnit cannot be empty. ")
	}
	return nil
}

func (r AddRequest) ResourcePath() string {
	return "/v2/heartbeats"
}

func (r AddRequest) Method() string {
	return http.MethodPost
}

type Unit string

const (
	Minutes Unit = "minutes"
	Hours   Unit = "hours"
	Days    Unit = "days"
)

type enableRequest struct {
	client.BaseRequest
	heartbeatName string
}

func (r enableRequest) Validate() error {
	if r.heartbeatName == "" {
		return errors.New("Invalid request. Name cannot be empty. ")
	}
	return nil
}

func (r enableRequest) ResourcePath() string {
	return "/v2/heartbeats/" + r.heartbeatName + "/enable"
}

func (r enableRequest) Method() string {
	return http.MethodPost
}

type disableRequest struct {
	client.BaseRequest
	heartbeatName string
}

func (r disableRequest) Validate() error {
	if r.heartbeatName == "" {
		return errors.New("Invalid request. Name cannot be empty. ")
	}
	return nil
}

func (r disableRequest) ResourcePath() string {
	return "/v2/heartbeats/" + r.heartbeatName + "/disable"
}

func (r disableRequest) Method() string {
	return http.MethodPost
}

type deleteRequest struct {
	client.BaseRequest
	HeartbeatName string
}

func (r deleteRequest) Validate() error {
	return nameValidation(r.HeartbeatName)
}

func (r deleteRequest) ResourcePath() string {
	return "/v2/heartbeats/" + r.HeartbeatName
}

func (r deleteRequest) Method() string {
	return http.MethodDelete
}
