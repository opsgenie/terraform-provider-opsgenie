package integration

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

func (c *Client) Get(context context.Context, request *GetRequest) (*GetResult, error) {
	result := &GetResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) List(context context.Context) (*ListResult, error) {
	request := listRequest{}
	result := &ListResult{}
	err := c.client.Exec(context, &request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) CreateApiBased(context context.Context, request *APIBasedIntegrationRequest) (*APIBasedIntegrationResult, error) {
	result := &APIBasedIntegrationResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) CreateWebhook(context context.Context, request *WebhookIntegrationRequest) (*WebhookIntegrationResult, error) {
	result := &WebhookIntegrationResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) CreateEmailBased(context context.Context, request *EmailBasedIntegrationRequest) (*EmailBasedIntegrationResult, error) {
	result := &EmailBasedIntegrationResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) ForceUpdateAllFields(context context.Context, request *UpdateIntegrationRequest) (*UpdateResult, error) {
	result := &UpdateResult{}
	if len(request.OtherFields) == 0 {
		request.OtherFields = map[string]interface{}{}
	}
	request.OtherFields["id"] = request.Id
	request.OtherFields["name"] = request.Name
	request.OtherFields["type"] = request.Type
	request.OtherFields["enabled"] = request.Enabled
	request.OtherFields["ignoreRespondersFromPayload"] = request.IgnoreRespondersFromPayload
	request.OtherFields["suppressNotifications"] = request.SuppressNotifications
	request.OtherFields["responders"] = request.Responders
	request.OtherFields["emailUsername"] = request.EmailUsername
	request.OtherFields["url"] = request.WebhookUrl
	request.OtherFields["addAlertDescription"] = request.AddAlertDescription
	request.OtherFields["addAlertDetails"] = request.AddAlertDetails
	request.OtherFields["headers"] = request.Headers
	err := c.client.Exec(context, request.OtherFields, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) Delete(context context.Context, request *DeleteIntegrationRequest) (*DeleteResult, error) {
	result := &DeleteResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) Enable(context context.Context, request *EnableIntegrationRequest) (*EnableResult, error) {
	result := &EnableResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) Disable(context context.Context, request *DisableIntegrationRequest) (*DisableResult, error) {
	result := &DisableResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) Authenticate(context context.Context, request *AuthenticateIntegrationRequest) (*AuthenticateResult, error) {
	result := &AuthenticateResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) GetActions(context context.Context, request *GetIntegrationActionsRequest) (*ActionsResult, error) {
	result := &ActionsResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) CreateActions(context context.Context, request *CreateIntegrationActionsRequest) (*ActionsResult, error) {
	result := &ActionsResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) UpdateAllActions(context context.Context, request *UpdateAllIntegrationActionsRequest) (*ActionsResult, error) {
	result := &ActionsResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
