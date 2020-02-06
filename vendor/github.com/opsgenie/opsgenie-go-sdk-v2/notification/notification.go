package notification

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
	return &Client{client: opsgenieClient}, nil
}
func (c *Client) CreateRuleStep(context context.Context, request *CreateRuleStepRequest) (*CreateRuleStepResult, error) {
	result := &CreateRuleStepResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) GetRuleStep(context context.Context, request *GetRuleStepRequest) (*GetRuleStepResult, error) {
	result := &GetRuleStepResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) UpdateRuleStep(context context.Context, request *UpdateRuleStepRequest) (*UpdateRuleStepResult, error) {
	result := &UpdateRuleStepResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) DeleteRuleStep(context context.Context, request *DeleteRuleStepRequest) (*DeleteRuleStepResult, error) {
	result := &DeleteRuleStepResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) ListRuleStep(context context.Context, request *ListRuleStepsRequest) (*ListRuleStepResult, error) {
	result := &ListRuleStepResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) EnableRuleStep(context context.Context, request *EnableRuleStepRequest) (*EnableRuleStepResult, error) {
	result := &EnableRuleStepResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) DisableRuleStep(context context.Context, request *DisableRuleStepRequest) (*DisableRuleStepResult, error) {
	result := &DisableRuleStepResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) CreateRule(context context.Context, request *CreateRuleRequest) (*CreateRuleResult, error) {
	result := &CreateRuleResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (c *Client) GetRule(context context.Context, request *GetRuleRequest) (*GetRuleResult, error) {
	result := &GetRuleResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) UpdateRule(context context.Context, request *UpdateRuleRequest) (*UpdateRuleResult, error) {
	result := &UpdateRuleResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) DeleteRule(context context.Context, request *DeleteRuleRequest) (*DeleteRuleResult, error) {
	result := &DeleteRuleResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) ListRule(context context.Context, request *ListRuleRequest) (*ListRuleResult, error) {
	result := &ListRuleResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) EnableRule(context context.Context, request *EnableRuleRequest) (*EnableRuleResult, error) {
	result := &EnableRuleResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) DisableRule(context context.Context, request *DisableRuleRequest) (*DisableRuleResult, error) {
	result := &DisableRuleResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) CopyRule(context context.Context, request *CopyNotificationRulesRequest) (*CopyNotificationRulesResult, error) {
	result := &CopyNotificationRulesResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
