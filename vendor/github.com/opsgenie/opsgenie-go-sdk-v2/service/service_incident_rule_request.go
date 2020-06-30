package service

import (
	"net/http"

	"github.com/opsgenie/opsgenie-go-sdk-v2/alert"
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/og"
	"github.com/pkg/errors"
)

type CreateIncidentRuleRequest struct {
	client.BaseRequest
	ServiceId          string
	Conditions         []og.Condition        `json:"conditions,omitempty"`
	ConditionMatchType og.ConditionMatchType `json:"conditionMatchType,omitempty"`
	IncidentProperties IncidentProperties    `json:"incidentProperties"`
}

func (r *CreateIncidentRuleRequest) Validate() error {
	err := validateServiceId(r.ServiceId)
	if err != nil {
		return err
	}

	err = og.ValidateConditions(r.Conditions)
	if err != nil {
		return err
	}

	err = validateIncidentProperties(r.IncidentProperties)
	if err != nil {
		return err
	}

	return nil
}

func (r *CreateIncidentRuleRequest) ResourcePath() string {
	return "/v1/services/" + r.ServiceId + "/incident-rules"
}

func (r *CreateIncidentRuleRequest) Method() string {
	return http.MethodPost
}

type UpdateIncidentRuleRequest struct {
	client.BaseRequest
	ServiceId          string
	IncidentRuleId     string
	Conditions         []og.Condition        `json:"conditions,omitempty"`
	ConditionMatchType og.ConditionMatchType `json:"conditionMatchType,omitempty"`
	IncidentProperties IncidentProperties    `json:"incidentProperties"`
}

func (r *UpdateIncidentRuleRequest) Validate() error {
	err := validateServiceId(r.ServiceId)
	if err != nil {
		return err
	}

	err = validateIncidentRuleId(r.IncidentRuleId)
	if err != nil {
		return err
	}

	err = og.ValidateConditions(r.Conditions)
	if err != nil {
		return err
	}

	err = validateIncidentProperties(r.IncidentProperties)
	if err != nil {
		return err
	}

	return nil
}

func (r *UpdateIncidentRuleRequest) ResourcePath() string {
	return "/v1/services/" + r.ServiceId + "/incident-rules/" + r.IncidentRuleId
}

func (r *UpdateIncidentRuleRequest) Method() string {
	return http.MethodPut
}

type DeleteIncidentRuleRequest struct {
	client.BaseRequest
	ServiceId      string
	IncidentRuleId string
}

func (r *DeleteIncidentRuleRequest) Validate() error {
	err := validateServiceId(r.ServiceId)
	if err != nil {
		return err
	}

	err = validateIncidentRuleId(r.IncidentRuleId)
	if err != nil {
		return err
	}
	return nil
}

func (r *DeleteIncidentRuleRequest) ResourcePath() string {
	return "/v1/services/" + r.ServiceId + "/incident-rules/" + r.IncidentRuleId
}

func (r *DeleteIncidentRuleRequest) Method() string {
	return http.MethodDelete
}

type GetIncidentRulesRequest struct {
	client.BaseRequest
	ServiceId string
}

func (r *GetIncidentRulesRequest) Validate() error {
	err := validateServiceId(r.ServiceId)
	if err != nil {
		return err
	}
	return nil
}

func (r *GetIncidentRulesRequest) ResourcePath() string {
	return "/v1/services/" + r.ServiceId + "/incident-rules"
}

func (r *GetIncidentRulesRequest) Method() string {
	return http.MethodGet
}

type IncidentProperties struct {
	Message               string                `json:"message"`
	Tags                  []string              `json:"tags,omitempty"`
	Details               map[string]string     `json:"details,omitempty"`
	Description           string                `json:"description,omitempty"`
	Priority              alert.Priority        `json:"priority"`
	StakeholderProperties StakeholderProperties `json:"stakeholderProperties"`
}

type StakeholderProperties struct {
	Enable      *bool  `json:"enable,omitempty"`
	Message     string `json:"message"`
	Description string `json:"description,omitempty"`
}

func validateServiceId(serviceId string) error {
	if serviceId == "" {
		return errors.New("Service Id cannot be empty.")
	} else if len(serviceId) > 130 {
		return errors.New("Service Id cannot be longer than 130 characters.")
	}
	return nil
}

func validateIncidentRuleId(incidentRuleId string) error {
	if incidentRuleId == "" {
		return errors.New("Incident Rule Id cannot be empty.")
	} else if len(incidentRuleId) > 130 {
		return errors.New("Incident Rule Id cannot be longer than 130 characters.")
	}
	return nil
}

func validateIncidentProperties(incidentProperties IncidentProperties) error {
	if incidentProperties.Message == "" {
		return errors.New("Message field of incident property cannot be empty.")
	} else if len(incidentProperties.Message) > 130 {
		return errors.New("Message field of incident property cannot be longer than 130 characters.")
	}
	if incidentProperties.Description != "" && len(incidentProperties.Description) > 10000 {
		return errors.New("Description field of incident property cannot be longer than 10000 characters.")
	}
	err := alert.ValidatePriority(incidentProperties.Priority)
	if err != nil {
		return err
	}
	err = validateStakeholderProperties(incidentProperties.StakeholderProperties)
	if err != nil {
		return err
	}
	return nil
}

func validateStakeholderProperties(stakeholderProperties StakeholderProperties) error {
	if stakeholderProperties.Message == "" {
		return errors.New("Message field of stakeholder property cannot be empty.")
	} else if len(stakeholderProperties.Message) > 130 {
		return errors.New("Message field of stakeholder property cannot be longer than 130 characters.")
	}
	if stakeholderProperties.Description != "" && len(stakeholderProperties.Description) > 10000 {
		return errors.New("Description field of stakeholder property cannot be longer than 10000 characters.")
	}
	return nil
}
