package contact

import (
	"net/http"

	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/pkg/errors"
)

type CreateRequest struct {
	client.BaseRequest
	UserIdentifier  string
	To              string     `json:"to"`
	MethodOfContact MethodType `json:"method"`
}

func (r *CreateRequest) Validate() error {
	if r.UserIdentifier == "" {
		return errors.New("User identifier cannot be empty.")
	}
	if r.To == "" {
		return errors.New("to cannot be empty.")
	}
	if r.MethodOfContact == "" {
		return errors.New("Method cannot be empty.")
	}

	return nil
}

func (r *CreateRequest) ResourcePath() string {
	return "/v2/users/" + r.UserIdentifier + "/contacts"
}

func (r *CreateRequest) Method() string {
	return http.MethodPost
}

type GetRequest struct {
	client.BaseRequest
	UserIdentifier    string
	ContactIdentifier string
}

func (r *GetRequest) Validate() error {
	err := validateIdentifier(r.UserIdentifier, r.ContactIdentifier)
	if err != nil {
		return err
	}
	return nil
}

func (r *GetRequest) ResourcePath() string {
	return "/v2/users/" + r.UserIdentifier + "/contacts/" + r.ContactIdentifier
}

func (r *GetRequest) Method() string {
	return http.MethodGet
}

type UpdateRequest struct {
	client.BaseRequest
	UserIdentifier    string
	ContactIdentifier string
	To                string `json:"to"`
}

func (r *UpdateRequest) Validate() error {

	err := validateIdentifier(r.UserIdentifier, r.ContactIdentifier)
	if err != nil {
		return err
	}

	if r.To == "" {
		return errors.New("to cannot be empty.")
	}

	return nil
}

func (r *UpdateRequest) ResourcePath() string {
	return "/v2/users/" + r.UserIdentifier + "/contacts/" + r.ContactIdentifier
}

func (r *UpdateRequest) Method() string {
	return http.MethodPatch
}

type DeleteRequest struct {
	client.BaseRequest
	UserIdentifier    string
	ContactIdentifier string
}

func (r *DeleteRequest) Validate() error {
	err := validateIdentifier(r.UserIdentifier, r.ContactIdentifier)
	if err != nil {
		return err
	}
	return nil
}
func (r *DeleteRequest) ResourcePath() string {
	return "/v2/users/" + r.UserIdentifier + "/contacts/" + r.ContactIdentifier
}

func (r *DeleteRequest) Method() string {
	return http.MethodDelete
}

type ListRequest struct {
	client.BaseRequest
	UserIdentifier string
}

func (r *ListRequest) Validate() error {
	if r.UserIdentifier == "" {
		return errors.New("User identifier cannot be empty.")
	}
	return nil
}
func (r *ListRequest) ResourcePath() string {
	return "/v2/users/" + r.UserIdentifier + "/contacts"
}

func (r *ListRequest) Method() string {
	return http.MethodGet
}

type EnableRequest struct {
	client.BaseRequest
	UserIdentifier    string
	ContactIdentifier string
}

func (r *EnableRequest) Validate() error {
	err := validateIdentifier(r.UserIdentifier, r.ContactIdentifier)
	if err != nil {
		return err
	}
	return nil
}
func (r *EnableRequest) ResourcePath() string {
	return "/v2/users/" + r.UserIdentifier + "/contacts/" + r.ContactIdentifier + "/enable"
}

func (r *EnableRequest) Method() string {
	return http.MethodPost
}

type DisableRequest struct {
	client.BaseRequest
	UserIdentifier    string
	ContactIdentifier string
}

func (r *DisableRequest) Validate() error {
	err := validateIdentifier(r.UserIdentifier, r.ContactIdentifier)
	if err != nil {
		return err
	}
	return nil
}
func (r *DisableRequest) ResourcePath() string {
	return "/v2/users/" + r.UserIdentifier + "/contacts/" + r.ContactIdentifier + "/disable"
}

func (r *DisableRequest) Method() string {
	return http.MethodPost
}

func validateIdentifier(userIdentifier string, contactIdentifier string) error {
	if userIdentifier == "" {
		return errors.New("User identifier cannot be empty.")
	}
	if contactIdentifier == "" {
		return errors.New("Contact identifier cannot be empty.")

	}
	return nil
}

type MethodType string

const (
	Sms   MethodType = "sms"
	Email MethodType = "email"
	Voice MethodType = "voice"
)
