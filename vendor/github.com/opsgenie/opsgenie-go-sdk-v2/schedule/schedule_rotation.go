package schedule

import "context"

func (c *Client) CreateRotation(context context.Context, request *CreateRotationRequest) (*CreateRotationResult, error) {
	result := &CreateRotationResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) GetRotation(context context.Context, request *GetRotationRequest) (*GetRotationResult, error) {
	result := &GetRotationResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) UpdateRotation(context context.Context, request *UpdateRotationRequest) (*UpdateRotationResult, error) {
	result := &UpdateRotationResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) DeleteRotation(context context.Context, request *DeleteRotationRequest) (*DeleteResult, error) {
	result := &DeleteResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) ListRotations(context context.Context, request *ListRotationsRequest) (*ListRotationsResult, error) {
	result := &ListRotationsResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
