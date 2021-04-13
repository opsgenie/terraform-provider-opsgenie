package policy

import (
	"github.com/opsgenie/opsgenie-go-sdk-v2/alert"
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
)

type CreateResult struct {
	client.ResultMetadata
	Id      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Enabled bool   `json:"enabled"`
}

type GetAlertPolicyResult struct {
	client.ResultMetadata
	MainFields
	Message                  string             `json:"message"`
	Continue                 bool               `json:"continue,omitempty"`
	Alias                    string             `json:"alias,omitempty"`
	AlertDescription         string             `json:"alertDescription,omitempty"`
	Entity                   string             `json:"entity,omitempty"`
	Source                   string             `json:"source,omitempty"`
	IgnoreOriginalDetails    bool               `json:"ignoreOriginalDetails,omitempty"`
	Actions                  []string           `json:"actions,omitempty"`
	IgnoreOriginalActions    bool               `json:"ignoreOriginalActions,omitempty"`
	Details                  interface{}        `json:"details,omitempty"`
	IgnoreOriginalResponders bool               `json:"ignoreOriginalResponders,omitempty"`
	Responders               *[]alert.Responder `json:"responders,omitempty"`
	IgnoreOriginalTags       bool               `json:"ignoreOriginalTags,omitempty"`
	Tags                     []string           `json:"tags,omitempty"`
	Priority                 alert.Priority     `json:"priority,omitempty"`
}

type GetNotificationPolicyResult struct {
	client.ResultMetadata
	MainFields
	AutoRestartAction         *AutoRestartAction   `json:"autoRestartAction,omitempty"`
	AutoCloseAction           *AutoCloseAction     `json:"autoCloseAction,omitempty"`
	DeDuplicationAction 	  *DeDuplicationAction `json:"deduplicationAction,omitempty"`
	DelayAction               *DelayAction         `json:"delayAction,omitempty"`
	Suppress                  bool                 `json:"suppress,omitempty"`
}

type PolicyResult struct {
	client.ResultMetadata
	Result string `json:"result,omitempty"`
}

type ListPolicyResult struct {
	client.ResultMetadata
	Policies []PolicyProps `json:"data,omitempty"`
}

type PolicyProps struct {
	Id      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Type    string `json:"type,omitempty"`
	Order   int    `json:"order,omitempty"`
	Enabled bool   `json:"enabled,omitempty"`
}
