package policy

import (
	"context"
	"errors"
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

func (c *Client) CreateAlertPolicy(context context.Context, request *CreateAlertPolicyRequest) (*CreateResult, error) {
	request.PolicyType = "alert"
	result := &CreateResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) CreateNotificationPolicy(context context.Context, request *CreateNotificationPolicyRequest) (*CreateResult, error) {
	request.PolicyType = "notification"
	result := &CreateResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) GetAlertPolicy(context context.Context, request *GetAlertPolicyRequest) (*GetAlertPolicyResult, error) {
	result := &GetAlertPolicyResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	if result.PolicyType != "alert" {
		return nil, errors.New("policy type is not alert")
	}
	return result, nil
}

func (c *Client) GetNotificationPolicy(context context.Context, request *GetNotificationPolicyRequest) (*GetNotificationPolicyResult, error) {
	result := &GetNotificationPolicyResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	if result.PolicyType != "notification" {
		return nil, errors.New("policy type is not notification")
	}
	return result, nil
}

func (c *Client) UpdateAlertPolicy(context context.Context, request *UpdateAlertPolicyRequest) (*PolicyResult, error) {
	request.PolicyType = "alert"
	result := &PolicyResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) UpdateNotificationPolicy(context context.Context, request *UpdateNotificationPolicyRequest) (*PolicyResult, error) {
	request.PolicyType = "notification"
	result := &PolicyResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) DeletePolicy(context context.Context, request *DeletePolicyRequest) (*PolicyResult, error) {
	result := &PolicyResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) DisablePolicy(context context.Context, request *DisablePolicyRequest) (*PolicyResult, error) {
	result := &PolicyResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) EnablePolicy(context context.Context, request *EnablePolicyRequest) (*PolicyResult, error) {
	result := &PolicyResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) ChangeOrder(context context.Context, request *ChangeOrderRequest) (*PolicyResult, error) {
	result := &PolicyResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) ListAlertPolicies(context context.Context, request *ListAlertPoliciesRequest) (*ListPolicyResult, error) {
	result := &ListPolicyResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) ListNotificationPolicies(context context.Context, request *ListNotificationPoliciesRequest) (*ListPolicyResult, error) {
	result := &ListPolicyResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
