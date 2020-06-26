package service

import "github.com/opsgenie/opsgenie-go-sdk-v2/client"

type UpdateAudienceTemplateResult struct {
	client.ResultMetadata
	Result string `json:"result"`
}

type GetAudienceTemplateResult struct {
	client.ResultMetadata
	Responder   ResponderOfAudience   `json:"responder"`
	Stakeholder StakeholderOfAudience `json:"stakeholder"`
}
