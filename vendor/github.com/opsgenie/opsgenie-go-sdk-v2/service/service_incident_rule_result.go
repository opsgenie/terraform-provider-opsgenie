package service

import (
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/og"
)

type CreateIncidentRuleResult struct {
	client.ResultMetadata
	Id string `json:"id"`
}

type UpdateIncidentRuleResult struct {
	client.ResultMetadata
	Id string `json:"id"`
}

type DeleteIncidentRuleResult struct {
	client.ResultMetadata
	Result string `json:"result"`
}

type GetIncidentRulesResult struct {
	client.ResultMetadata
	IncidentRule []IncidentRuleResult `json:"data,omitempty"`
}
type IncidentRuleResult struct {
	Id                 string                `json:"id"`
	Order              int                   `json:"order,omitempty"`
	ConditionMatchType og.ConditionMatchType `json:"conditionMatchType,omitempty"`
	Conditions         []og.Condition        `json:"conditions,omitempty"`
	IncidentProperties IncidentProperties    `json:"incidentProperties,omitempty"`
}
