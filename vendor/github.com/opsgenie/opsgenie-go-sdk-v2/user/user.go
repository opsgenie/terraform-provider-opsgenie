package user

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

func (c *Client) ListUserEscalations(context context.Context, request *ListUserEscalationsRequest) (*ListUserEscalationsResult, error) {
	result := &ListUserEscalationsResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) ListUserTeams(context context.Context, request *ListUserTeamsRequest) (*ListUserTeamsResult, error) {
	result := &ListUserTeamsResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) ListUserForwardingRules(context context.Context, request *ListUserForwardingRulesRequest) (*ListUserForwardingRulesResult, error) {
	result := &ListUserForwardingRulesResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) ListUserSchedules(context context.Context, request *ListUserSchedulesRequest) (*ListUserSchedulesResult, error) {
	result := &ListUserSchedulesResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (c *Client) GetSavedSearch(context context.Context, request *GetSavedSearchRequest) (*GetSavedSearchResult, error) {
	result := &GetSavedSearchResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (c *Client) ListSavedSearches(context context.Context, request *ListSavedSearchesRequest) (*ListSavedSearchesResult, error) {
	result := &ListSavedSearchesResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (c *Client) DeleteSavedSearch(context context.Context, request *DeleteSavedSearchRequest) (*DeleteSavedSearchResult, error) {
	result := &DeleteSavedSearchResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
