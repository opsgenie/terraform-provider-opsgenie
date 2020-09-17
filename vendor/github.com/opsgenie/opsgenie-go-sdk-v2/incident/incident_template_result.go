package incident

import "github.com/opsgenie/opsgenie-go-sdk-v2/client"

type TemplateIncident struct {
	Name                  string                `json:"name"`
	IncidentTemplateId    string                `json:"id"`
	Message               string                `json:"message"`
	Description           string                `json:"description,omitempty"`
	Tags                  []string              `json:"tags,omitempty"`
	Details               map[string]string     `json:"details,omitempty"`
	Priority              Priority              `json:"priority"`
	ImpactedServices      []string              `json:"impactedServices,omitempty"`
	StakeholderProperties StakeholderProperties `json:"stakeholderProperties"`
}

type CreateIncidentTemplateResult struct {
	client.ResultMetadata
	Result             string `json:"result"`
	IncidentTemplateId string `json:"id"`
}

type UpdateIncidentTemplateResult struct {
	client.ResultMetadata
	Result             string `json:"result"`
	IncidentTemplateId string `json:"id"`
}

type DeleteIncidentTemplateResult struct {
	client.ResultMetadata
	Result string `json:"result"`
}

type GetIncidentTemplateResult struct {
	client.ResultMetadata
	IncidentTemplates map[string][]TemplateIncident `json:"data"`
	Paging            Paging                        `json:"paging"`
}
