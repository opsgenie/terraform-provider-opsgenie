package notification

import (
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/og"
)

type Parent struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type RuleStep struct {
	Parent    Parent       `json:"_parent,omitempty"`
	Id        string       `json:"id,omitempty"`
	SendAfter og.SendAfter `json:"sendAfter,omitempty"`
	Contact   og.Contact   `json:"contact,omitempty"`
	Enabled   bool         `json:"enabled,omitempty"`
}

type CreateRuleStepResult struct {
	client.ResultMetadata
	Id string `json:"id,omitempty"`
}

type GetRuleStepResult struct {
	client.ResultMetadata
	RuleStep RuleStep `json:"data,omitempty"`
}
type UpdateRuleStepResult struct {
	client.ResultMetadata
	Id string `json:"id,omitempty"`
}

type DeleteRuleStepResult struct {
	client.ResultMetadata
	Result string `json:"result,omitempty"`
}

type ListRuleStepResult struct {
	client.ResultMetadata
	RuleSteps []RuleStep `json:"data,omitempty"`
}

type EnableRuleStepResult struct {
	client.ResultMetadata
	Id string `json:"id,omitempty"`
}

type DisableRuleStepResult struct {
	client.ResultMetadata
	Id string `json:"id,omitempty"`
}

type SimpleNotificationRuleResult struct {
	Id         string     `json:"id,omitempty"`
	Name       string     `json:"name,omitempty"`
	ActionType ActionType `json:"actionType,omitempty"`
	Order      uint32     `json:"order,omitempty"`
	Enabled    bool       `json:"bool,omitempty"`
}

type CreateRuleResult struct {
	client.ResultMetadata
	SimpleNotificationRule SimpleNotificationRuleResult `json:"data,omitempty"`
}

type GetRuleResult struct {
	client.ResultMetadata
	Id               string                 `json:"id,omitempty"`
	Name             string                 `json:"name,omitempty"`
	ActionType       ActionType             `json:"actionType,omitempty"`
	Order            uint32                 `json:"order,omitempty"`
	Enabled          bool                   `json:"enabled,omitempty"`
	NotificationTime []NotificationTimeType `json:"notificationTime,omitempty"`
	TimeRestriction  *og.TimeRestriction    `json:"timeRestriction,omitempty"`
	Steps            []*StepResult          `json:"steps,omitempty"`
	Schedules        []*Schedule            `json:"schedules,omitempty"`
}

type StepResult struct {
	Contact   og.Contact    `json:"contact,omitempty"`
	SendAfter *og.SendAfter `json:"sendAfter,omitempty"`
	Enabled   bool          `json:"enabled,omitempty"`
}

type UpdateRuleResult struct {
	client.ResultMetadata
	SimpleNotificationRule SimpleNotificationRuleResult `json:"data,omitempty"`
}
type DeleteRuleResult struct {
	client.ResultMetadata
	Result string `json:"result,omitempty"`
}

type ListRuleResult struct {
	client.ResultMetadata
	SimpleNotificationRules []SimpleNotificationRuleResult `json:"data,omitempty"`
}

type EnableRuleResult struct {
	client.ResultMetadata
	Id string `json:"id,omitempty"`
}

type DisableRuleResult struct {
	client.ResultMetadata
	Id string `json:"id,omitempty"`
}

type CopyNotificationRulesResult struct {
	client.ResultMetadata
	Result string `json:"result,omitempty"`
}
