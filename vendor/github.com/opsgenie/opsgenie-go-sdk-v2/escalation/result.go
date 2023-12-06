package escalation

import (
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/og"
)

type CreateResult struct {
	client.ResultMetadata
	Result  string            `json:"result,omitempty"`
	Message string            `json:"message,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
	Id      string            `json:"id,omitempty"`
	Name    string            `json:"name,omitempty"`
}

type UpdateResult struct {
	client.ResultMetadata
	Result  string            `json:"result,omitempty"`
	Message string            `json:"message,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
	Id      string            `json:"id,omitempty"`
	Name    string            `json:"name,omitempty"`
}

type DeleteResult struct {
	client.ResultMetadata
	Result  string `json:"result,omitempty"`
	Message string `json:"message,omitempty"`
}

type GetResult struct {
	client.ResultMetadata
	Escalation
}

type ListResult struct {
	client.ResultMetadata
	Escalations []Escalation `json:"data,omitempty"`
}

type Escalation struct {
	Id          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description,omitempty"`
	Rules       []Rule        `json:"rules,omitempty"`
	OwnerTeam   *og.OwnerTeam `json:"ownerTeam,omitempty"`
	Repeat      *Repeat       `json:"repeat,omitempty"`
}

type Repeat struct {
	WaitInterval         uint32 `json:"waitInterval,omitempty"`
	Count                uint32 `json:"count,omitempty"`
	ResetRecipientStates bool   `json:"resetRecipientStates,omitempty"`
	CloseAlertAfterAll   bool   `json:"closeAlertAfterAll,omitempty"`
}

type Rule struct {
	Condition  og.EscalationCondition `json:"condition,omitempty"`
	NotifyType og.NotifyType          `json:"notifyType,omitempty"`
	Recipient  og.Participant         `json:"recipient,omitempty"`
	Delay      EscalationDelay        `json:"delay,omitempty"`
}

type EscalationDelay struct {
	TimeUnit   og.TimeUnit `json:"timeUnit,omitempty"`
	TimeAmount uint32      `json:"timeAmount"`
}
