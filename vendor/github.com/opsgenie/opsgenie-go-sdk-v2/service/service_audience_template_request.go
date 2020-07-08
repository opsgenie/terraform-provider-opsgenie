package service

import (
	"net/http"

	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/og"
	"github.com/pkg/errors"
)

type GetAudienceTemplateRequest struct {
	client.BaseRequest
	ServiceId string
}

func (r *GetAudienceTemplateRequest) Validate() error {
	err := validateServiceId(r.ServiceId)
	if err != nil {
		return err
	}
	return nil
}

func (r *GetAudienceTemplateRequest) ResourcePath() string {
	return "/v1/services/" + r.ServiceId + "/audience-templates"
}

func (r *GetAudienceTemplateRequest) Method() string {
	return http.MethodGet
}

type UpdateAudienceTemplateRequest struct {
	client.BaseRequest
	ServiceId   string
	Responder   ResponderOfAudience   `json:"responder,omitempty"`
	Stakeholder StakeholderOfAudience `json:"stakeholder,omitempty"`
}

func (r *UpdateAudienceTemplateRequest) Validate() error {
	err := validateServiceId(r.ServiceId)
	if err != nil {
		return err
	}

	if &r.Responder != nil && (len(r.Responder.Teams) > 50 || len(r.Responder.Individuals) > 50) {
		return errors.New("You can set at most 50 team and 50 user to the template.")
	}
	if r.Stakeholder.ConditionMatchType == og.MatchAll {
		return errors.New("Condition match type can only be match-any-condition or match-all-conditions.")
	}
	for conditionIndex := range r.Stakeholder.Conditions {
		err = validateConditionOfStakeholder(r.Stakeholder.Conditions[conditionIndex])
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *UpdateAudienceTemplateRequest) ResourcePath() string {
	return "/v1/services/" + r.ServiceId + "/audience-templates"
}

func (r *UpdateAudienceTemplateRequest) Method() string {
	return http.MethodPatch
}

type ResponderOfAudience struct {
	Teams       []string `json:"teams,omitempty"`
	Individuals []string `json:"individuals,omitempty"`
}
type StakeholderOfAudience struct {
	Individuals        []string                 `json:"individuals,omitempty"`
	ConditionMatchType og.ConditionMatchType    `json:"conditionMatchType,omitempty"`
	Conditions         []ConditionOfStakeholder `json:"conditions,omitempty"`
}
type ConditionOfStakeholder struct {
	MatchField MatchField `json:"matchField,omitempty"`
	Key        string     `json:"key,omitempty"`
	Value      string     `json:"value,omitempty"`
}

func validateConditionOfStakeholder(condition ConditionOfStakeholder) error {
	if condition.MatchField == "" {
		return errors.New("Match field must be one of [country, state. city, zipCode, line, tag , customProperty].")
	}
	if condition.MatchField == CustomProperty && condition.Key == "" {
		return errors.New("Key field cannot be empty.")
	}
	if condition.Value == "" {
		return errors.New("Value field cannot be empty.")

	}
	return nil
}

type MatchField string

const (
	Country        MatchField = "country"
	State          MatchField = "state"
	City           MatchField = "city"
	ZipCode        MatchField = "zipCode"
	Line           MatchField = "line"
	Tag            MatchField = "tag"
	CustomProperty MatchField = "customProperty"
)
