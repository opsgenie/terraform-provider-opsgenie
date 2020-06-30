package service

import (
	"net/http"

	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/pkg/errors"
)

type CreateIncidentTemplateRequest struct {
	client.BaseRequest
	ServiceId        string
	IncidentTemplate IncidentTemplateRequest `json:"incidentTemplate"`
}

func (r *CreateIncidentTemplateRequest) Validate() error {
	err := validateServiceId(r.ServiceId)
	if err != nil {
		return err
	}

	err = validateIncidentTemplate(r.IncidentTemplate)
	if err != nil {
		return err
	}

	return nil
}

func (r *CreateIncidentTemplateRequest) ResourcePath() string {
	return "/v1/services/" + r.ServiceId + "/incident-templates"
}

func (r *CreateIncidentTemplateRequest) Method() string {
	return http.MethodPost
}

type UpdateIncidentTemplateRequest struct {
	client.BaseRequest
	ServiceId          string
	IncidentTemplateId string
	Name               string             `json:"name"`
	IncidentProperties IncidentProperties `json:"incidentProperties"`
}

func (r *UpdateIncidentTemplateRequest) Validate() error {
	err := validateServiceId(r.ServiceId)
	if err != nil {
		return err
	}

	if r.IncidentTemplateId == "" {
		return errors.New("Incident Template Id cannot be empty.")
	}

	if r.Name == "" {
		return errors.New("Name of incident template cannot be empty.")
	}

	err = validateIncidentProperties(r.IncidentProperties)
	if err != nil {
		return err
	}

	return nil
}

func (r *UpdateIncidentTemplateRequest) ResourcePath() string {
	return "/v1/services/" + r.ServiceId + "/incident-templates/" + r.IncidentTemplateId
}

func (r *UpdateIncidentTemplateRequest) Method() string {
	return http.MethodPut
}

type DeleteIncidentTemplateRequest struct {
	client.BaseRequest
	ServiceId          string
	IncidentTemplateId string
}

func (r *DeleteIncidentTemplateRequest) Validate() error {
	err := validateServiceId(r.ServiceId)
	if err != nil {
		return err
	}

	if r.IncidentTemplateId == "" {
		return errors.New("Incident Template Id cannot be empty.")
	}

	return nil
}

func (r *DeleteIncidentTemplateRequest) ResourcePath() string {
	return "/v1/services/" + r.ServiceId + "/incident-templates/" + r.IncidentTemplateId
}

func (r *DeleteIncidentTemplateRequest) Method() string {
	return http.MethodDelete
}

type GetIncidentTemplatesRequest struct {
	client.BaseRequest
	ServiceId string
}

func (r *GetIncidentTemplatesRequest) Validate() error {
	err := validateServiceId(r.ServiceId)
	if err != nil {
		return err
	}
	return nil
}

func (r *GetIncidentTemplatesRequest) ResourcePath() string {
	return "/v1/services/" + r.ServiceId + "/incident-templates"
}

func (r *GetIncidentTemplatesRequest) Method() string {
	return http.MethodGet
}

func validateIncidentTemplate(template IncidentTemplateRequest) error {
	if template.Name == "" {
		return errors.New("Name of incident template cannot be empty.")
	}
	err := validateIncidentProperties(template.IncidentProperties)
	if err != nil {
		return err
	}
	return nil
}

type IncidentTemplateRequest struct {
	Name               string             `json:"name"`
	IncidentProperties IncidentProperties `json:"incidentProperties"`
}
