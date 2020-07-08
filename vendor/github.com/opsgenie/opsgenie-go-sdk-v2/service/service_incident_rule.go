package service

import (
	"context"
)

func (c *Client) CreateIncidentRule(context context.Context, request *CreateIncidentRuleRequest) (*CreateIncidentRuleResult, error) {
	result := &CreateIncidentRuleResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) GetIncidentRules(context context.Context, request *GetIncidentRulesRequest) (*GetIncidentRulesResult, error) {
	result := &GetIncidentRulesResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) DeleteIncidentRule(context context.Context, request *DeleteIncidentRuleRequest) (*DeleteIncidentRuleResult, error) {
	result := &DeleteIncidentRuleResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) UpdateIncidentRule(context context.Context, request *UpdateIncidentRuleRequest) (*UpdateIncidentRuleResult, error) {
	result := &UpdateIncidentRuleResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
