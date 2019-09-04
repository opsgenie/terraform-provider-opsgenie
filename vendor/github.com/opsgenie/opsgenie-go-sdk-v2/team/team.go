package team

import (
	"context"
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
)

type Client struct {
	client *client.OpsGenieClient
}

func NewClient(config *client.Config) (*Client, error) {

	opsgenieClient, err := client.NewOpsGenieClient(config)

	if err != nil {
		return nil, err
	}

	return &Client{opsgenieClient}, nil
}

func (c *Client) Create(ctx context.Context, req *CreateTeamRequest) (*CreateTeamResult, error) {
	createTeamResponse := &CreateTeamResult{}

	err := c.client.Exec(ctx, req, createTeamResponse)
	if err != nil {
		return nil, err
	}

	return createTeamResponse, nil

}

func (c *Client) Get(ctx context.Context, req *GetTeamRequest) (*GetTeamResult, error) {

	getTeamResponse := &GetTeamResult{}

	err := c.client.Exec(ctx, req, getTeamResponse)
	if err != nil {
		return nil, err
	}

	return getTeamResponse, nil
}

func (c *Client) Update(ctx context.Context, req *UpdateTeamRequest) (*UpdateTeamResult, error) {

	updateTeamResponse := &UpdateTeamResult{}

	err := c.client.Exec(ctx, req, updateTeamResponse)
	if err != nil {
		return nil, err
	}

	return updateTeamResponse, nil
}

func (c *Client) Delete(ctx context.Context, req *DeleteTeamRequest) (*DeleteTeamResult, error) {

	deleteTeamResponse := &DeleteTeamResult{}

	err := c.client.Exec(ctx, req, deleteTeamResponse)
	if err != nil {
		return nil, err
	}

	return deleteTeamResponse, nil
}

func (c *Client) List(ctx context.Context, req *ListTeamRequest) (*ListTeamResult, error) {

	listTeamResponse := &ListTeamResult{}

	err := c.client.Exec(ctx, req, listTeamResponse)
	if err != nil {
		return nil, err
	}

	return listTeamResponse, nil
}

func (c *Client) ListTeamLogs(ctx context.Context, req *ListTeamLogsRequest) (*ListTeamLogsResult, error) {

	ListTeamLogsResponse := &ListTeamLogsResult{}

	err := c.client.Exec(ctx, req, ListTeamLogsResponse)
	if err != nil {
		return nil, err
	}

	return ListTeamLogsResponse, nil
}

//team role api
func (c *Client) CreateRole(ctx context.Context, req *CreateTeamRoleRequest) (*CreateTeamRoleResult, error) {

	createTeamRoleResponse := &CreateTeamRoleResult{}

	err := c.client.Exec(ctx, req, createTeamRoleResponse)
	if err != nil {
		return nil, err
	}

	return createTeamRoleResponse, nil
}

func (c *Client) GetRole(ctx context.Context, req *GetTeamRoleRequest) (*GetTeamRoleResult, error) {

	getTeamRoleResponse := &GetTeamRoleResult{}

	err := c.client.Exec(ctx, req, getTeamRoleResponse)
	if err != nil {
		return nil, err
	}

	return getTeamRoleResponse, nil
}

func (c *Client) UpdateRole(ctx context.Context, req *UpdateTeamRoleRequest) (*UpdateTeamRoleResult, error) {

	updateTeamRoleResponse := &UpdateTeamRoleResult{}

	err := c.client.Exec(ctx, req, updateTeamRoleResponse)
	if err != nil {
		return nil, err
	}

	return updateTeamRoleResponse, nil
}

func (c *Client) DeleteRole(ctx context.Context, req *DeleteTeamRoleRequest) (*DeleteTeamRoleResult, error) {

	deleteTeamRoleResponse := &DeleteTeamRoleResult{}

	err := c.client.Exec(ctx, req, deleteTeamRoleResponse)
	if err != nil {
		return nil, err
	}

	return deleteTeamRoleResponse, nil
}

func (c *Client) ListRole(ctx context.Context, req *ListTeamRoleRequest) (*ListTeamRoleResult, error) {

	listTeamRoleResponse := &ListTeamRoleResult{}

	err := c.client.Exec(ctx, req, listTeamRoleResponse)
	if err != nil {
		return nil, err
	}

	return listTeamRoleResponse, nil
}

//team member api
func (c *Client) AddMember(ctx context.Context, req *AddTeamMemberRequest) (*AddTeamMemberResult, error) {

	addTeamMemberResponse := &AddTeamMemberResult{}

	err := c.client.Exec(ctx, req, addTeamMemberResponse)
	if err != nil {
		return nil, err
	}

	return addTeamMemberResponse, nil
}

func (c *Client) RemoveMember(ctx context.Context, req *RemoveTeamMemberRequest) (*RemoveTeamMemberResult, error) {

	removeTeamMemberResponse := &RemoveTeamMemberResult{}

	err := c.client.Exec(ctx, req, removeTeamMemberResponse)
	if err != nil {
		return nil, err
	}

	return removeTeamMemberResponse, nil
}

//team routing rule api
func (c *Client) CreateRoutingRule(ctx context.Context, req *CreateRoutingRuleRequest) (*RoutingRuleResult, error) {

	createRoutingRuleResponse := &RoutingRuleResult{}

	err := c.client.Exec(ctx, req, createRoutingRuleResponse)
	if err != nil {
		return nil, err
	}

	return createRoutingRuleResponse, nil
}

func (c *Client) GetRoutingRule(ctx context.Context, req *GetRoutingRuleRequest) (*GetRoutingRuleResult, error) {

	getRoutingRuleResponse := &GetRoutingRuleResult{}

	err := c.client.Exec(ctx, req, getRoutingRuleResponse)
	if err != nil {
		return nil, err
	}

	return getRoutingRuleResponse, nil
}

func (c *Client) UpdateRoutingRule(ctx context.Context, req *UpdateRoutingRuleRequest) (*RoutingRuleResult, error) {

	updateRoutingRuleResponse := &RoutingRuleResult{}

	err := c.client.Exec(ctx, req, updateRoutingRuleResponse)
	if err != nil {
		return nil, err
	}

	return updateRoutingRuleResponse, nil
}

func (c *Client) DeleteRoutingRule(ctx context.Context, req *DeleteRoutingRuleRequest) (*DeleteRoutingRuleResult, error) {

	deleteRoutingRuleResponse := &DeleteRoutingRuleResult{}

	err := c.client.Exec(ctx, req, deleteRoutingRuleResponse)
	if err != nil {
		return nil, err
	}

	return deleteRoutingRuleResponse, nil
}

func (c *Client) ListRoutingRules(ctx context.Context, req *ListRoutingRulesRequest) (*ListRoutingRulesResult, error) {

	listRoutingRulesResponse := &ListRoutingRulesResult{}

	err := c.client.Exec(ctx, req, listRoutingRulesResponse)
	if err != nil {
		return nil, err
	}

	return listRoutingRulesResponse, nil
}

func (c *Client) ChangeRoutingRuleOrder(ctx context.Context, req *ChangeRoutingRuleOrderRequest) (*RoutingRuleResult, error) {

	changeRoutingRuleOrderResponse := &RoutingRuleResult{}

	err := c.client.Exec(ctx, req, changeRoutingRuleOrderResponse)
	if err != nil {
		return nil, err
	}

	return changeRoutingRuleOrderResponse, nil
}
