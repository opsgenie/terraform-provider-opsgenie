package team

import (
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/og"
)

type TeamMeta struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type ListedTeams struct {
	TeamMeta
	Description string `json:"description,omitempty"`
}

type RoutingRuleMeta struct {
	Id              string             `json:"id,omitempty"`
	Name            string             `json:"name,omitempty"`
	IsDefault       bool               `json:"isDefault,omitempty"`
	Criteria        og.Criteria        `json:"criteria,omitempty"`
	Timezone        string             `json:"timezone,omitempty"`
	TimeRestriction og.TimeRestriction `json:"timeRestriction,omitempty"`
	Notify          Notify             `json:"notify,omitempty"`
}

type CreateTeamResult struct {
	client.ResultMetadata
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type GetTeamResult struct {
	client.ResultMetadata
	TeamMeta
	Description string   `json:"description,omitempty"`
	Members     []Member `json:"members,omitempty"`
}

type UpdateTeamResult struct {
	client.ResultMetadata
	TeamMeta
}

type DeleteTeamResult struct {
	client.ResultMetadata
	Result string `json:"result"`
}

type ListTeamResult struct {
	client.ResultMetadata
	Teams []ListedTeams `json:"data"`
}

type LogEntry struct {
	Log         string `json:"log"`
	Owner       string `json:"owner"`
	CreatedDate string `json:"createdDate"`
}

type ListTeamLogsResult struct {
	client.ResultMetadata
	Offset string     `json:"offset,omitempty"`
	Logs   []LogEntry `json:logs,omitempty`
}

//team role api
type RoleMeta struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type RightMeta struct {
	Right   string `json:"right,omitempty"`
	Granted bool   `json:"granted,omitempty"`
}

type GetRoleInfo struct {
	RoleMeta
	Rights []Right `json:"rights"`
}

type CreateTeamRoleResult struct {
	client.ResultMetadata
	RoleMeta
}

type GetTeamRoleResult struct {
	client.ResultMetadata
	RoleMeta
	Rights []Right `json:"rights"`
}

type UpdateTeamRoleResult struct {
	client.ResultMetadata
	RoleMeta
}

type DeleteTeamRoleResult struct {
	client.ResultMetadata
	Result string `json:"result"`
}

type ListTeamRoleResult struct {
	client.ResultMetadata
	TeamRoles []GetRoleInfo `json:"data"`
}

type AddTeamMemberResult struct {
	client.ResultMetadata
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type RemoveTeamMemberResult struct {
	client.ResultMetadata
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type RoutingRuleResult struct {
	client.ResultMetadata
	Id string `json:"id,omitempty"`
}

type GetRoutingRuleResult struct {
	client.ResultMetadata
	RoutingRuleMeta
}

type DeleteRoutingRuleResult struct {
	client.ResultMetadata
	Result string `json:"result"`
}

type ListRoutingRulesResult struct {
	client.ResultMetadata
	RoutingRules []RoutingRuleMeta `json:"data"`
}
