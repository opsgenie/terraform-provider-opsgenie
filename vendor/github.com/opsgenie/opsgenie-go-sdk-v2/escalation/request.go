package escalation

import (
	"net/http"

	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/og"
	"github.com/pkg/errors"
)

type Identifier string

const (
	Name Identifier = "name"
	Id   Identifier = "id"
)

type RepeatRequest struct {
	WaitInterval         uint32 `json:"waitInterval,omitempty"`
	Count                uint32 `json:"count,omitempty"`
	ResetRecipientStates *bool  `json:"resetRecipientStates,omitempty"`
	CloseAlertAfterAll   *bool  `json:"closeAlertAfterAll,omitempty"`
}

type RuleRequest struct {
	Condition  og.EscalationCondition `json:"condition,omitempty"`
	NotifyType og.NotifyType          `json:"notifyType,omitempty"`
	Recipient  og.Participant         `json:"recipient,omitempty"`
	Delay      EscalationDelayRequest `json:"delay,omitempty"`
}

type EscalationDelayRequest struct {
	TimeAmount uint32 `json:"timeAmount"`
}

type CreateRequest struct {
	client.BaseRequest
	Name        string         `json:"name,omitempty"`
	Description string         `json:"description,omitempty"`
	Rules       []RuleRequest  `json:"rules,omitempty"`
	OwnerTeam   *og.OwnerTeam  `json:"ownerTeam,omitempty"`
	Repeat      *RepeatRequest `json:"repeat,omitempty"`
}

func (r *CreateRequest) Validate() error {
	if r.Name == "" {
		return errors.New("Name cannot be empty.")
	}
	if len(r.Rules) == 0 {
		return errors.New("Rules list cannot be empty.")
	}
	err := validateRules(r.Rules)
	if err != nil {
		return err
	}
	return nil
}

func (r *CreateRequest) ResourcePath() string {
	return "/v2/escalations"
}

func (r *CreateRequest) Method() string {
	return http.MethodPost
}

type GetRequest struct {
	client.BaseRequest
	IdentifierType Identifier
	Identifier     string
}

func (r *GetRequest) Validate() error {
	err := validateIdentifiers(r.Identifier, r.IdentifierType)
	if err != nil {
		return err
	}
	return nil
}

func (r *GetRequest) Method() string {
	return http.MethodGet
}

func (r *GetRequest) ResourcePath() string {
	return "/v2/escalations/" + r.Identifier
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
	Name           string         `json:"name,omitempty"`
	Description    string         `json:"description,omitempty"`
	Rules          []RuleRequest  `json:"rules,omitempty"`
	OwnerTeam      *og.OwnerTeam  `json:"ownerTeam,omitempty"`
	Repeat         *RepeatRequest `json:"repeat,omitempty"`
	IdentifierType Identifier
	Identifier     string
}

func (r *UpdateRequest) Validate() error {
	err := validateIdentifiers(r.Identifier, r.IdentifierType)
	if err != nil {
		return err
	}
	err = validateRules(r.Rules)
	if err != nil {
		return err
	}
	return nil
}

func (r *UpdateRequest) ResourcePath() string {
	return "/v2/escalations/" + r.Identifier
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

func (r *UpdateRequest) Method() string {
	return http.MethodPatch
}

type DeleteRequest struct {
	client.BaseRequest
	IdentifierType Identifier
	Identifier     string
}

func (r *DeleteRequest) Validate() error {
	err := validateIdentifiers(r.Identifier, r.IdentifierType)
	if err != nil {
		return err
	}
	return nil
}

func (r *DeleteRequest) Method() string {
	return http.MethodDelete
}

func (r *DeleteRequest) ResourcePath() string {
	return "/v2/escalations/" + r.Identifier
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

type listRequest struct {
	client.BaseRequest
}

func (r *listRequest) Validate() error {
	return nil
}

func (r *listRequest) Method() string {
	return http.MethodGet
}

func (r *listRequest) ResourcePath() string {
	return "/v2/escalations"
}

func validateRules(rules []RuleRequest) error {
	for _, rule := range rules {
		switch rule.Condition {
		case og.IfNotAcked, og.IfNotClosed:
			break
		default:
			return errors.New("Rule Condition should be one of these: 'if-not-acked', 'if-not-closed'.")
		}
		switch rule.NotifyType {
		case og.Next, og.Previous, og.Default, og.Users, og.Admins, og.All, og.Random:
			break
		default:
			return errors.New("Notify Type should be one of these: 'next', 'previous', 'default', 'users', 'admins', 'all'.")
		}
		err := validateRecipient(rule.Recipient)
		if err != nil {
			return err
		}
	}
	return nil
}

func validateRecipient(participant og.Participant) error {

	if participant.Type == "" {
		return errors.New("Recipient type cannot be empty.")
	}
	if participant.Type != og.User && participant.Type != og.Team && participant.Type != og.Schedule {
		return errors.New("Recipient type should be one of these: 'User', 'Team', 'Schedule'")
	}
	if participant.Type == og.User && participant.Username == "" && participant.Id == "" {
		return errors.New("For recipient type user either username or id must be provided.")
	}
	if (participant.Type == og.Team || participant.Type == og.Schedule) && participant.Name == "" && participant.Id == "" {
		return errors.New("For recipient type team and schedule either name or id must be provided.")
	}
	return nil
}

func validateIdentifiers(identifier string, identifierType Identifier) error {
	if identifierType != "" && identifierType != Name && identifierType != Id {
		return errors.New("Identifier Type should be one of this : 'id', 'name' or empty.")
	}

	if identifier == "" {
		return errors.New("Identifier cannot be empty.")
	}
	return nil
}
