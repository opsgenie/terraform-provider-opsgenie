package service

import "context"

func (c *Client) GetAudienceTemplate(context context.Context, request *GetAudienceTemplateRequest) (*GetAudienceTemplateResult, error) {
	result := &GetAudienceTemplateResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) UpdateAudienceTemplate(context context.Context, request *UpdateAudienceTemplateRequest) (*UpdateAudienceTemplateResult, error) {
	result := &UpdateAudienceTemplateResult{}
	err := c.client.Exec(context, request, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
