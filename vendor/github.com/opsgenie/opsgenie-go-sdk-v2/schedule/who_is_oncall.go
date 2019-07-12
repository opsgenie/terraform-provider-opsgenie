package schedule

import (
	"context"
	"os"
)

func (c *Client) GetOnCalls(context context.Context, request *GetOnCallsRequest) (*GetOnCallsResult, error) {
	result := &GetOnCallsResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) GetNextOnCall(context context.Context, request *GetNextOnCallsRequest) (*GetNextOnCallsResult, error) {
	result := &GetNextOnCallsResult{}
	err := c.client.Exec(context, request, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) ExportOnCallUser(context context.Context, request *ExportOnCallUserRequest) (*os.File, error) {
	result := &exportOncallUserResult{}

	file, err := os.Create(request.ExportedFilePath + request.getFileName())
	if err != nil {
		return nil, err
	}

	defer file.Close()

	err = c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}

	_, err = file.Write(result.FileContent)
	if err != nil {
		return nil, err
	}
	return file, nil
}
