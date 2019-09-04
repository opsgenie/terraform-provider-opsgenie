package schedule

import (
	"context"
)

func (c *Client) CreateScheduleOverride(context context.Context, request *CreateScheduleOverrideRequest) (*CreateScheduleOverrideResult, error) {
	result := &CreateScheduleOverrideResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) GetScheduleOverride(context context.Context, request *GetScheduleOverrideRequest) (*GetScheduleOverrideResult, error) {
	result := &GetScheduleOverrideResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) ListScheduleOverride(context context.Context, request *ListScheduleOverrideRequest) (*ListScheduleOverrideResult, error) {
	result := &ListScheduleOverrideResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) DeleteScheduleOverride(context context.Context, request *DeleteScheduleOverrideRequest) (*DeleteScheduleOverrideResult, error) {
	result := &DeleteScheduleOverrideResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) UpdateScheduleOverride(context context.Context, request *UpdateScheduleOverrideRequest) (*UpdateScheduleOverrideResult, error) {
	result := &UpdateScheduleOverrideResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
