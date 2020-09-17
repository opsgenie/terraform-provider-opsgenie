package incident

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

func (c *Client) GetRequestStatus(context context.Context, request *RequestStatusRequest) (*RequestStatusResult, error) {
	result := &RequestStatusResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) Create(context context.Context, request *CreateRequest) (*AsyncResult, error) {
	result := &AsyncResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) Delete(context context.Context, request *DeleteRequest) (*AsyncResult, error) {
	result := &AsyncResult{}
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

func (c *Client) List(context context.Context, request *ListRequest) (*ListResult, error) {
	result := &ListResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) Close(context context.Context, request *CloseRequest) (*AsyncResult, error) {
	result := &AsyncResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) AddNote(context context.Context, request *AddNoteRequest) (*AsyncResult, error) {
	result := &AsyncResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) AddResponder(context context.Context, request *AddResponderRequest) (*AsyncResult, error) {
	result := &AsyncResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) AddTags(context context.Context, request *AddTagsRequest) (*AsyncResult, error) {
	result := &AsyncResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) RemoveTags(context context.Context, request *RemoveTagsRequest) (*AsyncResult, error) {
	result := &AsyncResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) AddDetails(context context.Context, request *AddDetailsRequest) (*AsyncResult, error) {
	result := &AsyncResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) RemoveDetails(context context.Context, request *RemoveDetailsRequest) (*AsyncResult, error) {
	result := &AsyncResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) UpdatePriority(context context.Context, request *UpdatePriorityRequest) (*AsyncResult, error) {
	result := &AsyncResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) UpdateMessage(context context.Context, request *UpdateMessageRequest) (*AsyncResult, error) {
	result := &AsyncResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) UpdateDescription(context context.Context, request *UpdateDescriptionRequest) (*AsyncResult, error) {
	result := &AsyncResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) ListLogs(context context.Context, request *ListLogsRequest) (*ListLogsResult, error) {
	result := &ListLogsResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) ListNotes(context context.Context, request *ListNotesRequest) (*ListNotesResult, error) {
	result := &ListNotesResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
