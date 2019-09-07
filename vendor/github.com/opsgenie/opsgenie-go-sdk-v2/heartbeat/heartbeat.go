package heartbeat

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

func (c *Client) Ping(context context.Context, heartbeatName string) (*PingResult, error) {
	pingResult := &PingResult{}
	request := &pingRequest{HeartbeatName: heartbeatName}
	err := c.client.Exec(context, request, pingResult)
	if err != nil {
		return nil, err
	}
	return pingResult, nil
}

func (c *Client) Get(context context.Context, heartbeatName string) (*GetResult, error) {
	getResult := &GetResult{}
	request := &getRequest{HeartbeatName: heartbeatName}
	err := c.client.Exec(context, request, getResult)
	if err != nil {
		return nil, err
	}
	return getResult, nil
}

func (c *Client) List(context context.Context) (*ListResult, error) {
	listResult := &ListResult{}
	request := &listRequest{}
	err := c.client.Exec(context, request, listResult)
	if err != nil {
		return nil, err
	}
	return listResult, nil
}

func (c *Client) Update(context context.Context, request *UpdateRequest) (*HeartbeatInfo, error) {
	updateResult := &HeartbeatInfo{}
	err := c.client.Exec(context, request, updateResult)
	if err != nil {
		return nil, err
	}
	return updateResult, nil
}

func (c *Client) Add(context context.Context, request *AddRequest) (*AddResult, error) {
	result := &AddResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) Enable(context context.Context, heartbeatName string) (*HeartbeatInfo, error) {
	result := &HeartbeatInfo{}
	request := &enableRequest{heartbeatName: heartbeatName}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) Disable(context context.Context, heartbeatName string) (*HeartbeatInfo, error) {
	result := &HeartbeatInfo{}
	request := &disableRequest{heartbeatName: heartbeatName}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) Delete(context context.Context, heartbeatName string) (*DeleteResult, error) {
	deleteResult := &DeleteResult{}
	request := &deleteRequest{HeartbeatName: heartbeatName}
	err := c.client.Exec(context, request, deleteResult)
	if err != nil {
		return nil, err
	}
	return deleteResult, nil
}
