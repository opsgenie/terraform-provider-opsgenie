package contact

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
func (c *Client) Create(context context.Context, request *CreateRequest) (*CreateResult, error) {
	result := &CreateResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) Get(context context.Context, request *GetRequest) (*GetResult, error) {
	result := &GetResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) Update(context context.Context, request *UpdateRequest) (*UpdateResult, error) {
	result := &UpdateResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) Delete(context context.Context, request *DeleteRequest) (*DeleteResult, error) {
	result := &DeleteResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) List(context context.Context, request *ListRequest) (*ListResult, error) {
	result := &ListResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) Enable(context context.Context, request *EnableRequest) (*EnableResult, error) {
	result := &EnableResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) Disable(context context.Context, request *DisableRequest) (*DisableResult, error) {
	result := &DisableResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
