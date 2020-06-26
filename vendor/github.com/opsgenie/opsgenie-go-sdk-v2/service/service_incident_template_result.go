package service

import "github.com/opsgenie/opsgenie-go-sdk-v2/client"

type CreateIncidentTemplateResult struct {
	client.ResultMetadata
	Id string `json:"id"`
}

type UpdateIncidentTemplateResult struct {
	client.ResultMetadata
	Id string `json:"id"`
}

type DeleteIncidentTemplateResult struct {
	client.ResultMetadata
	Result string `json:"result"`
}

type GetIncidentTemplatesResult struct {
	client.ResultMetadata
	IncidentTemplates []IncidentTemplate `json:"data"`
}

type IncidentTemplate struct {
	Id                 string             `json:"id"`
	Name               string             `json:"name"`
	IncidentProperties IncidentProperties `json:"incidentProperties"`
}
