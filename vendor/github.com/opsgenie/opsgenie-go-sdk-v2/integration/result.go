package integration

import (
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/og"
)

type ListResult struct {
	client.ResultMetadata
	Integrations []GenericFields `json:"data"`
}

type GenericFields struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
	Type    string `json:"type"`
	TeamId  string `json:"teamId"`
}

type GetResult struct {
	client.ResultMetadata
	Data map[string]interface{} `json:"data"`
}

type APIBasedIntegrationResult struct {
	client.ResultMetadata
	GenericFields
	ApiKey string `json:"apiKey"`
}

type WebhookIntegrationResult struct {
	client.ResultMetadata
	GenericFields
	ApiKey string `json:"apiKey"`
}

type EmailBasedIntegrationResult struct {
	client.ResultMetadata
	GenericFields
	EmailAddress string `json:"emailAddress"`
}

type UpdateResult struct {
	client.ResultMetadata
	Data map[string]interface{} `json:"data"`
}

type DeleteResult struct {
	client.ResultMetadata
	Result string `json:"result"`
}

type EnableResult struct {
	client.ResultMetadata
	GenericFields
}

type DisableResult struct {
	client.ResultMetadata
	GenericFields
}

type AuthenticateResult struct {
	client.ResultMetadata
	Result string `json:"result"`
}

type ActionsResult struct {
	client.ResultMetadata
	Parent      ParentIntegration   `json:"_parent"`
	Ignore      []IntegrationAction `json:"ignore"`
	Create      []IntegrationAction `json:"create"`
	Close       []IntegrationAction `json:"close"`
	Acknowledge []IntegrationAction `json:"acknowledge"`
	AddNote     []IntegrationAction `json:"addNote"`
}

type ParentIntegration struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
	Type    string `json:"type"`
}

type GenericActionFields struct {
	Type   string       `json:"type"`
	Name   string       `json:"name"`
	Order  int          `json:"order"`
	Filter FilterResult `json:"filter"`
}

type FilterResult struct {
	ConditionMatchType og.ConditionMatchType `json:"conditionMatchType,omitempty"`
	Conditions         []ConditionResult     `json:"conditions,omitempty"`
}

type ConditionResult struct {
	Field         og.ConditionFieldType `json:"field,omitempty"`
	IsNot         bool                  `json:"not,omitempty"`
	Operation     og.ConditionOperation `json:"operation,omitempty"`
	ExpectedValue string                `json:"expectedValue,omitempty"`
	Key           string                `json:"key,omitempty"`
	Order         *int                  `json:"order,omitempty"`
}

type CreateAction struct {
	GenericActionFields
	User                             string            `json:"user"`
	Note                             string            `json:"note"`
	Alias                            string            `json:"alias"`
	Source                           string            `json:"source"`
	Message                          string            `json:"message"`
	Description                      string            `json:"description"`
	Entity                           string            `json:"entity"`
	AppendAttachments                bool              `json:"appendAttachments"`
	IgnoreAlertActionsFromPayload    bool              `json:"ignoreAlertActionsFromPayload"`
	IgnoreRespondersFromPayload      bool              `json:"ignoreRespondersFromPayload"`
	IgnoreTagsFromPayload            bool              `json:"ignoreTagsFromPayload"`
	IgnoreExtraPropertiesFromPayload bool              `json:"ignoreExtraPropertiesFromPayload"`
	AlertActions                     []string          `json:"alertActions"`
	Responders                       []Responder       `json:"responders"`
	Tags                             []string          `json:"tags"`
	ExtraProperties                  map[string]string `json:"extraProperties"`
}

type CloseAction struct {
	GenericActionFields
	User  string `json:"user"`
	Note  string `json:"note"`
	Alias string `json:"alias"`
}

type AcknowledgeAction struct {
	GenericActionFields
	User  string `json:"user"`
	Note  string `json:"note"`
	Alias string `json:"alias"`
}

type AddNoteAction struct {
	GenericActionFields
	User  string `json:"user"`
	Note  string `json:"note"`
	Alias string `json:"alias"`
}

type IgnoreAction struct {
	GenericActionFields
}

type ResponderType string
type ActionType string

const (
	User       ResponderType = "user"
	Team       ResponderType = "team"
	Escalation ResponderType = "escalation"
	Schedule   ResponderType = "schedule"

	Create      ActionType = "create"
	Close       ActionType = "close"
	Acknowledge ActionType = "acknowledge"
	AddNote     ActionType = "AddNote"
	Ignore      ActionType = "ignore"
)

type Responder struct {
	Type     ResponderType `json:"type, omitempty"`
	Name     string        `json:"name,omitempty"`
	Id       string        `json:"id,omitempty"`
	Username string        `json:"username, omitempty"`
}
